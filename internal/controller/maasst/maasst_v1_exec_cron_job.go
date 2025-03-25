package maasst

import (
	"context"

	"maa-backend/api/maasst/v1"
	"maa-backend/internal/dao"
	"maa-backend/internal/model/dto"
	"maa-backend/internal/service"
)

func (c *ControllerV1) ExecCronJob(ctx context.Context, req *v1.ExecCronJobReq) (res *v1.ExecCronJobRes, err error) {
	job := &dto.CronJob{}
	err = dao.CronJob.Ctx(ctx).WherePri(req.Id).Scan(job)
	if err != nil {
		return
	}
	err = service.Croner().ExecCronJob(job)
	return
}
