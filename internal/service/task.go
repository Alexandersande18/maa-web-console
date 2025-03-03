// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"maa-backend/internal/model/dto"
)

type (
	ITaskManager interface {
		Run(ctx context.Context) error
		// AddTask 为特定用户和设备添加任务
		AddTask(user string, device string, task dto.GUITask)
		// GetPendingTasks 获取未执行的任务
		GetPendingTasks(user string, device string) []dto.GUITask
		// MarkTaskExecuted 标记任务为已执行
		MarkTaskExecuted(user string, device string, taskID string)
	}
)

var (
	localTaskManager ITaskManager
)

func TaskManager() ITaskManager {
	if localTaskManager == nil {
		panic("implement not found for interface ITaskManager, forgot register?")
	}
	return localTaskManager
}

func RegisterTaskManager(i ITaskManager) {
	localTaskManager = i
}
