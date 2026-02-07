<template>
  <div class="min-h-screen bg-base-200/30 flex items-center justify-center py-12 px-4">
    <div class="max-w-md w-full">
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body text-center">
          <!-- Loading State -->
          <div v-if="loading" class="space-y-4">
            <div class="loading loading-spinner loading-lg text-primary mx-auto"></div>
            <h2 class="card-title justify-center">{{ $t('payment.verifying') }}</h2>
            <p class="text-base-content/60">{{ $t('payment.pleaseWait') }}</p>
          </div>

          <!-- Success State -->
          <div v-else-if="success" class="space-y-4">
            <div class="text-success">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h2 class="card-title justify-center text-2xl">{{ $t('payment.success') }}</h2>
            <p class="text-base-content/70">{{ $t('payment.successMessage') }}</p>
            
            <div v-if="orderInfo" class="bg-base-200 rounded-lg p-4 text-left space-y-2">
              <div class="flex justify-between text-sm">
                <span class="text-base-content/60">{{ $t('payment.orderId') }}:</span>
                <span class="font-mono">{{ orderInfo.orderId }}</span>
              </div>
              <div v-if="orderInfo.tradeNo" class="flex justify-between text-sm">
                <span class="text-base-content/60">{{ $t('payment.tradeNo') }}:</span>
                <span class="font-mono">{{ orderInfo.tradeNo }}</span>
              </div>
              <div v-if="orderInfo.amount" class="flex justify-between text-sm">
                <span class="text-base-content/60">{{ $t('payment.amount') }}:</span>
                <span class="font-semibold">${{ orderInfo.amount }}</span>
              </div>
            </div>

            <div class="card-actions justify-center mt-6">
              <button @click="goToOrders" class="btn btn-primary">
                {{ $t('payment.viewOrders') }}
              </button>
              <button @click="goToHome" class="btn btn-ghost">
                {{ $t('payment.backToHome') }}
              </button>
            </div>
          </div>

          <!-- Error State -->
          <div v-else-if="error" class="space-y-4">
            <div class="text-error">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h2 class="card-title justify-center text-2xl">{{ $t('payment.failed') }}</h2>
            <p class="text-base-content/70">{{ errorMessage || $t('payment.errorMessage') }}</p>
            
            <div class="card-actions justify-center mt-6">
              <button @click="goToOrders" class="btn btn-primary">
                {{ $t('payment.viewOrders') }}
              </button>
              <button @click="goToHome" class="btn btn-ghost">
                {{ $t('payment.backToHome') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from '@/utils/toast'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const success = ref(false)
const error = ref(false)
const errorMessage = ref('')
const orderInfo = ref(null)

onMounted(async () => {
  try {
    // 从 URL 参数获取支付信息
    const params = route.query
    
    // 检查支付状态
    if (params.trade_status === 'TRADE_SUCCESS') {
      success.value = true
      orderInfo.value = {
        orderId: params.out_trade_no,
        tradeNo: params.trade_no || params.api_trade_no,
        amount: params.money,
        type: params.type
      }
      
      toast.success('支付成功！')
      
      // 延迟3秒后自动跳转到订单页面
      setTimeout(() => {
        router.push('/orders')
      }, 3000)
    } else {
      error.value = true
      errorMessage.value = '支付未成功，请检查订单状态'
    }
  } catch (err) {
    console.error('Payment verification error:', err)
    error.value = true
    errorMessage.value = err.message || '验证支付状态时出错'
  } finally {
    loading.value = false
  }
})

const goToOrders = () => {
  router.push('/orders')
}

const goToHome = () => {
  router.push('/')
}
</script>
