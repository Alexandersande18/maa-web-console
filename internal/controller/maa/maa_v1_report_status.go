package maa

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"maa-backend/api/maa/v1"
)

func (c *ControllerV1) ReportStatus(ctx context.Context, req *v1.ReportStatusReq) (res *v1.ReportStatusRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
