<template>
  <div class="container mx-auto px-4 py-8">
    <div v-if="loading" class="flex justify-center py-16">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
    
    <div v-else-if="plugin" class="max-w-4xl mx-auto">
      <h1 class="text-4xl font-bold mb-4">{{ plugin.name }}</h1>
      <p class="text-xl mb-8">{{ plugin.description }}</p>
      
      <!-- Long Description -->
      <div v-if="plugin.long_description" class="card bg-base-100 shadow-xl mb-8">
        <div class="card-body">
          <h2 class="card-title">{{ $t('plugin.description') }}</h2>
          <div class="prose max-w-none">
            <p class="whitespace-pre-wrap">{{ plugin.long_description }}</p>
          </div>
        </div>
      </div>
      
      <div class="card bg-base-100 shadow-xl mb-8">
        <div class="card-body">
          <h2 class="card-title">{{ $t('plugin.details') }}</h2>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <strong>{{ $t('plugin.price') }}:</strong> ${{ plugin.price }}
            </div>
            <div>
              <strong>{{ $t('plugin.version') }}:</strong> {{ plugin.version || 'N/A' }}
            </div>
            <div v-if="plugin.category">
              <strong>{{ $t('plugin.category') }}:</strong> {{ plugin.category }}
            </div>
            <div>
              <strong>{{ $t('plugin.downloads') }}:</strong> {{ plugin.download_count || 0 }}
            </div>
          </div>
          <div class="card-actions justify-end mt-4">
            <button @click="purchase" class="btn btn-primary">
              {{ $t('plugin.purchase') }}
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="text-center py-16">
      <p class="text-lg">{{ $t('plugin.notFound') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const plugin = ref(null)
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const response = await api.get(`/plugins/${route.params.slug}`)
    plugin.value = response.data.plugin
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
