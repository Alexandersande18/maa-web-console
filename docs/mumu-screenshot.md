# MuMu 截图增强

每个 `maa.asst` 实例可配置：

```yaml
adbtype: MuMuEmulator12
mumu:
  enabled: true
  path: 'C:\Program Files\Netease\MuMu Player 12'
  # index: 0 # 网络桥接时显式指定；普通连接省略，按 ADB 端口推导。
```

`path` 是包含 `nx_device` 或 `shell` 的 MuMu 安装根目录，不是 MAA 资源目录。需要 Windows x64 版 MuMu 和提供 `AsstSetConnectionExtras` 的 MaaCore。后端检查 `external_renderer_ipc.dll` 是否存在，配置错误会在启动时报告。

后端在每次加载资源后的连接阶段设置 `MuMuEmulator12` extras，传递安装路径、客户端类型和可选实例编号。截图增强不改变触控模式。extras 是 MaaCore 进程级状态，因此每个账号连接前重新设置；未开启的 MuMu 账号会清空 extras，避免继承其他账号的设置。

修改配置后重新启动后端。运行任务后，在 `maa.userdir/debug/asst.log` 的 `FastestWayToScreencap` 回调中确认 `method` 为 `MumuExtras`；DLL 成功加载本身不代表截图增强已生效。如果核心回退到 ADB，仍可正常运行，需检查 MuMu 版本、路径和实例编号。网页截图预览继续使用独立的 ADB PNG 路径。

本机实测中，桌面用户会话可以使用增强，但以 `LocalSystem` 身份在 Session 0 运行时，MuMu SDK 报告 `screencap fail` 并回退到 ADB。部署时应在与 MuMu 相同的已登录用户会话内运行后端，例如使用登录触发的计划任务；该方式依赖用户保持登录。仅更改服务登录账号仍不能解决 Session 0 隔离。

当前本机部署已恢复为 Session 0 Windows 服务，两个实例均关闭 MuMu 增强，使用普通 ADB 截图；未安装登录启动任务。日常运行和更新见 [Windows 服务运维说明](service-operations.md)。

普通测试：`go test ./internal/logic/maa`。实际设备验证只连接并截图，不启动任务：

```powershell
$env:MAA_MUMU_SMOKE = '1'
$env:MAA_MUMU_SMOKE_CONFIG = 'C:\Sources\maa-backend\manifest\config\config.yaml'
$env:MAA_MUMU_SMOKE_OUTPUT = 'C:\Sources\maa-backend\.runtime\mumu-smoke'
go test ./internal/logic/maa -run '^TestMuMuNativeScreenshots$' -count=1 -timeout 180s -v
```

仅在账号队列空闲时运行实际设备验证。验证为每个账号分别测量关闭、开启增强时的截图，要求开启后实际使用 `MumuExtras`，检查 PNG 和设备 UUID，并输出 `screenshots.json` 与截图样本。输出目录包含游戏画面，请勿提交到版本库。

上游接口和参数参见 [AsstCaller.h](https://github.com/MaaAssistantArknights/MaaAssistantArknights/blob/v6.21.4/include/AsstCaller.h) 和 [MuMu12Extra.cs](https://github.com/MaaAssistantArknights/MaaAssistantArknights/blob/v6.21.4/src/MaaWpfGui/Models/EmulatorConnectionExtra/MuMu12Extra.cs)。
