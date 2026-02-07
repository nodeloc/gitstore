# 🛍️ gitstore

<div align="center">

**Plugin Store and License Management Platform**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](docker-compose.prod.yml)

Commercial distribution platform designed for Discourse and Flarum plugins, supporting innovative licensing model with permanent license + time-limited maintenance

[Live Demo](https://discourseplugin.com) · [Quick Deploy](#-quick-deployment) · [Full Documentation](DEPLOYMENT.md) · [中文文档](README.zh.md)

</div>

---

## 💡 Why Git-Store?

Are you an open-source developer creating plugins for Discourse, Flarum, or other platforms? Do you want to monetize your work while maintaining the open-source spirit?

**Git-Store solves the monetization challenge for plugin developers:**

- 📦 **Sell Your Plugins** - Transform your side projects into a sustainable income stream
- 🔓 **Keep Code Open** - Users own the code forever after purchase, no DRM lock-in
- ⚡ **Zero Integration Work** - No need to add license verification code to your plugins
- 🤝 **Fair Model** - Buyers get permanent access, you get paid for ongoing support
- 🎯 **Focus on Creating** - We handle payments, access control, and customer management

**Perfect for:**
- 👨‍💻 Independent developers with premium plugins
- 🏢 Teams offering commercial extensions
- 💼 Agencies providing custom forum solutions
- 🎓 Educators selling educational plugins

Stop worrying about piracy and license verification. Start earning from your creativity.

---

## ✨ Key Features

### 🎯 Unique Licensing Model
- **Permanent License** - One-time purchase, own the plugin code forever
- **Time-Limited Maintenance** - 12 months of updates and support by default
- **Zero Intrusive Design** - No license verification code needed in plugins
- **Flexible Access Control** - Manage plugin access through GitHub organizations

### 🔐 Deep GitHub Integration
- GitHub OAuth login, no additional account system required
- Support for both personal and organization accounts
- Direct integration with GitHub repositories
- Automated access management for purchased plugins

### 💳 Multiple Payment Methods
| Payment | Coverage | Status |
|---------|----------|--------|
| 💎 Stripe | Global | ✅ Integrated |
| 💙 PayPal | Global | 🚧 In Progress |
| 💚 Alipay | China | 🚧 In Progress |

### 📧 Smart Notification System
- ✉️ Instant purchase confirmation
- ⏰ Multi-level maintenance expiry reminders (30/7/1 days)
- 🔔 Maintenance expired notifications
- ✅ Renewal confirmation

### 📊 Data Analytics
- Real-time sales statistics and revenue trends
- User growth and activity analysis
- Plugin download and usage statistics
- Exportable visual reports

---

## 🏗️ Tech Stack

<table>
<tr>
<td><strong>Backend</strong></td>
<td>Go 1.21+ / Gin Framework / GORM</td>
</tr>
<tr>
<td><strong>Frontend</strong></td>
<td>Vue 3.4 / Vite 5 / TailwindCSS / DaisyUI</td>
</tr>
<tr>
<td><strong>Database</strong></td>
<td>PostgreSQL 14+</td>
</tr>
<tr>
<td><strong>Authentication</strong></td>
<td>GitHub OAuth + JWT</td>
</tr>
<tr>
<td><strong>Payment</strong></td>
<td>Stripe / PayPal / Alipay</td>
</tr>
<tr>
<td><strong>Deployment</strong></td>
<td>Docker / Nginx Proxy Manager</td>
</tr>
</table>

---

## 🚀 Quick Deployment

### Option 1: Docker One-Click Deploy (Recommended)

```bash
# 1. Clone repository
git clone https://github.com/nodeloc/gitstore.git /opt/gitstore
cd /opt/gitstore

# 2. Configure environment variables
cp .env.example .env
nano .env  # Modify necessary configurations

# 3. Run deployment script
export DOMAIN="your-domain.com"
sudo ./deploy.sh
```

After deployment, visit `http://YOUR_SERVER_IP:81` to configure Nginx Proxy Manager.

For detailed steps, see [NPM_SETUP.md](NPM_SETUP.md)

### Option 2: Local Development

```bash
# 1. Start database (using Docker)
docker-compose up -d postgres

# 2. Run migrations
psql -h localhost -p 5433 -U postgres -d git_store < migrations/001_initial_schema.sql

# 3. Start backend
go run main.go

# 4. Start frontend
cd frontend
npm install
npm run dev
```

Visit `http://localhost:3001`

---

## ⚙️ Configuration Guide

### 1️⃣ GitHub OAuth App

Visit [GitHub Developer Settings](https://github.com/settings/developers) to create an OAuth App:

- **Application name**: `gitstore`
- **Homepage URL**: `https://your-domain.com`
- **Callback URL**: `https://your-domain.com/api/auth/github/callback`

Get **Client ID** and **Client Secret**

### 2️⃣ GitHub Personal Access Token

To manage repository access for purchased plugins, you need a GitHub Personal Access Token:

1. Visit [GitHub Tokens](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Configure the token:
   - **Note**: `gitstore admin token`
   - **Expiration**: No expiration (or custom)
   - **Scopes**: Select `admin:org` → `write:org` (for managing organization members)
4. Generate and copy the token

Add to `.env`:
```env
GITHUB_ADMIN_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
```

**Note**: This token is used to automatically add/remove users to your GitHub organization when they purchase/expire plugins.

### 3️⃣ Stripe Payment Setup

Visit [Stripe Dashboard](https://dashboard.stripe.com/apikeys):

```env
STRIPE_SECRET_KEY=sk_live_xxx        # Use 'live' for production, 'test' for testing
STRIPE_PUBLISHABLE_KEY=pk_live_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx
```

**Create Webhook** (for receiving payment status):
- URL: `https://your-domain.com/api/webhooks/stripe`
- Events: `payment_intent.succeeded`, `payment_intent.payment_failed`

### 4️⃣ Environment Variables

Complete `.env` configuration example:

```env
# Application
APP_ENV=production
APP_PORT=8080
APP_URL=https://your-domain.com
FRONTEND_URL=https://your-domain.com

# Database
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

# Email (Optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@your-domain.com
```

---

## 📁 Project Structure

```
gitstore/
├── main.go                      # Application entry
├── internal/                    # Backend core code
│   ├── config/                  # Configuration management
│   ├── database/                # Database connection
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # Middlewares
│   ├── models/                  # Data models
│   ├── router/                  # Routing
│   ├── scheduler/               # Scheduled tasks
│   ├── services/                # Business logic
│   │   ├── payment_stripe.go    # Stripe payment
│   │   ├── payment_paypal.go    # PayPal payment
│   │   ├── payment_alipay.go    # Alipay
│   │   ├── github.go            # GitHub API
│   │   └── email.go             # Email sending
│   └── utils/                   # Utilities
├── frontend/                    # Vue 3 frontend
│   ├── src/
│   │   ├── components/          # Components
│   │   ├── views/               # Pages
│   │   ├── router/              # Routing
│   │   ├── stores/              # State management
│   │   ├── i18n/                # Internationalization
│   │   └── utils/               # Utilities
│   └── dist/                    # Build output
├── migrations/                  # Database migrations
├── docker-compose.prod.yml      # Production Docker config
├── Dockerfile                   # Backend image
├── deploy.sh                    # One-click deployment script
└── NPM_SETUP.md                 # Nginx Proxy Manager setup guide
```

---

## 🎮 Usage Flow

### Admin Side

1. **Login**: Sign in with GitHub account
2. **Create Categories**: Set categories for plugins (e.g., Feature Enhancement, Theme)
3. **Add Plugins**:
   - Fill in plugin name, description, price
   - Set GitHub repository URL
   - Configure default maintenance period (months)
4. **Manage Orders**: View order status and revenue statistics

### User Side

1. **Browse Plugins**: Browse available plugins in the store
2. **Purchase**: Select plugin, complete payment via Stripe
3. **Get Access**: Automatically receive GitHub private repository access
4. **Install Plugin**: Clone/download plugin code via GitHub
5. **Renew Maintenance**: Renew before expiry to extend update privileges

---

## 📝 API Documentation

### Core Endpoints

| Endpoint | Method | Description | Auth |
|----------|--------|-------------|------|
| `/api/health` | GET | Health check | No |
| `/api/auth/github` | GET | GitHub OAuth login | No |
| `/api/plugins` | GET | Get plugin list | No |
| `/api/plugins/:slug` | GET | Get plugin details | No |
| `/api/orders` | POST | Create order | Required |
| `/api/orders` | GET | Get my orders | Required |
| `/api/licenses` | GET | Get my licenses | Required |
| `/api/payment/stripe` | POST | Create Stripe payment | Required |
| `/api/webhooks/stripe` | POST | Stripe Webhook | No |

Full API Documentation: [View Swagger Docs](http://localhost:8080/swagger) (In Development)

---

## 🔧 Operations

### View Logs

```bash
# Backend logs
docker-compose -f docker-compose.prod.yml logs -f backend

# Database logs
docker-compose -f docker-compose.prod.yml logs -f postgres

# NPM logs
docker-compose -f docker-compose.prod.yml logs -f nginx-proxy-manager
```

### Database Backup

```bash
# Backup
docker exec gitstore-db-prod pg_dump -U postgres git_store > backup_$(date +%Y%m%d).sql

# Restore
cat backup_20260207.sql | docker exec -i gitstore-db-prod psql -U postgres git_store
```

### Update Deployment

```bash
cd /opt/gitstore
git pull origin main
sudo ./deploy.sh
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Create a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

---

## 🔗 Related Links

- [Full Deployment Guide](DEPLOYMENT.md)
- [Nginx Proxy Manager Setup](NPM_SETUP.md)
- [Quick Start Guide](QUICKSTART.md)
- [Project Structure](STRUCTURE.md)
- [Changelog](CHANGELOG.md)

---

## 💬 Support & Feedback

- 🐛 [Report Bug](https://github.com/nodeloc/gitstore/issues)
- 💡 [Feature Request](https://github.com/nodeloc/gitstore/issues)
- 📧 Email: james@nodeloc.com
- 💬 Forum: [Join Community](https://www.nodeloc.com)

---

<div align="center">

**⭐ If you find this helpful, please give it a Star ⭐**

Made with ❤️ by [NodeLoc](https://github.com/nodeloc)

</div>

</div>