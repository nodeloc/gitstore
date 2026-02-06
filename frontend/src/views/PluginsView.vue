<template>
  <div class="min-h-screen bg-base-200/30">
    <!-- Header Section with Search -->
    <div class="bg-gradient-to-br from-primary/5 via-base-100 to-secondary/5 border-b border-base-300">
      <div class="container mx-auto px-4 py-8">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div>
            <h1 class="text-3xl md:text-4xl font-bold mb-2 bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
              {{ currentCategory ? currentCategory.name : $t('plugins.title') }}
            </h1>
            <p class="text-base text-base-content/70">
              {{ currentCategory ? currentCategory.description : $t('plugins.subtitle') }}
            </p>
          </div>
          
          <!-- Search Box -->
          <div class="w-full md:w-96">
            <div class="relative">
              <input 
                v-model="searchQuery" 
                type="text" 
                :placeholder="$t('plugins.search')"
                class="input input-bordered w-full pl-12 bg-base-100 shadow-sm focus:shadow-md transition-shadow"
              />
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 absolute left-4 top-1/2 -translate-y-1/2 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Plugins Grid -->
    <div class="container mx-auto px-4 py-12">
      <!-- Loading State -->
      <div v-if="loading" class="flex flex-col items-center justify-center py-20">
        <span class="loading loading-spinner loading-lg text-primary"></span>
        <p class="mt-4 text-base-content/60">{{ $t('plugins.loading') }}</p>
      </div>
      
      <!-- Empty State -->
      <div v-else-if="filteredPlugins.length === 0" class="text-center py-20">
        <div class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-base-200 mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-base-content mb-2">{{ $t('plugins.noPlugins') }}</h3>
        <p class="text-base-content/60">{{ $t('plugins.noPluginsDesc') }}</p>
      </div>
      
      <!-- Plugins Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="plugin in filteredPlugins" 
             :key="plugin.id" 
             class="card bg-base-100 shadow-md hover:shadow-xl transition-all duration-300 border border-base-300 group">
          <div class="card-body">
            <!-- Plugin Icon & Title -->
            <div class="flex items-start gap-4 mb-3">
              <div class="w-14 h-14 rounded-xl bg-gradient-to-br from-primary/10 to-secondary/10 flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <h2 class="card-title text-lg mb-1 truncate">{{ plugin.name }}</h2>
                <div class="flex items-center gap-2 text-sm text-base-content/60">
                  <span v-if="plugin.category" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-primary/10 text-primary">{{ plugin.category }}</span>
                  <span v-if="plugin.version">v{{ plugin.version }}</span>
                </div>
              </div>
            </div>
            
            <!-- Description -->
            <p class="text-base-content/70 line-clamp-3 mb-4">
              {{ plugin.description || $t('plugins.noDescription') }}
            </p>
            
            <!-- Footer -->
            <div class="card-actions justify-between items-center pt-4 border-t border-base-300">
              <div class="text-2xl font-bold text-primary">
                ${{ plugin.price }}
              </div>
              <router-link 
                :to="`/plugins/${plugin.slug}`" 
                class="btn btn-primary btn-sm gap-2 shadow-sm hover:shadow-md transition-all">
                {{ $t('plugins.viewDetails') }}
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/utils/api'

const route = useRoute()
const plugins = ref([])
const categories = ref([])
const loading = ref(false)
const searchQuery = ref('')
const categoryFilter = ref('')

// Watch for URL query parameter changes
watch(() => route.query.category, (newCategory) => {
  categoryFilter.value = newCategory || ''
}, { immediate: true })

// Get current category info
const currentCategory = computed(() => {
  if (!categoryFilter.value || categories.value.length === 0) return null
  
  const categoryLower = categoryFilter.value.toLowerCase()
  return categories.value.find(cat => cat.slug?.toLowerCase() === categoryLower)
})

const filteredPlugins = computed(() => {
  let result = plugins.value
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(p => 
      p.name?.toLowerCase().includes(query) || 
      p.description?.toLowerCase().includes(query)
    )
  }
  
  // Filter by category (case-insensitive)
  if (categoryFilter.value) {
    const categoryLower = categoryFilter.value.toLowerCase()
    result = result.filter(p => {
      // Support both category slug string and category object
      if (typeof p.category === 'string') {
        return p.category.toLowerCase() === categoryLower
      } else if (p.category?.slug) {
        return p.category.slug.toLowerCase() === categoryLower
      } else if (p.category_slug) {
        return p.category_slug.toLowerCase() === categoryLower
      }
      return false
    })
  }
  
  return result
})

onMounted(async () => {
  loading.value = true
  try {
    const [pluginsRes, categoriesRes] = await Promise.all([
      api.get('/plugins'),
      api.get('/categories')
    ])
    plugins.value = pluginsRes.data.plugins || []
    categories.value = categoriesRes.data.categories || []
  } catch (error) {
    console.error('Failed to load data:', error)
  } finally {
    loading.value = false
  }
})
</script>