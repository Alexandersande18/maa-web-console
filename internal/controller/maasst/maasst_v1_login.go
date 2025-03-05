package maasst

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"maa-backend/api/maasst/v1"
	"github.com/golang-jwt/jwt/v5"
	"maa-backend/internal/middleware"
)

var mp = map[string]string{
	"REDACTED_CREDENTIAL_18": "my",
	"REDACTED_CREDENTIAL_17": "tb",
}

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	// In a real application, validate credentials against a database
	// This is just a demo with hardcoded values
	uname := "my"
	if v, ok := mp[req.Token]; !ok {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "Invalid credentials")
	} else {
		uname = v
	}

	// Create claims with user information
	claims := &middleware.JWTClaims{
		Username: uname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(middleware.JwtSecretKey))
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "Failed to generate token")
	}

	return &v1.LoginRes{
		Token: signedToken,
	}, nil
}
