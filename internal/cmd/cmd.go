package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"maa-backend/internal/controller/maasst"
	"maa-backend/internal/controller/hello"
	"maa-backend/internal/controller/maa"
)

const (
    MySwaggerUITemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI"/>
    <title>SwaggerUI</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui.min.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui-bundle.js" crossorigin></script>
<script>
    window.onload = () => {
        window.ui = SwaggerUIBundle({
            url:    '{SwaggerUIDocUrl}',
            dom_id: '#swagger-ui',
        });
    };
</script>
</body>
</html>
`
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
				group.Bind(
					maasst.NewV1(),
				)
			})
			s.SetSwaggerUITemplate(MySwaggerUITemplate)
			s.Run()
			return nil
		},
	}
)
