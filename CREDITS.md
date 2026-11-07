# Credits · 开源致谢

MAA Web Console（`maa-web-console`）建立在以下开源项目之上。感谢各项目的作者、维护者与贡献者。

本清单依据根目录与 `Golang/` 的 `go.mod`、`web/package.json`、`web/package-lock.json` 及实际导入情况整理，涵盖主要集成、直接依赖、构建工具和外部运行组件。传递依赖以对应依赖清单与锁文件为准。

## MAA 集成

| 项目 | 用途 | 许可证 / 上游声明 |
| --- | --- | --- |
| [MaaAssistantArknights](https://github.com/MaaAssistantArknights/MaaAssistantArknights) | MaaCore 自动化引擎、任务协议与 C 接口；核心二进制及资源由用户另行安装 | [AGPL-3.0-only](https://github.com/MaaAssistantArknights/MaaAssistantArknights/blob/dev-v2/LICENSE)，另见[上游声明](https://github.com/MaaAssistantArknights/MaaAssistantArknights#声明) |

MAA 自身集成的图像识别、OCR、推理运行时及其他组件，详见 [MAA 上游致谢](https://github.com/MaaAssistantArknights/MaaAssistantArknights#致谢)；本仓库不附带 MaaCore 或其第三方二进制。

## Go 后端

| 项目 | 用途 | 许可证 |
| --- | --- | --- |
| [GoFrame](https://github.com/gogf/gf) | HTTP 服务、配置、ORM、MySQL 驱动适配、定时任务调度与工具库 | [MIT](https://github.com/gogf/gf/blob/v2.9.0/LICENSE) |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT 签发与校验 | [MIT](https://github.com/golang-jwt/jwt/blob/v5.2.1/LICENSE) |
| [Logrus](https://github.com/sirupsen/logrus) | 日志记录与向 WebSocket 日志流转发 | [MIT](https://github.com/sirupsen/logrus/blob/v1.9.3/LICENSE) |
| [Gorilla WebSocket](https://github.com/gorilla/websocket) | WebSocket 日志传输 | [BSD-2-Clause](https://github.com/gorilla/websocket/blob/v1.5.3/LICENSE) |
| [Go MySQL Driver](https://github.com/go-sql-driver/mysql) | GoFrame MySQL 适配器使用的底层数据库驱动 | [MPL-2.0](https://github.com/go-sql-driver/mysql/blob/v1.9.1/LICENSE) |
| [Go supplementary system libraries](https://go.googlesource.com/sys/) | Windows 系统调用与原生 DLL 交互 | [BSD-3-Clause](https://go.googlesource.com/sys/+/refs/tags/v0.31.0/LICENSE) |
| [Gin](https://github.com/gin-gonic/gin) | `Golang/` 历史原型中的 HTTP 服务；仍列于依赖清单 | [MIT](https://github.com/gin-gonic/gin/blob/v1.10.0/LICENSE) |

## Vue 前端

| 项目 | 用途 | 许可证 |
| --- | --- | --- |
| [Vue](https://github.com/vuejs/core) | 界面组件与响应式状态 | [MIT](https://github.com/vuejs/core/blob/main/LICENSE) |
| [Vue Router](https://github.com/vuejs/router) | 登录页与控制台路由 | [MIT](https://github.com/vuejs/router/blob/main/LICENSE) |
| [Vant](https://github.com/vant-ui/vant) | 移动端 UI、弹窗、按钮与时间选择器 | [MIT](https://github.com/vant-ui/vant/blob/main/LICENSE) |
| [core-js](https://github.com/zloirock/core-js) | JavaScript 兼容性补充 | [MIT](https://github.com/zloirock/core-js/blob/master/LICENSE) |
| [Vuetify](https://github.com/vuetifyjs/vuetify) | 已声明的历史 UI 依赖；当前 `main.js` 中未启用 | [MIT](https://github.com/vuetifyjs/vuetify/blob/master/LICENSE.md) |
| [EventSource polyfill](https://github.com/Yaffle/EventSource) | 已声明的历史 SSE 依赖；当前日志使用 WebSocket | [MIT](https://github.com/Yaffle/EventSource/blob/master/LICENSE) |

## 构建与开发工具

| 项目 | 用途 | 许可证 |
| --- | --- | --- |
| [Go](https://go.dev/) | 后端语言、标准库与工具链 | [BSD-3-Clause](https://go.dev/LICENSE) |
| [Node.js](https://nodejs.org/) | 前端构建运行时 | [MIT，附带组件保留各自许可](https://github.com/nodejs/node/blob/main/LICENSE) |
| [Vue CLI](https://github.com/vuejs/vue-cli) | Vue CLI Service、Babel 与 ESLint 插件 | [MIT](https://github.com/vuejs/vue-cli/blob/dev/LICENSE) |
| [Babel](https://github.com/babel/babel) | JavaScript 编译与 ESLint 解析器 | [MIT](https://github.com/babel/babel/blob/main/LICENSE) |
| [ESLint](https://github.com/eslint/eslint) | JavaScript 静态检查 | [MIT](https://github.com/eslint/eslint/blob/main/LICENSE) |
| [eslint-plugin-vue](https://github.com/vuejs/eslint-plugin-vue) | Vue 单文件组件检查 | [MIT](https://github.com/vuejs/eslint-plugin-vue/blob/master/LICENSE) |
| [unplugin-vue-components](https://github.com/unplugin/unplugin-vue-components) | 组件自动导入与 Vant 组件解析 | [MIT](https://github.com/unplugin/unplugin-vue-components/blob/main/LICENSE) |

## 外部运行与部署组件

| 项目 | 用途 | 许可证 / 上游说明 |
| --- | --- | --- |
| [Android Debug Bridge / AOSP](https://android.googlesource.com/platform/packages/modules/adb/) | 连接模拟器、设备控制与截图；由用户另行安装 | [Apache-2.0](https://android.googlesource.com/platform/packages/modules/adb/+/refs/heads/main/NOTICE) |
| [MySQL Community Server](https://dev.mysql.com/) | 定时任务持久化；由用户另行部署 | [GPL-2.0，具体以所用发行版条款为准](https://github.com/mysql/mysql-server/blob/8.0/LICENSE) |
| [NGINX](https://nginx.org/) | 可选的前端静态服务与反向代理 | [BSD-2-Clause](https://github.com/nginx/nginx/blob/master/LICENSE) |

这里的 ADB 许可指 AOSP 源码；通过 Android SDK 获取的预编译 Platform-Tools 还应遵守其发行条款。MySQL、NGINX、MAA 及其资源均为独立安装的组件。

## 版权与许可边界

本项目的原创代码采用 [AGPL-3.0-only](LICENSE)。第三方代码、依赖及外部组件保留各自的许可证与版权声明；本项目的许可不对其重新授权。分发包含第三方代码的构建产物时，应随产物保留相应许可证与声明。

《明日方舟》的名称、商标、美术与游戏素材归各自权利人所有。MAA 标志与其他第三方标志也不由本项目重新授权。本项目是非官方的个人控制台，与 MAA 官方、鹰角网络、Hypergryph 无隶属关系。
