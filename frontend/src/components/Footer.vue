<template>
  <footer class="footer footer-center p-10 bg-base-300 text-base-content mt-auto">
    <div class="grid grid-flow-col gap-4">
      <router-link to="/pages/about" class="link link-hover">{{ $t('footer.about') }}</router-link>
      <router-link to="/pages/contact" class="link link-hover">{{ $t('footer.contact') }}</router-link>
      <router-link to="/pages/privacy-policy" class="link link-hover">{{ $t('footer.privacy') }}</router-link>
      <router-link to="/pages/terms-of-service" class="link link-hover">{{ $t('footer.terms') }}</router-link>
    </div>
    <div>
      <p class="font-bold">
        {{ siteName }}
      </p>
      <p v-if="siteSubtitle" class="text-sm text-base-content/70">{{ siteSubtitle }}</p>
    </div>
    <div>
      <div class="flex flex-col sm:flex-row items-center gap-2 sm:gap-4">
        <p>{{ $t('footer.copyright') }}</p>
        <p class="text-sm text-base-content/60">
          Powered by 
          <a href="https://github.com/nodeloc/gitstore" target="_blank" class="link link-hover text-primary">GitStore</a>
          from 
          <a href="https://www.nodeloc.com" target="_blank" class="link link-hover text-primary">Nodeloc</a>.
        </p>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const siteName = ref('GitStore')
const siteSubtitle = ref('')

onMounted(async () => {
  try {
    const response = await api.get('/settings/public')
    if (response.data.settings) {
      const siteNameSetting = response.data.settings.find(s => s.key === 'site_name')
      const siteSubtitleSetting = response.data.settings.find(s => s.key === 'site_subtitle')
      if (siteNameSetting) {
        siteName.value = siteNameSetting.value
      }
      if (siteSubtitleSetting) {
        siteSubtitle.value = siteSubtitleSetting.value
      }
    }
  } catch (error) {
    console.log('Failed to load site name, using default')
  }
})
</script>
