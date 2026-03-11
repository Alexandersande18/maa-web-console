package maasst

import (
	"context"
	v1 "maa-backend/api/maasst/v1"
	"maa-backend/internal/middleware"
	"maa-backend/internal/model/dto"
	"maa-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) LinkStart(ctx context.Context, req *v1.LinkStartReq) (res *v1.LinkStartRes, err error) {
	var taskList []dto.Tasker
	res = &v1.LinkStartRes{}
	taskList = append(taskList, &dto.TaskStartUp{
		Enable:           true,
		ClientType:       "Official",
		StartGameEnabled: true,
	})
	if req.TaskRecruit != nil && req.TaskRecruit.Enable {
		taskList = append(taskList, req.TaskRecruit)
	}
	if req.TaskInfrast != nil && req.TaskInfrast.Enable {
		taskList = append(taskList, req.TaskInfrast)
	}
	if req.TaskFight != nil && req.TaskFight.Enable {
		taskList = append(taskList, req.TaskFight)
	}
	if req.TaskRoguelike != nil && req.TaskRoguelike.Enable {
		taskList = append(taskList, req.TaskRoguelike)
	}
	if req.TaskMall != nil && req.TaskMall.Enable {
		taskList = append(taskList, req.TaskMall)
	}
	if req.TaskAward != nil && req.TaskAward.Enable {
		taskList = append(taskList, req.TaskAward)
	}
	taskList = append(taskList, &dto.TaskCloseDown{
		Enable:     true,
		ClientType: "Official",
	})

	uid := gconv.String(ctx.Value(middleware.CtxUsername))
	g.Log().Debug(context.TODO(), "uId: ", uid)

	if err = service.Maasst().LinkStart(uid, taskList); err != nil {

		return
	}
	return
}
