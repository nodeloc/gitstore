<template>
  <div class="bg-base-200/30">
    <div v-if="loading" class="flex justify-center py-16">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
    
    <div v-else-if="plugin">
      <!-- Header Section -->
      <div class="bg-gradient-to-br from-primary/5 via-base-100 to-secondary/5 border-b border-base-300">
        <div class="container mx-auto px-4 py-8">
          <button @click="$router.push('/plugins')" class="btn btn-ghost btn-sm mb-4">
            ← {{ $t('pluginDetail.backToList') }}
          </button>
          <h1 class="text-4xl font-bold mb-3">{{ plugin.name }}</h1>
          <p class="text-xl text-base-content/70 mb-4">{{ plugin.description }}</p>
          
          <!-- Tags/Badges -->
          <div class="flex gap-2 flex-wrap">
            <span v-if="plugin.category" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-primary/10 text-primary">{{ plugin.category }}</span>
            <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800">v{{ plugin.version || '1.0.0' }}</span>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="container mx-auto px-4 py-8">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Main Content (Left) -->
          <div class="lg:col-span-2 space-y-6">
            <!-- Long Description -->
            <div v-if="plugin.long_description" class="card bg-base-100 shadow-xl">
            <div class="card-body">
              <h2 class="card-title text-2xl mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ $t('pluginDetail.description') }}
              </h2>
              <MdPreview :modelValue="plugin.long_description" class="!bg-transparent" />
            </div>
          </div>
        </div>

        <!-- Sidebar (Right) -->
        <div class="lg:col-span-1">
          <!-- Purchase Card -->
          <div class="card bg-base-100 shadow-xl sticky top-20">
            <div class="card-body">
              <div class="text-center mb-4">
                <div class="text-4xl font-bold text-primary">${{ plugin.price }}</div>
                <div class="text-sm text-base-content/60 mt-1">{{ $t('pluginDetail.oneTimePurchase') }}</div>
              </div>

              <div class="divider"></div>

              <!-- Details -->
              <div class="space-y-3">
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <span class="text-sm">{{ $t('pluginDetail.permanentLicense') }}</span>
                </div>
                
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <span class="text-sm">{{ $t('pluginDetail.maintenanceIncluded', { months: plugin.default_maintenance_months || 12 }) }}</span>
                </div>
                
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <span class="text-sm">{{ $t('pluginDetail.githubAccess') }}</span>
                </div>
                
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <span class="text-sm">{{ $t('pluginDetail.instantDelivery') }}</span>
                </div>
              </div>

              <div class="divider"></div>

              <button @click="purchase" class="btn btn-primary btn-block btn-lg">
                {{ $t('pluginDetail.purchase') }}
              </button>

              <!-- Additional Info -->
              <div class="mt-4 text-xs text-base-content/60 text-center">
                {{ $t('pluginDetail.securePayment') }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    </div>
    
    <div v-else class="container mx-auto px-4 py-16 text-center">
      <p class="text-lg">{{ $t('plugin.notFound') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import api from '@/utils/api'
import { getPageSEO, updatePageSEO } from '@/utils/seo'

const route = useRoute()
const router = useRouter()
const plugin = ref(null)
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const response = await api.get(`/plugins/${route.params.slug}`)
    plugin.value = response.data.plugin
    
    // Update SEO with plugin data
    if (plugin.value) {
      // Get site settings
      let siteSettings = {}
      try {
        const settingsResponse = await api.get('/settings/public')
        if (settingsResponse.data.settings) {
          siteSettings = settingsResponse.data.settings.reduce((acc, setting) => {
            acc[setting.key] = setting.value
            return acc
          }, {})
        }
      } catch (err) {
        console.log('Failed to load site settings for SEO')
      }
      
      const seoData = getPageSEO('plugin-detail', {
        name: plugin.value.name,
        slug: plugin.value.slug,
        short_description: plugin.value.short_description,
        description: plugin.value.description,
        category: plugin.value.category,
        icon_url: plugin.value.icon_url,
        image_url: plugin.value.image_url
      }, siteSettings)
      
      updatePageSEO(seoData)
    }
  } catch (error) {
    console.error('Failed to load plugin:', error)
  } finally {
    loading.value = false
  }
})

const purchase = () => {
  router.push(`/purchase/${plugin.value.id}`)
}
</script>
