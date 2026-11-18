package maasst

import (
	"context"

	v1 "maa-backend/api/maasst/v1"
	"maa-backend/internal/dao"
	"maa-backend/internal/middleware"
	"maa-backend/internal/model/do"
	"maa-backend/internal/service"
)

func (c *ControllerV1) UpdateCronJob(ctx context.Context, req *v1.UpdateCronJobReq) (res *v1.UpdateCronJobRes, err error) {
	if _, err := req.TaskSelection.Build(req.IsTemplate == 0); err != nil {
		return nil, err
	}
	_, err = dao.CronJob.Ctx(ctx).Data(do.CronJob{
		Hour:            req.Hour,
		Minute:          req.Minute,
		TaskRecruit:     req.TaskRecruit,
		TaskInfrast:     req.TaskInfrast,
		TaskFight:       req.TaskFight,
		TaskRoguelike:   req.TaskRoguelike,
		TaskMall:        req.TaskMall,
		TaskAward:       req.TaskAward,
		TaskActivity:    req.TaskActivity,
		TaskReclamation: req.TaskReclamation,
		Name:            req.Name,
		IsTemplate:      req.IsTemplate,
	}).WherePri(req.Id).Where("username", ctx.Value(middleware.CtxUsername)).Update()
	if err != nil {
		return
	}

	service.Croner().SetupCronJobs()

	return
}
