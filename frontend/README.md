# GitStore Frontend

Vue 3 + Vite前端项目,使用DaisyUI和TailwindCSS构建。

## 特性

- ✅ Vue 3 + Composition API
- ✅ Vite 5 - 快速开发构建
- ✅ Vue Router 4 - 路由管理
- ✅ Pinia - 状态管理
- ✅ Vue I18n - 国际化支持 (中文/英文)
- ✅ DaisyUI - UI组件库
- ✅ TailwindCSS - 实用优先的CSS框架
- ✅ Axios - HTTP客户端
- ✅ 响应式设计
- ✅ 主题切换 (4个主题)
- ✅ 语言切换
- ✅ JWT认证集成

## 项目结构

```
frontend/
├── public/              # 静态资源
├── src/
│   ├── assets/         # 样式和静态文件
│   │   └── main.css    # 主样式文件
│   ├── components/     # Vue组件
│   │   ├── NavBar.vue
│   │   └── Footer.vue
│   ├── i18n/           # 国际化配置
│   │   ├── index.js
│   │   └── locales/
│   │       ├── en.json # 英文翻译
│   │       └── zh.json # 中文翻译
│   ├── router/         # 路由配置
│   │   └── index.js
│   ├── stores/         # Pinia状态管理
│   │   ├── auth.js     # 认证store
│   │   └── plugins.js  # 插件store
│   ├── utils/          # 工具函数
│   │   └── api.js      # Axios封装
│   ├── views/          # 页面组件
│   │   └── (待创建)
│   ├── App.vue         # 根组件
│   └── main.js         # 入口文件
├── index.html
├── package.json
├── vite.config.js
├── tailwind.config.js
└── postcss.config.js
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

访问: http://localhost:3000

### 生产构建

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 环境配置

复制 `.env.example` 为 `.env`:

```bash
cp .env.example .env
```

编辑 `.env` 文件:

```env
VITE_API_BASE_URL=http://localhost:8080/api
```

## API集成

### 使用方式

```js
import api from '@/utils/api'

// GET请求
const response = await api.get('/plugins')

// POST请求
const response = await api.post('/orders', { plugin_id: '123' })

// 自动添加Authorization header
// 自动处理错误
```

### API实例特性

- 自动添加JWT token到请求头
- 自动处理401(未授权)错误 - 跳转到登录
- 统一错误处理
- 30秒超时设置

## 状态管理

### Auth Store

```js
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 登录
await authStore.login()

// 获取用户信息
await authStore.fetchUser()

// 登出
authStore.logout()

// 状态
authStore.isAuthenticated  // 是否已登录
authStore.isAdmin          // 是否是管理员
authStore.user             // 用户信息
```

### Plugins Store

```js
import { usePluginsStore } from '@/stores/plugins'

const pluginsStore = usePluginsStore()

// 获取插件列表
await pluginsStore.fetchPlugins({ search: 'keyword', category: 'analytics' })

// 获取单个插件
await pluginsStore.fetchPlugin('plugin-slug')

// 状态
pluginsStore.plugins       // 插件列表
pluginsStore.currentPlugin // 当前插件
pluginsStore.loading       // 加载状态
```

## 国际化

### 切换语言

```js
import { useI18n } from 'vue-i18n'

const { locale, t } = useI18n()

// 切换到中文
locale.value = 'zh'

// 切换到英文
locale.value = 'en'

// 使用翻译
t('nav.home')  // 首页 或 Home
```

### 在模板中使用

```vue
<template>
  <div>
    {{ $t('nav.home') }}
    {{ $t('hero.title') }}
  </div>
</template>
```

### 添加新翻译

编辑 `src/i18n/locales/en.json` 和 `zh.json`:

```json
{
  "newSection": {
    "key": "Value"
  }
}
```

## 主题系统

### 可用主题

- light (默认)
- dark
- cupcake
- corporate

### 切换主题

```js
// 设置主题
document.documentElement.setAttribute('data-theme', 'dark')

// 保存到localStorage
localStorage.setItem('theme', 'dark')
```

### 在组件中使用

主题由DaisyUI提供,支持自动颜色切换。

## 路由

### 路由列表

| 路径 | 组件 | 描述 | 权限 |
|------|------|------|------|
| / | HomeView | 首页 | 公开 |
| /plugins | PluginsView | 插件列表 | 公开 |
| /plugins/:slug | PluginDetailView | 插件详情 | 公开 |
| /auth/callback | AuthCallbackView | OAuth回调 | 公开 |
| /dashboard | DashboardView | 用户Dashboard | 需登录 |
| /licenses | LicensesView | 授权管理 | 需登录 |
| /orders | OrdersView | 订单列表 | 需登录 |
| /purchase/:pluginId | PurchaseView | 购买页面 | 需登录 |
| /admin | AdminDashboard | 管理后台 | 需管理员 |

### 路由守卫

路由自动检查认证状态和管理员权限,未授权用户会被重定向。

## 组件开发

### 创建新组件

```vue
<template>
  <div class="component-name">
    <!-- 内容 -->
  </div>
</template>

<script setup>
import { ref } from 'vue'

// 组件逻辑
</script>

<style scoped>
/* 组件样式 */
</style>
```

### 使用DaisyUI组件

```vue
<template>
  <button class="btn btn-primary">按钮</button>
  <div class="card">卡片</div>
  <input class="input input-bordered" />
</template>
```

[DaisyUI文档](https://daisyui.com/components/)

## 待创建的视图

需要在 `src/views/` 创建以下组件:

1. **HomeView.vue** - 首页
   - Hero区域
   - 特性展示
   - 精选插件

2. **PluginsView.vue** - 插件列表
   - 搜索和筛选
   - 插件卡片网格
   - 分页

3. **PluginDetailView.vue** - 插件详情
   - 插件信息
   - 购买按钮
   - 安装说明

4. **AuthCallbackView.vue** - OAuth回调处理

5. **DashboardView.vue** - 用户Dashboard
   - 统计信息
   - 快捷操作

6. **LicensesView.vue** - 授权管理
   - 授权列表
   - 续费功能
   - GitHub仓库链接

7. **OrdersView.vue** - 订单列表
   - 订单历史
   - 订单详情

8. **PurchaseView.vue** - 购买页面
   - 支付方式选择
   - 订单确认

9. **NotFoundView.vue** - 404页面

10. **admin/AdminDashboard.vue** - 管理后台
    - 插件管理
    - 订单管理
    - 用户管理
    - 统计数据

参考 [FRONTEND_SETUP.md](../FRONTEND_SETUP.md) 获取视图组件模板。

## 部署

### 构建

```bash
npm run build
```

### 部署到Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
    }
}
```

### Docker部署

添加到主项目的Dockerfile或创建单独的前端Dockerfile。

## 开发规范

### 命名规范

- 组件文件: PascalCase (HomeView.vue)
- 普通文件: camelCase (api.js)
- CSS类: kebab-case (btn-primary)

### 代码风格

- 使用Composition API
- 使用 `<script setup>`
- 使用TailwindCSS和DaisyUI类
- 避免内联样式

### Git提交

```bash
git add .
git commit -m "feat: add plugin list view"
git push
```

## 故障排查

### 开发服务器无法启动

```bash
# 清除node_modules
rm -rf node_modules package-lock.json
npm install
```

### API请求失败

检查:
1. 后端服务是否运行在8080端口
2. `.env`文件中的API地址是否正确
3. 浏览器Console是否有CORS错误

### 热重载不工作

```bash
# 重启开发服务器
npm run dev
```

## 资源链接

- [Vue 3文档](https://vuejs.org/)
- [Vite文档](https://vitejs.dev/)
- [Vue Router文档](https://router.vuejs.org/)
- [Pinia文档](https://pinia.vuejs.org/)
- [DaisyUI文档](https://daisyui.com/)
- [TailwindCSS文档](https://tailwindcss.com/)
- [Vue I18n文档](https://vue-i18n.intlify.dev/)

## 许可证

MIT
