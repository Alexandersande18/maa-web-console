package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	dto "maa-backend/internal/model/dto"
)


// LoginReq is the request struct for login
type LoginReq struct {
	g.Meta   `path:"/maasst/login" method:"post" tags:"Auth" summary:"User login api"`
	Token string `v:"required#请输入Token" dc:"Token for authentication"`
}

// LoginRes is the response struct for login
type LoginRes struct {
	Token string `json:"token" dc:"JWT token for authentication"`
}

type LinkStartReq struct {
	g.Meta `path:"/maasst/linkstart" Tags:"maasst" method:"post" summary:"LinkStartReq"`
	TaskRecruit	*dto.TaskRecruit	`json:"taskRecruit" v:""`
	TaskInfrast *dto.TaskInfrast	`json:"taskInfrast" v:""`
	TaskFight 	*dto.TaskFight		`json:"taskFight" v:""`
	TaskMall 	*dto.TaskMall		`json:"taskMall" v:""`
	TaskAward 	*dto.TaskAward		`json:"taskAward" v:""`
}

type LinkStartRes struct {}

type StopTaskReq struct {
	g.Meta `path:"/maasst/stoptask" Tags:"maasst" method:"post" summary:"StopTaskReq"`
}

type StopTaskRes struct {}

type ScreenshotReq struct {
	g.Meta `path:"/maasst/screenshot" Tags:"maasst" method:"post" summary:"ScreenshotReq"`
	// UserID string `json:"userID" v:"required#用户ID为必填"`
}

type ScreenshotRes struct {
	ImageData string `json:"imageData" dc:"Base64 编码的 PNG 图片"`
	Timestamp int64  `json:"timestamp" dc:"截图时间戳"`
}

type LogsReq struct {
	g.Meta `path:"/maasst/logs" Tags:"maasst" method:"post" summary:"LogsReq"`
}

type LogsRes struct {}

type LogsWSReq struct {
	g.Meta `path:"/maasst/logs-ws" Tags:"maasst" method:"get" summary:"LogsGETReq"`
}

type LogsWSRes struct {}