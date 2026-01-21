<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center">
      <span class="loading loading-spinner loading-lg"></span>
      <p class="mt-4">{{ $t('auth.processing') }}</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

onMounted(async () => {
  try {
    await authStore.handleCallback()
    router.push('/dashboard')
  } catch (error) {
    console.error('Auth callback failed:', error)
    router.push('/')
  }
})
</script>
