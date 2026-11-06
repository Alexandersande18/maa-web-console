package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"maa-backend/internal/auth"
	"net/http"
	"strings"
)

const CtxUsername gctx.StrKey = "username"

func JWTAuth(r *ghttp.Request) {
	token := r.Header.Get("Authorization")
	// Browser WebSocket clients cannot set the Authorization header.
	if token == "" && r.URL.Path == "/maasst/logs-ws" {
		token = r.URL.Query().Get("token")
	}
	claims, err := auth.Parse(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		r.Response.Status = http.StatusUnauthorized
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "Invalid token"))
		return
	}
	r.SetCtxVar(CtxUsername, claims.Username)
	r.Middleware.Next()
}
