# 部署指南

本文档提供GitStore系统的完整部署指南。

## 目录

- [系统要求](#系统要求)
- [前置准备](#前置准备)
- [本地开发](#本地开发)
- [生产部署](#生产部署)
- [监控与维护](#监控与维护)
- [故障排查](#故障排查)

## 系统要求

### 硬件要求
- CPU: 2核心以上
- 内存: 4GB以上
- 磁盘: 20GB以上

### 软件要求
- Ubuntu 20.04+ / CentOS 8+ / macOS
- Go 1.21+
- PostgreSQL 14+
- Git

## 前置准备

### 1. 创建GitHub OAuth App

1. 访问 https://github.com/settings/developers
2. 点击 "New OAuth App"
3. 填写信息:
   - Application name: `GitStore`
   - Homepage URL: `https://your-domain.com`
   - Authorization callback URL: `https://your-domain.com/api/auth/github/callback`
4. 保存 Client ID 和 Client Secret

### 2. 创建GitHub App

1. 访问 https://github.com/settings/apps/new
2. 填写基本信息:
   - GitHub App name: `GitStore App`
   - Homepage URL: `https://your-domain.com`
   - Webhook: 留空
3. 设置权限:
   - Repository permissions:
     - Contents: Read-only
     - Metadata: Read-only
4. 创建后:
   - 记录 App ID
   - 生成并下载私钥 (.pem文件)
   - 安装到你的组织/账号
   - 记录 Installation ID

获取Installation ID:
```bash
# 使用GitHub CLI
gh api /app/installations

# 或者访问
https://github.com/settings/installations
```

### 3. 配置支付服务

#### Stripe
1. 访问 https://dashboard.stripe.com/register
2. 获取 API Keys: https://dashboard.stripe.com/apikeys
3. 创建 Webhook endpoint: https://dashboard.stripe.com/webhooks
   - Endpoint URL: `https://your-domain.com/api/webhooks/stripe`
   - Events: `payment_intent.succeeded`, `payment_intent.payment_failed`

#### PayPal
1. 访问 https://developer.paypal.com/home
2. 创建应用: https://developer.paypal.com/dashboard/applications
3. 获取 Client ID 和 Secret

#### 支付宝
1. 访问 https://open.alipay.com
2. 创建应用
3. 配置公私钥
4. 获取 App ID

### 4. 配置邮件服务

#### 使用Gmail
1. 开启2FA: https://myaccount.google.com/security
2. 生成应用专用密码: https://myaccount.google.com/apppasswords
3. 使用该密码作为SMTP密码

#### 使用SendGrid
1. 注册 https://sendgrid.com
2. 创建 API Key
3. 配置:
```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=your_sendgrid_api_key
```

## 本地开发

### 1. 克隆代码

```bash
git clone https://github.com/nodeloc/git-store.git
cd git-store
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env` 文件:

```env
# 应用配置
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000

# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=git_store

# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URL=http://localhost:8080/api/auth/github/callback

# GitHub App
GITHUB_APP_ID=your_app_id
GITHUB_APP_PRIVATE_KEY_PATH=./github-app-private-key.pem
GITHUB_APP_INSTALLATION_ID=your_installation_id
GITHUB_ORG_NAME=your-org-name

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this

# Stripe
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...

# 其他配置...
```

### 4. 创建数据库

```bash
# 创建数据库
createdb git_store

# 运行迁移
psql git_store < migrations/001_initial_schema.sql
```

### 5. 放置GitHub App私钥

```bash
# 将下载的私钥放到项目根目录
cp ~/Downloads/your-app-name.2024-01-01.private-key.pem ./github-app-private-key.pem
chmod 600 github-app-private-key.pem
```

### 6. 启动服务

```bash
# 方式1: 直接运行
go run main.go

# 方式2: 使用Make
make run

# 方式3: 使用Air(热重载)
air
```

访问 http://localhost:8080/api/health 验证服务运行。

## 生产部署

### 方案1: Docker部署(推荐)

#### 1. 准备环境

```bash
# 安装Docker
curl -fsSL https://get.docker.com | sh

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件,填入生产环境配置
```

#### 3. 启动服务

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

#### 4. 配置Nginx反向代理

```nginx
# /etc/nginx/sites-available/git-store
server {
    listen 80;
    server_name pluginstore.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用站点:
```bash
sudo ln -s /etc/nginx/sites-available/git-store /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

#### 5. 配置HTTPS (Let's Encrypt)

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d pluginstore.com
```

### 方案2: 系统服务部署

#### 1. 编译应用

```bash
# 在本地或CI/CD中编译
GOOS=linux GOARCH=amd64 go build -o git-store main.go

# 上传到服务器
scp git-store user@server:/opt/git-store/
scp -r migrations user@server:/opt/git-store/
scp .env user@server:/opt/git-store/
scp github-app-private-key.pem user@server:/opt/git-store/
```

#### 2. 创建systemd服务

```bash
sudo nano /etc/systemd/system/git-store.service
```

内容:
```ini
[Unit]
Description=GitStore
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/git-store
ExecStart=/opt/git-store/git-store
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

#### 3. 启动服务

```bash
sudo systemctl daemon-reload
sudo systemctl start git-store
sudo systemctl enable git-store
sudo systemctl status git-store
```

#### 4. 查看日志

```bash
sudo journalctl -u git-store -f
```

### 方案3: Kubernetes部署

#### 1. 创建ConfigMap

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: git-store-config
data:
  APP_ENV: production
  DB_HOST: postgres-service
  DB_PORT: "5432"
  DB_NAME: git_store
```

#### 2. 创建Secret

```bash
kubectl create secret generic git-store-secrets \
  --from-literal=DB_PASSWORD=your_password \
  --from-literal=JWT_SECRET=your_jwt_secret \
  --from-literal=GITHUB_CLIENT_SECRET=your_github_secret \
  --from-file=github-app-private-key.pem
```

#### 3. 创建Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: git-store
spec:
  replicas: 3
  selector:
    matchLabels:
      app: git-store
  template:
    metadata:
      labels:
        app: git-store
    spec:
      containers:
      - name: git-store
        image: your-registry/git-store:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: git-store-config
        - secretRef:
            name: git-store-secrets
        volumeMounts:
        - name: github-key
          mountPath: /app/github-app-private-key.pem
          subPath: github-app-private-key.pem
      volumes:
      - name: github-key
        secret:
          secretName: git-store-secrets
```

#### 4. 创建Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: git-store
spec:
  selector:
    app: git-store
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

#### 5. 部署

```bash
kubectl apply -f k8s/
```

## 监控与维护

### 1. 健康检查

```bash
# 检查服务状态
curl http://localhost:8080/api/health

# 预期响应
{"status":"ok"}
```

### 2. 日志管理

#### Docker日志
```bash
docker-compose logs -f app
docker-compose logs --tail=100 app
```

#### Systemd日志
```bash
sudo journalctl -u git-store -f
sudo journalctl -u git-store --since "1 hour ago"
```

### 3. 数据库备份

```bash
# 每日备份脚本
#!/bin/bash
BACKUP_DIR=/var/backups/git-store
DATE=$(date +%Y%m%d_%H%M%S)

pg_dump git_store > $BACKUP_DIR/backup_$DATE.sql
gzip $BACKUP_DIR/backup_$DATE.sql

# 保留最近7天的备份
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
```

添加到crontab:
```bash
# 每天凌晨3点备份
0 3 * * * /opt/scripts/backup-db.sh
```

### 4. 性能监控

推荐工具:
- **Prometheus** - 指标收集
- **Grafana** - 可视化面板
- **Alertmanager** - 告警管理

关键指标:
- API响应时间
- 数据库连接数
- 支付成功率
- 定时任务执行状态

### 5. 定时任务监控

```bash
# 查看定时任务日志
grep "maintenance expiry check" /var/log/git-store.log

# 检查授权回收是否正常
psql git_store -c "SELECT COUNT(*) FROM licenses WHERE status='expired' AND updated_at > NOW() - INTERVAL '1 day';"
```

## 故障排查

### 问题1: 服务无法启动

```bash
# 检查端口占用
sudo lsof -i :8080

# 检查数据库连接
psql -h localhost -U postgres -d git_store -c "SELECT 1;"

# 检查环境变量
env | grep -E "DB_|GITHUB_|STRIPE_"
```

### 问题2: GitHub App授权失败

```bash
# 检查私钥文件
ls -la github-app-private-key.pem
cat github-app-private-key.pem | head -1

# 验证App ID和Installation ID
curl -H "Authorization: Bearer $(generate-jwt-token)" \
  https://api.github.com/app/installations

# 测试仓库访问
# 使用GitHub API测试installation是否能访问仓库
```

### 问题3: 支付回调失败

```bash
# 检查Webhook日志
grep "webhook" /var/log/git-store.log

# 验证Webhook签名
# 在Stripe Dashboard查看webhook事件详情

# 手动重试webhook
curl -X POST http://localhost:8080/api/webhooks/stripe \
  -H "Content-Type: application/json" \
  -d @webhook-payload.json
```

### 问题4: 邮件发送失败

```bash
# 测试SMTP连接
telnet smtp.gmail.com 587

# 检查邮件日志
psql git_store -c "SELECT * FROM email_notifications WHERE status='failed' ORDER BY created_at DESC LIMIT 10;"

# 测试发送邮件
go run scripts/test-email.go
```

### 问题5: 定时任务未执行

```bash
# 检查cron是否运行
ps aux | grep cron

# 查看定时任务日志
grep "Running maintenance expiry check" /var/log/git-store.log

# 手动触发任务
curl -X POST http://localhost:8080/admin/cron/trigger
```

## 安全加固

### 1. 防火墙配置

```bash
# UFW (Ubuntu)
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 2. 限制数据库访问

```bash
# PostgreSQL配置
# /etc/postgresql/14/main/pg_hba.conf
host    git_store    postgres    127.0.0.1/32    md5

# 重启PostgreSQL
sudo systemctl restart postgresql
```

### 3. 启用速率限制

在Nginx中配置:
```nginx
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

location /api/ {
    limit_req zone=api burst=20 nodelay;
}
```

### 4. 定期更新

```bash
# 更新系统包
sudo apt update && sudo apt upgrade

# 更新Go依赖
go get -u ./...
go mod tidy

# 重新构建
make build
```

## 扩容方案

### 水平扩展

1. 使用负载均衡器(Nginx/HAProxy)
2. 部署多个应用实例
3. 使用Redis共享session
4. 数据库读写分离

### 垂直扩展

1. 升级服务器配置
2. 优化数据库查询
3. 添加缓存层(Redis)
4. 使用CDN加速静态资源

## 更新流程

### 1. 准备更新

```bash
# 备份数据库
pg_dump git_store > backup_before_update.sql

# 备份代码
cp -r /opt/git-store /opt/git-store.backup
```

### 2. 拉取更新

```bash
git pull origin main
go mod download
```

### 3. 运行迁移

```bash
# 如果有新的迁移文件
psql git_store < migrations/002_new_migration.sql
```

### 4. 重启服务

```bash
# Docker
docker-compose down
docker-compose up -d --build

# Systemd
sudo systemctl restart git-store
```

### 5. 验证

```bash
# 检查服务状态
curl http://localhost:8080/api/health

# 检查日志
docker-compose logs -f app
```

## 支持

如遇到问题:
1. 查看日志文件
2. 检查环境变量配置
3. 访问GitHub Issues
4. 联系技术支持

---

更新时间: 2024-01-21
