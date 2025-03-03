package maasst

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"maa-backend/api/maasst/v1"
)

func (c *ControllerV1) LinkStart(ctx context.Context, req *v1.LinkStartReq) (res *v1.LinkStartRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
