package maa

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"net"
	"os/exec"
	"time"
)

type screenshotADBRunner func(context.Context, string, ...string) ([]byte, error)

func runScreenshotADB(ctx context.Context, path string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.WaitDelay = time.Second
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ADB %v: %w: %s", args, err, stderr.String())
	}
	return data, nil
}

// Preview does not need a connected MaaCore handle or access to its shared
// resources. This also lets accounts take screenshots before their first task.
func takeADBScreenshot(ctx context.Context, cfg *Asst, run screenshotADBRunner) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	args := []string{"-s", cfg.ADBAddress, "exec-out", "screencap", "-p"}
	data, err := run(ctx, cfg.ADBPath, args...)
	if err != nil {
		// TCP emulators may not yet have been connected by a queued task.
		if _, _, addressErr := net.SplitHostPort(cfg.ADBAddress); addressErr == nil && ctx.Err() == nil {
			if _, connectErr := run(ctx, cfg.ADBPath, "connect", cfg.ADBAddress); connectErr != nil {
				return nil, connectErr
			}
			data, err = run(ctx, cfg.ADBPath, args...)
		}
		if err != nil {
			return nil, err
		}
	}
	if _, err := png.DecodeConfig(bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("ADB returned an invalid PNG screenshot: %w", err)
	}
	return data, nil
}
