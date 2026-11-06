package maa

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"maa-backend/internal/model/dto"
	"maa-backend/internal/service"
)

type sMaasst struct {
	instances map[string]*Maa
	queue     *taskQueue
}

var defaultService = &sMaasst{instances: make(map[string]*Maa)}

func init() { service.RegisterMaasst(defaultService) }

// Explicit startup keeps importing this package (including unit tests) free of
// native, ADB and database side effects. Start cron only after Initialize.
func Initialize() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	timeout, err := configuredDuration(cfg.TaskTimeout, 2*time.Hour)
	if err != nil {
		return fmt.Errorf("maa.tasktimeout: %w", err)
	}
	stopTimeout, err := configuredDuration(cfg.StopTimeout, 30*time.Second)
	if err != nil {
		return fmt.Errorf("maa.stoptimeout: %w", err)
	}
	if cfg.UserDir == "" {
		cfg.UserDir = filepath.Join(".runtime", "maa")
	}
	instances := make(map[string]*Maa)
	runners := make(map[string]queueRunner)
	devices := make(map[string]bool)
	client := ""
	for _, ast := range cfg.Asst {
		if ast.ID == "" || instances[ast.ID] != nil {
			return fmt.Errorf("missing or duplicate MAA instance ID: %s", ast.ID)
		}
		if ast.ADBAddress == "" || devices[ast.ADBAddress] {
			return fmt.Errorf("missing or duplicate ADB address: %s", ast.ADBAddress)
		}
		devices[ast.ADBAddress] = true
		if ast.ClientType == "" {
			ast.ClientType = "Official"
		}
		if gamePackages[ast.ClientType] == "" {
			return fmt.Errorf("unsupported clienttype: %s", ast.ClientType)
		}
		if client != "" && client != ast.ClientType {
			return fmt.Errorf("shared MaaCore runtime requires the same clienttype for every instance")
		}
		client = ast.ClientType
		if ast.ResourcePath == "" {
			ast.ResourcePath, err = os.Getwd()
			if err != nil {
				return err
			}
		}
		ast.ResourcePath, err = filepath.Abs(ast.ResourcePath)
		if err != nil {
			return err
		}
		ast.UserDir = cfg.UserDir
		instance, err := New(ast)
		if err != nil {
			return fmt.Errorf("initialize %s: %w", ast.ID, err)
		}
		instances[ast.ID] = instance
		runners[ast.ID] = instance
		log.WithFields(log.Fields{"instance": ast.ID, "version": instance.GetVersion()}).Info("MAA instance ready; ADB connects when its queued job starts")
	}
	defaultService.instances = instances
	defaultService.queue = newTaskQueue(runners, timeout, stopTimeout, func(err error) {
		// A native call may be hung. Releasing the queue would allow concurrent
		// execution again. Exit and let the installed Windows service restart us.
		log.WithError(err).Error("MAA did not stop within the grace period; restarting process")
		os.Exit(1)
	})
	return nil
}
func configuredDuration(value string, fallback time.Duration) (time.Duration, error) {
	if strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}
	if d <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}
	return d, nil
}
func (s *sMaasst) HasInstance(id string) bool { return s.instances[id] != nil }
func (s *sMaasst) getInstance(id string) *Maa { return s.instances[id] }
func (s *sMaasst) TakeScreenshot(ctx context.Context, id string) ([]byte, error) {
	m := s.instances[id]
	if m == nil {
		return nil, fmt.Errorf("instance not found")
	}
	return takeADBScreenshot(ctx, m.cfg, runScreenshotADB)
}
func (s *sMaasst) LinkStart(id string, tasks []dto.Tasker) error {
	if s.queue == nil {
		return fmt.Errorf("MAA is not initialized")
	}
	m := s.instances[id]
	if m == nil {
		return fmt.Errorf("instance not found")
	}
	// Timeout cleanup must target the same game package as the queued batch.
	for _, task := range tasks {
		if task == nil {
			return fmt.Errorf("nil task")
		}
		var params struct {
			ClientType string `json:"client_type"`
		}
		if err := json.Unmarshal([]byte(task.Serialize()), &params); err != nil {
			return err
		}
		if params.ClientType != "" && params.ClientType != m.cfg.ClientType {
			return fmt.Errorf("task client_type %q does not match configured clienttype %q", params.ClientType, m.cfg.ClientType)
		}
	}
	return s.queue.Submit(id, tasks)
}
func (s *sMaasst) StopTask(id string) error {
	if !s.HasInstance(id) {
		return fmt.Errorf("instance not found")
	}
	if s.queue == nil {
		return fmt.Errorf("MAA is not initialized")
	}
	return s.queue.Cancel(id)
}
func (s *sMaasst) NewSSEClient(mid, cid string) (*dto.SSEClientConn, error) {
	m := s.instances[mid]
	if m == nil {
		return nil, fmt.Errorf("instance not found")
	}
	m.NewClient(cid)
	if s.queue != nil {
		m.Notify(s.queue.Status(mid))
	}
	return m.GetClient(cid), nil
}
func (s *sMaasst) DeleteSSEClient(mid, cid string) error {
	m := s.instances[mid]
	if m == nil {
		return fmt.Errorf("instance not found")
	}
	m.DeleteClient(cid)
	return nil
}
func (s *sMaasst) GetSSEClient(mid, cid string) (*dto.SSEClientConn, error) {
	m := s.instances[mid]
	if m == nil {
		return nil, fmt.Errorf("instance not found")
	}
	return m.GetClient(cid), nil
}
