package maa

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"maa-backend/internal/model/dto"
)

type fakeRunner struct {
	started   chan struct{}
	release   chan struct{}
	exited    chan struct{}
	stops     atomic.Int32
	closes    atomic.Int32
	running   *atomic.Int32
	overlap   *atomic.Bool
	blockStop bool
	params    chan string
}

func fake() *fakeRunner {
	return &fakeRunner{started: make(chan struct{}), release: make(chan struct{}), exited: make(chan struct{}), params: make(chan string, 1)}
}
func (f *fakeRunner) Execute(ctx context.Context, tasks []dto.Tasker) error {
	if f.running != nil {
		if f.running.Add(1) != 1 {
			f.overlap.Store(true)
		}
		defer f.running.Add(-1)
	}
	defer close(f.exited)
	f.params <- tasks[0].Serialize()
	close(f.started)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-f.release:
		return nil
	}
}
func (f *fakeRunner) Stop(ctx context.Context) error {
	f.stops.Add(1)
	if f.blockStop {
		<-ctx.Done()
		return ctx.Err()
	}
	return nil
}
func (f *fakeRunner) CloseGame(context.Context) error { f.closes.Add(1); return nil }
func (f *fakeRunner) Notify(string)                   {}
func taskList() []dto.Tasker {
	return []dto.Tasker{&dto.TaskFight{Enable: true, Stage: "1-7", Times: 1, Series: 1}}
}
func await(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out")
	}
}
func jobDone(q *taskQueue, id string) <-chan struct{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.active != nil && q.active.id == id {
		return q.active.done
	}
	for _, j := range q.pending {
		if j.id == id {
			return j.done
		}
	}
	done := make(chan struct{})
	close(done)
	return done
}
func testQueue(r map[string]queueRunner, d time.Duration) *taskQueue {
	return newTaskQueue(r, d, time.Second, func(err error) { panic(err) })
}
func TestQueueSerializesCompleteRuns(t *testing.T) {
	a, b := fake(), fake()
	var active atomic.Int32
	var overlap atomic.Bool
	a.running, b.running = &active, &active
	a.overlap, b.overlap = &overlap, &overlap
	q := testQueue(map[string]queueRunner{"a": a, "b": b}, time.Minute)
	if err := q.Submit("a", taskList()); err != nil {
		t.Fatal(err)
	}
	await(t, a.started)
	aDone := jobDone(q, "a")
	if err := q.Submit("b", taskList()); err != nil {
		t.Fatal(err)
	}
	bDone := jobDone(q, "b")
	select {
	case <-b.started:
		t.Fatal("second account started while first runs")
	default:
	}
	close(a.release)
	await(t, aDone)
	await(t, b.started)
	close(b.release)
	await(t, bDone)
	if overlap.Load() {
		t.Fatal("native runs overlapped")
	}
}
func TestQueueRejectsDuplicateAndSnapshotsParameters(t *testing.T) {
	a, b := fake(), fake()
	q := testQueue(map[string]queueRunner{"a": a, "b": b}, time.Minute)
	_ = q.Submit("a", taskList())
	await(t, a.started)
	if err := q.Submit("a", taskList()); err == nil {
		t.Fatal("accepted duplicate running task")
	}
	tasks := taskList()
	_ = q.Submit("b", tasks)
	if err := q.Submit("b", taskList()); err == nil {
		t.Fatal("accepted duplicate pending task")
	}
	tasks[0].(*dto.TaskFight).Stage = "mutated"
	close(a.release)
	await(t, b.started)
	if value := <-b.params; value == "" || !containsStage(value, "1-7") {
		t.Fatalf("parameters were mutated: %s", value)
	}
	done := jobDone(q, "b")
	close(b.release)
	await(t, done)
}
func containsStage(value, stage string) bool { return strings.Contains(value, `"stage":"`+stage+`"`) }
func TestCancelPendingDoesNotStopActive(t *testing.T) {
	a, b := fake(), fake()
	q := testQueue(map[string]queueRunner{"a": a, "b": b}, time.Minute)
	_ = q.Submit("a", taskList())
	await(t, a.started)
	_ = q.Submit("b", taskList())
	if err := q.Cancel("b"); err != nil {
		t.Fatal(err)
	}
	if a.stops.Load() != 0 || b.stops.Load() != 0 {
		t.Fatal("canceling pending job touched a running instance")
	}
	if err := q.Cancel("a"); err != nil {
		t.Fatal(err)
	}
	if a.stops.Load() != 2 {
		t.Fatal("active task was not stopped")
	}
	select {
	case <-b.started:
		t.Fatal("canceled pending job started")
	default:
	}
}
func TestTimeoutClosesOnlyExpiredGameBeforeDispatch(t *testing.T) {
	a, b := fake(), fake()
	q := testQueue(map[string]queueRunner{"a": a, "b": b}, 80*time.Millisecond)
	_ = q.Submit("a", taskList())
	await(t, a.started)
	_ = q.Submit("b", taskList())
	done := jobDone(q, "b")
	await(t, b.started)
	if a.closes.Load() != 2 || a.stops.Load() != 2 {
		t.Fatal("expired account was not cleaned up")
	}
	if b.closes.Load() != 0 {
		t.Fatal("next account's game was closed")
	}
	// B waited almost its entire timeout in the queue, but gets a new deadline.
	close(b.release)
	await(t, done)
	if b.closes.Load() != 0 {
		t.Fatal("queue wait was counted as execution time")
	}
}
func TestUnconfirmedStopNeverReleasesQueue(t *testing.T) {
	a, b := fake(), fake()
	a.blockStop = true
	fatal := make(chan error, 1)
	q := newTaskQueue(map[string]queueRunner{"a": a, "b": b}, 30*time.Millisecond, 40*time.Millisecond, func(err error) { fatal <- err })
	_ = q.Submit("a", taskList())
	await(t, a.started)
	_ = q.Submit("b", taskList())
	select {
	case err := <-fatal:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("missing recovery")
	}
	select {
	case <-b.started:
		t.Fatal("dispatched despite failed stop")
	default:
	}
	if err := q.Submit("b", taskList()); err == nil {
		t.Fatal("accepted task after fatal failure")
	}
}
func TestConcurrentLogSubscribers(t *testing.T) {
	m := &Maa{cfg: &Asst{ID: "test"}, SSEClients: make(map[string]*dto.SSEClientConn)}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprint(i)
			for j := 0; j < 100; j++ {
				m.NewClient(id)
				_ = m.GetClient(id)
				m.AppendLog(dto.Maalog{})
				m.DeleteClient(id)
				m.DeleteClient(id)
			}
		}(i)
	}
	wg.Wait()
}
func TestDurationConfiguration(t *testing.T) {
	if d, err := configuredDuration("", 2*time.Hour); err != nil || d != 2*time.Hour {
		t.Fatal(d, err)
	}
	if d, err := configuredDuration("30m", 2*time.Hour); err != nil || d != 30*time.Minute {
		t.Fatal(d, err)
	}
	for _, v := range []string{"0s", "-1h", "invalid"} {
		if _, err := configuredDuration(v, time.Hour); err == nil {
			t.Fatalf("accepted %s", v)
		}
	}
}
