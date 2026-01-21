<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">{{ $t('licenses.title') }}</h1>
    
    <div class="overflow-x-auto">
      <table class="table w-full">
        <thead>
          <tr>
            <th>{{ $t('licenses.plugin') }}</th>
            <th>{{ $t('licenses.key') }}</th>
            <th>{{ $t('licenses.expiresAt') }}</th>
            <th>{{ $t('licenses.status') }}</th>
            <th>{{ $t('licenses.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="license in licenses" :key="license.id">
            <td>{{ license.plugin?.name }}</td>
            <td><code class="text-xs">{{ license.key }}</code></td>
            <td>{{ formatDate(license.expires_at) }}</td>
            <td>
              <span :class="getStatusClass(license.status)">
                {{ license.status }}
              </span>
            </td>
            <td>
              <button class="btn btn-sm btn-primary">
                {{ $t('licenses.renew') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <div v-if="licenses.length === 0" class="text-center py-8">
      <p class="text-lg">{{ $t('licenses.noLicenses') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const licenses = ref([])

onMounted(async () => {
  try {
    const response = await api.get('/api/user/licenses')
    licenses.value = response.data || []
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
