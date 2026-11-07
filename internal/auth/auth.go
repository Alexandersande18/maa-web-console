// Package auth loads deployment credentials without embedding them in source code.
package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type Config struct {
	JWTSecret              string `json:"jwtsecret"`
	Users                  []User `json:"users"`
	AllowLegacyCredentials bool   `json:"allowlegacycredentials"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// current is initialized before the HTTP server starts, and then stays immutable.
var current *Config

func Initialize(ctx context.Context) error {
	value, err := g.Cfg().Get(ctx, "auth")
	if err != nil {
		return fmt.Errorf("read auth configuration: %w", err)
	}
	cfg := &Config{}
	if err = value.Scan(cfg); err != nil {
		return fmt.Errorf("decode auth configuration: %w", err)
	}
	if secret, ok := os.LookupEnv("MAA_JWT_SECRET"); ok {
		cfg.JWTSecret = secret
	}
	if tokens, ok := os.LookupEnv("MAA_LOGIN_TOKENS"); ok {
		var users map[string]string
		if err = json.Unmarshal([]byte(tokens), &users); err != nil {
			return fmt.Errorf("MAA_LOGIN_TOKENS must be a JSON object of instance IDs and login tokens")
		}
		cfg.Users = nil
		for username, token := range users {
			cfg.Users = append(cfg.Users, User{Username: username, Token: token})
		}
	}
	if err = cfg.validate(); err != nil {
		return err
	}
	var instances []struct{ ID string }
	value, err = g.Cfg().Get(ctx, "maa.asst")
	if err != nil {
		return err
	}
	if err = value.Scan(&instances); err != nil {
		return err
	}
	for _, user := range cfg.Users {
		found := false
		for _, instance := range instances {
			found = found || instance.ID == user.Username
		}
		if !found {
			return fmt.Errorf("auth user %q must match a maa.asst instance ID", user.Username)
		}
	}
	current = cfg
	return nil
}

func (c *Config) validate() error {
	minimumLength := 32
	// An existing private deployment may explicitly retain its previous login
	// credentials and signing key. Public examples keep the stronger default.
	if c.AllowLegacyCredentials {
		minimumLength = 1
	}
	if len(c.JWTSecret) < minimumLength {
		return fmt.Errorf("auth.jwtsecret / MAA_JWT_SECRET must contain at least %d bytes", minimumLength)
	}
	if len(c.Users) == 0 {
		return fmt.Errorf("configure at least one auth.users entry / MAA_LOGIN_TOKENS")
	}
	names, tokens := map[string]bool{}, map[string]bool{}
	for _, user := range c.Users {
		if strings.TrimSpace(user.Username) == "" || len(user.Token) < minimumLength {
			return fmt.Errorf("each auth user needs a username and a login token of at least %d bytes", minimumLength)
		}
		if names[user.Username] || tokens[user.Token] || user.Token == c.JWTSecret {
			return fmt.Errorf("auth usernames, login tokens and JWT secret must be distinct")
		}
		names[user.Username], tokens[user.Token] = true, true
	}
	return nil
}

func (c *Config) authenticate(token string) (string, bool) {
	digest := sha256.Sum256([]byte(token))
	username := ""
	for _, user := range c.Users {
		expected := sha256.Sum256([]byte(user.Token))
		if subtle.ConstantTimeCompare(digest[:], expected[:]) == 1 {
			username = user.Username
		}
	}
	return username, username != ""
}

func Login(token string) (string, error) {
	if current == nil {
		return "", fmt.Errorf("authentication is not configured")
	}
	username, ok := current.authenticate(token)
	if !ok {
		return "", fmt.Errorf("invalid credentials")
	}
	now := time.Now()
	claims := &Claims{Username: username, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now), NotBefore: jwt.NewNumericDate(now),
	}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(current.JWTSecret))
}

func (c *Config) parse(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(c.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}
	for _, user := range c.Users {
		if claims.Username == user.Username {
			return claims, nil
		}
	}
	return nil, fmt.Errorf("unknown user")
}

func Parse(token string) (*Claims, error) {
	if current == nil {
		return nil, fmt.Errorf("authentication is not configured")
	}
	return current.parse(token)
}
