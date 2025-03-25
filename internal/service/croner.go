// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"maa-backend/internal/model/dto"
)

type (
	ICroner interface {
		// ExecCronJob 执行具体任务
		ExecCronJob(job *dto.CronJob) error
		// SetupCronJobs 初始化定时任务
		SetupCronJobs()
	}
)

var (
	localCroner ICroner
)

func Croner() ICroner {
	if localCroner == nil {
		panic("implement not found for interface ICroner, forgot register?")
	}
	return localCroner
}

func RegisterCroner(i ICroner) {
	localCroner = i
}
