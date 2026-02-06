<template>
  <nav class="navbar bg-base-100/80 backdrop-blur-lg shadow-sm border-b border-base-300 sticky top-0 z-50">
    <div class="container mx-auto px-4">
      <div class="flex-1">
        <div class="flex items-center gap-8">
          <!-- Logo -->
          <RouterLink to="/" class="flex items-center gap-3 hover:opacity-80 transition-opacity">
            <div v-if="logoUrl" class="w-10 h-10 flex items-center justify-center overflow-hidden">
              <img :src="logoUrl" :alt="siteName" class="w-full h-full object-contain" />
            </div>
            <div v-else class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center shadow-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <div class="hidden sm:block">
              <div class="font-bold text-xl bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                {{ siteName }}
              </div>
            </div>
          </RouterLink>

          <!-- Desktop Menu -->
          <div class="hidden lg:flex gap-1">
            <RouterLink to="/" 
                        class="px-4 py-2 rounded-lg hover:bg-base-200 transition-colors font-medium"
                        :class="{ 'bg-base-200 text-primary': $route.path === '/' }">
              {{ $t('nav.home') }}
            </RouterLink>
            
            <!-- Category Links -->
            <RouterLink 
              v-for="cat in categories" 
              :key="cat.slug"
              :to="`/plugins?category=${cat.slug}`" 
              class="px-4 py-2 rounded-lg hover:bg-base-200 transition-colors font-medium"
              :class="{ 'bg-base-200 text-primary': $route.query.category === cat.slug }">
              {{ cat.name }}
            </RouterLink>
          </div>
        </div>
      </div>

      <div class="flex-none">
        <div class="flex items-center gap-2">
          <!-- Language Switcher -->
          <div class="dropdown dropdown-end">
            <label tabindex="0" class="btn btn-ghost btn-sm gap-2 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129" />
              </svg>
              <span class="hidden md:inline font-medium">{{ locale.toUpperCase() }}</span>
            </label>
            <ul tabindex="0" class="dropdown-content menu p-2 shadow-lg bg-base-100 rounded-box w-40 mt-3 border border-base-300">
              <li><a @click="changeLanguage('en')" class="gap-2">
                <span class="text-lg">🇺🇸</span>
                <span>English</span>
              </a></li>
              <li><a @click="changeLanguage('zh')" class="gap-2">
                <span class="text-lg">🇨🇳</span>
                <span>中文</span>
              </a></li>
            </ul>
          </div>

          <!-- Theme Switcher -->
          <div class="dropdown dropdown-end">
            <label tabindex="0" class="btn btn-ghost btn-sm btn-circle rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
              </svg>
            </label>
            <ul tabindex="0" class="dropdown-content menu p-2 shadow-lg bg-base-100 rounded-box w-40 mt-3 border border-base-300">
              <li><a @click="changeTheme('light')" class="gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
                Light
              </a></li>
              <li><a @click="changeTheme('dark')" class="gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                </svg>
                Dark
              </a></li>
              <li><a @click="changeTheme('cupcake')" class="gap-2">
                <span>🧁</span>
                Cupcake
              </a></li>
              <li><a @click="changeTheme('corporate')" class="gap-2">
                <span>💼</span>
                Corporate
              </a></li>
            </ul>
          </div>

          <!-- Divider -->
          <div class="divider divider-horizontal mx-0"></div>

          <!-- User Menu -->
          <div v-if="authStore.isAuthenticated" class="dropdown dropdown-end">
            <label tabindex="0" class="btn btn-ghost btn-sm gap-2 rounded-lg">
              <div class="avatar">
                <div class="w-8 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                  <img :src="authStore.user?.avatar_url || '/default-avatar.png'" :alt="authStore.user?.name" />
                </div>
              </div>
              <span class="hidden md:inline font-medium max-w-[100px] truncate">
                {{ authStore.user?.name || authStore.user?.email }}
              </span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 opacity-60" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </label>
            <ul tabindex="0" class="dropdown-content menu p-2 shadow-lg bg-base-100 rounded-box w-60 mt-3 border border-base-300">
              <li class="menu-title px-4 py-2">
                <div class="flex items-center gap-2">
                  <div class="avatar">
                    <div class="w-10 rounded-full">
                      <img :src="authStore.user?.avatar_url || '/default-avatar.png'" :alt="authStore.user?.name" />
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="font-semibold truncate">{{ authStore.user?.name || 'User' }}</div>
                    <div class="text-xs opacity-60 truncate">{{ authStore.user?.email }}</div>
                  </div>
                </div>
              </li>
              <div class="divider my-1"></div>
              <li>
                <RouterLink to="/dashboard" class="gap-3">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                  </svg>
                  {{ $t('nav.dashboard') }}
                </RouterLink>
              </li>
              <li>
                <RouterLink to="/licenses" class="gap-3">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  {{ $t('nav.myLicenses') }}
                </RouterLink>
              </li>
              <li v-if="authStore.isAdmin">
                <RouterLink to="/admin" class="gap-3">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  {{ $t('nav.admin') }}
                </RouterLink>
              </li>
              <div class="divider my-1"></div>
              <li>
                <a @click="logout" class="gap-3 text-error">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
                  {{ $t('nav.logout') }}
                </a>
              </li>
            </ul>
          </div>

          <!-- Login Button -->
          <button v-else @click="authStore.login" class="btn btn-primary btn-sm rounded-lg gap-2 shadow-md hover:shadow-lg transition-all">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            {{ $t('nav.login') }}
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import api from '@/utils/api'

const { locale } = useI18n()
const authStore = useAuthStore()
const router = useRouter()
const categories = ref([])
const siteName = ref('GitStore')
const siteSubtitle = ref('Plugin Marketplace')
const logoUrl = ref('')

onMounted(async () => {
  try {
    const response = await api.get('/categories')
    categories.value = response.data.categories || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }

  // Load site settings
  try {
    const response = await api.get('/settings/public')
    if (response.data.settings) {
      const siteNameSetting = response.data.settings.find(s => s.key === 'site_name')
      const siteSubtitleSetting = response.data.settings.find(s => s.key === 'site_subtitle')
      const logoUrlSetting = response.data.settings.find(s => s.key === 'logo_url')
      if (siteNameSetting) siteName.value = siteNameSetting.value
      if (siteSubtitleSetting) siteSubtitle.value = siteSubtitleSetting.value
      if (logoUrlSetting && logoUrlSetting.value) logoUrl.value = logoUrlSetting.value
    }
  } catch (error) {
    console.log('Failed to load site settings, using defaults')
  }
})

const changeLanguage = (lang) => {
  locale.value = lang
  localStorage.setItem('language', lang)
}

const changeTheme = (theme) => {
  document.documentElement.setAttribute('data-theme', theme)
  localStorage.setItem('theme', theme)
}

const logout = () => {
  authStore.logout()
  router.push('/')
}

// Load theme on mount
const savedTheme = localStorage.getItem('theme') || 'light'
document.documentElement.setAttribute('data-theme', savedTheme)
</script>
