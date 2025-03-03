package maasst

import (
	"context"
	"encoding/base64"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"maa-backend/api/maasst/v1"
	"maa-backend/internal/service"
)

func (c *ControllerV1) Screenshot(ctx context.Context, req *v1.ScreenshotReq) (res *v1.ScreenshotRes, err error) {
	// 获取截图二进制数据
	imgData, err := service.Maasst().TakeScreenshot(req.UserID)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "截图失败: "+err.Error())
	}

	// 将二进制数据编码为 Base64
	base64Data := base64.StdEncoding.EncodeToString(imgData)

	// 构造 JSON 响应
	return &v1.ScreenshotRes{
		ImageData: base64Data,
		Timestamp: gtime.TimestampMilli(),
	}, nil
}
