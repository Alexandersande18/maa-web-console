package main

import (
	_ "maa-backend/internal/packed"

	_ "maa-backend/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"maa-backend/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
