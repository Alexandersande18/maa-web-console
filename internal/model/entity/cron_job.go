// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CronJob is the golang structure for table cron_job.
type CronJob struct {
	Id             int         `json:"id"             orm:"id"               description:"主键ID"`             // 主键ID
	Name           string      `json:"name"           orm:"name"             description:"任务名称（唯一）"`         // 任务名称（唯一）
	Username       string      `json:"username"       orm:"username"         description:"所属用户名"`            // 所属用户名
	IsTemplate     int         `json:"isTemplate"     orm:"is_template"      description:"是否为模板任务（0-否 1-是）"` // 是否为模板任务（0-否 1-是）
	Hour           uint        `json:"hour"           orm:"hour"             description:"小时（0-23）"`         // 小时（0-23）
	Minute         uint        `json:"minute"         orm:"minute"           description:"分钟（0-59）"`         // 分钟（0-59）
	TaskRecruit    string      `json:"taskRecruit"    orm:"task_recruit"     description:"公招任务配置"`           // 公招任务配置
	TaskInfrast    string      `json:"taskInfrast"    orm:"task_infrast"     description:"基建任务配置"`           // 基建任务配置
	TaskFight      string      `json:"taskFight"      orm:"task_fight"       description:"刷理智任务配置"`          // 刷理智任务配置
	TaskMall       string      `json:"taskMall"       orm:"task_mall"        description:"购物任务配置"`           // 购物任务配置
	TaskAward      string      `json:"taskAward"      orm:"task_award"       description:"领奖励任务配置"`          // 领奖励任务配置
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:"创建时间"`             // 创建时间
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:"更新时间"`             // 更新时间
	LastExecutedAt *gtime.Time `json:"lastExecutedAt" orm:"last_executed_at" description:"最近执行时间"`           // 最近执行时间
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"       description:"软删除时间"`            // 软删除时间
}
