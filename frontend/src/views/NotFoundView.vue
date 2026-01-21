<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center">
      <h1 class="text-6xl font-bold mb-4">404</h1>
      <p class="text-2xl mb-8">{{ $t('notFound.message') }}</p>
      <router-link to="/" class="btn btn-primary">
        {{ $t('notFound.goHome') }}
      </router-link>
    </div>
  </div>
</template>

<script setup>
</script>
