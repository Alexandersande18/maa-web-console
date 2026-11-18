package maa

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	"maa-backend/internal/model/dto"
)

func TestResourceOverlayOrderAndMigration(t *testing.T) {
	root := t.TempDir()
	for _, p := range []string{"resource", "cache/resource", "resource/global/YoStarJP/resource", "cache/resource/global/YoStarJP/resource"} {
		if err := os.MkdirAll(filepath.Join(root, p), 0755); err != nil {
			t.Fatal(err)
		}
	}
	data := []byte(`{"Event-8":{"algorithm":"JustReturn"}}`)
	if err := os.WriteFile(filepath.Join(root, "cache/resource/tasks.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	old := asstLoadResource
	defer func() { asstLoadResource = old }()
	var loaded []string
	asstLoadResource = func(a ...uintptr) (uintptr, uintptr, error) {
		loaded = append(loaded, CharPtr2String(a[0]))
		return 1, 0, nil
	}
	if err := loadResources(root, "YoStarJP"); err != nil {
		t.Fatal(err)
	}
	expected := []string{root, filepath.Join(root, "cache"), filepath.Join(root, "resource/global/YoStarJP"), filepath.Join(root, "cache/resource/global/YoStarJP")}
	if !reflect.DeepEqual(loaded, expected) {
		t.Fatal(loaded)
	}
	copied, err := os.ReadFile(filepath.Join(root, "cache/resource/tasks/tasks.json"))
	if err != nil || string(copied) != string(data) {
		t.Fatal(string(copied), err)
	}
	// Replacing an existing API cache must overwrite the old compatibility copy.
	changed := []byte(`{"Event-9":{}}`)
	_ = os.WriteFile(filepath.Join(root, "cache/resource/tasks.json"), changed, 0644)
	if err := prepareHotUpdate(filepath.Join(root, "cache")); err != nil {
		t.Fatal(err)
	}
	copied, _ = os.ReadFile(filepath.Join(root, "cache/resource/tasks/tasks.json"))
	if string(copied) != string(changed) {
		t.Fatal("stale task compatibility copy")
	}
}
func TestMissingOptionalCacheAndNativeFailure(t *testing.T) {
	root := t.TempDir()
	_ = os.Mkdir(filepath.Join(root, "resource"), 0755)
	old := asstLoadResource
	defer func() { asstLoadResource = old }()
	count := 0
	asstLoadResource = func(...uintptr) (uintptr, uintptr, error) { count++; return 1, 0, nil }
	if err := loadResources(root, "Official"); err != nil || count != 1 {
		t.Fatal(err, count)
	}
	asstLoadResource = func(...uintptr) (uintptr, uintptr, error) { return 0, 0, nil }
	if err := loadResources(root, "Official"); err == nil {
		t.Fatal("ignored native load failure")
	}
}
func TestNativeABI(t *testing.T) {
	if os.Getenv("MAA_NATIVE_SMOKE") != "1" {
		t.Skip("requires MAA_NATIVE_SMOKE=1; does not start tasks or connect a device")
	}
	cfg, err := LoadConfig()
	if err != nil || len(cfg.Asst) == 0 {
		t.Fatal(err)
	}
	c := *cfg.Asst[0]
	// MaaCore's process-global logger holds asst.log open until process exit;
	// testing.TempDir cleanup happens earlier on Windows.
	c.UserDir = filepath.Join(os.TempDir(), "maa-backend-native-smoke")
	c.ClientType = "Official"
	m, err := New(&c)
	if err != nil {
		t.Fatal(err)
	}
	defer m.AsstDestroy()
	t.Run("ToolboxResourceReload", func(t *testing.T) {
		overlay, err := filepath.Abs("../../../resource/maa")
		if err != nil {
			t.Fatal(err)
		}
		for round := 1; round <= 2; round++ {
			if err := loadResource(overlay); err != nil {
				t.Fatalf("round %d: load toolbox overlay: %v", round, err)
			}
			if err := reloadResources(context.Background(), c.ResourcePath, c.ClientType); err != nil {
				t.Fatalf("round %d: reload core resources after toolbox: %v", round, err)
			}
		}
	})
	task := &dto.TaskFight{Enable: true, Stage: "1-7", Times: 1, Series: 1}
	id, err := m.AsstAppendTask(dto.Fight, task)
	if err != nil || id <= 0 {
		t.Fatal(id, err)
	}
	if !m.AsstSetTaskParams(id, task) {
		t.Fatal("set params failed")
	}
	tasks, err := m.GetTasksList()
	if err != nil || len(tasks) != 1 || tasks[0] != id {
		t.Fatal(tasks, err)
	}
	if unsafe.Sizeof(dto.TaskID(0)) != 4 {
		t.Fatal("incorrect task ID width")
	}
	task.Series = 10
	if !m.AsstSetTaskParams(id, task) {
		t.Fatal("core rejected 10 consecutive fights")
	}
	for _, tc := range []struct {
		name     string
		mode     int
		strategy string
		utopia   bool
	}{
		{"automatic", 0, "", false},
		{"investment", 1, "", false},
		{"animals", 30001, "", false},
		{"burn", 0, "burn", false},
		{"burn_with_investment", 0, "burn_with_investment", false},
		{"ending1", 0, "ending1", false},
		{"ending2", 0, "ending2", false},
		{"ending3", 0, "ending3", false},
		{"start_run", 30003, "", false},
		{"utopia", 30003, "", true},
	} {
		t.Run("BlackFlow/"+tc.name, func(t *testing.T) {
			rogue := &dto.TaskRoguelike{
				Enable: true, Theme: "BlackFlow", Mode: tc.mode,
				Squad: "指挥分队", Roles: "取长补短", StartsCount: 3, Difficulty: 6,
				InvestmentEnabled: true, InvestmentsCount: 2147483647,
				BlackflowStrategy: tc.strategy, BlackflowCultivationTarget: "swaddled_cat",
				BlackflowStartAnimals: []string{"swaddled_cat", "swaddled_feathered_serpent"},
				BlackflowStartRewards: 65, BlackflowStartUtopia: tc.utopia,
				BlackflowStartRewardPriority: []string{"unnumbered_item", "pegasus"},
				CollectibleModeSquad:         "特勤分队", CollectibleModeRoles: "稳扎稳打",
				CollectibleModeCoreCharList: []dto.RoguelikeStartOper{{Name: "机械师"}},
			}
			rogueID, err := m.AsstAppendTask(dto.Roguelike, rogue)
			if err != nil || rogueID <= 0 {
				t.Fatal(rogueID, err)
			}
		})
	}
}
