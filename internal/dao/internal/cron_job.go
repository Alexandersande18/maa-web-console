// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CronJobDao is the data access object for the table cron_job.
type CronJobDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of the current DAO.
	columns CronJobColumns // columns contains all the column names of Table for convenient usage.
}

// CronJobColumns defines and stores column names for the table cron_job.
type CronJobColumns struct {
	Id             string // 主键ID
	Name           string // 任务名称（唯一）
	Username       string // 所属用户名
	IsTemplate     string // 是否为模板任务（0-否 1-是）
	Hour           string // 小时（0-23）
	Minute         string // 分钟（0-59）
	TaskRecruit    string // 公招任务配置
	TaskInfrast    string // 基建任务配置
	TaskFight      string // 刷理智任务配置
	TaskRoguelike  string // 肉鸽任务配置
	TaskMall       string // 购物任务配置
	TaskAward      string // 领奖励任务配置
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
	LastExecutedAt string // 最近执行时间
	DeletedAt      string // 软删除时间
}

// cronJobColumns holds the columns for the table cron_job.
var cronJobColumns = CronJobColumns{
	Id:             "id",
	Name:           "name",
	Username:       "username",
	IsTemplate:     "is_template",
	Hour:           "hour",
	Minute:         "minute",
	TaskRecruit:    "task_recruit",
	TaskInfrast:    "task_infrast",
	TaskFight:      "task_fight",
	TaskRoguelike:  "task_roguelike",
	TaskMall:       "task_mall",
	TaskAward:      "task_award",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	LastExecutedAt: "last_executed_at",
	DeletedAt:      "deleted_at",
}

// NewCronJobDao creates and returns a new DAO object for table data access.
func NewCronJobDao() *CronJobDao {
	return &CronJobDao{
		group:   "default",
		table:   "cron_job",
		columns: cronJobColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CronJobDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CronJobDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CronJobDao) Columns() CronJobColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CronJobDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CronJobDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *CronJobDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
