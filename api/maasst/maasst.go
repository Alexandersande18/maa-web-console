// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package maasst

import (
	"context"

	"maa-backend/api/maasst/v1"
)

type IMaasstV1 interface {
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	LinkStart(ctx context.Context, req *v1.LinkStartReq) (res *v1.LinkStartRes, err error)
	StopTask(ctx context.Context, req *v1.StopTaskReq) (res *v1.StopTaskRes, err error)
	Screenshot(ctx context.Context, req *v1.ScreenshotReq) (res *v1.ScreenshotRes, err error)
	Logs(ctx context.Context, req *v1.LogsReq) (res *v1.LogsRes, err error)
}
