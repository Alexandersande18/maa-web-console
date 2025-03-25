package middleware

import (
    "time"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
)

// LogMiddleware 日志中间件
func LogMiddleware(r *ghttp.Request) {
    start := time.Now()

    // 执行后续逻辑（包括其他中间件和路由处理函数）
    r.Middleware.Next()

    // 计算请求耗时
    latency := time.Since(start)

    // 记录访问日志
    g.Log().Info(r.Context(),
        "path:", r.Request.URL.Path,
        "method:", r.Method,
        "status:", r.Response.Status,
        "latency:", latency.String(),
        "client_ip:", r.GetClientIp(),
		"user:", r.GetCtxVar(CtxUsername),
    )
}