package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	dto "maa-backend/internal/model/dto"
)

type LinkStartReq struct {
	g.Meta `path:"/maasst/linkstart" Tags:"maasst" method:"post" summary:"LinkStartReq"`
	UserID string `json:"userID" v:"required#用户ID为必填"`
	TaskStartUp dto.TaskStartUp 	`json:"taskStartUp" v:"required#任务启动为必填"`
	TaskInfrast dto.TaskInfrast		`json:"taskInfrast" v:"required#任务基建为必填"`
	TaskFight dto.TaskFight		`json:"taskFight" v:"required#任务战斗为必填"`
	TaskMall dto.TaskMall		`json:"taskMall" v:"required#任务商城为必填"`
	TaskAward dto.TaskAward		`json:"taskAward" v:"required#任务领奖为必填"`
	// TaskRoguelike dto.TaskRoguelike	`json:"taskRoguelike" v:"required#任务罗德岛为必填"`
}

type LinkStartRes struct {
}

type ScreenshotReq struct {
	g.Meta `path:"/maasst/screenshot" Tags:"maasst" method:"post" summary:"ScreenshotReq"`
	UserID string `json:"userID" v:"required#用户ID为必填"`
}

type ScreenshotRes struct {
	g.Meta `mime:"application/json"` // 明确声明返回 JSON
	ImageData string `json:"imageData" dc:"Base64 编码的 PNG 图片"`
	Timestamp int64  `json:"timestamp" dc:"截图时间戳"`
}