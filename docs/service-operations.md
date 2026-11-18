# Windows 服务运行与更新

后端通过原有 Servy 包装的 `maa-backend` Windows 服务运行，账号为 `LocalSystem`，位于 Session 0，启动类型为自动。后端不依赖用户登录。登录启动计划任务和用户会话监督程序已移除。

两个 MuMu 实例使用 `adbtype: MuMuEmulator12`，`mumu.enabled: false`，任务截图使用普通 ADB 方式。网页截图预览也使用 ADB PNG。MuMu 截图增强在本机 Session 0 中不可用，不应在此服务配置中启用。

## 日常命令

从仓库根目录运行，需要管理服务时会弹出 UAC。维护检查依赖 Node.js 和 `web/node_modules` 中的 `js-yaml`、`ws`；后台服务日常运行不依赖 Node.js。

```powershell
# 查看登录、队列和定时任务状态
.\hack\manage-backend.ps1 -Action Status

# 重新启动 Windows 服务并验证
.\hack\manage-backend.ps1 -Action Restart

# 编译、更新程序并重新启动 Windows 服务
.\hack\manage-backend.ps1 -Action Update -Build

# 安装已编译好的程序
.\hack\manage-backend.ps1 -Action Update -Candidate 'C:\Downloads\maa-backend.exe'
```

管理脚本先检查所有账号队列均为空闲，再备份程序、停止服务、替换并校验程序、启动服务，验证原登录会话和定时任务；失败时尝试恢复旧程序。报告和备份位于 `.runtime/maintenance/<时间>/`。备份快照包含本地登录会话，请勿提交或对外分享。

上述更新只替换后端程序。配置、MAA DLL/资源或前端需要按各自方式更新。不要直接覆盖编译运行中的 `main.exe`。变更登录令牌或 JWT 密钥会使原会话检查失败，应单独安排迁移；数据库变更也需要对应的迁移和回退方案。

## 自动恢复和日志

继续使用原 Servy 服务包装器的进程恢复配置，未安装登录触发或用户会话恢复任务。队列检测到 MaaCore 无法在停止期限内退出时，已有的 `os.Exit(1)` 会退出后端，由服务包装器处理重启。

服务输出位于 `.logs/stdout.log`、`.logs/stderr.log`，MaaCore 日志位于 `.runtime/maa/debug/asst.log`。重启会重新加载定时任务；内存中的执行中和排队任务不会重放，也不会补跑停机期间错过的定时任务。
