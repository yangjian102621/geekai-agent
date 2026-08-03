# GeekAI-Agent

> 企业级 AI 智能体商业化与治理平台：统一接入智能体，提供用户、计费、支付、分账和运营能力。

GeekAI-Agent 用于将 OpenAI 兼容接口、Coze、Dify、阿里云百炼等平台上的智能体接入统一门户。项目由 Go API 和 Vue 3 Web 应用组成，支持私有化部署。

## 核心能力

- 多平台智能体统一接入和流式对话
- 工作流异步执行、进度轮询、失败退还积分
- 用户、积分、产品、订单和兑换码管理
- 创作者入驻、应用管理、收益分账和提现审核
- 支付宝、微信支付和易支付
- 本地、阿里云 OSS、七牛云和 MinIO 文件存储
- JWT、Redis 会话、接口限流和管理后台

## 技术栈

| 模块 | 技术 |
| --- | --- |
| API | Go 1.23、Gin、GORM、Uber Fx |
| Web | Vue 3、Vite、Pinia、Element Plus、Tailwind CSS |
| 数据 | MySQL 8.0、Redis 6.0 |
| 部署 | Docker、Docker Compose、Nginx |

## 项目结构

```text
geekai-agent/
├── api/                    # Go API：handler、service、store
├── web/                    # Vue 3 前端
├── database/               # 数据库初始化和升级 SQL
├── docker/                 # Compose、Nginx 和运行时配置
├── build/                  # API/Web 镜像构建脚本和 Dockerfile
├── CHANGELOG.md
└── README.md
```

## 本地开发

### 1. 环境要求

- Go 1.23.7 或兼容的 Go 1.23 工具链
- Node.js 20 LTS，npm 10 或兼容版本
- MySQL 8.0
- Redis 6.0 及以上
- Git

可选依赖：需要解析 Office、PDF 等文件内容时，再部署 Apache Tika。只开发基础对话和管理功能时可以暂不启动。

### 2. 获取源码

```bash
git clone https://github.com/yangjian102621/geekai-agent.git
cd geekai-agent
```

### 3. 初始化 MySQL

创建数据库并导入初始化脚本：

```bash
mysql -uroot -p -e "CREATE DATABASE IF NOT EXISTS geekai_agent CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -uroot -p geekai_agent < database/geekai_agent.sql
```

如果是从旧版本升级，请先备份数据库，再根据 `database/update.sql` 和 `CHANGELOG.md` 检查升级项。不要在生产库上直接重复导入全量初始化脚本。

### 4. 配置并启动 API

复制一份本地配置，避免把个人密钥提交到仓库：

```bash
cd api
cp config.toml config.local
```

至少修改以下配置：

```toml
Listen = "0.0.0.0:5678"
MysqlDns = "root:你的数据库密码@tcp(127.0.0.1:3306)/geekai_agent?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"

[Session]
SecretKey = "使用密码生成器创建的至少 64 位随机字符串"
MaxAge = 8640000

[AdminSession]
SecretKey = "使用另一个至少 64 位随机字符串"
MaxAge = 8640000

[Redis]
Host = "127.0.0.1"
Port = 6379
Password = ""
DB = 0
```

安装依赖并启动：

```bash
go mod download
GEEKAI_ADMIN_USERNAME=admin \
GEEKAI_ADMIN_PASSWORD='请设置至少12位强密码' \
CONFIG_FILE=config.local go run .
```

API 默认监听 `http://localhost:5678`。开发时也可以安装 [Fresh](https://github.com/gravityblast/fresh)，然后在 `api/` 目录运行 `fresh` 获得热更新。

### 5. 配置并启动 Web

另开一个终端：

```bash
cd web
npm ci
npm run dev
```

访问 `http://localhost:8888`。开发环境已将 `/api` 和 `/static` 代理到 `http://localhost:5678`；如需修改端口，请同步调整 `web/vite.config.js`。

### 6. 初始化后台

访问 `http://localhost:8888/#/admin`，使用启动 API 时通过 `GEEKAI_ADMIN_USERNAME` 和 `GEEKAI_ADMIN_PASSWORD` 创建的管理员登录。系统只会在管理员表为空时创建该账号，后续重启不会覆盖已有管理员。

## 配置说明

应用有两类配置：

1. `api/config.toml` 或 `CONFIG_FILE` 指定的 TOML 文件：监听地址、MySQL、Redis、会话密钥和基础服务地址。
2. 管理后台保存到 MySQL 的系统配置：站点信息、Coze、短信、对象存储、SMTP、支付和微信登录等。

敏感配置不要提交到 Git。建议本地使用已被 `.gitignore` 的 `api/config.local`，生产环境通过只读挂载的配置文件或密钥管理服务注入。会话密钥、数据库密码、Redis 密码和第三方密钥应分别生成，不要复用。

## 验证改动

### API

```bash
cd api
gofmt -w 你修改的.go文件
go test ./...
go build ./...
```

### Web

```bash
cd web
npm ci
npm run build
```

当前前端没有配置独立的 lint/test 脚本，以 `npm run build` 作为最低提交门槛。新增核心业务逻辑时，建议同步引入 Vitest 测试。

## Docker Compose 部署

仓库中的 Compose 包含 MySQL、Redis、API 和 Web，适合单机体验或作为生产部署模板。生产环境应使用外部数据库、密钥管理、HTTPS 和定期备份。

### 1. 准备服务器

建议最低配置为 2 核 CPU、4 GB 内存和 20 GB 可用磁盘，并安装 Docker Engine 与 Docker Compose v2：

```bash
docker --version
docker compose version
```

开放 Web 对外端口 `8081`。MySQL 的 `3308` 和 Redis 的 `6381` 默认也映射到宿主机；生产环境应删除这两个端口映射，或通过防火墙禁止公网访问。

### 2. 准备配置

```bash
cd docker
mkdir -p conf/mysql logs/mysql logs/nginx data/mysql/data data/redis static
cp .env.example .env
```

编辑 `.env`，分别设置 MySQL、Redis 密码、两个会话密钥和初始管理员账号。管理员密码不得少于 12 位。Compose 会通过环境变量把这些值注入 API，不需要把密钥写入受版本控制的 TOML 文件。

可使用 OpenSSL 生成会话密钥：

```bash
openssl rand -hex 32
openssl rand -hex 32
```

### 3. 选择镜像架构

`docker/docker-compose.yaml` 默认使用 `amd64` 镜像。如果服务器是 ARM64，请将 API 和 Web 镜像标签中的 `-amd64` 改为 `-arm64`。

```bash
uname -m
```

输出 `x86_64` 通常使用 `amd64`，输出 `aarch64` 或 `arm64` 使用 `arm64`。

### 4. 启动服务

```bash
docker compose pull
docker compose up -d
docker compose ps
docker compose logs -f geek-agent-api
```

首次启动时，MySQL 会执行 `docker/data/mysql/init.d/` 下的 SQL。该机制只在 MySQL 数据目录为空时执行；修改初始化 SQL 后，已有数据卷不会自动重新导入。

服务启动后访问：

- 用户端：`http://服务器IP:8081`
- 管理后台：`http://服务器IP:8081/#/admin`
- API：宿主机 `6789` 端口

登录后台后立即修改默认管理员密码，再配置站点信息和所需的第三方服务。

### 5. 查看日志和升级

```bash
cd docker
docker compose logs --tail=200 geek-agent-api
docker compose logs --tail=200 geek-agent-mysql
docker compose pull
docker compose up -d
```

升级前务必备份：

```bash
docker exec geek-agent-mysql mysqldump -uroot -p geekai_agent > geekai_agent-backup.sql
tar -czf geekai-agent-data-backup.tar.gz conf data static
```

`mysqldump` 会交互式询问密码，避免把密码写进 Shell 历史。升级完成后检查 API 日志、登录、对话、积分扣费、工作流和支付回调。

### 6. 停止服务

```bash
docker compose down
```

不要随意添加 `-v`，它会删除 Compose 管理的数据卷。当前项目主要使用绑定目录保存数据，清理 `docker/data/` 前仍须先备份。

## 从源码构建 Docker 镜像

构建脚本会先编译 Go API 和 Web，再生成指定架构的两个镜像：

```bash
cd web
npm ci
cd ../build
./build.sh v1.0.5 amd64
```

ARM64 构建：

```bash
cd build
./build.sh v1.0.5 arm64
```

传入第三个参数 `push` 会推送到脚本中写死的镜像仓库，仅仓库维护者应使用。普通贡献者应先把 `build/build.sh` 中的镜像命名空间改成自己的仓库，或直接执行 `docker build`。

## 生产环境检查清单

- 使用正式域名和 HTTPS，支付回调必须可从公网安全访问
- 替换全部默认密码、会话密钥和演示账号
- 不向公网暴露 MySQL、Redis 和 API 管理端口
- 使用独立低权限数据库账号，不让 API 使用 MySQL root
- 配置 MySQL、上传文件和对象存储的定期备份
- 限制上传文件类型和大小，校验反向代理请求体上限
- 第三方 API 密钥仅通过配置文件挂载或密钥管理服务注入
- 上线前回归登录、积分扣费、Coze/Dify/百炼工作流、支付和提现流程

## 常见问题

### Web 页面打开，但接口返回 502

确认 `geek-agent-api` 正常运行，并检查 Nginx 配置中的上游端口是否与 API 的 `Listen` 一致：

```bash
docker compose ps
docker compose logs --tail=200 geek-agent-api
```

### API 无法连接 MySQL 或 Redis

在 Docker 内应使用服务名 `geek-agent-mysql`、`geek-agent-redis`，不能使用 `localhost`。同时确认 `docker/.env` 中的密码已经设置，且容器已重新创建。

### 修改初始化 SQL 后没有生效

MySQL 官方镜像只会在空数据目录上执行初始化脚本。请通过迁移 SQL升级已有数据库；不要为了重跑脚本直接删除生产数据目录。

### 上传后的文件无法访问

使用本地存储时，确认 API 的 `./static/upload` 位于已持久化的 `docker/static` 挂载中，并检查 `/static/` 的 Nginx 反向代理。

## 贡献

提交改动前，请确保后端 `go test ./...`、`go build ./...` 和前端 `npm run build` 通过。Pull Request 需要说明改动背景、实现方式和验证结果；UI 变更请附截图，配置变更请同步更新示例配置。

## 文档与演示

- 在线演示：[https://agent.geekai.me](https://agent.geekai.me)
- 管理后台：[https://agent.geekai.me/#/admin](https://agent.geekai.me/#/admin)
- 使用文档：[https://docs.geekai.me/agent/](https://docs.geekai.me/agent/)

公开演示账号仅用于体验，部署自己的实例后请使用独立账号和强密码。

## 开源许可证

本项目采用 [Apache License 2.0](LICENSE)。使用和分发源码及衍生作品时，请遵守许可证要求。项目名称、Logo 及第三方素材的商标或其他权利不因代码许可证自动授予。

参与贡献前请阅读 [贡献指南](CONTRIBUTING.md)。安全漏洞请按照 [安全策略](SECURITY.md) 私密报告。
