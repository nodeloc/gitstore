<template>
  <div class="min-h-screen bg-base-200">
    <div class="container mx-auto p-6">
      <!-- Stats -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="stats shadow bg-gradient-to-br from-primary to-primary-focus text-primary-content">
          <div class="stat">
            <div class="stat-title text-primary-content opacity-80">总用户数</div>
            <div class="stat-value">{{ adminStore.dashboardStats.total_users }}</div>
            <div class="stat-desc text-primary-content">今日新增 +{{ adminStore.dashboardStats.new_users_today }}</div>
          </div>
        </div>

        <div class="stats shadow bg-gradient-to-br from-secondary to-secondary-focus text-secondary-content">
          <div class="stat">
            <div class="stat-title text-secondary-content opacity-80">总插件数</div>
            <div class="stat-value">{{ adminStore.dashboardStats.total_plugins }}</div>
            <div class="stat-desc text-secondary-content">已发布插件</div>
          </div>
        </div>

        <div class="stats shadow bg-gradient-to-br from-accent to-accent-focus text-accent-content">
          <div class="stat">
            <div class="stat-title text-accent-content opacity-80">总收入</div>
            <div class="stat-value text-2xl">${{ adminStore.dashboardStats.total_revenue?.toFixed(2) || '0.00' }}</div>
            <div class="stat-desc text-accent-content">今日 +${{ adminStore.dashboardStats.revenue_today?.toFixed(2) || '0.00' }}</div>
          </div>
        </div>

        <div class="stats shadow bg-gradient-to-br from-success to-success-focus text-success-content">
          <div class="stat">
            <div class="stat-title text-success-content opacity-80">活跃授权</div>
            <div class="stat-value">{{ adminStore.dashboardStats.active_licenses }}</div>
            <div class="stat-desc text-success-content">未过期授权</div>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="tabs tabs-boxed bg-base-100 shadow mb-6">
        <a class="tab" :class="{ 'tab-active': activeTab === 'plugins' }" @click="switchTab('plugins')">
          插件管理
        </a>
        <a class="tab" :class="{ 'tab-active': activeTab === 'users' }" @click="switchTab('users')">
          用户管理
        </a>
        <a class="tab" :class="{ 'tab-active': activeTab === 'orders' }" @click="switchTab('orders')">
          订单管理
        </a>
        <a class="tab" :class="{ 'tab-active': activeTab === 'licenses' }" @click="switchTab('licenses')">
          授权管理
        </a>
        <a class="tab" :class="{ 'tab-active': activeTab === 'settings' }" @click="switchTab('settings')">
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
              <div class="flex gap-2">
                <input type="text" v-model="pluginSearch" placeholder="搜索插件..."
                       class="input input-bordered input-sm" @input="debouncedSearchPlugins" />
                <select v-model="pluginStatusFilter" class="select select-bordered select-sm" @change="loadPlugins(1)">
                  <option value="">全部状态</option>
                  <option value="published">已发布</option>
                  <option value="draft">草稿</option>
                  <option value="archived">已归档</option>
                </select>
                <button class="btn btn-primary btn-sm" @click="openPluginModal()">创建插件</button>
              </div>
            </div>
            <div v-if="adminStore.loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else-if="adminStore.plugins.length === 0" class="alert alert-info">
              <span>暂无插件，点击"创建插件"按钮添加</span>
            </div>
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>名称</th>
                    <th>Slug</th>
                    <th>价格</th>
                    <th>状态</th>
                    <th>下载量</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="plugin in adminStore.plugins" :key="plugin.id">
                    <td>
                      <div class="flex items-center gap-2">
                        <div class="avatar placeholder" v-if="plugin.icon_url">
                          <div class="w-8 rounded">
                            <img :src="plugin.icon_url" :alt="plugin.name" />
                          </div>
                        </div>
                        <span class="font-medium">{{ plugin.name }}</span>
                      </div>
                    </td>
                    <td class="text-sm opacity-70">{{ plugin.slug }}</td>
                    <td>${{ plugin.price?.toFixed(2) || '0.00' }}</td>
                    <td>
                      <div class="badge" :class="getStatusBadgeClass(plugin.status)">
                        {{ getStatusText(plugin.status) }}
                      </div>
                    </td>
                    <td>{{ plugin.download_count || 0 }}</td>
                    <td>
                      <button class="btn btn-ghost btn-xs" @click="openPluginModal(plugin)">编辑</button>
                      <button class="btn btn-ghost btn-xs text-error" @click="confirmDeletePlugin(plugin)">删除</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <!-- Pagination -->
            <div class="flex justify-center mt-4" v-if="adminStore.pluginsPagination.total_pages > 1">
              <div class="join">
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.pluginsPagination.page <= 1"
                        @click="loadPlugins(adminStore.pluginsPagination.page - 1)">«</button>
                <button class="join-item btn btn-sm">
                  {{ adminStore.pluginsPagination.page }} / {{ adminStore.pluginsPagination.total_pages }}
                </button>
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.pluginsPagination.page >= adminStore.pluginsPagination.total_pages"
                        @click="loadPlugins(adminStore.pluginsPagination.page + 1)">»</button>
              </div>
            </div>
          </div>

          <!-- Users Tab -->
          <div v-if="activeTab === 'users'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="card-title">用户管理</h2>
              <div class="flex gap-2">
                <input type="text" v-model="userSearch" placeholder="搜索用户..."
                       class="input input-bordered input-sm" @input="debouncedSearchUsers" />
                <select v-model="userRoleFilter" class="select select-bordered select-sm" @change="loadUsers(1)">
                  <option value="">全部角色</option>
                  <option value="admin">管理员</option>
                  <option value="user">普通用户</option>
                </select>
              </div>
            </div>
            <div v-if="adminStore.loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else-if="adminStore.users.length === 0" class="alert alert-info">
              <span>暂无用户数据</span>
            </div>
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>用户</th>
                    <th>邮箱</th>
                    <th>角色</th>
                    <th>状态</th>
                    <th>注册时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="u in adminStore.users" :key="u.id">
                    <td>
                      <div class="flex items-center gap-3">
                        <div class="avatar placeholder">
                          <div class="w-10 rounded-full bg-neutral-focus text-neutral-content">
                            <span>{{ u.name?.[0] || u.email?.[0] || 'U' }}</span>
                          </div>
                        </div>
                        <div>{{ u.name || '未设置' }}</div>
                      </div>
                    </td>
                    <td>{{ u.email }}</td>
                    <td>
                      <div class="badge" :class="u.role === 'admin' ? 'badge-error' : 'badge-ghost'">
                        {{ u.role === 'admin' ? '管理员' : '用户' }}
                      </div>
                    </td>
                    <td>
                      <div class="badge" :class="u.is_active ? 'badge-success' : 'badge-warning'">
                        {{ u.is_active ? '活跃' : '已停用' }}
                      </div>
                    </td>
                    <td>{{ formatDate(u.created_at) }}</td>
                    <td>
                      <button class="btn btn-ghost btn-xs" @click="openUserModal(u)">编辑</button>
                      <button class="btn btn-ghost btn-xs text-error" @click="confirmDeleteUser(u)"
                              v-if="u.id !== user?.id">停用</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <!-- Pagination -->
            <div class="flex justify-center mt-4" v-if="adminStore.usersPagination.total_pages > 1">
              <div class="join">
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.usersPagination.page <= 1"
                        @click="loadUsers(adminStore.usersPagination.page - 1)">«</button>
                <button class="join-item btn btn-sm">
                  {{ adminStore.usersPagination.page }} / {{ adminStore.usersPagination.total_pages }}
                </button>
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.usersPagination.page >= adminStore.usersPagination.total_pages"
                        @click="loadUsers(adminStore.usersPagination.page + 1)">»</button>
              </div>
            </div>
          </div>

          <!-- Orders Tab -->
          <div v-if="activeTab === 'orders'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="card-title">订单管理</h2>
              <div class="flex gap-2">
                <input type="text" v-model="orderSearch" placeholder="搜索订单号..."
                       class="input input-bordered input-sm" @input="debouncedSearchOrders" />
                <select v-model="orderStatusFilter" class="select select-bordered select-sm" @change="loadOrders(1)">
                  <option value="">全部状态</option>
                  <option value="pending">待支付</option>
                  <option value="paid">已支付</option>
                  <option value="refunded">已退款</option>
                  <option value="failed">失败</option>
                </select>
              </div>
            </div>
            <div v-if="adminStore.loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else-if="adminStore.orders.length === 0" class="alert alert-info">
              <span>暂无订单数据</span>
            </div>
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>订单号</th>
                    <th>用户</th>
                    <th>插件</th>
                    <th>金额</th>
                    <th>支付方式</th>
                    <th>状态</th>
                    <th>创建时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="order in adminStore.orders" :key="order.id">
                    <td class="font-mono text-sm">{{ order.order_number }}</td>
                    <td>{{ order.User?.name || order.User?.email || '-' }}</td>
                    <td>{{ order.Plugin?.name || '-' }}</td>
                    <td>${{ order.amount?.toFixed(2) }}</td>
                    <td>
                      <div class="badge badge-outline">{{ getPaymentMethodText(order.payment_method) }}</div>
                    </td>
                    <td>
                      <div class="badge" :class="getOrderStatusBadgeClass(order.payment_status)">
                        {{ getOrderStatusText(order.payment_status) }}
                      </div>
                    </td>
                    <td>{{ formatDate(order.created_at) }}</td>
                    <td>
                      <button class="btn btn-ghost btn-xs text-warning"
                              v-if="order.payment_status === 'paid'"
                              @click="confirmRefundOrder(order)">退款</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <!-- Pagination -->
            <div class="flex justify-center mt-4" v-if="adminStore.ordersPagination.total_pages > 1">
              <div class="join">
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.ordersPagination.page <= 1"
                        @click="loadOrders(adminStore.ordersPagination.page - 1)">«</button>
                <button class="join-item btn btn-sm">
                  {{ adminStore.ordersPagination.page }} / {{ adminStore.ordersPagination.total_pages }}
                </button>
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.ordersPagination.page >= adminStore.ordersPagination.total_pages"
                        @click="loadOrders(adminStore.ordersPagination.page + 1)">»</button>
              </div>
            </div>
          </div>

          <!-- Licenses Tab -->
          <div v-if="activeTab === 'licenses'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="card-title">授权管理</h2>
              <div class="flex gap-2">
                <select v-model="licenseStatusFilter" class="select select-bordered select-sm" @change="loadLicenses(1)">
                  <option value="">全部状态</option>
                  <option value="active">活跃</option>
                  <option value="expired">已过期</option>
                  <option value="revoked">已撤销</option>
                </select>
              </div>
            </div>
            <div v-if="adminStore.loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else-if="adminStore.licenses.length === 0" class="alert alert-info">
              <span>暂无授权数据</span>
            </div>
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>用户</th>
                    <th>插件</th>
                    <th>类型</th>
                    <th>维护到期</th>
                    <th>状态</th>
                    <th>创建时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="license in adminStore.licenses" :key="license.id">
                    <td>{{ license.User?.name || license.User?.email || '-' }}</td>
                    <td>{{ license.Plugin?.name || '-' }}</td>
                    <td>
                      <div class="badge badge-outline">
                        {{ license.license_type === 'permanent' ? '永久' : '试用' }}
                      </div>
                    </td>
                    <td>{{ formatDate(license.maintenance_until) }}</td>
                    <td>
                      <div class="badge" :class="getLicenseStatusBadgeClass(license.status)">
                        {{ getLicenseStatusText(license.status) }}
                      </div>
                    </td>
                    <td>{{ formatDate(license.created_at) }}</td>
                    <td>
                      <button class="btn btn-ghost btn-xs"
                              v-if="license.status === 'active'"
                              @click="openExtendModal(license)">延期</button>
                      <button class="btn btn-ghost btn-xs text-error"
                              v-if="license.status === 'active'"
                              @click="confirmRevokeLicense(license)">撤销</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <!-- Pagination -->
            <div class="flex justify-center mt-4" v-if="adminStore.licensesPagination.total_pages > 1">
              <div class="join">
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.licensesPagination.page <= 1"
                        @click="loadLicenses(adminStore.licensesPagination.page - 1)">«</button>
                <button class="join-item btn btn-sm">
                  {{ adminStore.licensesPagination.page }} / {{ adminStore.licensesPagination.total_pages }}
                </button>
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.licensesPagination.page >= adminStore.licensesPagination.total_pages"
                        @click="loadLicenses(adminStore.licensesPagination.page + 1)">»</button>
              </div>
            </div>
          </div>

          <!-- Settings Tab -->
          <div v-if="activeTab === 'settings'">
            <h2 class="card-title mb-4">系统设置</h2>
            <div v-if="adminStore.loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else class="space-y-4 max-w-xl">
              <div v-for="setting in adminStore.settings" :key="setting.key" class="form-control">
                <label class="label">
                  <span class="label-text">{{ setting.description || setting.key }}</span>
                </label>
                <input type="text" v-model="setting.value" class="input input-bordered" />
              </div>
              <button class="btn btn-primary" @click="saveSettings">保存设置</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Plugin Modal -->
    <dialog ref="pluginModal" class="modal">
      <div class="modal-box w-11/12 max-w-3xl">
        <h3 class="font-bold text-lg">{{ editingPlugin ? '编辑插件' : '创建插件' }}</h3>
        <form @submit.prevent="savePlugin" class="space-y-4 mt-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label"><span class="label-text">名称 *</span></label>
              <input type="text" v-model="pluginForm.name" class="input input-bordered" required />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">Slug *</span></label>
              <input type="text" v-model="pluginForm.slug" class="input input-bordered" required />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">价格 (USD)</span></label>
              <input type="number" v-model="pluginForm.price" step="0.01" min="0" class="input input-bordered" />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">状态</span></label>
              <select v-model="pluginForm.status" class="select select-bordered">
                <option value="draft">草稿</option>
                <option value="published">已发布</option>
                <option value="archived">已归档</option>
              </select>
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">分类</span></label>
              <input type="text" v-model="pluginForm.category" class="input input-bordered" />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">版本</span></label>
              <input type="text" v-model="pluginForm.version" class="input input-bordered" placeholder="1.0.0" />
            </div>
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">简短描述</span></label>
            <textarea v-model="pluginForm.description" class="textarea textarea-bordered" rows="2"></textarea>
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">详细描述</span></label>
            <textarea v-model="pluginForm.long_description" class="textarea textarea-bordered" rows="4"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="form-control col-span-2">
              <label class="label">
                <span class="label-text">GitHub 仓库</span>
                <button type="button" class="btn btn-sm btn-ghost" @click="loadGitHubRepos">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                  刷新列表
                </button>
              </label>
              <select v-if="githubRepos.length > 0" v-model="selectedGitHubRepo" class="select select-bordered" @change="onGitHubRepoSelected">
                <option :value="null">手动输入或选择仓库</option>
                <option v-for="repo in githubRepos" :key="repo.id" :value="repo">
                  {{ repo.full_name }} {{ repo.private ? '(私有)' : '' }}
                </option>
              </select>
              <div v-if="loadingRepos" class="flex items-center justify-center py-2">
                <span class="loading loading-spinner loading-sm mr-2"></span>
                <span class="text-sm">加载仓库列表...</span>
              </div>
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">GitHub 仓库 URL</span></label>
              <input type="text" v-model="pluginForm.github_repo_url" class="input input-bordered" />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">图标 URL</span></label>
              <input type="text" v-model="pluginForm.icon_url" class="input input-bordered" />
            </div>
          </div>
          <div class="modal-action">
            <button type="button" class="btn" @click="closePluginModal">取消</button>
            <button type="submit" class="btn btn-primary">{{ editingPlugin ? '保存' : '创建' }}</button>
          </div>
        </form>
      </div>
    </dialog>

    <!-- User Modal -->
    <dialog ref="userModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">编辑用户</h3>
        <form @submit.prevent="saveUser" class="space-y-4 mt-4">
          <div class="form-control">
            <label class="label"><span class="label-text">用户名</span></label>
            <input type="text" v-model="userForm.name" class="input input-bordered" />
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">角色</span></label>
            <select v-model="userForm.role" class="select select-bordered">
              <option value="user">普通用户</option>
              <option value="admin">管理员</option>
            </select>
          </div>
          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">账户状态</span>
              <input type="checkbox" v-model="userForm.is_active" class="toggle toggle-success" />
            </label>
          </div>
          <div class="modal-action">
            <button type="button" class="btn" @click="closeUserModal">取消</button>
            <button type="submit" class="btn btn-primary">保存</button>
          </div>
        </form>
      </div>
    </dialog>

    <!-- Extend License Modal -->
    <dialog ref="extendModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">延长授权</h3>
        <form @submit.prevent="extendLicense" class="space-y-4 mt-4">
          <div class="form-control">
            <label class="label"><span class="label-text">延长月数</span></label>
            <select v-model="extendMonths" class="select select-bordered">
              <option :value="1">1 个月</option>
              <option :value="3">3 个月</option>
              <option :value="6">6 个月</option>
              <option :value="12">12 个月</option>
              <option :value="24">24 个月</option>
            </select>
          </div>
          <div class="modal-action">
            <button type="button" class="btn" @click="closeExtendModal">取消</button>
            <button type="submit" class="btn btn-primary">确认延期</button>
          </div>
        </form>
      </div>
    </dialog>

    <!-- Confirm Dialog -->
    <dialog ref="confirmModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">{{ confirmTitle }}</h3>
        <p class="py-4">{{ confirmMessage }}</p>
        <div class="form-control" v-if="confirmNeedsReason">
          <label class="label"><span class="label-text">原因</span></label>
          <textarea v-model="confirmReason" class="textarea textarea-bordered" rows="2" required></textarea>
        </div>
        <div class="modal-action">
          <button class="btn" @click="closeConfirmModal">取消</button>
          <button class="btn btn-error" @click="executeConfirm">确认</button>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAdminStore } from '@/stores/admin'

const router = useRouter()
const authStore = useAuthStore()
const adminStore = useAdminStore()
const user = ref(authStore.user)

const activeTab = ref('plugins')

// Search and filters
const pluginSearch = ref('')
const pluginStatusFilter = ref('')
const userSearch = ref('')
const userRoleFilter = ref('')
const orderSearch = ref('')
const orderStatusFilter = ref('')
const licenseStatusFilter = ref('')

// Plugin form
const pluginModal = ref(null)
const editingPlugin = ref(null)
const pluginForm = ref({
  name: '',
  slug: '',
  description: '',
  long_description: '',
  price: 0,
  status: 'draft',
  category: '',
  version: '',
  github_repo_url: '',
  github_repo_id: 0,
  github_repo_name: '',
  icon_url: ''
})

// GitHub repos
const githubRepos = ref([])
const selectedGitHubRepo = ref(null)
const loadingRepos = ref(false)

// User form
const userModal = ref(null)
const editingUser = ref(null)
const userForm = ref({
  name: '',
  role: 'user',
  is_active: true
})

// Extend modal
const extendModal = ref(null)
const extendingLicense = ref(null)
const extendMonths = ref(12)

// Confirm modal
const confirmModal = ref(null)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmNeedsReason = ref(false)
const confirmReason = ref('')
const confirmAction = ref(null)

// Debounce helper
function debounce(fn, delay) {
  let timer
  return (...args) => {
    clearTimeout(timer)
    timer = setTimeout(() => fn(...args), delay)
  }
}

const debouncedSearchPlugins = debounce(() => loadPlugins(1), 300)
const debouncedSearchUsers = debounce(() => loadUsers(1), 300)
const debouncedSearchOrders = debounce(() => loadOrders(1), 300)

onMounted(async () => {
  await adminStore.fetchDashboardStats()
  await loadPlugins()
})

function switchTab(tab) {
  activeTab.value = tab
  if (tab === 'plugins') loadPlugins()
  else if (tab === 'users') loadUsers()
  else if (tab === 'orders') loadOrders()
  else if (tab === 'licenses') loadLicenses()
  else if (tab === 'settings') adminStore.fetchSettings()
}

async function loadPlugins(page = 1) {
  await adminStore.fetchPlugins({
    page,
    page_size: 10,
    search: pluginSearch.value,
    status: pluginStatusFilter.value
  })
}

async function loadUsers(page = 1) {
  await adminStore.fetchUsers({
    page,
    page_size: 10,
    search: userSearch.value,
    role: userRoleFilter.value
  })
}

async function loadOrders(page = 1) {
  await adminStore.fetchOrders({
    page,
    page_size: 10,
    search: orderSearch.value,
    status: orderStatusFilter.value
  })
}

async function loadLicenses(page = 1) {
  await adminStore.fetchLicenses({
    page,
    page_size: 10,
    status: licenseStatusFilter.value
  })
}

// GitHub repos functions
async function loadGitHubRepos() {
  loadingRepos.value = true
  try {
    const response = await adminStore.fetchGitHubRepos()
    githubRepos.value = response.repositories || []
    if (githubRepos.value.length === 0) {
      alert('未找到可用的 GitHub 仓库，请检查 GitHub App 配置')
    }
  } catch (err) {
    alert('加载 GitHub 仓库失败: ' + (err.response?.data?.error || err.message))
  } finally {
    loadingRepos.value = false
  }
}

function onGitHubRepoSelected() {
  if (selectedGitHubRepo.value) {
    const repo = selectedGitHubRepo.value
    pluginForm.value.github_repo_url = repo.html_url
    pluginForm.value.github_repo_id = repo.id
    pluginForm.value.github_repo_name = repo.full_name
    
    // 自动填充名称和描述（如果为空）
    if (!pluginForm.value.name) {
      pluginForm.value.name = repo.name
    }
    if (!pluginForm.value.slug) {
      pluginForm.value.slug = repo.name.toLowerCase().replace(/[^a-z0-9-]/g, '-')
    }
    if (!pluginForm.value.description && repo.description) {
      pluginForm.value.description = repo.description
    }
  }
}

// Plugin CRUD
function openPluginModal(plugin = null) {
  editingPlugin.value = plugin
  selectedGitHubRepo.value = null
  
  if (plugin) {
    pluginForm.value = { ...plugin }
  } else {
    pluginForm.value = {
      name: '',
      slug: '',
      description: '',
      long_description: '',
      price: 0,
      status: 'draft',
      category: '',
      version: '',
      github_repo_url: '',
      github_repo_id: 0,
      github_repo_name: '',
      icon_url: ''
    }
    // 自动加载 GitHub 仓库列表（仅在创建新插件时）
    if (githubRepos.value.length === 0) {
      loadGitHubRepos()
    }
  }
  pluginModal.value?.showModal()
}

function closePluginModal() {
  pluginModal.value?.close()
  editingPlugin.value = null
  selectedGitHubRepo.value = null
}

async function savePlugin() {
  try {
    if (editingPlugin.value) {
      await adminStore.updatePlugin(editingPlugin.value.id, pluginForm.value)
    } else {
      await adminStore.createPlugin(pluginForm.value)
    }
    closePluginModal()
    await loadPlugins()
  } catch (err) {
    alert('保存失败: ' + (err.response?.data?.error || err.message))
  }
}

function confirmDeletePlugin(plugin) {
  confirmTitle.value = '确认删除'
  confirmMessage.value = `确定要删除插件 "${plugin.name}" 吗？此操作不可撤销。`
  confirmNeedsReason.value = false
  confirmAction.value = async () => {
    try {
      await adminStore.deletePlugin(plugin.id)
      await loadPlugins()
    } catch (err) {
      alert('删除失败: ' + (err.response?.data?.error || err.message))
    }
  }
  confirmModal.value?.showModal()
}

// User CRUD
function openUserModal(u) {
  editingUser.value = u
  userForm.value = {
    name: u.name || '',
    role: u.role || 'user',
    is_active: u.is_active !== false
  }
  userModal.value?.showModal()
}

function closeUserModal() {
  userModal.value?.close()
  editingUser.value = null
}

async function saveUser() {
  try {
    await adminStore.updateUser(editingUser.value.id, userForm.value)
    closeUserModal()
    await loadUsers()
  } catch (err) {
    alert('保存失败: ' + (err.response?.data?.error || err.message))
  }
}

function confirmDeleteUser(u) {
  confirmTitle.value = '确认停用'
  confirmMessage.value = `确定要停用用户 "${u.name || u.email}" 吗？`
  confirmNeedsReason.value = false
  confirmAction.value = async () => {
    try {
      await adminStore.deleteUser(u.id)
      await loadUsers()
    } catch (err) {
      alert('操作失败: ' + (err.response?.data?.error || err.message))
    }
  }
  confirmModal.value?.showModal()
}

// Order actions
function confirmRefundOrder(order) {
  confirmTitle.value = '确认退款'
  confirmMessage.value = `确定要退款订单 "${order.order_number}" ($${order.amount?.toFixed(2)}) 吗？相关授权将被撤销。`
  confirmNeedsReason.value = false
  confirmAction.value = async () => {
    try {
      await adminStore.refundOrder(order.id)
      await loadOrders()
    } catch (err) {
      alert('退款失败: ' + (err.response?.data?.error || err.message))
    }
  }
  confirmModal.value?.showModal()
}

// License actions
function openExtendModal(license) {
  extendingLicense.value = license
  extendMonths.value = 12
  extendModal.value?.showModal()
}

function closeExtendModal() {
  extendModal.value?.close()
  extendingLicense.value = null
}

async function extendLicense() {
  try {
    await adminStore.extendLicense(extendingLicense.value.id, extendMonths.value)
    closeExtendModal()
    await loadLicenses()
  } catch (err) {
    alert('延期失败: ' + (err.response?.data?.error || err.message))
  }
}

function confirmRevokeLicense(license) {
  confirmTitle.value = '确认撤销授权'
  confirmMessage.value = `确定要撤销 "${license.User?.name || license.User?.email}" 对 "${license.Plugin?.name}" 的授权吗？`
  confirmNeedsReason.value = true
  confirmReason.value = ''
  confirmAction.value = async () => {
    if (!confirmReason.value) {
      alert('请填写撤销原因')
      return
    }
    try {
      await adminStore.revokeLicense(license.id, confirmReason.value)
      await loadLicenses()
    } catch (err) {
      alert('撤销失败: ' + (err.response?.data?.error || err.message))
    }
  }
  confirmModal.value?.showModal()
}

// Settings
async function saveSettings() {
  try {
    await adminStore.updateSettings(adminStore.settings)
    alert('设置已保存')
  } catch (err) {
    alert('保存失败: ' + (err.response?.data?.error || err.message))
  }
}

// Confirm modal
function closeConfirmModal() {
  confirmModal.value?.close()
  confirmAction.value = null
  confirmReason.value = ''
}

async function executeConfirm() {
  if (confirmAction.value) {
    await confirmAction.value()
  }
  closeConfirmModal()
}

// Helpers
function formatDate(date) {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

function getStatusBadgeClass(status) {
  switch (status) {
    case 'published': return 'badge-success'
    case 'draft': return 'badge-warning'
    case 'archived': return 'badge-ghost'
    default: return 'badge-ghost'
  }
}

function getStatusText(status) {
  switch (status) {
    case 'published': return '已发布'
    case 'draft': return '草稿'
    case 'archived': return '已归档'
    default: return status
  }
}

function getOrderStatusBadgeClass(status) {
  switch (status) {
    case 'paid': return 'badge-success'
    case 'pending': return 'badge-warning'
    case 'refunded': return 'badge-info'
    case 'failed': return 'badge-error'
    default: return 'badge-ghost'
  }
}

function getOrderStatusText(status) {
  switch (status) {
    case 'paid': return '已支付'
    case 'pending': return '待支付'
    case 'refunded': return '已退款'
    case 'failed': return '失败'
    default: return status
  }
}

function getPaymentMethodText(method) {
  switch (method) {
    case 'stripe': return 'Stripe'
    case 'paypal': return 'PayPal'
    case 'alipay': return '支付宝'
    default: return method
  }
}

function getLicenseStatusBadgeClass(status) {
  switch (status) {
    case 'active': return 'badge-success'
    case 'expired': return 'badge-warning'
    case 'revoked': return 'badge-error'
    default: return 'badge-ghost'
  }
}

function getLicenseStatusText(status) {
  switch (status) {
    case 'active': return '活跃'
    case 'expired': return '已过期'
    case 'revoked': return '已撤销'
    default: return status
  }
}

function logout() {
  authStore.logout()
  router.push('/')
}
</script>
