package maa

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/png"
	"reflect"
	"testing"
	"time"
)

func TestScreenshotBeforeFirstTaskUsesAccountDevice(t *testing.T) {
	var img bytes.Buffer
	if err := png.Encode(&img, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	cfg := &Asst{ADBPath: "adb.exe", ADBAddress: "127.0.0.1:16416"}
	var calls [][]string
	run := func(ctx context.Context, path string, args ...string) ([]byte, error) {
		deadline, ok := ctx.Deadline()
		if !ok || time.Until(deadline) > 15*time.Second || path != cfg.ADBPath {
			t.Fatal("missing screenshot timeout or wrong ADB binary")
		}
		calls = append(calls, args)
		if len(calls) == 1 {
			return nil, errors.New("device not found")
		}
		if len(calls) == 2 {
			return []byte("connected"), nil
		}
		return img.Bytes(), nil
	}
	got, err := takeADBScreenshot(context.Background(), cfg, run)
	want := [][]string{
		{"-s", cfg.ADBAddress, "exec-out", "screencap", "-p"},
		{"connect", cfg.ADBAddress},
		{"-s", cfg.ADBAddress, "exec-out", "screencap", "-p"},
	}
	if err != nil || !bytes.Equal(got, img.Bytes()) || !reflect.DeepEqual(calls, want) {
		t.Fatal(err, calls)
	}
}

func TestScreenshotRejectsInvalidImageAndRespectsCancellation(t *testing.T) {
	cfg := &Asst{ADBAddress: "127.0.0.1:16416"}
	if _, err := takeADBScreenshot(context.Background(), cfg, func(context.Context, string, ...string) ([]byte, error) {
		return []byte("not a screenshot"), nil
	}); err == nil {
		t.Fatal("accepted invalid screenshot")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	_, err := takeADBScreenshot(ctx, cfg, func(ctx context.Context, _ string, _ ...string) ([]byte, error) {
		calls++
		return nil, ctx.Err()
	})
	if !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatal("canceled screenshot attempted to reconnect", err, calls)
	}
	if _, err := (&sMaasst{}).TakeScreenshot(context.Background(), "missing"); err == nil {
		t.Fatal("accepted unknown account")
	}
}
