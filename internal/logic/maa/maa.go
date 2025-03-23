package maa

import (
	// "MaaAssistantArknights/conf"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	dto "maa-backend/internal/model/dto"

	"github.com/gogf/gf/v2/frame/g"
	log "github.com/sirupsen/logrus"
)

var (
	ErrorMAAFalse          = errors.New("MAA return false")
	ErrorAsstCreateEx      = errors.New("AsstCreateEx failed")
	ErrorAsstLoadSource    = errors.New("AsstLoadSource failed")
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
)

type AsstCallback func(message int, details uintptr, arg uintptr) uintptr

type Maa struct {
	asstHandle uintptr
	version    string
	hooks      []Hook
	mu         sync.Mutex
	cfg        *Asst
	SSEClients    map[string]*dto.SSEClientConn
}

func New(cfg *Asst) (maa *Maa, err error) {
	if len(cfg.ResourcePath) == 0 {
		cfg.ResourcePath, err = os.Getwd()
		if err != nil {
			return
		}
	}
	maa = &Maa{
		hooks: make([]Hook, 0),
		cfg:   cfg,
		SSEClients:  make(map[string]*dto.SSEClientConn), // 确保map初始化
	}
	err = maa.initasst()
	if err != nil {
		return
	}
	maa.setDefaultHook()
	return
}

func (m *Maa) initasst() (err error) {
	cfg := m.cfg
	err = newSyscall(cfg.ResourcePath)
	if err != nil {
		return
	}
	err = m.asstLoadResource(cfg.ResourcePath)
	if err != nil {
		return
	}
	err = m.AsstCreateEx()
	if err != nil {
		return
	}
	err = m.AsstConnect(cfg.ADBPath, cfg.ADBAddress, cfg.ADBType, cfg.UserDir)
	if err != nil {
		return
	}

	m.SetVersion()
	return
}

func (m *Maa) callback(message int, details uintptr, arg uintptr) uintptr {
	m.mu.Lock()
	hooks := m.hooks
	m.mu.Unlock()
	for _, hook := range hooks {
		hook.Do(message, CharPtr2String(details), arg)
	}
	return 0
}

func (m *Maa) AppendHook(hooks ...Hook) {
	if len(hooks) == 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks = append(m.hooks, hooks...)
}

func (m *Maa) asstLoadResource(path string) error {
	a1, err := UintPtrFromString(path)
	if err != nil {
		return err
	}
	r1, _, err := asstLoadResource(a1)
	log.WithFields(log.Fields{
		"path": path,
		"r1":   r1,
		"err":  err,
	}).Println("asstLoadResource")
	if !errors.Is(err, syscall.Errno(0)) {
		return ErrorAsstLoadSource
	}
	if r1 == 0 {
		return ErrorMAAFalse
	}
	return nil
}

func (m *Maa) AsstCreateEx() error {
	a1 := syscall.NewCallback(m.callback)
	a2 := uintptr(0) // 这里的arg之后看有没有需要再更新
	r1, _, err := asstCreateEx(a1, a2)

	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("AsstCreateEx")

	if !errors.Is(err, syscall.Errno(0)) {
		return ErrorAsstCreateEx
	}
	m.asstHandle = r1
	return nil
}

func (m *Maa) AsstCreate() error {
	r1, _, err := asstCreate()
	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("AsstCreate")
	if !errors.Is(err, syscall.Errno(0)) {
		return ErrorAsstCreate
	}
	m.asstHandle = r1
	return nil
}

func (m *Maa) AsstConnect(path string, address string, config string, usrdir string) error {
	fmt.Println("adb start connect")
	a1, err := UintPtrFromString(path)
	if err != nil {
		return err
	}
	a2, err := UintPtrFromString(address)
	if err != nil {
		return err
	}
	a3, err := UintPtrFromString(config)
	if err != nil {
		return err
	}
	a4, err := UintPtrFromString(usrdir)
	if err != nil {
		return err
	}
	r0, _, err := asstSetUserDir(m.asstHandle, a4)
	log.WithFields(log.Fields{
		"r0":  r0,
		"err": err,
	}).Println("AsstSetUserDir")
	r1, _, err := asstConnect(m.asstHandle, a1, a2, a3)
	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("AsstConnect")

	if r1 == 0 {
		return ErrorAsstConnect
	}
	return nil
}

func (m *Maa) setDefaultHook() {
	m.hooks = append(m.hooks, &TaskCallback{
		Instance: m,
	})
}

func (m *Maa) SetVersion() {
	r1, _, err := asstGetVersion(m.asstHandle)
	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("SetVersion")
	if !errors.Is(err, syscall.Errno(0)) {
		return
	}
	if r1 > 0 {
		m.version = CharPtr2String(r1)
	}
}

func (m *Maa) GetVersion() string {
	return m.version
}

func (m *Maa) AsstAppendTask(typ dto.TaskType, task dto.Tasker) (dto.TaskID, error) {
	a1, err := UintPtrFromString(typ.String())
	if err != nil {
		return 0, err
	}
	
	a2, err := UintPtrFromString(task.Serialize())
	if err != nil {
		return 0, err
	}
	r1, _, err := asstAppendTask(m.asstHandle, a1, a2)
	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("asstAppendTask")
	if r1 == 0 {
		return 0, ErrorAsstAppendTask
	}

	return dto.TaskID(r1), nil
}

func (m *Maa) AsstSetTaskParams(taskID dto.TaskID, task dto.Tasker) bool {
	a1 := uintptr(taskID)
	a2, err := UintPtrFromString(task.Serialize())
	if err != nil {
		return false
	}
	r1, _, err := asstSetTaskParams(m.asstHandle, a1, a2)
	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("AsstSetTaskParams")

	return r1 > 0
}

func (m *Maa) AsstStart() (bool, error) {
	r1, _, err := asstStart(m.asstHandle)
	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("AsstStart")
	return r1 > 0, err
}

func (m *Maa) AsstStop() (bool, error) {
	r1, _, err := asstStop(m.asstHandle)
	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("AsstStop")

	return r1 > 0, err
}

func (m *Maa) AsstRunning() bool {
	r1, _, err := asstRunning(m.asstHandle)
	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("AsstRunning")
	return r1 > 0
}

func (m *Maa) AsstDestroy() error {
	r1, _, err := asstDestroy(m.asstHandle)
	log.WithFields(log.Fields{
		"r1":  r1,
		"err": err,
	}).Println("AsstDestroy")
	if r1 == 0 {
		return ErrorAsstDestroy
	}
	return nil
}


func (m *Maa) GetImage() ([]byte, error) {
    const (
        initialBufferSize = 2 * 1920 * 1080 * 4
        maxBufferSize     = 10 * initialBufferSize
    )
    
    buffSize := initialBufferSize
    m.GetImageAsync()
    for {
        if buffSize > maxBufferSize {
            return nil, ErrTooLargeAllocation
        }
        
        buffer := make([]byte, buffSize)
        bufferPtr := unsafe.Pointer(&buffer[0])
        
        actualSize, _, err := asstGetImage(
            m.asstHandle,
            uintptr(bufferPtr),
            uintptr(buffSize),
        )
        g.Log().Debugf(context.TODO(), "AsstGetImage attempt: buffSize: %d, actualSize: %d, error: %v", buffSize, actualSize, err)
        
        // 处理可能的系统错误
        if err != nil && !errors.Is(err, syscall.Errno(0)) {
            return nil, fmt.Errorf("system error: %w", err)
        }
        
        // 假设返回0表示缓冲区不足
        if actualSize == 0 {
            buffSize *= 2
            continue
        }
        
        // 返回有效数据（自动截断多余容量）
        return buffer[:actualSize], nil
    }
}

func (m *Maa) GetImageAsync() {
	callId, _, err := asstAsyncScreencap(m.asstHandle, uintptr(1))
	log.WithFields(log.Fields{
		"r1":  callId,
		"err": err,
	}).Println("AsstAsyncScreencap")
	g.Log().Debugf(context.TODO(), "AsstAsyncScreencap callId: %d", callId)
}


func (m *Maa) GetTasksList() ([]dto.TaskID, error) {
    // 第一次调用获取任务数量
    requiredSize, _, err := asstGetTasksList(m.asstHandle, 0, 0)
    if requiredSize == 0 {
        return nil, fmt.Errorf("failed to get tasks count: %v", err)
    }

    // 分配任务ID数组
    tasks := make([]dto.TaskID, requiredSize)
    tasksPtr := uintptr(unsafe.Pointer(&tasks[0]))

    // 第二次调用获取实际任务列表
    actualSize, _, err := asstGetTasksList(m.asstHandle, tasksPtr, uintptr(requiredSize))
    if actualSize == 0 {
        return nil, fmt.Errorf("failed to get tasks list: %v", err)
    }

    return tasks[:actualSize], nil
}

func (m *Maa) NewClient(id string) {
	ch := make(chan dto.Maalog, 10)
	conn := &dto.SSEClientConn{
		Ch:        ch,
		CreatedAt: time.Now(),
	}
	m.SSEClients[id] = conn
	firstMsg := dto.Maalog{
		Level:    "warn",
		Timestamp: time.Now().Unix(),
	}
	if m.AsstRunning(){
		firstMsg.Message = "当前有正在运行的任务"
	} else {
		firstMsg.Message = "没有正在运行的任务"
	}
	conn.Ch <- firstMsg
}

func (m *Maa) GetClient(id string) *dto.SSEClientConn {
	return m.SSEClients[id]
}

func (m *Maa) DeleteClient(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	close(m.SSEClients[id].Ch)
	delete(m.SSEClients, id)
}

func (m *Maa) AppendLog(maalog dto.Maalog) {
	m.mu.Lock()
	defer m.mu.Unlock()
    for k, v := range m.SSEClients {
        if v == nil {
            g.Log().Error(context.TODO(), "Client is nil", k)
            continue
        }
        // 使用闭包处理panic防止通道关闭导致崩溃
        func() {
            defer func() {
                if r := recover(); r != nil {
                    g.Log().Errorf(context.TODO(), "Send log to client %s failed: channel already closed", k)
                }
            }()
            select {
            case v.Ch <- maalog:
            default:
                g.Log().Errorf(context.TODO(), "Send log to client %s failed: channel full", k)
            }
        }()
    }
}

