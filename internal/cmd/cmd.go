package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"maa-backend/internal/auth"
	"maa-backend/internal/controller/maasst"
	"maa-backend/internal/dao"
	maalogic "maa-backend/internal/logic/maa"
	"maa-backend/internal/middleware"
	"maa-backend/internal/service"
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
			if err = auth.Initialize(ctx); err != nil {
				return err
			}
			if err = dao.CheckCronJobSchema(ctx); err != nil {
				return err
			}
			if err = maalogic.Initialize(); err != nil {
				return err
			}
			service.Croner().SetupCronJobs()
			s := g.Server()
			// 设置静态文件服务
			s.AddStaticPath("/", "resource/public/html")

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(maasst.NewV1().Login)
				group.Middleware(middleware.JWTAuth)
				group.Middleware(middleware.LogMiddleware)
				group.Bind(
					maasst.NewV1().LinkStart,
					maasst.NewV1().Screenshot,
					maasst.NewV1().StopTask,
					maasst.NewV1().AddCronJob,
					maasst.NewV1().DeleteCronJob,
					maasst.NewV1().GetCronJobList,
					maasst.NewV1().UpdateCronJob,
					maasst.NewV1().ExecCronJob,
				)
			})
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.JWTAuth)
				group.Middleware(middleware.LogMiddleware)
				group.Bind(
					maasst.NewV1().Logs,
					maasst.NewV1().LogsWS,
				)
			})
			s.SetSwaggerUITemplate(MySwaggerUITemplate)
			s.Run()
			return nil
		},
	}
)
