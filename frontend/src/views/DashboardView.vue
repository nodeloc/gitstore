<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">{{ $t('dashboard.title') }}</h1>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="stats shadow">
        <div class="stat">
          <div class="stat-title">{{ $t('dashboard.totalLicenses') }}</div>
          <div class="stat-value">{{ stats.licenses }}</div>
        </div>
      </div>
      <div class="stats shadow">
        <div class="stat">
          <div class="stat-title">{{ $t('dashboard.totalOrders') }}</div>
          <div class="stat-value">{{ stats.orders }}</div>
        </div>
      </div>
      <div class="stats shadow">
        <div class="stat">
          <div class="stat-title">{{ $t('dashboard.activeLicenses') }}</div>
          <div class="stat-value">{{ stats.active }}</div>
        </div>
      </div>
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <router-link to="/licenses" class="btn btn-primary btn-lg">
        {{ $t('dashboard.viewLicenses') }}
      </router-link>
      <router-link to="/orders" class="btn btn-secondary btn-lg">
        {{ $t('dashboard.viewOrders') }}
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const stats = ref({
  licenses: 0,
  orders: 0,
  active: 0
})

onMounted(() => {
  // Load stats from API
})
</script>
