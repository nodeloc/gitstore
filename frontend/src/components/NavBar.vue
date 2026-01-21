<template>
  <div class="navbar bg-base-100 shadow-lg sticky top-0 z-50">
    <div class="container mx-auto">
      <div class="flex-1">
        <RouterLink to="/" class="btn btn-ghost normal-case text-xl">
          🔌 GitStore
        </RouterLink>

        <!-- Desktop Menu -->
        <div class="hidden lg:flex ml-8">
          <RouterLink to="/" class="btn btn-ghost">
            {{ $t('nav.home') }}
          </RouterLink>
          <RouterLink to="/plugins" class="btn btn-ghost">
            {{ $t('nav.plugins') }}
          </RouterLink>
        </div>
      </div>

      <div class="flex-none gap-2">
        <!-- Language Switcher -->
        <div class="dropdown dropdown-end">
          <label tabindex="0" class="btn btn-ghost btn-circle">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-5 h-5 stroke-current">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129" />
            </svg>
          </label>
          <ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-32 mt-4">
            <li><a @click="changeLanguage('en')">English</a></li>
            <li><a @click="changeLanguage('zh')">中文</a></li>
          </ul>
        </div>

        <!-- Theme Switcher -->
        <div class="dropdown dropdown-end">
          <label tabindex="0" class="btn btn-ghost btn-circle">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-5 h-5 stroke-current">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01" />
            </svg>
          </label>
          <ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-32 mt-4">
            <li><a @click="changeTheme('light')">Light</a></li>
            <li><a @click="changeTheme('dark')">Dark</a></li>
            <li><a @click="changeTheme('cupcake')">Cupcake</a></li>
            <li><a @click="changeTheme('corporate')">Corporate</a></li>
          </ul>
        </div>

        <!-- User Menu -->
        <div v-if="authStore.isAuthenticated" class="dropdown dropdown-end">
          <label tabindex="0" class="btn btn-ghost btn-circle avatar">
            <div class="w-10 rounded-full">
              <img :src="authStore.user?.avatar_url || '/default-avatar.png'" :alt="authStore.user?.name" />
            </div>
          </label>
          <ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52 mt-4">
            <li class="menu-title">
              <span>{{ authStore.user?.name || authStore.user?.email }}</span>
            </li>
            <li><RouterLink to="/dashboard">{{ $t('nav.dashboard') }}</RouterLink></li>
            <li><RouterLink to="/licenses">{{ $t('nav.myLicenses') }}</RouterLink></li>
            <li v-if="authStore.isAdmin"><RouterLink to="/admin">{{ $t('nav.admin') }}</RouterLink></li>
            <li><a @click="logout">{{ $t('nav.logout') }}</a></li>
          </ul>
        </div>

        <button v-else @click="authStore.login" class="btn btn-primary">
          {{ $t('nav.login') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const { locale } = useI18n()
const authStore = useAuthStore()
const router = useRouter()

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
