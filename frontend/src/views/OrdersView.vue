<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">{{ $t('orders.title') }}</h1>
    
    <div class="overflow-x-auto">
      <table class="table w-full">
        <thead>
          <tr>
            <th>{{ $t('orders.id') }}</th>
            <th>{{ $t('orders.plugin') }}</th>
            <th>{{ $t('orders.amount') }}</th>
            <th>{{ $t('orders.status') }}</th>
            <th>{{ $t('orders.date') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in orders" :key="order.id">
            <td>{{ order.id }}</td>
            <td>{{ order.plugin?.name }}</td>
            <td>${{ order.amount }}</td>
            <td>
              <span :class="getStatusClass(order.status)">
                {{ order.status }}
              </span>
            </td>
            <td>{{ formatDate(order.created_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <div v-if="orders.length === 0" class="text-center py-8">
      <p class="text-lg">{{ $t('orders.noOrders') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const orders = ref([])

onMounted(async () => {
  try {
    const response = await api.get('/api/user/orders')
    orders.value = response.data || []
  } catch (error) {
    console.error('Failed to load orders:', error)
  }
})

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

const getStatusClass = (status) => {
  return {
    'badge badge-success': status === 'completed',
    'badge badge-error': status === 'failed',
    'badge badge-warning': status === 'pending'
  }
}
</script>
