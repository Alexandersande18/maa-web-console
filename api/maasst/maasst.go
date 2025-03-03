// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package maasst

import (
	"context"

	"maa-backend/api/maasst/v1"
)

type IMaasstV1 interface {
	LinkStart(ctx context.Context, req *v1.LinkStartReq) (res *v1.LinkStartRes, err error)
	Screenshot(ctx context.Context, req *v1.ScreenshotReq) (res *v1.ScreenshotRes, err error)
}
