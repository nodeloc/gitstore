<template>
  <div id="app" class="min-h-screen flex flex-col">
    <NavBar />
    <main class="flex-grow" :class="{ 'bg-base-200': $route.path.startsWith('/admin') }">
      <RouterView v-slot="{ Component, route }">
        <Transition name="fade" mode="out-in">
          <component :is="Component" :key="route.path" />
        </Transition>
      </RouterView>
    </main>
    <Footer />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavBar from '@/components/NavBar.vue'
import Footer from '@/components/Footer.vue'
import api from '@/utils/api'

const authStore = useAuthStore()

onMounted(async () => {
  // Fetch user data if token exists
  if (authStore.token) {
    authStore.fetchUser()
  }

  // Load site settings and update page title and favicon
  try {
    const response = await api.get('/settings/public')
    if (response.data.settings) {
      const siteNameSetting = response.data.settings.find(s => s.key === 'site_name')
      const siteSubtitleSetting = response.data.settings.find(s => s.key === 'site_subtitle')
      const logoUrlSetting = response.data.settings.find(s => s.key === 'logo_url')
      
      // Update title
      if (siteNameSetting && siteSubtitleSetting) {
        document.title = `${siteNameSetting.value} - ${siteSubtitleSetting.value}`
      } else if (siteNameSetting) {
        document.title = siteNameSetting.value
      }
      
      // Update favicon
      if (logoUrlSetting && logoUrlSetting.value) {
        let link = document.querySelector("link[rel~='icon']")
        if (!link) {
          link = document.createElement('link')
          link.rel = 'icon'
          document.head.appendChild(link)
        }
        link.href = logoUrlSetting.value
      }
    }
  } catch (error) {
    console.log('Failed to load site settings for title and favicon')
  }
})
</script>

<style scoped>
#app {
  min-height: 100vh;
}
</style>
