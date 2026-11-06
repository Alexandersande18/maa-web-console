package maa

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"maa-backend/internal/model/dto"
)

func TestExecutePreservesCoreFailureAfterCloseDown(t *testing.T) {
	for _, tc := range []struct {
		name    string
		message int
		delay   bool
		failed  bool
	}{
		{"successful batch", dto.TaskChainCompleted, false, false},
		{"roguelike failure followed by successful shutdown", dto.TaskChainError, false, true},
		{"callback arrives after core becomes idle", dto.TaskChainError, true, true},
		{"recoverable subtask error", dto.SubTaskError, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.Mkdir(filepath.Join(root, "resource"), 0755); err != nil {
				t.Fatal(err)
			}
			oldLoad, oldConnect, oldAppend, oldStart, oldRunning := asstLoadResource, asstConnect, asstAppendTask, asstStart, asstRunning
			defer func() {
				asstLoadResource, asstConnect, asstAppendTask, asstStart, asstRunning = oldLoad, oldConnect, oldAppend, oldStart, oldRunning
			}()
			ok := func(...uintptr) (uintptr, uintptr, error) { return 1, 0, nil }
			asstLoadResource, asstConnect = ok, ok
			var id uintptr
			asstAppendTask = func(...uintptr) (uintptr, uintptr, error) { id++; return id, 0, nil }
			var running, closed atomic.Bool
			m := &Maa{cfg: &Asst{ID: "test", ResourcePath: root, ClientType: "Official"}}
			hook := &TaskCallback{Instance: m}
			asstStart = func(...uintptr) (uintptr, uintptr, error) {
				running.Store(true)
				go func() {
					if tc.delay {
						running.Store(false)
						time.Sleep(400 * time.Millisecond)
					}
					hook.Do(tc.message, `{"taskchain":"Roguelike","taskid":2}`, nil)
					// A successful CloseDown and AllTasksCompleted do not erase failure.
					hook.Do(dto.TaskChainCompleted, `{"taskchain":"CloseDown","taskid":3}`, nil)
					closed.Store(true)
					hook.Do(dto.AllTasksCompleted, `{"finished_tasks":[1,2,3],"taskchain":"CloseDown","taskid":3}`, nil)
					running.Store(false)
				}()
				return 1, 0, nil
			}
			asstRunning = func(...uintptr) (uintptr, uintptr, error) {
				if running.Load() {
					return 1, 0, nil
				}
				return 0, 0, nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err := m.Execute(ctx, []dto.Tasker{
				&dto.TaskStartUp{}, &dto.TaskRoguelike{Theme: "BlackFlow"}, &dto.TaskCloseDown{},
			})
			if tc.failed {
				if err == nil || !strings.Contains(err.Error(), "Roguelike (task_id=2)") {
					t.Fatalf("lost core failure: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if !closed.Load() {
				t.Fatal("returned before CloseDown and the final callback")
			}
		})
	}
}

func TestTaskResultIgnoresPreviousBatchCallbacks(t *testing.T) {
	m := &Maa{}
	old := m.beginTaskRun()
	old.ids[1] = true
	m.recordTaskCallback(dto.TaskChainError, `{"taskchain":"Roguelike","taskid":1}`)
	m.endTaskRun(old)
	run := m.beginTaskRun()
	run.ids[2] = true
	m.recordTaskCallback(dto.TaskChainError, `{"taskchain":"Roguelike","taskid":1}`)
	m.recordTaskCallback(dto.AllTasksCompleted, `{"taskchain":"Roguelike","taskid":1}`)
	if complete, err := m.taskRunResult(run); complete || err != nil {
		t.Fatalf("previous batch contaminated new result: complete=%v err=%v", complete, err)
	}
	m.recordTaskCallback(dto.AllTasksCompleted, `{"taskchain":"Roguelike","taskid":2}`)
	if complete, err := m.taskRunResult(run); !complete || err != nil {
		t.Fatalf("new batch did not finish: complete=%v err=%v", complete, err)
	}
}
