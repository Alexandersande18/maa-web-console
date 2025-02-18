// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package maa

import (
	"context"

	"maa-backend/api/maa/v1"
)

type IMaaV1 interface {
	GetTask(ctx context.Context, req *v1.GetTaskReq) (res *v1.GetTaskRes, err error)
	ReportStatus(ctx context.Context, req *v1.ReportStatusReq) (res *v1.ReportStatusRes, err error)
}
