import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { toast } from '@/utils/toast'
import i18n from '@/i18n'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue')
    },
    {
      path: '/plugins',
      name: 'plugins',
      component: () => import('@/views/PluginsView.vue')
    },
    {
      path: '/plugins/:slug',
      name: 'plugin-detail',
      component: () => import('@/views/PluginDetailView.vue')
    },
    {
      path: '/auth/callback',
      name: 'auth-callback',
      component: () => import('@/views/AuthCallbackView.vue')
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/licenses',
      name: 'licenses',
      component: () => import('@/views/LicensesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/licenses/:id',
      name: 'license-detail',
      component: () => import('@/views/LicenseDetailView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/orders',
      name: 'orders',
      component: () => import('@/views/OrdersView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/purchase/:pluginId',
      name: 'purchase',
      component: () => import('@/views/PurchaseView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/admin/AdminDashboard.vue'),
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/admin/plugins/create',
      name: 'admin-plugin-create',
      component: () => import('@/views/admin/PluginForm.vue'),
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/admin/plugins/:id/edit',
      name: 'admin-plugin-edit',
      component: () => import('@/views/admin/PluginForm.vue'),
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/admin/pages',
      name: 'admin-pages',
      component: () => import('@/views/admin/PageManagement.vue'),
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/pages/:slug',
      name: 'page',
      component: () => import('@/views/PageView.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue')
    }
  ]
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  console.log('🔍 Route guard check:', {
    path: to.path,
    requiresAuth: to.meta.requiresAuth,
    requiresAdmin: to.meta.requiresAdmin,
    isAuthenticated: authStore.isAuthenticated,
    isAdmin: authStore.isAdmin,
    hasToken: !!authStore.token,
    hasUser: !!authStore.user,
    userRole: authStore.user?.role
  })

  // Check if route requires authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    console.log('❌ Auth required but not authenticated, redirecting to home')
    toast.error(i18n.global.t('auth.loginRequired'))
    next({ name: 'home' })
    return
  }

  // Check if route requires admin
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    console.log('❌ Admin required but not admin, redirecting to dashboard')
    toast.error(i18n.global.t('auth.adminRequired'))
    next({ name: 'dashboard' })
    return
  }

  console.log('✅ Route guard passed')
  next()
})

export default router
