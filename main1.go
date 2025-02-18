package main

import (
	"net/http"
	"sync"
	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)
// Task 表示一个任务的结构
type Task struct {
	ID     string      `json:"id"`
	Type   string      `json:"type"`
	Params interface{} `json:"params,omitempty"`
}
// GetTaskRequest 表示获取任务请求的结构
type GetTaskRequest struct {
	User   string `json:"user" binding:"required"`
	Device string `json:"device" binding:"required"`
	// 可选参数
	// Add other fields if needed
}
// GetTaskResponse 表示获取任务响应的结构
type GetTaskResponse struct {
	Tasks []Task `json:"tasks"`
	// 可选返回值
	// Add other fields if needed
}
// ReportStatusRequest 表示汇报任务状态请求的结构
type ReportStatusRequest struct {
	User    string `json:"user" binding:"required"`
	Device  string `json:"device" binding:"required"`
	TaskID  string `json:"task"`
	Status  string `json:"status" binding:"required,oneof=SUCCESS FAILED"`
	Payload string `json:"payload"`
	// 可选参数
	// Add other fields if needed
}
// TaskManager 管理任务的结构
type TaskManager struct {
	sync.Mutex
	// key: user_device, value: list of tasks
	tasks map[string][]Task
	// key: user_device_taskID, value: bool (executed)
	executedTasks map[string]bool
}
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:          make(map[string][]Task),
		executedTasks: make(map[string]bool),
	}
}
// AddTask 为特定用户和设备添加任务
func (tm *TaskManager) AddTask(user, device string, task Task) {
	tm.Lock()
	defer tm.Unlock()
	key := user + "_" + device
	tm.tasks[key] = append(tm.tasks[key], task)
}
// GetPendingTasks 获取未执行的任务
func (tm *TaskManager) GetPendingTasks(user, device string) []Task {
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
func (tm *TaskManager) MarkTaskExecuted(user, device, taskID string) {
	tm.Lock()
	defer tm.Unlock()
	key := user + "_" + device
	execKey := key + "_" + taskID
	tm.executedTasks[execKey] = true
}
// Initialize sample tasks (for demonstration purposes)
func initializeSampleTasks(tm *TaskManager) {
	user := "ea6c39eb-a45f-4d82-9ecc-33a7bf2ae4dc"
	device := "f7cd9682-3de9-4eef-9137-ec124ea9e9ec"
	// 生成一些示例任务
	tasks := []Task{
		// {
		// 	ID:   uuid.New().String(),
		// 	Type: "CaptureImage",
		// },
		// {
		// 	ID:   uuid.New().String(),
		// 	Type: "LinkStart",
		// },
		// {
		// 	ID:   uuid.New().String(),
		// 	Type: "LinkStart-Recruiting",
		// },
		// {
		// 	ID:   uuid.New().String(),
		// 	Type: "Toolbox-GachaOnce",
		// },
		// {
		// 	ID:     uuid.New().String(),
		// 	Type:   "Settings-ConnectionAddress",
		// 	Params: "new_address_value",
		// },
		// {
		// 	ID:   uuid.New().String(),
		// 	Type: "CaptureImageNow",
		// },
		// {
		// 	ID:   uuid.New().String(),
		// 	Type: "StopTask",
		// },
		// {
		// 	ID:   uuid.New().String(),
		// 	Type: "HeartBeat",
		// },
	}
	for _, task := range tasks {
		tm.AddTask(user, device, task)
	}
}
func main1() {
	r := gin.Default()
	taskManager := NewTaskManager()
	// 初始化一些示例任务
	initializeSampleTasks(taskManager)
	// 获取任务端点
	r.POST("/maa/getTask", func(c *gin.Context) {
		var req GetTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pendingTasks := taskManager.GetPendingTasks(req.User, req.Device)
		// 标记这些任务为已发送，以避免重复发送
		for _, task := range pendingTasks {
			taskManager.MarkTaskExecuted(req.User, req.Device, task.ID)
		}
		response := GetTaskResponse{
			Tasks: pendingTasks,
		}
		c.JSON(http.StatusOK, response)
	})
	// 汇报任务状态端点
	r.POST("/maa/reportStatus", func(c *gin.Context) {
		var req ReportStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 这里可以根据需要处理汇报的状态，比如记录到数据库、日志等
		// 这里只是简单打印
		c.JSON(http.StatusOK, gin.H{"message": "Status received"})
		// 示例：标记任务为已完成（已执行）
		// 你可以根据不同的状态进一步处理
		taskManager.MarkTaskExecuted(req.User, req.Device, req.TaskID)
	})
	// 启动服务器
	r.Run(":8080") // 可以根据需要更改端口
}