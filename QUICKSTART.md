# 快速开始指南

这是一个5分钟快速启动指南,帮助你快速运行GitStore。

## 前置条件

确保你已安装:
- Docker & Docker Compose
- Git
- 文本编辑器

## 步骤1: 克隆代码

```bash
git clone https://github.com/nodeloc/git-store.git
cd git-store
```

## 步骤2: 配置环境变量

```bash
cp .env.example .env
```

**最小配置** (用于本地测试):

```env
# 数据库配置
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=git_store

# JWT密钥 (随机生成)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# GitHub OAuth (暂时使用测试值)
GITHUB_CLIENT_ID=test
GITHUB_CLIENT_SECRET=test
GITHUB_REDIRECT_URL=http://localhost:8080/api/auth/github/callback

# GitHub App (暂时使用测试值)
GITHUB_APP_ID=test
GITHUB_APP_PRIVATE_KEY_PATH=./github-app-private-key.pem
GITHUB_APP_INSTALLATION_ID=test
GITHUB_ORG_NAME=test

# 其他保持默认
```

## 步骤3: 创建GitHub App私钥占位文件

```bash
touch github-app-private-key.pem
```

## 步骤4: 启动服务

```bash
docker-compose up -d
```

等待30秒让数据库初始化完成。

## 步骤5: 验证运行

```bash
# 检查健康状态
curl http://localhost:8080/api/health

# 预期输出
{"status":"ok"}
```

## 步骤6: 访问前端

在浏览器打开:
```
http://localhost:8080/
```

你应该能看到GitStore主页!

## 下一步

### 配置真实的GitHub OAuth

1. 访问 https://github.com/settings/developers
2. 创建新的 OAuth App
3. 更新 `.env` 中的 `GITHUB_CLIENT_ID` 和 `GITHUB_CLIENT_SECRET`
4. 重启服务: `docker-compose restart app`

### 配置GitHub App

1. 访问 https://github.com/settings/apps/new
2. 创建GitHub App并下载私钥
3. 将私钥保存为 `github-app-private-key.pem`
4. 更新 `.env` 中的 GitHub App 配置
5. 重启服务: `docker-compose restart app`

### 配置支付方式

根据需要配置 Stripe / PayPal / Alipay,详见 [README.md](README.md)

## 常用命令

```bash
# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down

# 重启服务
docker-compose restart app

# 进入数据库
docker-compose exec postgres psql -U postgres -d git_store

# 查看表
\dt

# 查看用户
SELECT * FROM users;
```

## 故障排查

### 问题: 端口已被占用

```bash
# 修改 docker-compose.yml 中的端口映射
ports:
  - "8081:8080"  # 改为8081
```

### 问题: 数据库连接失败

```bash
# 检查PostgreSQL状态
docker-compose ps postgres

# 重启数据库
docker-compose restart postgres
```

### 问题: 无法访问前端

```bash
# 确认应用正在运行
docker-compose ps

# 检查日志
docker-compose logs app
```

## 测试API

### 测试健康检查
```bash
curl http://localhost:8080/api/health
```

### 测试插件列表
```bash
curl http://localhost:8080/api/plugins
```

### 测试GitHub OAuth流程
```bash
curl http://localhost:8080/api/auth/github
```

## 下一步阅读

- [完整文档](README.md)
- [部署指南](DEPLOYMENT.md)
- [项目结构](STRUCTURE.md)
- [项目总结](PROJECT_SUMMARY.md)

## 获取帮助

遇到问题?
1. 查看 [README.md](README.md) 的常见问题部分
2. 查看 [DEPLOYMENT.md](DEPLOYMENT.md) 的故障排查章节
3. 提交 GitHub Issue

---

**祝你使用愉快!** 🎉
