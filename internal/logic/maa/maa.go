package maa

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"

	log "github.com/sirupsen/logrus"
	"maa-backend/internal/model/dto"
)

var (
	ErrorMAAFalse          = errors.New("MAA return false")
	ErrorAsstCreateEx      = errors.New("AsstCreateEx failed")
	ErrorAsstLoadSource    = errors.New("AsstLoadResource failed")
	ErrorAsstConnect       = errors.New("AsstConnect failed")
	ErrorAsstCreate        = errors.New("AsstCreate failed")
	ErrorAsstGetVersion    = errors.New("AsstGetVersion failed")
	ErrorAsstAppendTask    = errors.New("AsstAppendTask failed")
	ErrorAsstSetTaskParams = errors.New("AsstSetTaskParams failed")
	ErrorAsstStart         = errors.New("AsstStart failed")
	ErrorAsstStop          = errors.New("AsstStop failed")
	ErrorAsstRunning       = errors.New("AsstRunning failed")
	ErrorAsstDestroy       = errors.New("AsstDestroy failed")
	ErrTooLargeAllocation  = errors.New("buffer allocation exceeds maximum limit")
	// Native control calls are serialized across handles. This does not replace
	// the queue's exclusion over the entire lifetime of a running task.
	nativeMu sync.Mutex
)

const nativeNullSize = ^uintptr(0)

type AsstCallback func(message int, details uintptr, arg uintptr) uintptr

type Maa struct {
	asstHandle uintptr
	version    string
	hooks      []Hook
	mu         sync.Mutex // hooks and log subscribers only; never held across native calls
	cfg        *Asst
	SSEClients map[string]*dto.SSEClientConn
	run        *taskRun // guarded by mu; callbacks are delivered on a native message thread
}

func New(cfg *Asst) (*Maa, error) {
	if err := validateMuMuConfig(cfg); err != nil {
		return nil, err
	}
	if err := initializeRuntime(cfg.ResourcePath, cfg.UserDir, cfg.ClientType); err != nil {
		return nil, err
	}
	if muMuEnabled(cfg) && asstSetConnectionExtras == nil {
		return nil, errors.New("MuMu screenshot enhancement requires MaaCore AsstSetConnectionExtras")
	}
	m := &Maa{cfg: cfg, SSEClients: make(map[string]*dto.SSEClientConn)}
	m.setDefaultHook()
	if err := m.AsstCreateEx(); err != nil {
		return nil, err
	}
	m.SetVersion()
	return m, nil
}

func (m *Maa) callback(message int, details uintptr, arg uintptr) (result uintptr) {
	defer func() {
		if r := recover(); r != nil {
			log.WithField("instance", m.cfg.ID).Errorf("MAA callback panic: %v", r)
		}
	}()
	// MaaCore owns details and frees it when this callback returns.
	detail := CharPtr2String(details)
	m.mu.Lock()
	hooks := append([]Hook(nil), m.hooks...)
	m.mu.Unlock()
	for _, hook := range hooks {
		hook.Do(message, detail, arg)
	}
	return 0
}

func (m *Maa) AppendHook(hooks ...Hook) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks = append(m.hooks, hooks...)
}
func (m *Maa) setDefaultHook() { m.AppendHook(&TaskCallback{Instance: m}) }

// callStrings retains typed heap pointers until the DLL call finishes. Never
// return just a uintptr to Go memory from a conversion helper.
func callStrings(fn procCall, prefix []uintptr, values ...string) (uintptr, error) {
	pointers := make([]*byte, len(values))
	args := append([]uintptr(nil), prefix...)
	for i, value := range values {
		p, err := bytePtrFromString(value)
		if err != nil {
			return 0, err
		}
		pointers[i] = p
		args = append(args, uintptr(unsafe.Pointer(p)))
	}
	r, _, _ := fn(args...)
	runtime.KeepAlive(pointers)
	return r, nil // MaaCore uses its return value, not Windows GetLastError.
}

func (m *Maa) asstLoadResource(path string) error {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	return loadResource(path)
}
func (m *Maa) AsstCreateEx() error {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	r, _, _ := asstCreateEx(syscall.NewCallback(m.callback), 0)
	if r == 0 {
		return ErrorAsstCreateEx
	}
	m.asstHandle = r
	return nil
}
func (m *Maa) AsstCreate() error {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	r, _, _ := asstCreate()
	if r == 0 {
		return ErrorAsstCreate
	}
	m.asstHandle = r
	return nil
}
func (m *Maa) AsstConnect(path, address, config, usrdir string) error {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	if err := m.configureConnectionExtras(config); err != nil {
		return err
	}
	// AsstSetUserDir is global and is set once, before resource loading.
	r, err := callStrings(asstConnect, []uintptr{m.asstHandle}, path, address, config)
	if err != nil {
		return err
	}
	if r == 0 {
		return ErrorAsstConnect
	}
	return nil
}
func (m *Maa) SetVersion() {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	r, _, _ := asstGetVersion()
	m.version = CharPtr2String(r)
}
func (m *Maa) GetVersion() string { return m.version }
func (m *Maa) AsstAppendTask(typ dto.TaskType, task dto.Tasker) (dto.TaskID, error) {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	params := task.Serialize()
	r, err := callStrings(asstAppendTask, []uintptr{m.asstHandle}, typ.String(), params)
	if err != nil {
		return 0, err
	}
	if int32(r) <= 0 {
		return 0, fmt.Errorf("%w: type=%s params=%s resource=%s version=%s", ErrorAsstAppendTask, typ, params, m.cfg.ResourcePath, m.version)
	}
	return dto.TaskID(int32(r)), nil
}
func (m *Maa) AsstSetTaskParams(id dto.TaskID, task dto.Tasker) bool {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	r, err := callStrings(asstSetTaskParams, []uintptr{m.asstHandle, uintptr(id)}, task.Serialize())
	return err == nil && r != 0
}
func (m *Maa) AsstStart() (bool, error) {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	r, _, _ := asstStart(m.asstHandle)
	if r == 0 {
		return false, ErrorAsstStart
	}
	return true, nil
}
func (m *Maa) AsstStop() (bool, error) {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	r, _, _ := asstStop(m.asstHandle)
	if r == 0 {
		return false, ErrorAsstStop
	}
	return true, nil
}
func (m *Maa) AsstRunning() bool {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	r, _, _ := asstRunning(m.asstHandle)
	return r != 0
}
func (m *Maa) AsstDestroy() error {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	if m.asstHandle != 0 {
		asstDestroy(m.asstHandle)
		m.asstHandle = 0
	}
	return nil
}
func (m *Maa) GetImage() ([]byte, error) {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	r, _, _ := asstAsyncScreencap(m.asstHandle, 1)
	if int32(r) <= 0 {
		return nil, errors.New("AsstAsyncScreencap failed")
	}
	for size := 8 * 1024 * 1024; size <= 64*1024*1024; size *= 2 {
		buffer := make([]byte, size)
		n, _, _ := asstGetImage(m.asstHandle, uintptr(unsafe.Pointer(&buffer[0])), uintptr(size))
		runtime.KeepAlive(buffer)
		if n == nativeNullSize {
			continue
		}
		if n == 0 {
			return nil, errors.New("MAA image is empty")
		}
		if n > uintptr(len(buffer)) {
			return nil, errors.New("invalid image size returned by MaaCore")
		}
		return buffer[:int(n)], nil
	}
	return nil, ErrTooLargeAllocation
}
func (m *Maa) GetImageAsync() {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	asstAsyncScreencap(m.asstHandle, 1)
}
func (m *Maa) GetTasksList() ([]dto.TaskID, error) {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	// A null buffer returns AsstGetNullSize(), not the required count.
	for size := 32; size <= 65536; size *= 2 {
		tasks := make([]dto.TaskID, size) // AsstTaskId is int32, including on amd64.
		n, _, _ := asstGetTasksList(m.asstHandle, uintptr(unsafe.Pointer(&tasks[0])), uintptr(size))
		runtime.KeepAlive(tasks)
		if n == nativeNullSize {
			continue
		}
		if n > uintptr(size) {
			return nil, errors.New("invalid task count returned by MaaCore")
		}
		return tasks[:int(n)], nil
	}
	return nil, ErrTooLargeAllocation
}
func (m *Maa) NewClient(id string) {
	conn := &dto.SSEClientConn{Ch: make(chan dto.Maalog, 64), CreatedAt: time.Now()}
	m.mu.Lock()
	defer m.mu.Unlock()
	if old := m.SSEClients[id]; old != nil {
		close(old.Ch)
	}
	m.SSEClients[id] = conn
}
func (m *Maa) GetClient(id string) *dto.SSEClientConn {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.SSEClients[id]
}
func (m *Maa) DeleteClient(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c := m.SSEClients[id]; c != nil {
		delete(m.SSEClients, id)
		close(c.Ch)
	}
}
func (m *Maa) AppendLog(entry dto.Maalog) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.SSEClients {
		if c != nil {
			select {
			case c.Ch <- entry:
			default:
			}
		}
	}
}
func (m *Maa) Notify(message string) {
	log.WithField("instance", m.cfg.ID).Info(message)
	m.AppendLog(dto.Maalog{Level: "info", Message: message, Timestamp: time.Now().Unix()})
}

// Execute is called exclusively by the global task queue. All work, including
// connect/resource reload, is included in the task deadline.
func (m *Maa) Execute(ctx context.Context, tasks []dto.Tasker) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := reloadResources(ctx, m.cfg.ResourcePath, m.cfg.ClientType); err != nil {
		return err
	}
	batches, err := m.prepareBatches(ctx, tasks)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := m.AsstConnect(m.cfg.ADBPath, m.cfg.ADBAddress, m.cfg.ADBType, ""); err != nil {
		return err
	}
	return m.executeBatches(ctx, batches)
}

// drained reports whether the native queue finished and can accept another batch.
func (m *Maa) executeBatch(ctx context.Context, tasks []dto.Tasker) (drained bool, err error) {
	run := m.beginTaskRun()
	defer m.endTaskRun(run)
	for _, task := range tasks {
		if err := ctx.Err(); err != nil {
			return false, err
		}
		id, err := m.AsstAppendTask(task.GetType(), task)
		if err != nil {
			return false, err
		}
		m.mu.Lock()
		run.ids[id] = true
		m.mu.Unlock()
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if _, err := m.AsstStart(); err != nil {
		return false, err
	}
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	var idleSince time.Time
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-ticker.C:
			if !m.AsstRunning() {
				completed, err := m.taskRunResult(run)
				if completed {
					return true, err
				}
				// AsstRunning can become false before the message thread delivers
				// TaskChainError and AllTasksCompleted. Drain that result first.
				if idleSince.IsZero() {
					idleSince = time.Now()
				} else if time.Since(idleSince) >= 5*time.Second {
					return false, errors.Join(err, errors.New("MaaCore stopped without an AllTasksCompleted callback"))
				}
			} else {
				idleSince = time.Time{}
			}
		}
	}
}
func (m *Maa) Stop(ctx context.Context) error {
	if _, err := m.AsstStop(); err != nil {
		return err
	}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		if !m.AsstRunning() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
