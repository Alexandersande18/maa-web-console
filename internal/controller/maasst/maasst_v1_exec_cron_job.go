package maasst

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"maa-backend/api/maasst/v1"
	"maa-backend/internal/dao"
	"maa-backend/internal/middleware"
	"maa-backend/internal/model/dto"
	"maa-backend/internal/service"
)

func (c *ControllerV1) ExecCronJob(ctx context.Context, req *v1.ExecCronJobReq) (res *v1.ExecCronJobRes, err error) {
	job := &dto.CronJob{}
	err = dao.CronJob.Ctx(ctx).WherePri(req.Id).Where("username", ctx.Value(middleware.CtxUsername)).Scan(job)
	if err != nil {
		return
	}
	if job.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "Task not found")
	}
	err = service.Croner().ExecCronJob(job)
	return
}
