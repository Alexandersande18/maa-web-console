package maa

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

var gamePackages = map[string]string{
	"Official": "com.hypergryph.arknights", "Bilibili": "com.hypergryph.arknights.bilibili",
	"YoStarEN": "com.YoStarEN.Arknights", "YoStarJP": "com.YoStarJP.Arknights",
	"YoStarKR": "com.YoStarKR.Arknights", "txwy": "tw.txwy.and.arknights",
}

func (m *Maa) CloseGame(ctx context.Context) error {
	pkg := gamePackages[m.cfg.ClientType]
	if pkg == "" {
		return fmt.Errorf("unsupported client type: %s", m.cfg.ClientType)
	}
	cmd := exec.CommandContext(ctx, m.cfg.ADBPath, "-s", m.cfg.ADBAddress, "shell", "am", "force-stop", pkg)
	cmd.WaitDelay = time.Second
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ADB force-stop failed for %s: %w: %s", m.cfg.ADBAddress, err, output)
	}
	return nil
}
