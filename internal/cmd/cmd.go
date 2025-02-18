package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"maa-backend/internal/controller/hello"
	"maa-backend/internal/controller/maa"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// s.Use(func (r *ghttp.Request) {
			// 	r.Middleware.Next()
			// 	if err := r.GetError(); err != nil {
			// 		// 记录到自定义错误日志文件
			// 		g.Log().Error(ctx, err)
			// 		//返回固定的友好信息
			// 		r.Response.ClearBuffer()
			// 		r.Response.Writeln("服务器居然开小差了，请稍后再试吧！")
			// 	}
			// })
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
				)
				group.Bind(
					maa.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
