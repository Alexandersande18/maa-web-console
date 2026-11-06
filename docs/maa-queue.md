# MAA 资源与任务队列

手动启动、立即执行定时任务和自动定时任务共用一个 FIFO 队列。全进程最多运行一批任务；一批任务包含开始唤醒、选中的各任务和关闭游戏。`LinkStart` 成功表示已受理排队，实际开始、失败、超时及结束通过账号日志通知。

同一账号已有执行中或排队任务时，重复提交返回错误。停止接口会取消该账号的排队任务，或停止该账号当前任务，不会停止其他账号。定时任务不再抢占正在执行的任务。

## 配置

在 `manifest/config/config.yaml` 设置，修改后重启服务：

```yaml
maa:
  tasktimeout: 2h       # 默认 2 小时，也可以写 90m、3h；必须大于 0
  stoptimeout: 30s      # 默认 30 秒，停止任务和退出游戏的宽限时间
  userdir: .runtime/maa # 进程级用户目录，不能给每个句柄切换不同目录
  asst:
    - id: example
      adbaddress: 127.0.0.1:16384
      adbpath: C:\platform-tools\adb.exe
      adbtype: BlueStacks
      resourcepath: C:\MAA
      clienttype: Official # 可省略，默认 Official
```

计时从任务出队开始，涵盖资源加载、连接模拟器、添加任务和实际执行，不包含等待排队时间。超时会请求 MaaCore 停止，同时用 `adb -s <该账号地址> shell am force-stop <游戏包名>` 退出游戏。停止确认后再次确保游戏关闭，再放行下一账号。不会退出模拟器程序。

如果停止超过 `stoptimeout` 仍无法确认核心线程结束，程序退出，由已安装的 Windows 服务管理器重启。即使 DLL 调用卡住，ADB 退出游戏仍独立执行。进程内的 `recover` 无法恢复原生崩溃，因此这个方案仍不是进程隔离。

队列保存在内存，服务重启后不会自动恢复或重放任务。账号的 ADB 地址必须互不相同；所有账号必须使用同一个 MaaCore 路径和客户端类型，避免混用进程级资源。

## 热更新加载

首次启动及每批任务出队后、执行前，会依次加载：

1. `<resourcepath>/resource`
2. `<resourcepath>/cache/resource`（存在时）

API 兼容路径 `cache/resource/tasks.json` 会同步到 `cache/resource/tasks/tasks.json`。海外服再加载对应的 `global/<clienttype>/resource` 资源。传给 `AsstLoadResource` 的都是这些 `resource` 目录的上级目录。

程序读取本机已有热更新文件，不主动下载 OTA。更新器应先完整写好资源，再提交任务。仅更新 `version.json` 或地图数据并不保证存在活动关卡导航定义。更换 MaaCore DLL 或资源根目录后需要重启服务。

## 验证

```powershell
go test ./...
go vet ./...
```

队列测试使用模拟核心，覆盖完整生命周期串行、重复提交、取消、超时退出游戏、停止失败不放行、排队不消耗执行时限，以及并发日志订阅。

真实 DLL 的无任务启动验证（加载资源、创建/销毁句柄、添加任务、修改参数、读取任务列表；不连接设备）：

```powershell
$env:MAA_NATIVE_SMOKE = '1'
go test ./internal/logic/maa -run '^TestNativeABI$' -count=1 -v
Remove-Item Env:MAA_NATIVE_SMOKE
```

旧的模拟器操作测试默认跳过，只有显式设置 `MAA_INTEGRATION_TEST=1` 才会执行。

## 截图

预览截图使用当前账号配置的 ADB 地址执行 `exec-out screencap -p`，不依赖 MaaCore 句柄已经连接设备，所以服务重启后、首次任务前也可以截图。TCP 设备未连接时会尝试连接后重试；请求最长 15 秒，并验证返回的 PNG。截图不进入任务队列、不启动游戏，也不使用其他账号的句柄或设备。
