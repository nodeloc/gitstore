# GitStore - GitHub App授权商业分发平台

基于GitHub App + 私有仓库的插件商店系统,专为Discourse插件设计的商业分发与授权平台。

## 核心特性

### ✨ 授权模型
- **永久许可证** - 一次购买,插件永久可用
- **时间限定维护** - 默认赠送12个月更新权限
- **零侵入式授权** - 无需在插件代码中集成License校验
- **自动权限回收** - 维护到期自动移除GitHub仓库访问权限

### 🔐 GitHub App集成
- 基于GitHub App统一授权管理
- 细粒度仓库访问控制
- 支持个人账号和组织账号
- 自动token管理和刷新

### 💳 多支付方式
- **Stripe** - 国际信用卡支付
- **PayPal** - 全球在线支付
- **支付宝** - 国内主流支付方式

### 📧 自动化通知
- 购买成功邮件
- 维护到期提醒(30天/7天/1天)
- 维护已到期通知
- 续费成功确认

### 📊 数据统计
- 实时销售统计
- 用户增长分析
- 收入趋势报表
- 插件下载统计

## 技术栈

- **后端**: Go 1.21+
- **框架**: Gin
- **数据库**: PostgreSQL 14+
- **前端**: DaisyUI (Tailwind CSS)
- **支付**: Stripe, PayPal, Alipay
- **邮件**: SMTP (支持Gmail, SendGrid等)
- **定时任务**: Cron
- **认证**: GitHub OAuth + JWT

## 快速开始

### 前置要求

- Go 1.21+
- PostgreSQL 14+
- GitHub OAuth App
- GitHub App (用于仓库访问控制)
- Stripe/PayPal/Alipay 账号 (至少一个)

### 安装步骤

1. **克隆仓库**
```bash
git clone https://github.com/nodeloc/git-store.git
cd git-store
```

2. **安装依赖**
```bash
go mod download
```

3. **配置环境变量**
```bash
cp .env.example .env
# 编辑 .env 文件,填入你的配置
```

4. **创建数据库**
```bash
createdb git_store
psql git_store < migrations/001_initial_schema.sql
```

5. **运行服务**
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

## 配置说明

### 1. GitHub OAuth配置

前往 https://github.com/settings/developers 创建OAuth App:

- **Application name**: GitStore
- **Homepage URL**: `http://localhost:8080`
- **Authorization callback URL**: `http://localhost:8080/api/auth/github/callback`

获取 Client ID 和 Client Secret,填入 `.env`:

```env
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
```

### 2. GitHub App配置

前往 https://github.com/settings/apps/new 创建GitHub App:

**基本信息:**
- **GitHub App name**: GitStore App
- **Homepage URL**: `http://localhost:8080`
- **Webhook**: 可选

**权限设置:**
- Repository permissions:
  - Contents: Read-only
  - Metadata: Read-only

创建后:
1. 下载私钥文件 (`.pem`)
2. 安装到你的组织/账号
3. 获取 App ID 和 Installation ID

填入 `.env`:
```env
GITHUB_APP_ID=your_app_id
GITHUB_APP_PRIVATE_KEY_PATH=./github-app-private-key.pem
GITHUB_APP_INSTALLATION_ID=your_installation_id
GITHUB_ORG_NAME=your-org-name
```

### 3. Stripe配置

前往 https://dashboard.stripe.com/apikeys 获取密钥:

```env
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
```

### 4. PayPal配置

前往 https://developer.paypal.com/dashboard/applications 创建应用:

```env
PAYPAL_CLIENT_ID=your_client_id
PAYPAL_CLIENT_SECRET=your_client_secret
PAYPAL_MODE=sandbox  # 生产环境改为 live
```

### 5. 支付宝配置

前往 https://open.alipay.com 获取配置:

```env
ALIPAY_APP_ID=your_app_id
ALIPAY_PRIVATE_KEY=your_private_key
ALIPAY_PUBLIC_KEY=alipay_public_key
```

### 6. 邮件配置 (SMTP)

使用Gmail示例:

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@pluginstore.com
```

## 核心业务流程

### 购买流程

```
1. 用户浏览插件商店
   ↓
2. 选择插件,点击购买
   ↓
3. 选择支付方式 (Stripe/PayPal/Alipay)
   ↓
4. 完成支付
   ↓
5. 系统创建订单和授权记录
   ↓
6. 自动授予GitHub App仓库访问权限
   ↓
7. 发送购买成功邮件 (含安装教程)
   ↓
8. 用户可通过 git clone 下载插件
```

### 到期回收流程

```
定时任务 (每日凌晨2点):

1. 扫描所有授权记录
   ↓
2. 发现维护到期的授权
   ↓
3. 调用GitHub API移除仓库访问权限
   ↓
4. 更新授权状态为 "expired"
   ↓
5. 发送到期通知邮件
   ↓
6. 用户本地插件继续可用
   ↓
7. 用户无法 git pull 获取更新
```

### 续费流程

```
1. 用户在后台查看已购插件
   ↓
2. 点击"续费"按钮
   ↓
3. 选择续费时长 (6/12/24个月)
   ↓
4. 完成支付
   ↓
5. 系统延长维护到期时间
   ↓
6. 自动恢复GitHub仓库访问权限
   ↓
7. 发送续费成功邮件
```

## API文档

### 公开接口

#### 用户认证
- `GET /api/auth/github` - GitHub OAuth登录
- `GET /api/auth/github/callback` - OAuth回调
- `GET /api/auth/me` - 获取当前用户信息
- `POST /api/auth/logout` - 登出

#### 插件浏览
- `GET /api/plugins` - 获取插件列表
- `GET /api/plugins/:slug` - 获取插件详情

#### 教程
- `GET /api/tutorials/public` - 获取公开教程
- `GET /api/tutorials/:slug` - 获取教程详情

### 需认证接口

#### 用户中心
- `GET /api/user/licenses` - 我的授权
- `GET /api/user/orders` - 我的订单
- `GET /api/user/github-accounts` - GitHub账号列表

#### 订单管理
- `POST /api/orders` - 创建订单
- `GET /api/orders/:id` - 订单详情

#### 支付
- `POST /api/payments/stripe/create-intent` - 创建Stripe支付意图
- `POST /api/payments/paypal/create-order` - 创建PayPal订单
- `POST /api/payments/paypal/capture-order` - 确认PayPal支付
- `POST /api/payments/alipay/create` - 创建支付宝支付

#### 授权管理
- `GET /api/licenses/:id` - 授权详情
- `POST /api/licenses/:id/renew` - 续费授权
- `GET /api/licenses/:id/history` - 授权历史

### 管理员接口

#### 插件管理
- `GET /api/admin/plugins` - 所有插件
- `POST /api/admin/plugins` - 创建插件
- `PUT /api/admin/plugins/:id` - 更新插件
- `DELETE /api/admin/plugins/:id` - 删除插件
- `POST /api/admin/plugins/sync-repos` - 同步GitHub仓库

#### 订单管理
- `GET /api/admin/orders` - 所有订单
- `POST /api/admin/orders/:id/refund` - 退款

#### 授权管理
- `GET /api/admin/licenses` - 所有授权
- `POST /api/admin/licenses/:id/revoke` - 撤销授权
- `POST /api/admin/licenses/:id/extend` - 延长授权

#### 统计数据
- `GET /api/admin/statistics/dashboard` - 仪表盘统计
- `GET /api/admin/statistics/revenue` - 收入统计
- `GET /api/admin/statistics/users` - 用户统计
- `GET /api/admin/statistics/plugins` - 插件统计

## Discourse集成

### 安装插件

用户购买后,在Discourse服务器上执行:

1. 编辑 `app.yml`:
```yaml
hooks:
  after_code:
    - git clone https://github.com/your-org/your-plugin.git
```

2. 重建容器:
```bash
./launcher rebuild app
```

### 更新插件

维护期内:
```bash
cd plugins/your-plugin
git pull
```

维护到期后,git pull会失败(403),但插件继续运行。

## 定时任务

系统使用cron调度以下任务:

| 任务 | 时间 | 说明 |
|------|------|------|
| 维护到期检查 | 每日02:00 | 检查并回收到期授权 |
| 统计数据汇总 | 每日01:00 | 汇总前一天统计数据 |
| 到期提醒 | 每日02:00 | 发送30/7/1天到期提醒 |

配置cron时间:
```env
CRON_MAINTENANCE_CHECK=0 2 * * *
```

## 数据库Schema

核心表结构:

- **users** - 用户表
- **github_accounts** - GitHub账号绑定
- **plugins** - 插件商品
- **orders** - 订单记录
- **licenses** - 授权记录 (核心)
- **license_history** - 授权历史 (审计)
- **tutorials** - 安装教程
- **email_notifications** - 邮件通知记录
- **statistics** - 统计数据
- **system_settings** - 系统设置

详见 [migrations/001_initial_schema.sql](migrations/001_initial_schema.sql)

## 安全说明

### 授权安全
- 所有GitHub仓库必须为Private
- 使用GitHub App细粒度权限控制
- 不在插件代码中做License校验
- 授权操作记录完整audit trail

### 数据安全
- 密码使用JWT token认证
- 敏感配置存储在环境变量
- GitHub access token加密存储
- 支付信息不存储到数据库

### API安全
- 所有用户接口需JWT认证
- 管理接口需admin角色
- CORS跨域保护
- Webhook签名验证

## 生产部署

### 使用Docker

```bash
# 构建镜像
docker build -t git-store .

# 运行容器
docker run -d \
  --name git-store \
  -p 8080:8080 \
  --env-file .env \
  git-store
```

### 使用systemd

```bash
# 编译二进制
go build -o git-store

# 复制systemd service文件
sudo cp deploy/git-store.service /etc/systemd/system/

# 启动服务
sudo systemctl start git-store
sudo systemctl enable git-store
```

### 反向代理 (Nginx)

```nginx
server {
    listen 80;
    server_name pluginstore.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 监控与日志

### 日志位置
- 应用日志: stdout
- 错误日志: stderr
- 定时任务日志: stdout

### 推荐监控指标
- API响应时间
- 支付成功率
- 授权回收成功率
- 邮件发送成功率

## 常见问题

### Q: 用户卸载了GitHub App怎么办?
A: 系统会检测到installation失效,暂停所有相关授权,并告警给管理员。

### Q: 仓库被删除了怎么办?
A: 插件自动下架,已购用户本地代码不受影响。

### Q: 支付成功但授权失败?
A: 系统会记录订单状态,管理员可在后台手动补授权。

### Q: 如何支持试用?
A: 创建LicenseType=trial的授权,设置较短的maintenance_until时间。

## 贡献指南

欢迎提交Issue和Pull Request!

## 开源协议

MIT License

## 联系方式

- Email: support@pluginstore.com
- GitHub: https://github.com/nodeloc/git-store
- 文档: https://docs.pluginstore.com

---

**这是一个正确的系统** ✓

- GitHub官方推荐的商业分发方式
- Discourse生态最干净的方案
- 后期维护成本最低的设计
