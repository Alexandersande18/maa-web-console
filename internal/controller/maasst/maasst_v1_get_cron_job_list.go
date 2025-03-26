package maasst

import (
	"context"

	"maa-backend/api/maasst/v1"
	"maa-backend/internal/dao"
	"maa-backend/internal/middleware"
	// "maa-backend/internal/model/do"
	// "maa-backend/internal/model/dto"

	// "github.com/gogf/gf/v2/frame/g"
	// "github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) GetCronJobList(ctx context.Context, req *v1.GetCronJobListReq) (res *v1.GetCronJobListRes, err error) {
	res = &v1.GetCronJobListRes{}
	uid := gconv.String(ctx.Value(middleware.CtxUsername))
	err = dao.CronJob.Ctx(ctx).Where("username = ?", uid).Scan(&res.CronJobList)

	return 
}
