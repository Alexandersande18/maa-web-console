package task

import (
	"context"
	"sync"
)

// Task 表示一个任务的结构
type Task struct {
    ID     string      `json:"id"`
    Type   string      `json:"type"`
    Params interface{} `json:"params,omitempty"`
}

// sTaskManager 管理任务的结构
type sTaskManager struct {
    sync.Mutex
    // key: user_device, value: list of tasks
    tasks map[string][]Task
    // key: user_device_taskID, value: bool (executed)
    executedTasks map[string]bool
}

// NewTaskManager 创建一个新的 TaskManager 实例
func NewTaskManager() *sTaskManager {
    return &sTaskManager{
        tasks:          make(map[string][]Task),
        executedTasks: make(map[string]bool),
    }
}

func init() {
    InitializeSampleTasks()
}

func (tm *sTaskManager) Run(ctx context.Context) error {
    return nil
}


var TaskManagerInstance = NewTaskManager()

// AddTask 为特定用户和设备添加任务
func (tm *sTaskManager) AddTask(user, device string, task Task) {
    tm.Lock()
    defer tm.Unlock()
    key := user + "_" + device
    tm.tasks[key] = append(tm.tasks[key], task)
}

// GetPendingTasks 获取未执行的任务
func (tm *sTaskManager) GetPendingTasks(user, device string) []Task {
    tm.Lock()
    defer tm.Unlock()
    key := user + "_" + device
    allTasks := tm.tasks[key]
    var pending []Task
    for _, task := range allTasks {
        execKey := key + "_" + task.ID
        if !tm.executedTasks[execKey] {
            pending = append(pending, task)
        }
    }
    return pending
}

// MarkTaskExecuted 标记任务为已执行
func (tm *sTaskManager) MarkTaskExecuted(user, device, taskID string) {
    tm.Lock()
    defer tm.Unlock()
    key := user + "_" + device
    execKey := key + "_" + taskID
    tm.executedTasks[execKey] = true
    
}

// InitializeSampleTasks 初始化一些示例任务（用于演示）
func InitializeSampleTasks() {
    user := "ea6c39eb-a45f-4d82-9ecc-33a7bf2ae4dc"
    device := "f7cd9682-3de9-4eef-9137-ec124ea9e9ec"
    // 生成一些示例任务
    tasks := []Task{
        {ID: "task1", Type: "CaptureImage"},
        {ID: "task2", Type: "LinkStart"},
        {ID: "task3", Type: "LinkStart-Recruiting"},
        {ID: "task4", Type: "Toolbox-GachaOnce"},
        {ID: "task5", Type: "Settings-ConnectionAddress", Params: "new_address_value"},
        {ID: "task6", Type: "CaptureImageNow"},
        {ID: "task7", Type: "StopTask"},
        {ID: "task8", Type: "HeartBeat"},
    }
    for _, task := range tasks {
        TaskManagerInstance.AddTask(user, device, task)
    }
}