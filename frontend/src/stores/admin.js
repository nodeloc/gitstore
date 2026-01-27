import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/api'

export const useAdminStore = defineStore('admin', () => {
  // Dashboard Stats
  const dashboardStats = ref({
    total_users: 0,
    total_plugins: 0,
    total_orders: 0,
    total_revenue: 0,
    active_licenses: 0,
    new_users_today: 0,
    new_orders_today: 0,
    revenue_today: 0
  })

  // Plugins
  const plugins = ref([])
  const pluginsPagination = ref({ page: 1, page_size: 10, total: 0, total_pages: 0 })

  // Users
  const users = ref([])
  const usersPagination = ref({ page: 1, page_size: 10, total: 0, total_pages: 0 })

  // Orders
  const orders = ref([])
  const ordersPagination = ref({ page: 1, page_size: 10, total: 0, total_pages: 0 })

  // Licenses
  const licenses = ref([])
  const licensesPagination = ref({ page: 1, page_size: 10, total: 0, total_pages: 0 })

  // Tutorials
  const tutorials = ref([])
  const tutorialsPagination = ref({ page: 1, page_size: 10, total: 0, total_pages: 0 })

  // Settings
  const settings = ref([])

  // Loading state
  const loading = ref(false)
  const error = ref(null)

  // ==================== Dashboard ====================
  async function fetchDashboardStats() {
    try {
      loading.value = true
      const response = await api.get('/admin/statistics/dashboard')
      dashboardStats.value = response.data
    } catch (err) {
      error.value = err.message
      console.error('Failed to fetch dashboard stats:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchRevenueStats(period = '30d') {
    try {
      const response = await api.get('/admin/statistics/revenue', { params: { period } })
      return response.data
    } catch (err) {
      console.error('Failed to fetch revenue stats:', err)
      return { period, data: [] }
    }
  }

  async function fetchUserStats(period = '30d') {
    try {
      const response = await api.get('/admin/statistics/users', { params: { period } })
      return response.data
    } catch (err) {
      console.error('Failed to fetch user stats:', err)
      return { period, data: [] }
    }
  }

  async function fetchPluginStats() {
    try {
      const response = await api.get('/admin/statistics/plugins')
      return response.data.plugins || []
    } catch (err) {
      console.error('Failed to fetch plugin stats:', err)
      return []
    }
  }

  // ==================== Plugins ====================
  async function fetchPlugins(params = {}) {
    try {
      loading.value = true
      const response = await api.get('/admin/plugins', { params })
      plugins.value = response.data.plugins || []
      pluginsPagination.value = response.data.pagination || pluginsPagination.value
    } catch (err) {
      error.value = err.message
      console.error('Failed to fetch plugins:', err)
    } finally {
      loading.value = false
    }
  }

  async function createPlugin(pluginData) {
    const response = await api.post('/admin/plugins', pluginData)
    return response.data.plugin
  }

  async function getPlugin(id) {
    const response = await api.get(`/admin/plugins/${id}`)
    return response.data.plugin
  }

  async function updatePlugin(id, pluginData) {
    const response = await api.put(`/admin/plugins/${id}`, pluginData)
    return response.data.plugin
  }

  async function deletePlugin(id) {
    await api.delete(`/admin/plugins/${id}`)
  }

  async function syncGitHubRepos() {
    const response = await api.post('/admin/plugins/sync-repos')
    return response.data.repositories || []
  }

  // ==================== Users ====================
  async function fetchUsers(params = {}) {
    try {
      loading.value = true
      const response = await api.get('/admin/users', { params })
      users.value = response.data.users || []
      usersPagination.value = response.data.pagination || usersPagination.value
    } catch (err) {
      error.value = err.message
      console.error('Failed to fetch users:', err)
    } finally {
      loading.value = false
    }
  }

  async function getUser(id) {
    const response = await api.get(`/admin/users/${id}`)
    return response.data.user
  }

  async function updateUser(id, userData) {
    const response = await api.put(`/admin/users/${id}`, userData)
    return response.data.user
  }

  async function deleteUser(id) {
    await api.delete(`/admin/users/${id}`)
  }

  // ==================== Orders ====================
  async function fetchOrders(params = {}) {
    try {
      loading.value = true
      const response = await api.get('/admin/orders', { params })
      orders.value = response.data.orders || []
      ordersPagination.value = response.data.pagination || ordersPagination.value
    } catch (err) {
      error.value = err.message
      console.error('Failed to fetch orders:', err)
    } finally {
      loading.value = false
    }
  }

  async function getOrder(id) {
    const response = await api.get(`/admin/orders/${id}`)
    return response.data.order
  }

  async function updateOrderStatus(id, paymentStatus) {
    const response = await api.put(`/admin/orders/${id}/status`, { payment_status: paymentStatus })
    return response.data
  }

  async function refundOrder(id) {
    const response = await api.post(`/admin/orders/${id}/refund`)
    return response.data
  }

  // ==================== Licenses ====================
  async function fetchLicenses(params = {}) {
    try {
      loading.value = true
      const response = await api.get('/admin/licenses', { params })
      licenses.value = response.data.licenses || []
      licensesPagination.value = response.data.pagination || licensesPagination.value
    } catch (err) {
      error.value = err.message
      console.error('Failed to fetch licenses:', err)
    } finally {
      loading.value = false
    }
  }

  async function getLicense(id) {
    const response = await api.get(`/admin/licenses/${id}`)
    return response.data.license
  }

  async function revokeLicense(id, reason) {
    const response = await api.post(`/admin/licenses/${id}/revoke`, { reason })
    return response.data
  }

  async function extendLicense(id, months) {
    const response = await api.post(`/admin/licenses/${id}/extend`, { months })
    return response.data
  }

  // ==================== Tutorials ====================
  async function fetchTutorials(params = {}) {
    try {
      loading.value = true
      const response = await api.get('/admin/tutorials', { params })
      tutorials.value = response.data.tutorials || []
      tutorialsPagination.value = response.data.pagination || tutorialsPagination.value
    } catch (err) {
      error.value = err.message
      console.error('Failed to fetch tutorials:', err)
    } finally {
      loading.value = false
    }
  }

  async function createTutorial(tutorialData) {
    const response = await api.post('/admin/tutorials', tutorialData)
    return response.data.tutorial
  }

  async function getTutorial(id) {
    const response = await api.get(`/admin/tutorials/${id}`)
    return response.data.tutorial
  }

  async function updateTutorial(id, tutorialData) {
    const response = await api.put(`/admin/tutorials/${id}`, tutorialData)
    return response.data.tutorial
  }

  async function deleteTutorial(id) {
    await api.delete(`/admin/tutorials/${id}`)
  }

  // ==================== Settings ====================
  async function fetchSettings() {
    try {
      loading.value = true
      const response = await api.get('/admin/settings')
      settings.value = response.data.settings || []
    } catch (err) {
      error.value = err.message
      console.error('Failed to fetch settings:', err)
    } finally {
      loading.value = false
    }
  }

  async function updateSettings(settingsData) {
    const response = await api.put('/admin/settings', { settings: settingsData })
    return response.data
  }

  // ==================== GitHub ====================
  async function fetchGitHubRepos() {
    const response = await api.get('/admin/github/repositories')
    return response.data
  }

  return {
    // State
    dashboardStats,
    plugins,
    pluginsPagination,
    users,
    usersPagination,
    orders,
    ordersPagination,
    licenses,
    licensesPagination,
    tutorials,
    tutorialsPagination,
    settings,
    loading,
    error,

    // Dashboard Actions
    fetchDashboardStats,
    fetchRevenueStats,
    fetchUserStats,
    fetchPluginStats,

    // Plugin Actions
    fetchPlugins,
    createPlugin,
    getPlugin,
    updatePlugin,
    deletePlugin,
    syncGitHubRepos,

    // User Actions
    fetchUsers,
    getUser,
    updateUser,
    deleteUser,

    // Order Actions
    fetchOrders,
    getOrder,
    updateOrderStatus,
    refundOrder,

    // License Actions
    fetchLicenses,
    getLicense,
    revokeLicense,
    extendLicense,

    // Tutorial Actions
    fetchTutorials,
    createTutorial,
    getTutorial,
    updateTutorial,
    deleteTutorial,

    // Settings Actions
    fetchSettings,
    updateSettings,
    
    // GitHub Actions
    fetchGitHubRepos
  }
})
