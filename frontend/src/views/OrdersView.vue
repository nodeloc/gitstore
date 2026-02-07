<template>
  <div class="bg-base-200/30">
    <!-- Header -->
    <div class="bg-gradient-to-br from-secondary/5 via-base-100 to-accent/5 border-b border-base-300">
      <div class="container mx-auto px-4 py-12">
        <h1 class="text-3xl md:text-4xl font-bold bg-gradient-to-r from-secondary to-accent bg-clip-text text-transparent mb-2">
          {{ $t('orders.title') }}
        </h1>
        <p class="text-base-content/60">{{ $t('orders.subtitle') }}</p>
      </div>
    </div>

    <div class="container mx-auto px-4 py-8">
      <!-- Orders Grid -->
      <div v-if="orders.length > 0" class="space-y-4">
        <div 
          v-for="order in orders" 
          :key="order.id"
          class="card bg-base-100 shadow-md hover:shadow-xl transition-all border border-base-300"
        >
          <div class="card-body">
            <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
              <!-- Left: Order Info -->
              <div class="flex-1">
                <div class="flex items-start gap-4">
                  <!-- Plugin Icon -->
                  <div class="w-16 h-16 rounded-xl bg-gradient-to-br from-secondary to-accent flex items-center justify-center text-white text-2xl font-bold flex-shrink-0">
                    {{ order.plugin?.name?.charAt(0) || 'P' }}
                  </div>
                  
                  <!-- Order Details -->
                  <div class="flex-1 min-w-0">
                    <h3 class="text-lg font-bold text-base-content mb-1">
                      {{ order.plugin?.name || 'N/A' }}
                    </h3>
                    <div class="flex flex-wrap gap-3 text-sm text-base-content/60">
                      <div class="flex items-center gap-1">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
                        </svg>
                        <code class="text-xs">{{ order.order_number }}</code>
                      </div>
                      <div class="flex items-center gap-1">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                        {{ formatDate(order.created_at) }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Center: Amount & Status -->
              <div class="flex items-center gap-6">
                <div class="text-right">
                  <p class="text-sm text-base-content/60 mb-1">{{ $t('orders.amount') }}</p>
                  <p class="text-2xl font-bold text-base-content">
                    {{ order.currency }} {{ order.amount }}
                  </p>
                </div>
                
                <div>
                  <span 
                    :class="[
                      'px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full',
                      order.payment_status === 'paid' ? 'bg-green-100 text-green-800' :
                      order.payment_status === 'failed' ? 'bg-red-100 text-red-800' :
                      order.payment_status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                      'bg-gray-100 text-gray-800'
                    ]"
                  >
                    {{ $t(`orders.${order.payment_status}`) }}
                  </span>
                </div>
              </div>

              <!-- Right: Actions -->
              <div class="flex gap-2 lg:flex-col">
                <button 
                  v-if="order.payment_status === 'pending'" 
                  @click="retryPayment(order)"
                  class="btn btn-sm btn-primary gap-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
                  </svg>
                  {{ $t('orders.payNow') }}
                </button>
                <router-link
                  v-if="order.payment_status === 'paid' && order.license?.id"
                  :to="`/licenses/${order.license.id}`"
                  class="btn btn-sm btn-success gap-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  {{ $t('orders.viewLicense') }}
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-else class="text-center py-16">
        <div class="w-24 h-24 mx-auto mb-6 rounded-full bg-base-300 flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
          </svg>
        </div>
        <h3 class="text-xl font-bold mb-2">{{ $t('orders.noOrders') }}</h3>
        <p class="text-base-content/60 mb-6">{{ $t('orders.noOrdersDesc') }}</p>
        <router-link to="/plugins" class="btn btn-primary gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          {{ $t('orders.browsePlugins') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'

const router = useRouter()
const orders = ref([])

onMounted(async () => {
  try {
    const response = await api.get('/user/orders')
    orders.value = response.data.orders || []
  } catch (error) {
    console.error('Failed to load orders:', error)
  }
})

const formatDate = (date) => {
  if (!date) return 'N/A'
  try {
    return new Date(date).toLocaleDateString()
  } catch (e) {
    return 'Invalid Date'
  }
}

const getStatusClass = (status) => {
  return {
    'badge badge-success': status === 'paid',
    'badge badge-error': status === 'failed',
    'badge badge-warning': status === 'pending'
  }
}

const retryPayment = async (order) => {
  try {
    // 重新发起支付
    const response = await api.post('/payments/alipay/create', {
      order_id: order.id
    })
    
    if (response.data.pay_url) {
      // 跳转到支付页面
      window.location.href = response.data.pay_url
    } else if (response.data.qrcode) {
      // 显示二维码
      toast.info('Please scan the QR code: ' + response.data.qrcode)
    }
  } catch (error) {
    console.error('Failed to retry payment:', error)
    toast.error('Failed to initiate payment')
  }
}
</script>
