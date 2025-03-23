package maasst

import (
	"context"
	v1 "maa-backend/api/maasst/v1"
	"maa-backend/internal/middleware"
	"maa-backend/internal/service"
	"net/http"
	"time"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
)

func (c *ControllerV1) LogsWS(ctx context.Context, req *v1.LogsWSReq) (res *v1.LogsWSRes, err error) {
	r := g.RequestFromCtx(ctx)
	g.Log().Info(ctx, "WebSocket 连接成功")
	// 升级 HTTP 请求为 WebSocket 连接
	upgrader := websocket.Upgrader{
		// 此处简单允许所有来源，实际生产环境请根据需要进行校验
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		g.Log().Error(ctx, "WebSocket 升级失败: ", err)
		return
	}
	defer conn.Close()
	// 创建客户端连接

	username, ok := ctx.Value(middleware.CtxUsername).(string)
	if !ok {
		// Handle error case when username is not a string
		g.Log().Error(ctx, "Username context value is not a string")
		r.Response.WriteStatus(http.StatusInternalServerError)
		return
	}
	
	// 创建日志客户端，复用你原来 SSE 客户端相关逻辑
	clientId := gconv.String(gtime.TimestampNano())
	client, err := service.Maasst().NewSSEClient(username, clientId)
	if err != nil {
		g.Log().Error(ctx, "创建日志客户端失败: ", err)
		conn.WriteMessage(websocket.CloseMessage, []byte("创建日志客户端失败"))
		return
	}
	// 确保连接关闭时删除客户端
	defer service.Maasst().DeleteSSEClient(username, clientId)
	g.Log().Info(ctx, "WebSocket 客户端创建成功: ", clientId)

	// 定义一个定时器发送心跳
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// 循环监听客户端日志消息和心跳定时器
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			// 发送心跳包，心跳包格式可根据前后端约定进行调整
			heartbeat := map[string]interface{}{
				"event": "heartbeat",
				"data":  "{}",
			}
			if err = conn.WriteJSON(heartbeat); err != nil {
				g.Log().Error(ctx, "发送心跳包失败: ", err)
				return
			}
		case msg, ok := <-client.Ch:
			if !ok {
				// 如果日志管道关闭则退出循环
				g.Log().Info(ctx, "日志管道已关闭")
				return
			}
			// 直接将日志对象以 JSON 方式发送给前端
			if err = conn.WriteJSON(msg); err != nil {
				g.Log().Error(ctx, "发送日志消息失败: ", err)
				return
			}
		}
	}
}
