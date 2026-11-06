package dto

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gogf/gf/v2/util/gvalid"
)

func TestTaskLimits(t *testing.T) {
	for _, tc := range []struct {
		name  string
		data  interface{}
		valid bool
	}{
		{"fight10", TaskFight{Series: 10, Server: "CN", ClientType: "Official"}, true},
		{"fight11", TaskFight{Series: 11, Server: "CN", ClientType: "Official"}, false},
		{"rogue3", TaskRoguelike{Theme: "BlackFlow", Mode: 30003, StartsCount: 3}, true},
		{"rogue4", TaskRoguelike{Theme: "BlackFlow", Mode: 30003, StartsCount: 4}, false},
		{"animals", TaskRoguelike{Theme: "BlackFlow", Mode: 30001, StartsCount: 3}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := gvalid.New().Data(tc.data).Run(context.Background())
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v, err=%v", tc.valid, err)
			}
		})
	}
}

func TestRoguelikeLegacyLimitAndStartRunSerialization(t *testing.T) {
	task := TaskRoguelike{Theme: "BlackFlow", Mode: 30003, StartsCount: 8}
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(task.Serialize()), &params); err != nil {
		t.Fatal(err)
	}
	if params["starts_count"] != float64(3) || task.StartsCount != 8 {
		t.Fatal("legacy task must be capped without mutating its saved configuration", params)
	}
	if animals, ok := params["blackflow_start_animals"].([]interface{}); !ok || len(animals) != 0 {
		t.Fatal("start-run animals must be an empty array, not null", params)
	}
	if _, exists := params["collectible_mode_core_char_list"]; exists {
		t.Fatal("unspecified preparation operators must preserve core-char fallback")
	}
}
