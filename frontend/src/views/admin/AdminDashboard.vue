<template>
  <div class="min-h-screen bg-base-200">
    <!-- Admin Header -->
    <div class="navbar bg-base-100 shadow-lg">
      <div class="flex-1">
        <a class="btn btn-ghost normal-case text-xl">🔧 GitStore 管理后台</a>
      </div>
      <div class="flex-none gap-2">
        <div class="badge badge-success">管理员</div>
        <div class="dropdown dropdown-end">
          <label tabindex="0" class="btn btn-ghost btn-circle avatar">
            <div class="w-10 rounded-full bg-primary text-primary-content flex items-center justify-center">
              {{ user?.name?.[0] || 'A' }}
            </div>
          </label>
          <ul tabindex="0" class="mt-3 p-2 shadow menu menu-compact dropdown-content bg-base-100 rounded-box w-52">
            <li><a>{{ user?.email }}</a></li>
            <li><a @click="logout">退出登录</a></li>
          </ul>
        </div>
      </div>
    </div>

    <div class="container mx-auto p-6">
      <!-- Stats -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="stats shadow bg-gradient-to-br from-primary to-primary-focus text-primary-content">
          <div class="stat">
            <div class="stat-title text-primary-content opacity-80">总用户数</div>
            <div class="stat-value">{{ stats.totalUsers }}</div>
            <div class="stat-desc text-primary-content">注册用户总数</div>
          </div>
        </div>
        
        <div class="stats shadow bg-gradient-to-br from-secondary to-secondary-focus text-secondary-content">
          <div class="stat">
            <div class="stat-title text-secondary-content opacity-80">总插件数</div>
            <div class="stat-value">{{ stats.totalPlugins }}</div>
            <div class="stat-desc text-secondary-content">已发布插件</div>
          </div>
        </div>
        
        <div class="stats shadow bg-gradient-to-br from-accent to-accent-focus text-accent-content">
          <div class="stat">
            <div class="stat-title text-accent-content opacity-80">总收入</div>
            <div class="stat-value text-2xl">${{ stats.totalRevenue }}</div>
            <div class="stat-desc text-accent-content">累计收入</div>
          </div>
        </div>
        
        <div class="stats shadow bg-gradient-to-br from-success to-success-focus text-success-content">
          <div class="stat">
            <div class="stat-title text-success-content opacity-80">活跃授权</div>
            <div class="stat-value">{{ stats.activeLicenses }}</div>
            <div class="stat-desc text-success-content">未过期授权</div>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="tabs tabs-boxed bg-base-100 shadow mb-6">
        <a class="tab" :class="{ 'tab-active': activeTab === 'plugins' }" @click="activeTab = 'plugins'">
          插件管理
        </a>
        <a class="tab" :class="{ 'tab-active': activeTab === 'users' }" @click="activeTab = 'users'">
          用户管理
        </a>
        <a class="tab" :class="{ 'tab-active': activeTab === 'orders' }" @click="activeTab = 'orders'">
          订单管理
        </a>
        <a class="tab" :class="{ 'tab-active': activeTab === 'licenses' }" @click="activeTab = 'licenses'">
          授权管理
        </a>
        <a class="tab" :class="{ 'tab-active': activeTab === 'settings' }" @click="activeTab = 'settings'">
          系统设置
        </a>
      </div>

      <!-- Content -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <!-- Plugins Tab -->
          <div v-if="activeTab === 'plugins'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="card-title">插件管理</h2>
              <button class="btn btn-primary" @click="createPlugin">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                创建插件
              </button>
            </div>
            <div v-if="loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else-if="plugins.length === 0" class="alert alert-info">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              <span>暂无插件，点击"创建插件"按钮添加第一个插件</span>
            </div>
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>名称</th>
                    <th>价格</th>
                    <th>状态</th>
                    <th>销量</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="plugin in plugins" :key="plugin.id">
                    <td>{{ plugin.name }}</td>
                    <td>${{ plugin.price }}</td>
                    <td>
                      <div class="badge" :class="plugin.status === 'active' ? 'badge-success' : 'badge-warning'">
                        {{ plugin.status }}
                      </div>
                    </td>
                    <td>{{ plugin.sales || 0 }}</td>
                    <td>
                      <button class="btn btn-ghost btn-xs">编辑</button>
                      <button class="btn btn-ghost btn-xs text-error">删除</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Users Tab -->
          <div v-if="activeTab === 'users'">
            <h2 class="card-title mb-4">用户管理</h2>
            <div v-if="loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else-if="users.length === 0" class="alert alert-info">
              <span>暂无用户数据</span>
            </div>
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>用户</th>
                    <th>邮箱</th>
                    <th>角色</th>
                    <th>注册时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="u in users" :key="u.id">
                    <td>
                      <div class="flex items-center gap-3">
                        <div class="avatar placeholder">
                          <div class="w-10 rounded-full bg-neutral-focus text-neutral-content">
                            <span>{{ u.name?.[0] || 'U' }}</span>
                          </div>
                        </div>
                        <div>{{ u.name }}</div>
                      </div>
                    </td>
                    <td>{{ u.email }}</td>
                    <td>
                      <div class="badge" :class="u.role === 'admin' ? 'badge-error' : 'badge-ghost'">
                        {{ u.role }}
                      </div>
                    </td>
                    <td>{{ formatDate(u.created_at) }}</td>
                    <td>
                      <button class="btn btn-ghost btn-xs">详情</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Orders Tab -->
          <div v-if="activeTab === 'orders'">
            <h2 class="card-title mb-4">订单管理</h2>
            <div class="alert alert-info">
              <span>订单管理功能开发中...</span>
            </div>
          </div>

          <!-- Licenses Tab -->
          <div v-if="activeTab === 'licenses'">
            <h2 class="card-title mb-4">授权管理</h2>
            <div class="alert alert-info">
              <span>授权管理功能开发中...</span>
            </div>
          </div>

          <!-- Settings Tab -->
          <div v-if="activeTab === 'settings'">
            <h2 class="card-title mb-4">系统设置</h2>
            <div class="space-y-4">
              <div class="form-control">
                <label class="label">
                  <span class="label-text">网站名称</span>
                </label>
                <input type="text" placeholder="GitStore" class="input input-bordered" />
              </div>
              <div class="form-control">
                <label class="label">
                  <span class="label-text">默认维护期限（月）</span>
                </label>
                <input type="number" placeholder="12" class="input input-bordered" />
              </div>
              <button class="btn btn-primary">保存设置</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'

const router = useRouter()
const authStore = useAuthStore()
const user = ref(authStore.user)

const activeTab = ref('plugins')
const loading = ref(false)

const stats = ref({
  totalUsers: 0,
  totalPlugins: 0,
  totalRevenue: 0,
  activeLicenses: 0
})

const plugins = ref([])
const users = ref([])

onMounted(async () => {
  await loadData()
})

async function loadData() {
  loading.value = true
  try {
    // 加载统计数据
    try {
      const statsRes = await api.get('/admin/statistics/dashboard')
      if (statsRes.data && typeof statsRes.data === 'object') {
        stats.value = { ...stats.value, ...statsRes.data }
      }
    } catch (e) {
      console.log('Stats not available:', e.message)
    }

    // 加载插件列表
    try {
      const pluginsRes = await api.get('/admin/plugins')
      if (Array.isArray(pluginsRes.data)) {
        plugins.value = pluginsRes.data
      } else if (pluginsRes.data?.plugins) {
        plugins.value = pluginsRes.data.plugins
      }
    } catch (e) {
      console.log('Plugins not available:', e.message)
    }

    // 加载用户列表
    try {
      const usersRes = await api.get('/admin/users')
      if (Array.isArray(usersRes.data)) {
        users.value = usersRes.data
      } else if (usersRes.data?.users) {
        users.value = usersRes.data.users
      }
    } catch (e) {
      console.log('Users not available:', e.message)
    }
  } finally {
    loading.value = false
  }
}

function createPlugin() {
  alert('创建插件功能开发中...')
}

function formatDate(date) {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

function logout() {
  authStore.logout()
  router.push('/')
}
</script>
