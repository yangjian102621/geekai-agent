# Repository Guidelines

## 项目结构与模块划分
- `api/`：Go 语言后端，按照 `handler/`（HTTP 层）、`service/`（业务逻辑）、`store/`（MySQL、Redis 适配器）分层，`config.toml` 存放默认配置，`static/` 与 `res/` 提供内置素材。
- `web/`：基于 Vue 3 + Vite，`src/views` 覆盖用户端与后台页面，`src/stores` 由 Pinia 维护全局状态。
- `database/` 包含初始化 SQL（如 `geekai_agent.sql`）；`docker/` 下的 `conf/config.toml` 与 `build/` 中的 Dockerfile 用于交付；`build/build.sh` 统一打包镜像。

## 构建、测试与开发命令
- 后端：`cd api && make amd64` 或 `make arm64` 生成无 CGO 的 Linux 二进制；`make clean` 清理由 `bin/` 产物。
- 前端：首次运行需 `cd web && npm install`，本地调试用 `npm run dev -- --host`，发布前执行 `npm run build` 并可用 `npm run preview` 自查。
- 一体化构建：`cd build && ./build.sh <version> [amd64|arm64] [push]` 将同时编译 API 与 Web 并生成 `registry.cn-shenzhen.aliyuncs.com/geekmaster/*` 镜像。

## 代码风格与命名
- Go 代码统一 `gofmt`，保持包级函数名使用 CamelCase，接口/结构体以业务语义命名（如 `ChatService`）；配置键遵循 snake_case，与 `config.toml` 对齐。
- Vue 组件中模板缩进 2 空格，组合式 API + `<script setup>` 为默认；Pinia store 命名 `useXxxStore`，路由文件使用下划线文件夹区分后台与用户端。
- 提交前建议运行 `golangci-lint run`（若已安装）及 `npm run lint`（根据个人 ESLint 配置）。

## 测试指南
- 当前仓库尚未附带自动化测试，新增功能请补充：后端使用 `go test ./...` 并将测试文件命名为 `xxx_test.go`，可借助 testify 断言；前端推荐引入 Vitest，对核心 store 与组件建立最小快照或交互测试。
- 回归检查需覆盖：积分扣费流程、Coze/Dify 导入流程以及支付日志列表，必要时通过演示环境模拟完整对话链路。

## 提交与 Pull Request
- 历史提交采用中英文混排、动词开头的方式（例如 “支持给智能体设置灵活的扣费模式” 或 “refactor creator and workflow modules”）；请控制在 72 字符以内并描述受影响模块。
- PR 必须说明改动背景、主要实现点、测试/回归结果，并链接相关 Issue；若变更 UI，请附录前后截图或短视频，涉及配置的要更新 `docker/conf/config.toml` 示例。
- 在合并前确保分支与 `main` 同步，解决冲突后再请求评审。

## 配置与安全提示
- 默认配置位于 `api/config.toml` 及 `docker/conf/config.toml`，提交前勿包含生产密钥；敏感凭据请通过环境变量或外部密钥管理注入。
- 运行队列、存储等依赖（Redis、MySQL）建议使用 docker-compose 或云托管服务，确认网络访问策略仅开放必要端口。
