<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">{{ $t('plugins.title') }}</h1>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="plugin in plugins" :key="plugin.id" class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">{{ plugin.name }}</h2>
          <p>{{ plugin.description }}</p>
          <div class="card-actions justify-end">
            <router-link :to="`/plugins/${plugin.slug}`" class="btn btn-primary btn-sm">
              {{ $t('plugins.viewDetails') }}
            </router-link>
          </div>
        </div>
      </div>
    </div>
    
    <div v-if="loading" class="flex justify-center py-8">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
    
    <div v-if="!loading && plugins.length === 0" class="text-center py-8">
      <p class="text-lg">{{ $t('plugins.noPlugins') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const plugins = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const response = await api.get('/plugins')
    plugins.value = response.data.plugins || []
  } catch (error) {
    console.error('Failed to load plugins:', error)
  } finally {
    loading.value = false
  }
})
</script>
