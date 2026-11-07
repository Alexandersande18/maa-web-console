package auth

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/golang-jwt/jwt/v5"
)

func testConfig() *Config {
	return &Config{JWTSecret: strings.Repeat("signing-test-only-", 4), Users: []User{
		{Username: "alice", Token: strings.Repeat("alice-test-only-", 4)},
		{Username: "bob", Token: strings.Repeat("bob-test-only-", 4)},
	}}
}

func TestCredentialValidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*Config)
	}{
		{"empty secret", func(c *Config) { c.JWTSecret = "" }},
		{"short secret", func(c *Config) { c.JWTSecret = "short" }},
		{"no users", func(c *Config) { c.Users = nil }},
		{"empty username", func(c *Config) { c.Users[0].Username = " " }},
		{"short login token", func(c *Config) { c.Users[0].Token = "short" }},
		{"duplicate user", func(c *Config) { c.Users[1].Username = c.Users[0].Username }},
		{"duplicate token", func(c *Config) { c.Users[1].Token = c.Users[0].Token }},
		{"shared signing secret", func(c *Config) { c.Users[0].Token = c.JWTSecret }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := testConfig()
			tc.change(cfg)
			if cfg.validate() == nil {
				t.Fatal("unsafe configuration accepted")
			}
		})
	}
	if err := testConfig().validate(); err != nil {
		t.Fatal(err)
	}
}

func TestLoginUsesConfiguredIdentity(t *testing.T) {
	cfg := testConfig()
	for _, user := range cfg.Users {
		name, ok := cfg.authenticate(user.Token)
		if !ok || name != user.Username {
			t.Fatal("wrong user for configured credential")
		}
	}
	for _, token := range []string{"", "wrong", cfg.Users[0].Token + "x", cfg.JWTSecret} {
		if _, ok := cfg.authenticate(token); ok {
			t.Fatal("invalid login accepted")
		}
	}
}

func TestLegacyCredentialsRequireExplicitOptIn(t *testing.T) {
	cfg := &Config{JWTSecret: "legacy-signing-fixture", Users: []User{
		{Username: "alice", Token: "old-login"},
	}}
	if cfg.validate() == nil {
		t.Fatal("short credentials accepted by default")
	}
	cfg.AllowLegacyCredentials = true
	if err := cfg.validate(); err != nil {
		t.Fatal(err)
	}
	oldSession, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username: "alice", RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}).SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cfg.parse(oldSession); err != nil {
		t.Fatal("existing unexpired session was rejected", err)
	}
	if name, ok := cfg.authenticate("old-login"); !ok || name != "alice" {
		t.Fatal("legacy login no longer maps to its original instance")
	}
	cfg.JWTSecret = ""
	if cfg.validate() == nil {
		t.Fatal("legacy compatibility must still reject empty signing keys")
	}
	cfg.JWTSecret = "legacy-signing-fixture"
	cfg.Users[0].Token = ""
	if cfg.validate() == nil {
		t.Fatal("legacy compatibility must still reject empty login tokens")
	}
}

func TestJWTRejectsForgedExpiredAndUnknownIdentity(t *testing.T) {
	cfg := testConfig()
	for _, tc := range []struct {
		name, username string
		method         jwt.SigningMethod
		key            interface{}
		expiry         *jwt.NumericDate
		valid          bool
	}{
		{"valid", "alice", jwt.SigningMethodHS256, []byte(cfg.JWTSecret), jwt.NewNumericDate(time.Now().Add(time.Hour)), true},
		{"wrong signature", "alice", jwt.SigningMethodHS256, []byte("other-key"), jwt.NewNumericDate(time.Now().Add(time.Hour)), false},
		{"wrong algorithm", "alice", jwt.SigningMethodHS512, []byte(cfg.JWTSecret), jwt.NewNumericDate(time.Now().Add(time.Hour)), false},
		{"unsigned", "alice", jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, jwt.NewNumericDate(time.Now().Add(time.Hour)), false},
		{"expired", "alice", jwt.SigningMethodHS256, []byte(cfg.JWTSecret), jwt.NewNumericDate(time.Now().Add(-time.Hour)), false},
		{"missing expiry", "alice", jwt.SigningMethodHS256, []byte(cfg.JWTSecret), nil, false},
		{"unknown user", "mallory", jwt.SigningMethodHS256, []byte(cfg.JWTSecret), jwt.NewNumericDate(time.Now().Add(time.Hour)), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			signed, err := jwt.NewWithClaims(tc.method, &Claims{Username: tc.username, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: tc.expiry}}).SignedString(tc.key)
			if err != nil {
				t.Fatal(err)
			}
			claims, err := cfg.parse(signed)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v, error=%v", tc.valid, err)
			}
			if tc.valid && claims.Username != tc.username {
				t.Fatal("wrong claims")
			}
		})
	}
}

func TestInitializeAndEnvironmentOverrides(t *testing.T) {
	originalAdapter, originalConfig := g.Cfg().GetAdapter(), current
	t.Cleanup(func() { g.Cfg().SetAdapter(originalAdapter); current = originalConfig })
	// Isolate these tests from credentials in the developer's environment.
	for _, name := range []string{"MAA_JWT_SECRET", "MAA_LOGIN_TOKENS"} {
		value, exists := os.LookupEnv(name)
		os.Unsetenv(name)
		t.Cleanup(func() {
			if exists {
				os.Setenv(name, value)
			} else {
				os.Unsetenv(name)
			}
		})
	}
	cfg := testConfig()
	content, err := json.Marshal(map[string]interface{}{
		"auth": cfg, "maa": map[string]interface{}{"asst": []map[string]string{{"id": "alice"}, {"id": "bob"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	adapter, err := gcfg.NewAdapterContent(string(content))
	if err != nil {
		t.Fatal(err)
	}
	g.Cfg().SetAdapter(adapter)
	if err = Initialize(context.Background()); err != nil {
		t.Fatal(err)
	}
	signed, err := Login(cfg.Users[0].Token)
	if err != nil {
		t.Fatal(err)
	}
	if claims, err := Parse(signed); err != nil || claims.Username != "alice" {
		t.Fatal("login round trip failed", err)
	}
	t.Setenv("MAA_JWT_SECRET", strings.Repeat("new-signing-test-only-", 3))
	t.Setenv("MAA_LOGIN_TOKENS", `{"bob":"`+strings.Repeat("new-login-test-only-", 3)+`"}`)
	if err = Initialize(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err = Parse(signed); err == nil {
		t.Fatal("rotated signing key accepted old JWT")
	}
	if _, err = Login(cfg.Users[0].Token); err == nil {
		t.Fatal("removed user can still log in")
	}
	t.Setenv("MAA_LOGIN_TOKENS", `{"unknown":"`+strings.Repeat("new-login-test-only-", 3)+`"}`)
	if err = Initialize(context.Background()); err == nil {
		t.Fatal("unknown instance accepted")
	}
	t.Setenv("MAA_LOGIN_TOKENS", "not-json")
	if err = Initialize(context.Background()); err == nil {
		t.Fatal("malformed override accepted")
	}
}

func TestExampleConfigParsesWithoutUsableCredentials(t *testing.T) {
	adapter, err := gcfg.NewAdapterFile("../../manifest/config/config.example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	c := gcfg.NewWithAdapter(adapter)
	value, err := c.Get(context.Background(), "auth")
	if err != nil {
		t.Fatal(err)
	}
	var cfg Config
	if err = value.Scan(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.validate() == nil {
		t.Fatal("public example must not contain usable credentials")
	}
	value, err = c.Get(context.Background(), "maa.asst.0.adbpath")
	if err != nil || value.String() != `C:\platform-tools\adb.exe` {
		t.Fatal("invalid example ADB path")
	}
}
