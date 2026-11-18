package maasst

import (
	"context"
	v1 "maa-backend/api/maasst/v1"
	"maa-backend/internal/middleware"
	"maa-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) LinkStart(ctx context.Context, req *v1.LinkStartReq) (res *v1.LinkStartRes, err error) {
	taskList, err := req.TaskSelection.Build(false)
	if err != nil {
		return nil, err
	}

	uid := gconv.String(ctx.Value(middleware.CtxUsername))
	g.Log().Debug(context.TODO(), "uId: ", uid)

	if err = service.Maasst().LinkStart(uid, taskList); err != nil {

		return
	}
	return &v1.LinkStartRes{Message: "任务已加入队列，将按顺序执行"}, nil
}
