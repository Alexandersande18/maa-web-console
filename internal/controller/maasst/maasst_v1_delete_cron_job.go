package maasst

import (
	"context"
	"time"

	"maa-backend/api/maasst/v1"
	"maa-backend/internal/dao"
	"maa-backend/internal/model/do"
	"maa-backend/internal/service"
	"fmt"
)

func (c *ControllerV1) DeleteCronJob(ctx context.Context, req *v1.DeleteCronJobReq) (res *v1.DeleteCronJobRes, err error) {
	_, err = dao.CronJob.Ctx(ctx).Data(do.CronJob{
		Name: fmt.Sprint(time.Now().UnixNano()),
	}).WherePri(req.Id).Update()
	if err != nil {
		return nil, err
	}
	_, err = dao.CronJob.Ctx(ctx).WherePri(req.Id).Delete()
	if err != nil {
		return nil, err
	}
	service.Croner().SetupCronJobs()
	return 
}
