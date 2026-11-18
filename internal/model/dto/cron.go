package dto

import "github.com/gogf/gf/v2/os/gtime"

type CronJob struct {
	TaskSelection
	Hour           int         `json:"hour" v:"required|between:0,23"`
	Minute         int         `json:"minute" v:"required|between:0,59"`
	Id             int         `json:"id" v:""`
	Name           string      `json:"name" v:""`
	UserName       string      `json:"userName" v:""`
	IsTemplate     int         `json:"isTemplate" v:""`
	LastExecutedAt *gtime.Time `json:"lastExecutedAt" v:""`
}
