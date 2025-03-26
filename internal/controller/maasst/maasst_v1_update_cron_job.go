package maasst

import (
	"context"

	"maa-backend/api/maasst/v1"
	"maa-backend/internal/dao"
	"maa-backend/internal/model/do"
	"maa-backend/internal/service"
)

func (c *ControllerV1) UpdateCronJob(ctx context.Context, req *v1.UpdateCronJobReq) (res *v1.UpdateCronJobRes, err error) {
	_, err = dao.CronJob.Ctx(ctx).Data(do.CronJob{
		Hour: req.Hour,
		Minute: req.Minute,
		TaskRecruit: req.TaskRecruit,
		TaskInfrast: req.TaskInfrast,
		TaskFight: req.TaskFight,
		TaskMall: req.TaskMall,
		TaskAward: req.TaskAward,
		Name: req.Name,
		IsTemplate: req.IsTemplate,
	}).WherePri(req.Id).Update()
	if err != nil {
		return
	}

	service.Croner().SetupCronJobs()
	
	return 
}
