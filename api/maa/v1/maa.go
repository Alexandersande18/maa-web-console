// api/task.go

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Task 表示一个任务的结构
type Task struct {
	ID     string      `json:"id"    v:"required"`
	Type   string      `json:"type"  v:"required"`
	Params interface{} `json:"params,omitempty"`
}

// GetTaskReq 表示获取任务的请求结构体
type GetTaskReq struct {
	g.Meta `path:"/maa/getTask" method:"post" tags:"Task" summary:"获取任务"`
	User   string `json:"user"    v:"required#用户ID为必填"`
	Device string `json:"device"  v:"required#设备ID为必填"`
	// 可选字段可以在此添加
}

// GetTaskRes 表示获取任务的响应结构体
type GetTaskRes struct {
	Tasks []Task `json:"tasks"`
	// 可选返回值可以在此添加
}

// ReportStatusReq 表示汇报任务状态的请求结构体
type ReportStatusReq struct {
	g.Meta `path:"/maa/reportStatus" method:"post" tags:"Task" summary:"汇报任务状态"`
	User    string `json:"user"    v:"required#用户ID为必填"`
	Device  string `json:"device"  v:"required#设备ID为必填"`
	TaskID  string `json:"task"`
	Status  string `json:"status"  v:"required|in:SUCCESS,FAILED#状态为必填，且必须为 SUCCESS 或 FAILED"`
	Payload string `json:"payload"`
	// 可选字段可以在此添加
}

// ReportStatusRes 表示汇报任务状态的响应结构体
type ReportStatusRes struct {
	Message string `json:"message"`
}