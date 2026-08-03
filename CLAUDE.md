# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GeekAI-Agent is a full-stack AI agent platform that allows monetization of AI agents from platforms like Coze and Dify. Built with Go backend and Vue 3 frontend, it includes user management, payment systems, and admin features.

## Development Commands

### Backend (Go)

```bash
# Development with hot reload
fresh                              # Start with auto-reload (uses fresh.conf)

# Manual build
go build -o bin/geek-agent main.go

# Production builds
make all                           # Build for both amd64 and arm64
make amd64                         # Build for Linux AMD64
make arm64                         # Build for Linux ARM64
make clean                         # Clean build artifacts
```

### Frontend (Vue 3)

```bash
# Development server (port 8888)
npm run dev

# Production build
npm run build

# Preview production build
npm run preview
```

### Docker Deployment

```bash
# Full stack deployment
docker-compose up -d
```

## Architecture

### Backend Structure (`/api`)

- **Entry Point**: `main.go` - Uses Uber FX for dependency injection
- **Core**: `core/app_server.go` - HTTP server initialization
- **Handlers**: `handler/` - HTTP request handlers organized by feature
- **Services**: `service/` - Business logic layer
- **Store**: `store/model/` (database models) and `store/vo/` (view objects)
- **Modules**: `modules/` - Feature modules (e.g., `creator/` for creator platform)
- **Config**: `config.toml` - Main configuration file

### Frontend Structure (`/web/src`)

- **Entry Point**: `main.js` - Vue 3 application bootstrap
- **Router**: `router.js` - Route definitions
- **Views**: `views/` - Page components (includes `admin/` for admin panel)
- **Components**: `components/` - Reusable Vue components
- **API Layer**: `js/action/` - API interaction logic
- **Styling**: Uses TailwindCSS + Element Plus

### Technology Stack

- **Backend**: Go 1.23.0, Gin framework, GORM (MySQL), Redis, WebSocket
- **Frontend**: Vue 3 Composition API, Element Plus UI, Vite, TailwindCSS
- **Database**: MySQL 8.0 (primary), Redis 6.0 (caching)
- **Storage**: Multiple providers (Aliyun OSS, Qiniu, MinIO, Local)

## Key Development Notes

### Module Architecture

The backend follows a modular pattern where each feature (like `creator`) has its own:

- `handler/` - HTTP endpoints
- `model/` - Database models
- `service/` - Business logic

### Configuration

- Backend configuration in `api/config.toml`
- Frontend uses Vite with proxy to backend during development
- Hot reload configured for both frontend (Vite HMR) and backend (Fresh)

### Database

- Schema located in `database/geekai_agent.sql`
- Uses GORM for ORM with MySQL
- Redis for caching and sessions

### API Integration

- Supports multiple LLM APIs (OpenAI, Anthropic, etc.)
- Coze agent import functionality
- Planned Dify integration

### Business Features

- User management with points/credits system
- Payment integration (Alipay, WeChat Pay)
- Creator platform for monetization
- Admin dashboard for system management

### 开发总原则

1. 所有的功能开发都要遵循 `MVC` 架构模式，即 `Model`、`View`、`Controller` 分离。
2. 所有的功能开发都要遵循 `DRY` 原则，即 `Don't Repeat Yourself`，不要重复造轮子。
3. 所有的功能开发都要遵循 `KISS` 原则，即 `Keep It Simple Stupid`，保持简单。
4. 所有的功能开发都要遵循 `YAGNI` 原则，即 `You Aren't Gonna Need It`，不要过度设计。
5. 所有 Vue 文件的 JS 代码都要写在 `<script setup>` 标签中，不要写在 `<script>` 标签中。
6. 所有的 javascript 代码都要写在 `web/src/js/action` 目录下，不要写在视图文件中。
7. 无论是 golang 代码还是 javascript 代码，工具函数尽量调用现在已有的工具函数，不要自己造轮子。
8. Model 属性的命名统一采取驼峰命名法，不要使用下划线命名法，也不要全部大写，比如 Id 不要写成 ID,Pid 不要写成 PID。

### 新增一个管理后台功能的开发流程如下：

1. 在 `api/handler/admin` 目录下创建新的控制器文件，控制器文件的命名为 `{feature_name}_handler.go`,可以参考 `api/handler/admin/user_handler.go` 文件
2. 如果有封装组件，组件文件写在 `web/src/components/admin` 目录下，文件的命名为 `{feature_name}.vue` 首字母大写。
3. 在 `web/src/js/action/admin` 目录下创建新的业务逻辑文件，文件的命名为 `{feature_name}.js` 所有的 Javascript 功能代码都写在这个文件中，不要在视图文件中写任何业务逻辑代码。可以参考 `web/src/js/action/admin/user.js` 文件
4. 在 `routers.js` 文件中添加新的路由配置
5. 如果有引入新的组件，在 `main.js` 文件中添加新的组件配置
6. 在 `main.go` 文件中添加新的控制器实例
7. 在 `web/src/components/admin/SideBar.vue` 文件中添加菜单配置

### 新增一个前台功能的开发流程如下：

1. 在 `api/handler` 目录下创建新的控制器文件，控制器文件的命名为 `{feature_name}_handler.go`,可以参考 `api/handler/admin/user_handler.go` 文件
2. css 样式文件写在 `web/src/assets/css/` 目录下，文件的命名为 `{feature_name}.css` 名称全部小写。
3. 如果有封装组件，组件文件写在 `web/src/components` 目录下，文件的命名为 `{feature_name}.vue` 首字母大写。
4. 在 `web/src/js/action` 目录下创建新的业务逻辑文件，文件的命名为 `{feature_name}.js`，文件名全部小写， 所有的 Javascript 功能代码都写在这个文件中，不要在视图文件中写任何业务逻辑代码。可以参考 `web/src/js/action/admin/user.js` 文件
5. 在 `web/src/router.js` 文件中添加新的路由配置
6. 如果有引入新的组件，在 `main.js` 文件中添加新的组件配置
7. 在 `main.go` 文件中添加新的控制器实例
