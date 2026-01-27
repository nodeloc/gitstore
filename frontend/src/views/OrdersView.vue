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
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in orders" :key="order.id">
            <td>{{ order.order_number }}</td>
            <td>{{ order.plugin?.name || 'N/A' }}</td>
            <td>{{ order.currency }} {{ order.amount }}</td>
            <td>
              <span :class="getStatusClass(order.payment_status)">
                {{ order.payment_status }}
              </span>
            </td>
            <td>{{ formatDate(order.created_at) }}</td>
            <td>
              <button 
                v-if="order.payment_status === 'pending'" 
                @click="retryPayment(order)"
                class="btn btn-sm btn-primary"
              >
                Pay Now
              </button>
            </td>
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
      alert('Please scan the QR code: ' + response.data.qrcode)
    }
  } catch (error) {
    console.error('Failed to retry payment:', error)
    alert('Failed to initiate payment')
  }
}
</script>
