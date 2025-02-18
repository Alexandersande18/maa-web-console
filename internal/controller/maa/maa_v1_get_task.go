package maa

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"maa-backend/api/maa/v1"
	// "maa-backend/internal/logic/task"
)

func (c *ControllerV1) GetTask(ctx context.Context, req *v1.GetTaskReq) (res *v1.GetTaskRes, err error) {
	// tasks := task.TaskManagerInstance.GetPendingTasks(req.User, req.Device)
	res = &v1.GetTaskRes{
		Tasks: []v1.Task{},
	}
	return res, gerror.NewCode(gcode.CodeOK, "success")
}
