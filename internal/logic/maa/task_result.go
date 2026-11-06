package maa

import (
	"encoding/json"
	"errors"
	"fmt"

	"maa-backend/internal/model/dto"
)

// AllTasksCompleted means the batch drained, including tasks that failed.
// Keep TaskChainError until the final callback; CloseDown must still run.
type taskRun struct {
	ids       map[dto.TaskID]bool
	completed bool
	err       error
}

func (m *Maa) beginTaskRun() *taskRun {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.run = &taskRun{ids: make(map[dto.TaskID]bool)}
	return m.run
}

func (m *Maa) endTaskRun(run *taskRun) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.run == run {
		m.run = nil
	}
}

func (m *Maa) taskRunResult(run *taskRun) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return run.completed, run.err
}

func (m *Maa) recordTaskCallback(message int, details string) {
	switch message {
	case dto.TaskChainError, dto.AllTasksCompleted, dto.InternalError, dto.InitFailed:
	default:
		return
	}
	var event struct {
		TaskChain string     `json:"taskchain"`
		TaskID    dto.TaskID `json:"taskid"`
	}
	if err := json.Unmarshal([]byte(details), &event); err != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	run := m.run
	if run == nil || run.completed {
		return
	}
	if message == dto.InternalError || message == dto.InitFailed {
		run.err = errors.Join(run.err, fmt.Errorf("MaaCore %s: %s", dto.GetAsstMsgDescription(message), details))
		return
	}
	// Ignore delayed callbacks from an earlier batch on the same handle.
	if !run.ids[event.TaskID] {
		return
	}
	if message == dto.TaskChainError {
		run.err = errors.Join(run.err, fmt.Errorf("MaaCore task failed: %s (task_id=%d)", event.TaskChain, event.TaskID))
	} else {
		run.completed = true
	}
}
