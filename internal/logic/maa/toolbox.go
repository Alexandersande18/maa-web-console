package maa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"maa-backend/internal/model/dto"
)

type activityRoute struct {
	Name            string    `json:"name"`
	Clients         []string  `json:"clients"`
	StartsAt        time.Time `json:"starts_at"`
	EndsAt          time.Time `json:"ends_at"`
	Navigation      string    `json:"navigation"`
	FromCurrentPage string    `json:"from_current_page"`
}

func (m *Maa) resolveActivity(task *dto.TaskActivity, now time.Time) (*dto.TaskCustom, error) {
	data, err := os.ReadFile(filepath.Join(m.cfg.ToolboxResourcePath, "activities.json"))
	if err != nil {
		return nil, fmt.Errorf("读取活动目录失败: %w", err)
	}
	var routes map[string]activityRoute
	if err := json.Unmarshal(data, &routes); err != nil {
		return nil, fmt.Errorf("活动目录格式有误: %w", err)
	}
	route, exists := routes[task.ActivityID]
	if !exists {
		return nil, fmt.Errorf("暂无支持的活动路线: %s", task.ActivityID)
	}
	if !slices.Contains(route.Clients, m.cfg.ClientType) {
		return nil, fmt.Errorf("%s 暂不支持客户端 %s", route.Name, m.cfg.ClientType)
	}
	if now.Before(route.StartsAt) || (!route.EndsAt.IsZero() && !now.Before(route.EndsAt)) {
		return nil, fmt.Errorf("%s 当前未开放", route.Name)
	}
	entry := route.Navigation
	if task.CurrentPage {
		entry = route.FromCurrentPage
	}
	if entry == "" {
		return nil, fmt.Errorf("%s 缺少所选启动方式的路线", route.Name)
	}
	return &dto.TaskCustom{TaskNames: []string{entry}}, nil
}

type nativeBatch struct {
	tasks        []dto.Tasker
	afterSuccess dto.Tasker
}

func (m *Maa) prepareBatches(ctx context.Context, tasks []dto.Tasker) ([]nativeBatch, error) {
	var batches []nativeBatch
	var pending []dto.Tasker
	needsOverlay := false
	for _, task := range tasks {
		switch task.GetType() {
		case dto.Activity:
			var config dto.TaskActivity
			if err := json.Unmarshal([]byte(task.Serialize()), &config); err != nil {
				return nil, err
			}
			custom, err := m.resolveActivity(&config, time.Now())
			if err != nil {
				return nil, err
			}
			task = custom
			needsOverlay = true
		case dto.Reclamation:
			var config dto.TaskReclamation
			if err := json.Unmarshal([]byte(task.Serialize()), &config); err != nil {
				return nil, err
			}
			if err := config.Validate(); err != nil {
				return nil, err
			}
			if !config.CurrentPage || config.Theme == dto.ReclamationAnchor {
				if !config.CurrentPage && m.cfg.ClientType != "Official" && m.cfg.ClientType != "Bilibili" {
					return nil, fmt.Errorf("生息演算自动导航目前适用于简中客户端，请使用从当前游戏页面开始")
				}
				if len(pending) > 0 {
					batches = append(batches, nativeBatch{tasks: pending})
					pending = nil
				}
				gate := &dto.TaskCustom{TaskNames: []string{reclamationNavigation(config)}}
				batches = append(batches, nativeBatch{tasks: []dto.Tasker{gate}, afterSuccess: task})
				needsOverlay = true
				continue
			}
		}
		pending = append(pending, task)
	}
	if len(pending) > 0 {
		batches = append(batches, nativeBatch{tasks: pending})
	}
	if needsOverlay {
		nativeMu.Lock()
		defer nativeMu.Unlock()
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := loadResource(m.cfg.ToolboxResourcePath); err != nil {
			return nil, fmt.Errorf("加载牛杂导航资源失败: %w", err)
		}
	}
	return batches, nil
}

func reclamationNavigation(task dto.TaskReclamation) string {
	if task.Theme == dto.ReclamationAnchor {
		stage := map[int]string{
			16: "RA1",
			32: "RA15",
			48: "RA4",
		}[task.Mode]
		if task.CurrentPage {
			return "WebNav@Anchor@Panel" + stage
		}
		return "WebNav@Anchor@" + stage
	}
	if task.Mode == 1 {
		return "WebNav@Tales@InSave"
	}
	return "WebNav@Tales@NoSave"
}

func (m *Maa) executeBatches(ctx context.Context, batches []nativeBatch) error {
	var result error
	for _, batch := range batches {
		drained, err := m.executeBatch(ctx, batch.tasks)
		if err != nil && batch.afterSuccess != nil && ctx.Err() == nil {
			err = fmt.Errorf("生息演算未到达所选模式的开始页面，请检查主题、存档或关卡解锁状态: %w", err)
		}
		result = errors.Join(result, err)
		if !drained {
			return result
		}
		if batch.afterSuccess != nil && err == nil {
			drained, err = m.executeBatch(ctx, []dto.Tasker{batch.afterSuccess})
			result = errors.Join(result, err)
			if !drained {
				return result
			}
		}
	}
	return result
}
