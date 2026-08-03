# 贡献指南

感谢你参与 GeekAI-Agent。提交代码前请先搜索现有 Issue，较大的功能或兼容性变更应先创建 Issue 说明用户问题、方案和影响范围。

## 开发流程

1. Fork 仓库，从最新的 `main` 创建功能分支。
2. 只修改与目标直接相关的代码，不提交密钥、生产数据、构建产物或编辑器配置。
3. 后端遵循 Handler、Service、Store 分层，Go 代码运行 `gofmt`。
4. 前端使用 Vue 3 `<script setup>`，业务请求和状态逻辑放在 `web/src/js` 对应模块。
5. 新功能和缺陷修复应提供能证明行为的测试；暂时无法自动化时，在 PR 中写明人工回归步骤。

## 本地验证

```bash
cd api
go test ./...
go build ./...

cd ../web
npm ci
npm run build
```

涉及 Docker 配置时还要执行：

```bash
cd docker
cp .env.example .env
docker compose config --quiet
```

不要将生成的 `docker/.env` 提交到仓库。

## 提交和 Pull Request

提交信息建议遵循 Conventional Commits，例如 `feat(workflow): support new provider`、`fix(payment): verify callback signature`。

Pull Request 必须包含：

- 要解决的用户问题和改动范围
- 关键实现与取舍
- 测试命令及结果
- 配置、数据库迁移和兼容性影响
- UI 变更的截图或录屏

数据库结构变更必须同时更新模型、初始化 SQL、升级迁移和部署文档。禁止通过删除生产字段或重建数据目录来代替迁移。

## 许可证

提交贡献即表示你有权提交相关代码，并同意贡献内容按照仓库的 Apache License 2.0 进行许可。第三方代码和素材必须保留许可证及归属信息。
