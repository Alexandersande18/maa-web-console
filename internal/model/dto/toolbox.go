package dto

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TaskActivity struct {
	Enable      bool   `json:"enable"`
	ActivityID  string `json:"activity_id"`
	CurrentPage bool   `json:"current_page"`
}

func (t *TaskActivity) GetType() TaskType { return Activity }

func (t *TaskActivity) Serialize() string {
	data, _ := json.Marshal(t)
	return string(data)
}

type TaskCustom struct {
	TaskNames []string `json:"task_names"`
}

func (t *TaskCustom) GetType() TaskType { return Custom }

func (t *TaskCustom) Serialize() string {
	data, _ := json.Marshal(t)
	return string(data)
}

const (
	ReclamationTales  = "Tales"
	ReclamationAnchor = "RelaunchAnchor"
)

type TaskReclamation struct {
	Enable          bool     `json:"enable"`
	Theme           string   `json:"theme"`
	Mode            int      `json:"mode"`
	ToolsToCraft    []string `json:"tools_to_craft"`
	IncrementMode   int      `json:"increment_mode"`
	NumCraftBatches int      `json:"num_craft_batches"`
	ClearStore      bool     `json:"clear_store"`
	CurrentPage     bool     `json:"current_page"`
}

func (t *TaskReclamation) GetType() TaskType { return Reclamation }

func (t *TaskReclamation) Validate() error {
	switch t.Theme {
	case "", ReclamationTales:
		switch t.Mode {
		case 0:
			return nil
		case 1:
			if t.IncrementMode < 0 || t.IncrementMode > 1 || t.NumCraftBatches < 0 {
				return fmt.Errorf("生息演算制造参数无效")
			}
			return nil
		}
	case ReclamationAnchor:
		switch t.Mode {
		case 16, 32, 48:
			return nil
		}
	}
	return fmt.Errorf("生息演算主题和模式不匹配")
}

func (t *TaskReclamation) Serialize() string {
	config := *t
	if config.Theme == "" {
		config.Theme = ReclamationTales
	}
	// MaaCore v6.22.0 requires a real array in Tales' crafting mode.
	config.ToolsToCraft = []string{}
	if config.Theme == ReclamationTales && config.Mode == 1 {
		for _, tool := range t.ToolsToCraft {
			if tool = strings.TrimSpace(tool); tool != "" {
				config.ToolsToCraft = append(config.ToolsToCraft, tool)
			}
		}
		if len(config.ToolsToCraft) == 0 {
			config.ToolsToCraft = append(config.ToolsToCraft, "荧光棒")
		}
		if config.NumCraftBatches == 0 {
			config.NumCraftBatches = 1
		}
	} else {
		config.IncrementMode = 0
		config.NumCraftBatches = 0
	}
	config.ClearStore = config.ClearStore && config.Theme == ReclamationTales && config.Mode == 0
	data, _ := json.Marshal(config)
	return string(data)
}
