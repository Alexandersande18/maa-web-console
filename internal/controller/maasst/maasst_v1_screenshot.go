package maasst

import (
	"context"
	"encoding/base64"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"maa-backend/api/maasst/v1"
	"maa-backend/internal/middleware"
	"maa-backend/internal/service"
)

func (c *ControllerV1) Screenshot(ctx context.Context, req *v1.ScreenshotReq) (res *v1.ScreenshotRes, err error) {
	res = &v1.ScreenshotRes{}
	// 获取截图二进制数据
	uid := gconv.String(ctx.Value(middleware.CtxUsername))
	g.Log().Debug(context.TODO(), "uId: ", uid)
	imgData, err := service.Maasst().TakeScreenshot(uid)
	if err != nil {
		err = gerror.WrapCode(gcode.CodeInternalError, err, "截图失败")
		return
	}

	// 将二进制数据编码为 Base64
	base64Data := base64.StdEncoding.EncodeToString(imgData)
	// 构造 JSON 响应
	res.ImageData = base64Data
	res.Timestamp = gtime.TimestampMilli()
	return
}
