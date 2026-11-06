package maasst

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"maa-backend/api/maasst/v1"
	"maa-backend/internal/auth"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRes, error) {
	token, err := auth.Login(req.Token)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "Invalid credentials")
	}
	return &v1.LoginRes{Token: token}, nil
}
