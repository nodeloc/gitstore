<template>
  <div class="min-h-screen bg-base-200/30">
    <!-- Header -->
    <div class="bg-gradient-to-br from-primary/5 via-base-100 to-secondary/5 border-b border-base-300">
      <div class="container mx-auto px-4 py-12">
        <h1 class="text-3xl md:text-4xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent mb-2">
          {{ $t('licenses.title') }}
        </h1>
        <p class="text-base-content/60">{{ $t('licenses.subtitle') }}</p>
      </div>
    </div>

    <div class="container mx-auto px-4 py-8">
      <!-- Licenses Grid -->
      <div v-if="licenses.length > 0" class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div 
          v-for="license in licenses" 
          :key="license.id"
          class="card bg-base-100 shadow-md hover:shadow-xl transition-all border border-base-300"
        >
          <div class="card-body">
            <!-- Plugin Info -->
            <div class="flex items-start gap-4 mb-4">
              <div class="w-16 h-16 rounded-xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white text-2xl font-bold flex-shrink-0">
                {{ license.plugin?.name?.charAt(0) || 'P' }}
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="text-xl font-bold text-base-content mb-1">
                  {{ license.plugin?.name || 'N/A' }}
                </h3>
                <div class="flex items-center gap-2">
                  <span 
                    class="badge"
                    :class="{
                      'badge-success': license.status === 'active',
                      'badge-error': license.status === 'expired',
                      'badge-warning': license.status === 'pending'
                    }"
                  >
                    {{ license.status }}
                  </span>
                </div>
              </div>
            </div>

            <!-- License Key -->
            <div class="bg-base-200 rounded-lg p-3 mb-3">
              <p class="text-xs text-base-content/60 mb-1">{{ $t('licenses.key') }}</p>
              <code class="text-sm font-mono break-all">{{ license.id }}</code>
            </div>

            <!-- Maintenance Period -->
            <div class="flex items-center gap-2 mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-base-content/60" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <div>
                <p class="text-sm text-base-content/60">{{ $t('licenses.maintenanceUntil') }}</p>
                <p class="text-base font-semibold">{{ formatDate(license.maintenance_until) }}</p>
              </div>
            </div>

            <!-- Actions -->
            <div class="card-actions justify-end">
              <router-link :to="`/licenses/${license.id}`" class="btn btn-primary btn-sm gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                {{ $t('licenses.viewDetails') }}
              </router-link>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-else class="text-center py-16">
        <div class="w-24 h-24 mx-auto mb-6 rounded-full bg-base-300 flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <h3 class="text-xl font-bold mb-2">{{ $t('licenses.noLicenses') }}</h3>
        <p class="text-base-content/60 mb-6">{{ $t('licenses.noLicensesDesc') }}</p>
        <router-link to="/plugins" class="btn btn-primary gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          {{ $t('licenses.browsePlugins') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const licenses = ref([])

onMounted(async () => {
  try {
    const response = await api.get('/user/licenses')
    licenses.value = response.data.licenses || []
  } catch (error) {
    console.error('Failed to load licenses:', error)
  }
})

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

const getStatusClass = (status) => {
  return {
    'badge badge-success': status === 'active',
    'badge badge-error': status === 'expired',
    'badge badge-warning': status === 'pending'
  }
}
</script>
