# 项目结构

```
git-store/
├── .env.example                      # 环境变量配置示例
├── .gitignore                        # Git忽略文件
├── README.md                         # 项目说明文档
├── DEPLOYMENT.md                     # 部署指南
├── PROJECT_SUMMARY.md                # 项目总结
├── STRUCTURE.md                      # 项目结构说明(本文件)
├── Dockerfile                        # Docker镜像构建文件
├── docker-compose.yml                # Docker编排配置
├── Makefile                          # 构建脚本
├── go.mod                            # Go依赖管理
├── main.go                           # 应用入口文件
│
├── migrations/                       # 数据库迁移文件
│   └── 001_initial_schema.sql       # 初始化数据库Schema
│
├── internal/                         # 内部代码(不对外暴露)
│   │
│   ├── config/                       # 配置管理
│   │   └── config.go                # 配置加载和定义
│   │
│   ├── database/                     # 数据库
│   │   └── database.go              # 数据库连接和迁移
│   │
│   ├── models/                       # 数据模型
│   │   └── models.go                # GORM模型定义
│   │
│   ├── middleware/                   # HTTP中间件
│   │   └── auth.go                  # JWT认证 + CORS + 角色权限
│   │
│   ├── utils/                        # 工具函数
│   │   ├── jwt.go                   # JWT生成和验证
│   │   └── helpers.go               # 通用辅助函数
│   │
│   ├── services/                     # 业务服务层
│   │   ├── github_oauth.go          # GitHub OAuth认证
│   │   ├── github_app.go            # GitHub App授权管理
│   │   ├── payment_stripe.go        # Stripe支付集成
│   │   ├── payment_paypal.go        # PayPal支付集成
│   │   ├── payment_alipay.go        # 支付宝支付集成
│   │   └── email.go                 # 邮件发送服务
│   │
│   ├── handlers/                     # HTTP处理器(Controller层)
│   │   ├── auth.go                  # 认证相关处理器
│   │   └── handlers_all.go          # 其他处理器(插件/订单/支付/授权/教程/管理/统计)
│   │
│   ├── router/                       # 路由配置
│   │   └── router.go                # 路由注册和分组
│   │
│   └── scheduler/                    # 定时任务
│       └── scheduler.go             # Cron任务(到期回收/统计汇总)
│
└── frontend/                         # 前端页面
    └── index.html                    # 主页面(DaisyUI + Tailwind CSS)
```

## 代码统计

### 文件数量
- Go源代码: 14个文件
- 配置文件: 5个
- 文档文件: 4个
- 前端文件: 1个
- 数据库文件: 1个
- **总计: 25个文件**

### 代码行数(估算)
- Go代码: ~3500行
- SQL: ~400行
- HTML/JS: ~300行
- 文档: ~2000行
- **总计: ~6200行**

## 模块说明

### 1. 配置管理 (config)
- 加载环境变量
- 提供统一配置接口
- 支持默认值

### 2. 数据库 (database)
- PostgreSQL连接
- GORM自动迁移
- 连接池管理

### 3. 数据模型 (models)
- User - 用户
- GitHubAccount - GitHub账号绑定
- Plugin - 插件商品
- Order - 订单
- License - 授权(核心)
- LicenseHistory - 授权历史
- Tutorial - 教程
- EmailNotification - 邮件通知
- Statistic - 统计数据
- SystemSetting - 系统设置

### 4. 中间件 (middleware)
- JWT认证
- CORS跨域
- 管理员权限检查

### 5. 工具函数 (utils)
- JWT生成
- 订单号生成
- Slug转换
- 日期计算

### 6. 业务服务 (services)

#### GitHub OAuth Service
- 用户登录
- 获取用户信息
- 获取用户邮箱
- 获取组织列表

#### GitHub App Service
- Installation认证
- 授予仓库访问权限
- 撤销仓库访问权限
- 列出可访问仓库
- 创建Installation Token

#### Payment Services
- **Stripe**: 创建支付意图、确认支付、Webhook验证
- **PayPal**: 创建订单、捕获支付、查询订单
- **Alipay**: 创建支付、回调验证

#### Email Service
- 发送购买成功邮件
- 发送到期提醒邮件
- 发送到期通知邮件
- 发送续费成功邮件
- 模板渲染

### 7. HTTP处理器 (handlers)

#### Auth Handler
- GitHub OAuth登录
- 回调处理
- 获取当前用户
- 登出
- 获取GitHub账号列表

#### Plugin Handler
- 列表查询
- 详情查询

#### Order Handler
- 创建订单
- 查询订单
- 用户订单列表

#### Payment Handler
- Stripe支付
- PayPal支付
- Alipay支付
- Webhook处理

#### License Handler
- 用户授权列表
- 授权详情
- 授权续费
- 授权历史

#### Tutorial Handler
- 公开教程列表
- 教程详情
- 授权后教程访问

#### Admin Handler
- 插件管理(CRUD)
- 订单管理(查询/退款)
- 授权管理(查询/撤销/延长)
- 教程管理(CRUD)
- 用户管理(CRUD)
- 系统设置
- GitHub仓库同步

#### Dashboard Handler
- 仪表盘统计
- 收入统计
- 用户统计
- 插件统计

### 8. 路由 (router)
- 公开路由 (12个)
- 认证路由 (11个)
- 管理路由 (30+个)

### 9. 定时任务 (scheduler)
- 维护到期检查
- 到期授权回收
- 到期提醒邮件
- 统计数据汇总

### 10. 前端 (frontend)
- DaisyUI组件库
- Tailwind CSS样式
- Axios API调用
- GitHub OAuth集成
- 插件列表展示
- 购买流程

## 数据流向

### 购买流程
```
用户 → Frontend → API Handler → Order Service
  ↓
Payment Service (Stripe/PayPal/Alipay)
  ↓
Webhook → Payment Handler → Create License
  ↓
GitHub App Service → Grant Repository Access
  ↓
Email Service → Send Success Email
```

### 到期回收流程
```
Cron Scheduler → Check Expired Licenses
  ↓
GitHub App Service → Revoke Repository Access
  ↓
Update License Status
  ↓
Email Service → Send Expiry Email
```

## API分层架构

```
HTTP Request
    ↓
Middleware (Auth, CORS)
    ↓
Router
    ↓
Handler (Controller)
    ↓
Service (Business Logic)
    ↓
Model (Data Access)
    ↓
Database
```

## 依赖关系

```
main.go
  ├── config
  ├── database
  ├── router
  │     ├── handlers
  │     │     ├── services
  │     │     ├── models
  │     │     └── utils
  │     └── middleware
  └── scheduler
        ├── services
        └── models
```

## 关键文件说明

### main.go
- 应用入口
- 加载配置
- 初始化数据库
- 注册路由
- 启动定时任务
- 启动HTTP服务器

### config/config.go
- 30+配置项
- 环境变量加载
- 默认值设置

### models/models.go
- 10个数据模型
- GORM注解
- 关联关系定义
- UUID生成钩子

### router/router.go
- 路由分组
- 中间件应用
- Handler注册
- API版本管理

### scheduler/scheduler.go
- Cron表达式配置
- 到期检查逻辑
- 统计汇总逻辑
- 错误处理

## 代码规范

### 命名规范
- 文件名: snake_case
- 包名: 小写单词
- 结构体: PascalCase
- 函数: PascalCase (公开) / camelCase (私有)
- 变量: camelCase

### 项目约定
- 所有业务逻辑在 Service 层
- Handler 只负责请求响应
- 使用 GORM 进行数据库操作
- 错误统一返回 JSON 格式
- 日志使用标准库 log

### 注释规范
- 每个导出函数必须有注释
- 复杂逻辑必须有行内注释
- 每个包必须有包注释

## 部署文件

### Dockerfile
- 多阶段构建
- Alpine基础镜像
- 最小化镜像体积

### docker-compose.yml
- PostgreSQL服务
- 应用服务
- 数据持久化
- 健康检查

### Makefile
- build: 构建二进制
- run: 运行应用
- test: 运行测试
- docker-build: 构建镜像
- docker-up: 启动容器

## 文档文件

### README.md (主文档)
- 项目介绍
- 快速开始
- 配置说明
- API文档
- 常见问题

### DEPLOYMENT.md (部署文档)
- 系统要求
- 前置准备
- 本地开发
- 生产部署
- 监控维护
- 故障排查

### PROJECT_SUMMARY.md (项目总结)
- 功能清单
- 技术架构
- 业务流程
- 成本估算
- 后续开发

### STRUCTURE.md (本文件)
- 目录结构
- 模块说明
- 代码统计
- 依赖关系

## 扩展点

### 1. 添加新支付方式
在 `internal/services/` 创建 `payment_xxx.go`

### 2. 添加新通知渠道
在 `internal/services/email.go` 扩展或创建新文件

### 3. 添加新的定时任务
在 `internal/scheduler/scheduler.go` 添加新的 cron job

### 4. 添加新的API端点
1. 在 `internal/handlers/` 创建handler方法
2. 在 `internal/router/router.go` 注册路由

### 5. 添加新的数据模型
1. 在 `internal/models/models.go` 定义结构体
2. 运行数据库迁移

## 测试建议

### 单元测试结构
```
internal/
├── config/
│   └── config_test.go
├── services/
│   ├── github_oauth_test.go
│   ├── github_app_test.go
│   ├── payment_stripe_test.go
│   └── email_test.go
├── handlers/
│   └── auth_test.go
└── utils/
    ├── jwt_test.go
    └── helpers_test.go
```

### 集成测试
- API端点测试
- 数据库操作测试
- 第三方服务Mock测试

## 性能考虑

### 当前瓶颈
1. 数据库直连
2. 同步API调用
3. 内存状态管理

### 优化方案
1. 添加Redis缓存
2. 异步任务队列
3. 数据库连接池优化
4. CDN加速

## 安全检查清单

- [x] JWT认证
- [x] 角色权限控制
- [x] SQL注入防护 (GORM)
- [x] XSS防护
- [x] CORS配置
- [x] 敏感信息加密
- [x] Webhook签名验证
- [x] 最小权限原则 (GitHub App)
- [ ] API速率限制 (待实现)
- [ ] HTTPS强制 (生产环境配置)

## 监控指标

### 应用层
- API响应时间
- 错误率
- 并发连接数

### 业务层
- 订单量
- 支付成功率
- 授权回收成功率
- 邮件发送成功率

### 基础设施
- CPU使用率
- 内存使用率
- 数据库连接数
- 磁盘使用率

---

**更新时间**: 2024-01-21
