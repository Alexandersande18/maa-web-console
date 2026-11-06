package maa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"maa-backend/internal/model/dto"
)

type queueRunner interface {
	Execute(context.Context, []dto.Tasker) error
	Stop(context.Context) error
	CloseGame(context.Context) error
	Notify(string)
}
type frozenTask struct {
	typ    dto.TaskType
	params string
}

func (t frozenTask) GetType() dto.TaskType { return t.typ }
func (t frozenTask) Serialize() string     { return t.params }

type queuedJob struct {
	id     string
	tasks  []dto.Tasker
	cancel context.CancelFunc
	done   chan struct{}
}
type taskQueue struct {
	mu                   sync.Mutex
	runners              map[string]queueRunner
	pending              []*queuedJob
	active               *queuedJob
	wake                 chan struct{}
	timeout, stopTimeout time.Duration
	halted               bool
	fatal                func(error)
}

func newTaskQueue(runners map[string]queueRunner, timeout, stopTimeout time.Duration, fatal func(error)) *taskQueue {
	q := &taskQueue{runners: runners, timeout: timeout, stopTimeout: stopTimeout, wake: make(chan struct{}, 1), fatal: fatal}
	go q.loop()
	return q
}
func (q *taskQueue) Submit(id string, tasks []dto.Tasker) error {
	if q.runners[id] == nil {
		return fmt.Errorf("instance not found: %s", id)
	}
	if len(tasks) == 0 {
		return errors.New("task list is empty")
	}
	snapshot := make([]dto.Tasker, 0, len(tasks))
	for _, t := range tasks {
		if t == nil {
			return errors.New("nil task")
		}
		params := t.Serialize()
		if t.GetType().String() == "" || !json.Valid([]byte(params)) {
			return errors.New("invalid task parameters")
		}
		snapshot = append(snapshot, frozenTask{typ: t.GetType(), params: params})
	}
	q.mu.Lock()
	if q.halted {
		q.mu.Unlock()
		return errors.New("task queue is shutting down")
	}
	if q.active != nil && q.active.id == id {
		q.mu.Unlock()
		return errors.New("该账号已有正在执行的任务，请先停止后再提交")
	}
	for _, j := range q.pending {
		if j.id == id {
			q.mu.Unlock()
			return errors.New("该账号已有排队任务，请先取消后再提交")
		}
	}
	j := &queuedJob{id: id, tasks: snapshot, done: make(chan struct{})}
	q.pending = append(q.pending, j)
	position := len(q.pending)
	if q.active != nil {
		position++
	}
	// Announce before waking the worker, so accepted always precedes started.
	q.runners[id].Notify(fmt.Sprintf("任务已加入队列，当前顺位 %d；最长执行时间 %s（从实际执行开始计时）", position, q.timeout))
	q.mu.Unlock()
	select {
	case q.wake <- struct{}{}:
	default:
	}
	return nil
}
func (q *taskQueue) Cancel(id string) error {
	q.mu.Lock()
	for i, j := range q.pending {
		if j.id == id {
			q.pending = append(q.pending[:i], q.pending[i+1:]...)
			close(j.done)
			q.mu.Unlock()
			q.runners[id].Notify("排队任务已取消")
			return nil
		}
	}
	j := q.active
	if j == nil || j.id != id {
		q.mu.Unlock()
		return nil
	}
	j.cancel()
	q.mu.Unlock()
	select {
	case <-j.done:
		return nil
	case <-time.After(q.stopTimeout + time.Second):
		return errors.New("停止任务超时，正在执行故障恢复")
	}
}
func (q *taskQueue) Status(id string) string {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.halted {
		return "任务队列已暂停，等待服务恢复"
	}
	if q.active != nil && q.active.id == id {
		return "当前任务正在执行"
	}
	for i, j := range q.pending {
		if j.id == id {
			return fmt.Sprintf("任务正在排队，等待顺位 %d", i+1)
		}
	}
	return "没有正在执行或排队的任务"
}
func (q *taskQueue) loop() {
	for range q.wake {
		for {
			q.mu.Lock()
			if q.halted || len(q.pending) == 0 {
				q.mu.Unlock()
				break
			}
			j := q.pending[0]
			q.pending = q.pending[1:]
			ctx, cancel := context.WithTimeout(context.Background(), q.timeout)
			j.cancel = cancel
			q.active = j
			q.mu.Unlock()
			r := q.runners[j.id]
			r.Notify("任务开始执行")
			err, cleanupErr := q.execute(ctx, r, j.tasks)
			cancel()
			if cleanupErr != nil {
				r.Notify(fmt.Sprintf("无法确认 MaaCore 已停止，暂停队列并重启服务：%v", cleanupErr))
				q.mu.Lock()
				q.halted = true
				close(j.done)
				q.mu.Unlock()
				q.fatal(cleanupErr)
				return // Never dispatch another job while a native thread may still run.
			}
			if errors.Is(err, context.DeadlineExceeded) {
				r.Notify("任务执行超时，已停止任务并关闭游戏")
			} else if errors.Is(err, context.Canceled) {
				r.Notify("任务已停止")
			} else if err != nil {
				r.Notify(fmt.Sprintf("任务执行失败：%v", err))
			} else {
				r.Notify("本次任务执行结束")
			}
			q.mu.Lock()
			q.active = nil
			close(j.done)
			q.mu.Unlock()
		}
	}
}
func (q *taskQueue) execute(ctx context.Context, r queueRunner, tasks []dto.Tasker) (error, error) {
	result := make(chan error, 1)
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		defer func() {
			if v := recover(); v != nil {
				result <- fmt.Errorf("task panic: %v", v)
			}
		}()
		result <- r.Execute(ctx, tasks)
	}()
	var err error
	select {
	case err = <-result:
	case <-ctx.Done():
		err = ctx.Err()
	}
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	cleanupCtx, cancel := context.WithTimeout(context.Background(), q.stopTimeout)
	defer cancel()
	cleaned := make(chan error, 1)
	go func() {
		// ADB force-stop runs independently of the ABI mutex, even if a DLL
		// call is stuck. Only the assigned emulator and configured package are used.
		stopped := make(chan error, 1)
		go func() { stopped <- r.Stop(cleanupCtx) }()
		var closeErr error
		timedOut := errors.Is(err, context.DeadlineExceeded)
		if timedOut {
			closeErr = r.CloseGame(cleanupCtx)
		}
		var stopErr error
		select {
		case stopErr = <-stopped:
		case <-cleanupCtx.Done():
			stopErr = cleanupCtx.Err()
		}
		if stopErr == nil {
			select {
			case <-finished:
			case <-cleanupCtx.Done():
				stopErr = cleanupCtx.Err()
			}
		}
		// A canceled Execute may have already been waiting for the ABI lock
		// when the first Stop ran. Drain it, then clear/stop once more so no
		// late append/start can survive into this account's next queued batch.
		if err != nil && stopErr == nil {
			stopErr = r.Stop(cleanupCtx)
		}
		// Core may have been restarting the game when the first force-stop ran.
		// Once no task can execute any more, enforce the final closed state.
		if timedOut && stopErr == nil {
			closeErr = r.CloseGame(cleanupCtx)
		}
		cleaned <- errors.Join(closeErr, stopErr)
	}()
	select {
	case cleanupErr := <-cleaned:
		return err, cleanupErr
	case <-cleanupCtx.Done():
		return err, cleanupCtx.Err()
	}
}
