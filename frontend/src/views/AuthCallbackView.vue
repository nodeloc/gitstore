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
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

onMounted(async () => {
  try {
    const token = route.query.token
    if (token) {
      // Save token and fetch user info
      localStorage.setItem('token', token)
      await authStore.fetchUser()
      router.push('/dashboard')
    } else {
      throw new Error('No token received')
    }
  } catch (error) {
    console.error('Auth callback failed:', error)
    router.push('/')
  }
})
</script>
