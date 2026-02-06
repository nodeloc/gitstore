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
            <label class="label cursor-pointer hover:bg-base-200 rounded-lg p-4 transition-colors">
              <span class="label-text flex items-center gap-3">
                <svg class="w-12 h-8" viewBox="0 0 60 25" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M59.64 14.28h-8.06c.19 1.93 1.6 2.55 3.2 2.55 1.64 0 2.96-.37 4.05-.95v3.32a8.33 8.33 0 0 1-4.56 1.1c-4.01 0-6.83-2.5-6.83-7.48 0-4.19 2.39-7.52 6.3-7.52 3.92 0 5.96 3.28 5.96 7.5 0 .4-.04 1.26-.06 1.48zm-5.92-5.62c-1.03 0-2.17.73-2.17 2.58h4.25c0-1.85-1.07-2.58-2.08-2.58zM40.95 20.3c-1.44 0-2.32-.6-2.9-1.04l-.02 4.63-4.12.87V5.57h3.76l.08 1.02a4.7 4.7 0 0 1 3.23-1.29c2.9 0 5.62 2.6 5.62 7.4 0 5.23-2.7 7.6-5.65 7.6zM40 8.95c-.95 0-1.54.34-1.97.81l.02 6.12c.4.44.98.78 1.95.78 1.52 0 2.54-1.65 2.54-3.87 0-2.15-1.04-3.84-2.54-3.84zM28.24 5.57h4.13v14.44h-4.13V5.57zm0-4.7L32.37 0v3.36l-4.13.88V.87zm-4.32 9.35v9.79H19.8V5.57h3.7l.12 1.22c1-1.77 3.07-1.41 3.62-1.22v3.79c-.52-.17-2.29-.43-3.32.86zm-8.55 4.72c0 2.43 2.6 1.68 3.12 1.46v3.36c-.55.3-1.54.54-2.89.54a4.15 4.15 0 0 1-4.27-4.24l.01-13.17 4.02-.86v3.54h3.14V9.1h-3.13v5.85zm-4.91.7c0 2.97-2.31 4.66-5.73 4.66a11.2 11.2 0 0 1-4.46-.93v-3.93c1.38.75 3.1 1.31 4.46 1.31.92 0 1.53-.24 1.53-1C6.26 13.77 0 14.51 0 9.95 0 7.04 2.28 5.3 5.62 5.3c1.36 0 2.72.2 4.09.75v3.88a9.23 9.23 0 0 0-4.1-1.06c-.86 0-1.44.25-1.44.9 0 1.85 6.29.97 6.29 5.88z" fill="#635BFF"/>
                </svg>
                <span class="text-base font-medium">Stripe</span>
              </span>
              <input type="radio" name="payment" value="stripe" v-model="paymentMethod" class="radio radio-primary" />
            </label>
          </div>
          <div class="form-control">
            <label class="label cursor-pointer hover:bg-base-200 rounded-lg p-4 transition-colors">
              <span class="label-text flex items-center gap-3">
                <svg class="w-12 h-8" viewBox="0 0 124 33" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M46.211 6.749h-6.839a.95.95 0 0 0-.939.802l-2.766 17.537a.57.57 0 0 0 .564.658h3.265a.95.95 0 0 0 .939-.803l.746-4.73a.95.95 0 0 1 .938-.803h2.165c4.505 0 7.105-2.18 7.784-6.5.306-1.89.013-3.375-.872-4.415-.972-1.142-2.696-1.746-4.985-1.746zm.789 6.405c-.374 2.454-2.249 2.454-4.062 2.454h-1.032l.724-4.583a.57.57 0 0 1 .563-.481h.473c1.235 0 2.4 0 3.002.704.359.42.468 1.044.332 1.906zM66.654 13.075h-3.275a.57.57 0 0 0-.563.481l-.145.916-.229-.332c-.709-1.029-2.29-1.373-3.868-1.373-3.619 0-6.71 2.741-7.312 6.586-.313 1.918.132 3.752 1.22 5.031.998 1.176 2.426 1.666 4.125 1.666 2.916 0 4.533-1.875 4.533-1.875l-.146.91a.57.57 0 0 0 .562.66h2.95a.95.95 0 0 0 .939-.803l1.77-11.209a.568.568 0 0 0-.561-.658zm-4.565 6.374c-.316 1.871-1.801 3.127-3.695 3.127-.951 0-1.711-.305-2.199-.883-.484-.574-.668-1.391-.514-2.301.295-1.855 1.805-3.152 3.67-3.152.93 0 1.686.309 2.184.892.499.589.697 1.411.554 2.317zM84.096 13.075h-3.291a.954.954 0 0 0-.787.417l-4.539 6.686-1.924-6.425a.953.953 0 0 0-.912-.678h-3.234a.57.57 0 0 0-.541.754l3.625 10.638-3.408 4.811a.57.57 0 0 0 .465.9h3.287a.949.949 0 0 0 .781-.408l10.946-15.8a.57.57 0 0 0-.468-.895z" fill="#253B80"/>
                  <path d="M94.992 6.749h-6.84a.95.95 0 0 0-.938.802l-2.766 17.537a.569.569 0 0 0 .562.658h3.51a.665.665 0 0 0 .656-.562l.785-4.971a.95.95 0 0 1 .938-.803h2.164c4.506 0 7.105-2.18 7.785-6.5.307-1.89.012-3.375-.873-4.415-.971-1.142-2.694-1.746-4.983-1.746zm.789 6.405c-.373 2.454-2.248 2.454-4.062 2.454h-1.031l.725-4.583a.568.568 0 0 1 .562-.481h.473c1.234 0 2.4 0 3.002.704.359.42.468 1.044.331 1.906zM115.434 13.075h-3.273a.567.567 0 0 0-.562.481l-.145.916-.23-.332c-.709-1.029-2.289-1.373-3.867-1.373-3.619 0-6.709 2.741-7.311 6.586-.312 1.918.131 3.752 1.219 5.031 1 1.176 2.426 1.666 4.125 1.666 2.916 0 4.533-1.875 4.533-1.875l-.146.91a.57.57 0 0 0 .564.66h2.949a.95.95 0 0 0 .938-.803l1.771-11.209a.571.571 0 0 0-.565-.658zm-4.565 6.374c-.314 1.871-1.801 3.127-3.695 3.127-.949 0-1.711-.305-2.199-.883-.484-.574-.666-1.391-.514-2.301.297-1.855 1.805-3.152 3.67-3.152.93 0 1.686.309 2.184.892.501.589.699 1.411.554 2.317zM119.295 7.23l-2.807 17.858a.569.569 0 0 0 .562.658h2.822c.469 0 .867-.34.938-.803l2.768-17.536a.57.57 0 0 0-.562-.659h-3.16a.571.571 0 0 0-.561.482z" fill="#179BD7"/>
                  <path d="M7.266 29.154l.523-3.322-1.165-.027H1.061L4.927 1.292a.316.316 0 0 1 .314-.268h9.38c3.114 0 5.263.648 6.385 1.927.526.6.861 1.227 1.023 1.917.17.724.173 1.589.007 2.644l-.012.077v.676l.526.298a3.69 3.69 0 0 1 1.065.812c.45.513.741 1.165.864 1.938.127.795.085 1.741-.123 2.812-.24 1.232-.628 2.305-1.152 3.183a6.547 6.547 0 0 1-1.825 2c-.696.494-1.523.869-2.458 1.109-.906.236-1.939.355-3.072.355h-.73c-.522 0-1.029.188-1.427.525a2.21 2.21 0 0 0-.744 1.328l-.055.299-.924 5.855-.042.215c-.011.068-.03.102-.058.125a.155.155 0 0 1-.096.035H7.266z" fill="#253B80"/>
                  <path d="M23.048 7.667c-.028.179-.06.362-.096.55-1.237 6.351-5.469 8.545-10.874 8.545H9.326c-.661 0-1.218.48-1.321 1.132L6.596 26.83l-.399 2.533a.704.704 0 0 0 .695.814h4.881c.578 0 1.069-.42 1.16-.99l.048-.248.919-5.832.059-.32c.09-.572.582-.992 1.16-.992h.73c4.729 0 8.431-1.92 9.513-7.476.452-2.321.218-4.259-.978-5.622a4.667 4.667 0 0 0-1.336-1.03z" fill="#179BD7"/>
                  <path d="M21.754 7.151a9.757 9.757 0 0 0-1.203-.267 15.284 15.284 0 0 0-2.426-.177h-7.352a1.172 1.172 0 0 0-1.159.992L8.05 17.605l-.045.289a1.336 1.336 0 0 1 1.321-1.132h2.752c5.405 0 9.637-2.195 10.874-8.545.037-.188.068-.371.096-.55a6.594 6.594 0 0 0-1.017-.429 9.045 9.045 0 0 0-.277-.087z" fill="#222D65"/>
                  <path d="M9.614 7.699a1.169 1.169 0 0 1 1.159-.991h7.352c.871 0 1.684.057 2.426.177a9.757 9.757 0 0 1 1.481.353c.365.121.704.264 1.017.429.368-2.347-.003-3.945-1.272-5.392C20.378.682 17.853 0 14.622 0h-9.38c-.66 0-1.223.48-1.325 1.133L.01 25.898a.806.806 0 0 0 .795.932h5.791l1.454-9.225 1.564-9.906z" fill="#253B80"/>
                </svg>
                <span class="text-base font-medium">PayPal</span>
              </span>
              <input type="radio" name="payment" value="paypal" v-model="paymentMethod" class="radio radio-primary" />
            </label>
          </div>
          <div class="form-control">
            <label class="label cursor-pointer hover:bg-base-200 rounded-lg p-4 transition-colors">
              <span class="label-text flex items-center gap-3">
                <svg class="h-8" viewBox="0 0 149.369 37.663" xmlns="http://www.w3.org/2000/svg">
                  <path d="M31.365 0H5.982C2.677 0 0 2.7 0 6.034V31.63c0 3.331 2.677 6.033 5.982 6.033h25.383c3.305 0 5.98-2.702 5.98-6.033V6.034C37.345 2.7 34.67 0 31.365 0" fill="#1677ff"/>
                  <path d="M10.091 28.964c-5.81 0-7.528-4.616-4.656-7.14.958-.854 2.709-1.27 3.642-1.364 3.451-.344 6.646.984 10.416 2.84-2.65 3.486-6.025 5.664-9.402 5.664m20.657-5.314c-1.495-.505-3.5-1.277-5.733-2.092 1.34-2.352 2.412-5.03 3.116-7.94H20.77v-2.674h9.018V9.45H20.77V4.994h-3.68c-.646 0-.646.642-.646.642v3.815h-9.12v1.493h9.12v2.674h-7.53v1.492h14.604a26.38 26.38 0 0 1-2.103 5.185c-4.74-1.577-9.796-2.855-12.973-2.069-2.031.505-3.34 1.406-4.108 2.35-3.529 4.33-.998 10.905 6.453 10.905 4.405 0 8.65-2.477 11.939-6.56 4.906 2.38 14.62 6.463 14.62 6.463v-5.82s-1.22-.098-6.598-1.914" fill="#fff"/>
                  <path d="M81.092 5.514c0 1.948 1.423 3.26 3.41 3.26s3.41-1.312 3.41-3.26c0-1.911-1.424-3.26-3.41-3.26s-3.41 1.349-3.41 3.26" fill="#1677ff"/>
                  <path d="M71.574 30.433h5.92V3.078h-5.92zM52.575 20.952l3.522-12.178h.15l3.335 12.178zm8.619-16.937h-7.945l-8.881 26.418h5.471l1.499-5.171h9.406l1.423 5.171h7.008zM81.541 30.433h5.92V10.348h-5.92zM149.332 10.385l.037-.037h-5.583l-3.522 12.216h-.188l-4.047-12.216h-6.633l7.982 20.16-3.335 6.146v.15h5.209zM98.442 26.986a6.458 6.458 0 0 1-2.023-.3V15.819c1.237-.862 2.249-1.274 3.522-1.274 2.212 0 3.973 1.761 3.973 5.508 0 4.797-2.586 6.933-5.472 6.933m3.748-17.05c-2.173 0-3.86.824-5.77 2.398v-1.986h-5.922v26.455h5.921v-6.557c1.124.3 2.173.45 3.448.45 5.283 0 10.043-3.898 10.043-10.83 0-6.22-3.449-9.93-7.72-9.93M122.314 25.562c-1.574.861-2.474 1.199-3.522 1.199-1.424 0-2.324-.937-2.324-2.436 0-.562.113-1.124.562-1.574.712-.712 2.099-1.236 5.284-1.986zm5.92-.15v-8.394c0-4.572-2.697-7.082-7.456-7.082-3.036 0-5.134.524-8.956 1.686l1.048 4.609c3.485-1.574 5.022-2.249 6.634-2.249 1.948 0 2.81 1.387 2.81 3.523v.15c-6.783 1.274-8.881 1.986-10.193 3.298-.974.974-1.386 2.36-1.386 3.971 0 3.86 2.998 5.922 5.809 5.922 2.098 0 3.784-.788 6.07-2.512l.412 2.1h5.92z"/>
                </svg>
                <span class="text-base font-medium">支付宝 (Alipay)</span>
              </span>
              <input type="radio" name="payment" value="alipay" v-model="paymentMethod" class="radio radio-primary" />
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
