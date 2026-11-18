package maa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"maa-backend/internal/model/dto"
)

func TestReclamationNavigationGatesNativeTask(t *testing.T) {
	for _, tc := range []struct {
		name string
		fail bool
		want []string
	}{
		{"arrived", false, []string{"StartUp", "Custom", "Reclamation", "CloseDown"}},
		{"wrong page", true, []string{"StartUp", "Custom", "CloseDown"}},
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
			asstRunning = func(...uintptr) (uintptr, uintptr, error) { return 0, 0, nil }
			var appended []string
			asstAppendTask = func(args ...uintptr) (uintptr, uintptr, error) {
				appended = append(appended, CharPtr2String(args[1]))
				return uintptr(len(appended)), 0, nil
			}
			m := &Maa{cfg: &Asst{ResourcePath: root, ToolboxResourcePath: root, ClientType: "Official"}}
			asstStart = func(...uintptr) (uintptr, uintptr, error) {
				chain := appended[len(appended)-1]
				detail := fmt.Sprintf(`{"taskchain":%q,"taskid":%d}`, chain, len(appended))
				if chain == "Custom" && tc.fail {
					m.recordTaskCallback(dto.TaskChainError, detail)
				}
				m.recordTaskCallback(dto.AllTasksCompleted, detail)
				return 1, 0, nil
			}
			err := m.Execute(context.Background(), []dto.Tasker{
				&dto.TaskStartUp{}, &dto.TaskReclamation{Enable: true, Theme: dto.ReclamationTales}, &dto.TaskCloseDown{},
			})
			if tc.fail {
				if err == nil || !strings.Contains(err.Error(), "未到达所选模式的开始页面") {
					t.Fatalf("incorrect navigation failure: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(appended, tc.want) {
				t.Fatalf("appended %v, want %v", appended, tc.want)
			}
		})
	}
}
