<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">{{ $t('purchase.title') }}</h1>
    
    <div class="max-w-2xl mx-auto">
      <!-- Plugin Info -->
      <div v-if="plugin" class="card bg-base-100 shadow-xl mb-6">
        <div class="card-body">
          <h2 class="card-title">{{ plugin.name }}</h2>
          <p>{{ plugin.description }}</p>
          <div class="divider"></div>
          <div class="flex justify-between items-center">
            <span class="text-lg font-semibold">{{ $t('purchase.total') }}</span>
            <span class="text-2xl font-bold">${{ plugin.price }}</span>
          </div>
        </div>
      </div>

      <!-- Payment Method Selection -->
      <div class="card bg-base-100 shadow-xl" v-if="!showStripePayment">
        <div class="card-body">
          <h2 class="card-title">{{ $t('purchase.selectPayment') }}</h2>
          
          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">💳 Stripe (Credit Card)</span>
              <input type="radio" name="payment" value="stripe" v-model="paymentMethod" class="radio" />
            </label>
          </div>
          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">💰 PayPal</span>
              <input type="radio" name="payment" value="paypal" v-model="paymentMethod" class="radio" />
            </label>
          </div>
          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">💸 Alipay (支付宝)</span>
              <input type="radio" name="payment" value="alipay" v-model="paymentMethod" class="radio" />
            </label>
          </div>
          
          <div v-if="error" class="alert alert-error mt-4">
            <span>{{ error }}</span>
          </div>
          
          <div class="card-actions justify-end mt-4">
            <button @click="processPurchase" class="btn btn-primary" :disabled="!paymentMethod || processing">
              <span v-if="processing" class="loading loading-spinner loading-sm mr-2"></span>
              {{ processing ? $t('purchase.processing') : $t('purchase.continue') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Stripe Payment Form -->
      <div v-if="showStripePayment" class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">💳 {{ $t('purchase.enterCard') }}</h2>
          
          <div id="stripe-card-element" class="border border-base-300 rounded-lg p-4 my-4"></div>
          
          <div v-if="error" class="alert alert-error mt-4">
            <span>{{ error }}</span>
          </div>
          
          <div class="card-actions justify-between mt-4">
            <button @click="cancelStripePayment" class="btn btn-ghost" :disabled="processing">
              {{ $t('purchase.back') }}
            </button>
            <button @click="confirmStripePayment" class="btn btn-primary" :disabled="processing">
              <span v-if="processing" class="loading loading-spinner loading-sm mr-2"></span>
              {{ processing ? $t('purchase.processing') : $t('purchase.pay') + ' $' + plugin?.price }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { loadStripe } from '@stripe/stripe-js'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const paymentMethod = ref('alipay')
const plugin = ref(null)
const processing = ref(false)
const error = ref(null)
const showStripePayment = ref(false)
const stripePromise = ref(null)
const stripe = ref(null)
const cardElement = ref(null)
const currentOrder = ref(null)
const clientSecret = ref(null)

onMounted(async () => {
  // Load plugin details
  const pluginId = route.params.pluginId
  try {
    const response = await api.get(`/plugins/id/${pluginId}`)
    plugin.value = response.data.plugin
    
    // Initialize Stripe with publishable key from environment
    const stripeKey = import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY || 'pk_test_51QamdqRxiUU8ECKhCaJ0yCyH6QdmfxSXUxqLIUgdTmhAmqjBhWN1b9SXfVTUIhbv5UqNXOiT8Xjw4jN1uo3D2cBj00y6KK7gjD'
    stripePromise.value = loadStripe(stripeKey)
  } catch (err) {
    console.error('Failed to load plugin:', err)
    error.value = 'Failed to load plugin details'
  }
})

onUnmounted(() => {
  // Cleanup Stripe elements
  if (cardElement.value) {
    cardElement.value.destroy()
  }
})

const processPurchase = async () => {
  processing.value = true
  error.value = null
  
  try {
    // Step 1: Create order
    const orderResponse = await api.post('/orders', {
      plugin_id: route.params.pluginId,
      payment_method: paymentMethod.value
    })
    
    const order = orderResponse.data.order
    console.log('Order created:', order)
    currentOrder.value = order
    
    // Step 2: Process payment based on method
    if (paymentMethod.value === 'alipay') {
      const paymentResponse = await api.post('/payments/alipay/create', {
        order_id: order.id
      })
      
      // Redirect to Alipay payment page
      if (paymentResponse.data.pay_url) {
        window.location.href = paymentResponse.data.pay_url
      } else {
        throw new Error('No payment URL returned')
      }
    } else if (paymentMethod.value === 'stripe') {
      const paymentResponse = await api.post('/payments/stripe/create-intent', {
        order_id: order.id
      })
      
      // Store client secret and show Stripe payment form
      clientSecret.value = paymentResponse.data.client_secret
      showStripePayment.value = true
      processing.value = false
      
      // Initialize Stripe Elements
      await initStripeElements()
    } else if (paymentMethod.value === 'paypal') {
      const paymentResponse = await api.post('/payments/paypal/create-order', {
        order_id: order.id
      })
      
      // Redirect to PayPal
      if (paymentResponse.data.approve_url) {
        window.location.href = paymentResponse.data.approve_url
      } else {
        throw new Error('No PayPal URL returned')
      }
    }
  } catch (err) {
    console.error('Purchase failed:', err)
    error.value = err.response?.data?.error || err.message || 'Purchase failed. Please try again.'
    processing.value = false
  }
}

const initStripeElements = async () => {
  try {
    stripe.value = await stripePromise.value
    const elements = stripe.value.elements()
    cardElement.value = elements.create('card', {
      style: {
        base: {
          fontSize: '16px',
          color: '#424770',
          '::placeholder': {
            color: '#aab7c4',
          },
        },
        invalid: {
          color: '#9e2146',
        },
      },
    })
    
    // Wait for DOM to update
    await new Promise(resolve => setTimeout(resolve, 100))
    const cardElementDiv = document.getElementById('stripe-card-element')
    if (cardElementDiv) {
      cardElement.value.mount('#stripe-card-element')
    }
  } catch (err) {
    console.error('Failed to initialize Stripe:', err)
    error.value = 'Failed to initialize payment form'
  }
}

const confirmStripePayment = async () => {
  processing.value = true
  error.value = null
  
  try {
    const { error: stripeError, paymentIntent } = await stripe.value.confirmCardPayment(
      clientSecret.value,
      {
        payment_method: {
          card: cardElement.value,
        },
      }
    )
    
    if (stripeError) {
      throw new Error(stripeError.message)
    }
    
    if (paymentIntent.status === 'succeeded') {
      // Payment successful, redirect to orders page
      router.push('/orders')
    } else {
      throw new Error('Payment was not successful')
    }
  } catch (err) {
    console.error('Stripe payment failed:', err)
    error.value = err.message || 'Payment failed. Please try again.'
  } finally {
    processing.value = false
  }
}

const cancelStripePayment = () => {
  showStripePayment.value = false
  if (cardElement.value) {
    cardElement.value.destroy()
    cardElement.value = null
  }
  clientSecret.value = null
  currentOrder.value = null
  error.value = null
}
</script>
