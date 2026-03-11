// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CronJob is the golang structure of table cron_job for DAO operations like Where/Data.
type CronJob struct {
	g.Meta         `orm:"table:cron_job, do:true"`
	Id             interface{} // 主键ID
	Name           interface{} // 任务名称（唯一）
	Username       interface{} // 所属用户名
	IsTemplate     interface{} // 是否为模板任务（0-否 1-是）
	Hour           interface{} // 小时（0-23）
	Minute         interface{} // 分钟（0-59）
	TaskRecruit    interface{} // 公招任务配置
	TaskInfrast    interface{} // 基建任务配置
	TaskFight      interface{} // 刷理智任务配置
	TaskRoguelike  interface{} // 肉鸽任务配置
	TaskMall       interface{} // 购物任务配置
	TaskAward      interface{} // 领奖励任务配置
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
	LastExecutedAt *gtime.Time // 最近执行时间
	DeletedAt      *gtime.Time // 软删除时间
}
