package maasst

import (
	"context"

	"maa-backend/api/maasst/v1"
	"maa-backend/internal/dao"
	"maa-backend/internal/service"
)

func (c *ControllerV1) DeleteCronJob(ctx context.Context, req *v1.DeleteCronJobReq) (res *v1.DeleteCronJobRes, err error) {
	_, err = dao.CronJob.Ctx(ctx).WherePri(req.Id).Delete()
	service.Croner().SetupCronJobs()
	return 
}
