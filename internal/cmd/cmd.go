package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"maa-backend/internal/controller/maasst"
	"maa-backend/internal/middleware"
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
            // 设置静态文件服务
            s.AddStaticPath("/", "resource/public/html")


            s.Group("/", func(group *ghttp.RouterGroup) {
                group.Middleware(ghttp.MiddlewareHandlerResponse)
                group.Bind(maasst.NewV1().Login)
                group.Middleware(middleware.JWTAuth)
                group.Bind(
                    maasst.NewV1().LinkStart,
                    maasst.NewV1().Screenshot,
                    maasst.NewV1().StopTask,
                )
            })
            s.Group("/", func(group *ghttp.RouterGroup) {
                group.Middleware(middleware.JWTAuth)
                group.Bind(
                    maasst.NewV1().Logs,
                )
            })
            s.SetSwaggerUITemplate(MySwaggerUITemplate)
            s.Run()
            return nil
        },
    }
)
