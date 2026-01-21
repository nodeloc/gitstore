# GitStore Frontend Setup Guide

## 项目结构

已创建的Vue 3 + Vite前端项目,包含以下功能:

### ✅ 已完成

1. **项目配置**
   - package.json - 依赖管理
   - vite.config.js - Vite配置
   - tailwind.config.js - TailwindCSS配置
   - postcss.config.js - PostCSS配置

2. **国际化 (i18n)**
   - 英文 (en.json)
   - 中文 (zh.json)
   - 支持语言切换

3. **路由配置**
   - 首页、插件列表、插件详情
   - 用户Dashboard、授权管理、订单管理
   - 管理后台
   - 认证保护

4. **状态管理 (Pinia)**
   - auth store - 用户认证
   - plugins store - 插件数据

5. **核心组件**
   - NavBar - 导航栏(含语言切换、主题切换、用户菜单)
   - Footer - 页脚

6. **工具类**
   - api.js - Axios封装,自动添加Token

### 📝 需要手动创建的视图组件

在 `src/views/` 目录下创建以下文件:

1. **HomeView.vue** - 首页
2. **PluginsView.vue** - 插件列表页
3. **PluginDetailView.vue** - 插件详情页
4. **AuthCallbackView.vue** - OAuth回调页
5. **DashboardView.vue** - 用户Dashboard
6. **LicensesView.vue** - 授权管理页
7. **OrdersView.vue** - 订单列表页
8. **PurchaseView.vue** - 购买页面
9. **NotFoundView.vue** - 404页面
10. **admin/AdminDashboard.vue** - 管理后台

## 安装步骤

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 配置环境变量

创建 `.env` 文件:

```env
VITE_API_BASE_URL=http://localhost:8080/api
```

### 3. 启动开发服务器

```bash
npm run dev
```

访问: http://localhost:3000

### 4. 构建生产版本

```bash
npm run build
```

构建文件将生成在 `dist/` 目录

## 视图组件模板

### HomeView.vue 模板

```vue
<template>
  <div>
    <!-- Hero Section -->
    <div class="hero min-h-screen bg-base-200">
      <div class="hero-content text-center">
        <div class="max-w-md">
          <h1 class="text-5xl font-bold">{{ $t('hero.title') }}</h1>
          <p class="py-6">{{ $t('hero.subtitle') }}</p>
          <RouterLink to="/plugins" class="btn btn-primary">
            {{ $t('hero.browseButton') }}
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Features Section -->
    <div class="container mx-auto px-4 py-12">
      <h2 class="text-3xl font-bold text-center mb-8">
        {{ $t('features.title') }}
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- Feature cards -->
      </div>
    </div>
  </div>
</template>

<script setup>
</script>
```

### PluginsView.vue 模板

```vue
<template>
  <div class="container mx-auto px-4 py-12">
    <h1 class="text-4xl font-bold mb-8">{{ $t('plugins.title') }}</h1>

    <!-- Search and Filter -->
    <div class="mb-8 flex gap-4">
      <input 
        v-model="searchQuery" 
        type="text" 
        :placeholder="$t('plugins.search')" 
        class="input input-bordered flex-1"
      />
      <select v-model="selectedCategory" class="select select-bordered">
        <option value="">{{ $t('plugins.allCategories') }}</option>
        <!-- Categories -->
      </select>
    </div>

    <!-- Plugin Grid -->
    <div v-if="loading" class="text-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
      <p>{{ $t('plugins.loading') }}</p>
    </div>

    <div v-else-if="plugins.length === 0" class="text-center py-12">
      <p>{{ $t('plugins.noPlugins') }}</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <PluginCard 
        v-for="plugin in plugins" 
        :key="plugin.id" 
        :plugin="plugin" 
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { usePluginsStore } from '@/stores/plugins'
import PluginCard from '@/components/PluginCard.vue'

const pluginsStore = usePluginsStore()
const searchQuery = ref('')
const selectedCategory = ref('')

const { plugins, loading } = pluginsStore

onMounted(() => {
  pluginsStore.fetchPlugins()
})

watch([searchQuery, selectedCategory], () => {
  pluginsStore.fetchPlugins({
    search: searchQuery.value,
    category: selectedCategory.value
  })
})
</script>
```

## API集成说明

所有API调用都通过 `src/utils/api.js` 进行,会自动:
- 添加Authorization header
- 处理401/403错误
- 统一错误处理

示例:
```js
import api from '@/utils/api'

// GET请求
const response = await api.get('/plugins')

// POST请求
const response = await api.post('/orders', { plugin_id: '123' })
```

## 主题和语言

### 切换语言
```js
import { useI18n } from 'vue-i18n'
const { locale } = useI18n()
locale.value = 'zh' // 或 'en'
```

### 切换主题
```js
document.documentElement.setAttribute('data-theme', 'dark')
```

可用主题: light, dark, cupcake, corporate

## 下一步

1. 根据上面的模板创建所有视图组件
2. 创建PluginCard组件用于插件展示
3. 测试与后端API的集成
4. 优化用户体验和样式

## 文件清单

前端项目已创建的文件:
- ✅ package.json
- ✅ vite.config.js
- ✅ tailwind.config.js
- ✅ postcss.config.js
- ✅ index.html
- ✅ src/main.js
- ✅ src/App.vue
- ✅ src/assets/main.css
- ✅ src/router/index.js
- ✅ src/stores/auth.js
- ✅ src/stores/plugins.js
- ✅ src/utils/api.js
- ✅ src/i18n/index.js
- ✅ src/i18n/locales/en.json
- ✅ src/i18n/locales/zh.json
- ✅ src/components/NavBar.vue
- ✅ src/components/Footer.vue

需要创建的视图文件可以参考上面的模板,或者根据需求自定义。
