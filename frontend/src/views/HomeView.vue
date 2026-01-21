<template>
  <div class="min-h-screen">
    <!-- Hero Section -->
    <div class="hero min-h-[60vh] bg-base-200">
      <div class="hero-content text-center">
        <div class="max-w-md">
          <h1 class="text-5xl font-bold">{{ $t('home.title') }}</h1>
          <p class="py-6">{{ $t('home.subtitle') }}</p>
          <router-link to="/plugins" class="btn btn-primary">
            {{ $t('home.browsePlugins') }}
          </router-link>
        </div>
      </div>
    </div>

    <!-- Features Section -->
    <div class="container mx-auto px-4 py-16">
      <h2 class="text-3xl font-bold text-center mb-12">{{ $t('home.features') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body items-center text-center">
            <h3 class="card-title">{{ $t('home.feature1Title') }}</h3>
            <p>{{ $t('home.feature1Description') }}</p>
          </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body items-center text-center">
            <h3 class="card-title">{{ $t('home.feature2Title') }}</h3>
            <p>{{ $t('home.feature2Description') }}</p>
          </div>
        </div>
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body items-center text-center">
            <h3 class="card-title">{{ $t('home.feature3Title') }}</h3>
            <p>{{ $t('home.feature3Description') }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
</script>
