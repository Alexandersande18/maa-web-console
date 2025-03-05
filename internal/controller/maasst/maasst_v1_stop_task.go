package maasst

import (
	"context"
	"maa-backend/api/maasst/v1"
	"maa-backend/internal/middleware"
	"maa-backend/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) StopTask(ctx context.Context, req *v1.StopTaskReq) (res *v1.StopTaskRes, err error) {
	uid := gconv.String(ctx.Value(middleware.CtxUsername))
	if err = service.Maasst().StopTask(uid); err != nil {
		return
	}
	return
}
