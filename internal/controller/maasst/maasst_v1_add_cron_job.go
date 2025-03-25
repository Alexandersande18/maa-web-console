package maasst

import (
	"context"

	"maa-backend/api/maasst/v1"
	"maa-backend/internal/dao"
	"maa-backend/internal/middleware"
	"maa-backend/internal/model/do"
	"maa-backend/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) AddCronJob(ctx context.Context, req *v1.AddCronJobReq) (res *v1.AddCronJobRes, err error) {
	uid := gconv.String(ctx.Value(middleware.CtxUsername))
	insertId, err := dao.CronJob.Ctx(ctx).Data(do.CronJob{
		Username: uid,
		Name: req.Name,
		IsTemplate: func() int {
			if req.IsTemplate == 0 {
				return 0
			}
			return 1
		}(),
		Hour: req.Hour,
		Minute: req.Minute,
		TaskRecruit: req.TaskRecruit,
		TaskInfrast: req.TaskInfrast,
		TaskFight: req.TaskFight,
		TaskMall: req.TaskMall,
		TaskAward: req.TaskAward,
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	res = &v1.AddCronJobRes{
		Id: int(insertId),
	}
	
	service.Croner().SetupCronJobs()

	return
}
