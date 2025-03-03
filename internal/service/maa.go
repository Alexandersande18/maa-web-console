// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	dto "maa-backend/internal/model/dto"
)

type (
	IMaasst interface {
		HasInstance(id string) bool
		TakeScreenshot(id string) ([]byte, error)
		LinkStart(id string, tasks []dto.Tasker) error
	}
)

var (
	localMaasst IMaasst
)

func Maasst() IMaasst {
	if localMaasst == nil {
		panic("implement not found for interface IMaasst, forgot register?")
	}
	return localMaasst
}

func RegisterMaasst(i IMaasst) {
	localMaasst = i
}
