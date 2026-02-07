# 🛍️ gitstore

<div align="center">

**插件商店和许可证管理平台**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](docker-compose.prod.yml)

专为 Discourse 和 Flarum 插件设计的商业分发平台，支持永久许可证 + 时限维护的创新授权模式

[在线演示](https://discourseplugin.com) · [快速部署](#-快速部署) · [完整文档](DEPLOYMENT.md) · [English](README.md)

</div>

---

## 💡 为什么选择 Git-Store？

你是一位为 Discourse、Flarum 或其他平台开发插件的开源开发者吗？想要将你的创意项目变现，同时保持开源精神？

**Git-Store 解决了插件开发者的变现难题：**

- 📦 **出售你的插件** - 将业余项目转化为可持续的收入来源
- 🔓 **保持代码开放** - 用户购买后永久拥有代码，无 DRM 锁定
- ⚡ **零集成工作** - 无需在插件中添加授权验证代码
- 🤝 **公平模式** - 买家获得永久访问权，你获得持续支持的报酬
- 🎯 **专注创作** - 我们处理支付、访问控制和客户管理

**完美适用于：**
- 👨‍💻 拥有高级插件的独立开发者
- 🏢 提供商业扩展的团队
- 💼 提供定制论坛解决方案的代理商
- 🎓 销售教育类插件的教育工作者

不再担心盗版和许可证验证。开始从你的创意中获得收益。

---

## ✨ 核心特性

### 🎯 独特的授权模式
- **永久许可证** - 一次购买，插件代码永久拥有
- **时限维护期** - 默认12个月技术支持和更新
- **零侵入设计** - 无需在代码中集成授权验证逻辑
- **灵活访问控制** - 通过 GitHub 组织管理插件访问权限

### 🔐 GitHub 深度集成
- GitHub OAuth 登录，无需额外账号体系
- 支持个人和组织账号
- 直接集成 GitHub 仓库
- 已购插件的自动访问管理

### 💳 多种支付方式
| 支付方式 | 覆盖地区 | 状态 |
|---------|---------|------|
| 💎 Stripe | 全球 | ✅ 已集成 |
| 💙 PayPal | 全球 | 🚧 开发中 |
| 💚 支付宝 | 中国 | 🚧 开发中 |

### 📧 智能通知系统
- ✉️ 购买成功即时通知
- ⏰ 维护到期多级提醒（30天/7天/1天）
- 🔔 维护已到期通知
- ✅ 续费成功确认

### 📊 数据分析
- 实时销售统计和收入趋势
- 用户增长和活跃度分析
- 插件下载和使用统计
- 可视化报表导出

---

## 🏗️ 技术架构

<table>
<tr>
<td><strong>后端</strong></td>
<td>Go 1.21+ / Gin Framework / GORM</td>
</tr>
<tr>
<td><strong>前端</strong></td>
<td>Vue 3.4 / Vite 5 / TailwindCSS / DaisyUI</td>
</tr>
<tr>
<td><strong>数据库</strong></td>
<td>PostgreSQL 14+</td>
</tr>
<tr>
<td><strong>认证</strong></td>
<td>GitHub OAuth + JWT</td>
</tr>
<tr>
<td><strong>支付</strong></td>
<td>Stripe / PayPal / Alipay</td>
</tr>
<tr>
<td><strong>部署</strong></td>
<td>Docker / Nginx Proxy Manager</td>
</tr>
</table>

---

## 🚀 快速部署

### 方式一：Docker 一键部署（推荐）

```bash
# 1. 克隆代码
git clone https://github.com/nodeloc/gitstore.git /opt/gitstore
cd /opt/gitstore

# 2. 配置环境变量
cp .env.example .env
nano .env  # 修改必要的配置

# 3. 执行部署脚本
export DOMAIN="your-domain.com"
sudo ./deploy.sh
```

部署完成后访问 `http://YOUR_SERVER_IP:81` 配置 Nginx Proxy Manager。

详细步骤请参考 [NPM_SETUP.md](NPM_SETUP.md)

### 方式二：本地开发

```bash
# 1. 启动数据库（使用 Docker）
docker-compose up -d postgres

# 2. 运行迁移
psql -h localhost -p 5433 -U postgres -d git_store < migrations/001_initial_schema.sql

# 3. 启动后端
go run main.go

# 4. 启动前端
cd frontend
npm install
npm run dev
```

访问 `http://localhost:3001`

---

## ⚙️ 配置指南

### 1️⃣ GitHub OAuth 应用

访问 [GitHub Developer Settings](https://github.com/settings/developers) 创建 OAuth App：

- **Application name**: `gitstore`
- **Homepage URL**: `https://your-domain.com`
- **Callback URL**: `https://your-domain.com/api/auth/github/callback`

获取 **Client ID** 和 **Client Secret**

### 2️⃣ GitHub 个人访问令牌

为了管理已购插件的仓库访问权限，你需要一个 GitHub 个人访问令牌：

1. 访问 [GitHub Tokens](https://github.com/settings/tokens)
2. 点击「Generate new token (classic)」
3. 配置令牌：
   - **Note**: `gitstore admin token`
   - **Expiration**: No expiration（或自定义）
   - **Scopes**: 选择 `admin:org` → `write:org`（用于管理组织成员）
4. 生成并复制令牌

添加到 `.env`：
```env
GITHUB_ADMIN_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
```

**说明**：此令牌用于在用户购买/过期插件时，自动添加/移除用户到你的 GitHub 组织。

### 3️⃣ Stripe 支付配置

访问 [Stripe Dashboard](https://dashboard.stripe.com/apikeys)：

```env
STRIPE_SECRET_KEY=sk_live_xxx        # 生产环境用 live，测试用 test
STRIPE_PUBLISHABLE_KEY=pk_live_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx
```

**创建 Webhook**（用于接收支付状态）：
- URL: `https://your-domain.com/api/webhooks/stripe`
- Events: `payment_intent.succeeded`, `payment_intent.payment_failed`

### 4️⃣ 环境变量配置

完整的 `.env` 配置示例：

```env
# 应用配置
APP_ENV=production
APP_PORT=8080
APP_URL=https://your-domain.com
FRONTEND_URL=https://your-domain.com

# 数据库
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_strong_password
DB_NAME=git_store

# GitHub OAuth
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
GITHUB_REDIRECT_URL=https://your-domain.com/api/auth/github/callback
GITHUB_ORG_NAME=your-org-name
GITHUB_ADMIN_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx

# JWT
JWT_SECRET=your-random-64-character-secret-key
JWT_EXPIRY_HOURS=720

# Stripe
STRIPE_SECRET_KEY=sk_live_xxx
STRIPE_PUBLISHABLE_KEY=pk_live_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx

# 邮件（可选）
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@your-domain.com
```

---

## 📁 项目结构

```
gitstore/
├── main.go                      # 应用入口
├── internal/                    # 后端核心代码
│   ├── config/                  # 配置管理
│   ├── database/                # 数据库连接
│   ├── handlers/                # HTTP 处理器
│   ├── middleware/              # 中间件
│   ├── models/                  # 数据模型
│   ├── router/                  # 路由配置
│   ├── scheduler/               # 定时任务
│   ├── services/                # 业务逻辑
│   │   ├── payment_stripe.go    # Stripe 支付
│   │   ├── payment_paypal.go    # PayPal 支付
│   │   ├── payment_alipay.go    # 支付宝
│   │   ├── github.go            # GitHub API
│   │   └── email.go             # 邮件发送
│   └── utils/                   # 工具函数
├── frontend/                    # Vue 3 前端
│   ├── src/
│   │   ├── components/          # 组件
│   │   ├── views/               # 页面
│   │   ├── router/              # 路由
│   │   ├── stores/              # 状态管理
│   │   ├── i18n/                # 国际化
│   │   └── utils/               # 工具
│   └── dist/                    # 构建输出
├── migrations/                  # 数据库迁移
├── docker-compose.prod.yml      # 生产环境 Docker 配置
├── Dockerfile                   # 后端镜像
├── deploy.sh                    # 一键部署脚本
└── NPM_SETUP.md                 # Nginx Proxy Manager 配置指南
```

---

## 🎮 使用流程

### 管理员端

1. **登录后台**：使用 GitHub 账号登录
2. **创建分类**：为插件设置分类（如：功能增强、主题美化）
3. **添加插件**：
   - 填写插件名称、描述、价格
   - 设置 GitHub 仓库地址
   - 配置默认维护期（月）
4. **管理订单**：查看订单状态、收入统计

### 用户端

1. **浏览插件**：在商店浏览可用插件
2. **购买插件**：选择插件，使用 Stripe 完成支付
3. **获取访问权**：自动获得 GitHub 私有仓库访问权限
4. **安装插件**：通过 GitHub 克隆/下载插件代码
5. **维护续费**：到期前可续费延长更新权限

---

## 📝 API 文档

### 核心端点

| 端点 | 方法 | 说明 | 认证 |
|------|-----|------|-----|
| `/api/health` | GET | 健康检查 | 无 |
| `/api/auth/github` | GET | GitHub OAuth 登录 | 无 |
| `/api/plugins` | GET | 获取插件列表 | 无 |
| `/api/plugins/:slug` | GET | 获取插件详情 | 无 |
| `/api/orders` | POST | 创建订单 | 需要 |
| `/api/orders` | GET | 获取我的订单 | 需要 |
| `/api/licenses` | GET | 获取我的许可证 | 需要 |
| `/api/payment/stripe` | POST | 创建 Stripe 支付 | 需要 |
| `/api/webhooks/stripe` | POST | Stripe Webhook | 无 |

完整 API 文档：[查看 Swagger 文档](http://localhost:8080/swagger)（开发中）

---

## 🔧 运维管理

### 查看日志

```bash
# 后端日志
docker-compose -f docker-compose.prod.yml logs -f backend

# 数据库日志
docker-compose -f docker-compose.prod.yml logs -f postgres

# NPM 日志
docker-compose -f docker-compose.prod.yml logs -f nginx-proxy-manager
```

### 数据库备份

```bash
# 备份
docker exec gitstore-db-prod pg_dump -U postgres git_store > backup_$(date +%Y%m%d).sql

# 恢复
cat backup_20260207.sql | docker exec -i gitstore-db-prod psql -U postgres git_store
```

### 更新部署

```bash
cd /opt/gitstore
git pull origin main
sudo ./deploy.sh
```

---

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 🔗 相关链接

- [完整部署文档](DEPLOYMENT.md)
- [Nginx Proxy Manager 配置](NPM_SETUP.md)
- [快速开始指南](QUICKSTART.md)
- [项目结构说明](STRUCTURE.md)
- [更新日志](CHANGELOG.md)

---

## 💬 支持与反馈

- 🐛 [提交 Bug](https://github.com/nodeloc/gitstore/issues)
- 💡 [功能建议](https://github.com/nodeloc/gitstore/issues)
- 📧 Email: james@nodeloc.com
- 💬 Forum: [加入社区](https://www.nodeloc.com)

---

<div align="center">

**⭐ 如果觉得有帮助，请给个 Star ⭐**

Made with ❤️ by [NodeLoc](https://github.com/nodeloc)

</div>
