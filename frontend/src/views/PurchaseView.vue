<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">{{ $t('purchase.title') }}</h1>
    
    <div class="max-w-2xl mx-auto">
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">{{ $t('purchase.selectPayment') }}</h2>
          
          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">Stripe</span>
              <input type="radio" name="payment" value="stripe" v-model="paymentMethod" class="radio" />
            </label>
          </div>
          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">PayPal</span>
              <input type="radio" name="payment" value="paypal" v-model="paymentMethod" class="radio" />
            </label>
          </div>
          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">Alipay</span>
              <input type="radio" name="payment" value="alipay" v-model="paymentMethod" class="radio" />
            </label>
          </div>
          
          <div class="card-actions justify-end mt-4">
            <button @click="processPurchase" class="btn btn-primary" :disabled="!paymentMethod">
              {{ $t('purchase.continue') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const paymentMethod = ref('')

const processPurchase = () => {
  console.log('Processing purchase with', paymentMethod.value)
  // Handle payment processing
}
</script>
