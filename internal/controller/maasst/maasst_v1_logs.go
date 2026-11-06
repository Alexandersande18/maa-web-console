package maasst

import (
	"context"
	v1 "maa-backend/api/maasst/v1"
	"maa-backend/internal/middleware"
	"maa-backend/internal/service"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

func Logs(ctx context.Context, r *ghttp.Request) {

}
func (c *ControllerV1) Logs(ctx context.Context, req *v1.LogsReq) (res *v1.LogsRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("X-Accel-Buffering", "no")
	// 创建客户端连接
	username, ok := ctx.Value(middleware.CtxUsername).(string)
	if !ok {
		// Handle error case when username is not a string
		g.Log().Error(ctx, "Username context value is not a string")
		r.Response.WriteStatus(http.StatusInternalServerError)
		return
	}
	clientId := gconv.String(gtime.TimestampNano())
	client, err := service.Maasst().NewSSEClient(username, clientId)
	if err != nil {
		// Handle error case when client creation fails
		g.Log().Error(ctx, "Failed to create SSE client: ", err)
		r.Response.WriteStatus(http.StatusInternalServerError)
		return
	}
	defer service.Maasst().DeleteSSEClient(username, clientId)
	// 保持连接
	g.Log().Info(ctx, "SSE client created: ", clientId)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			// 发送心跳
			r.Response.Write("event: heartbeat\ndata: {}\n\n")
			r.Response.Flush()
		case msg, ok := <-client.Ch:
			if !ok {
				return
			}
			// 发送日志消息
			data, _ := gjson.Encode(msg)
			r.Response.Writef("data: %s\n\n", data)
			r.Response.Flush()
		}
	}
}
