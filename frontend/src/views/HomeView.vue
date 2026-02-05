<template>
  <div class="min-h-screen">
    <!-- Hero Section -->
    <div class="relative overflow-hidden bg-gradient-to-br from-primary/10 via-base-100 to-secondary/10">
      <div class="absolute inset-0 bg-grid-pattern opacity-5"></div>
      <div class="relative container mx-auto px-4 py-20 md:py-32">
        <div class="max-w-4xl mx-auto text-center">
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 text-primary text-sm font-medium mb-6">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            {{ $t('home.badge') }}
          </div>
          <h1 class="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
            {{ $t('home.title') }}
          </h1>
          <p class="text-xl md:text-2xl text-base-content/70 mb-8 max-w-2xl mx-auto">
            {{ $t('home.subtitle') }}
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <router-link to="/plugins" class="btn btn-primary btn-lg gap-2 shadow-lg hover:shadow-xl transition-all">
              {{ $t('home.browsePlugins') }}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
            </router-link>
            <a href="#features" class="btn btn-outline btn-lg gap-2">
              {{ $t('home.learnMore') }}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Section -->
    <div class="bg-base-200 py-12">
      <div class="container mx-auto px-4">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-8 max-w-4xl mx-auto">
          <div class="text-center">
            <div class="text-4xl font-bold text-primary">100+</div>
            <div class="text-sm text-base-content/60 mt-2">{{ $t('home.stats.plugins') }}</div>
          </div>
          <div class="text-center">
            <div class="text-4xl font-bold text-secondary">1000+</div>
            <div class="text-sm text-base-content/60 mt-2">{{ $t('home.stats.users') }}</div>
          </div>
          <div class="text-center">
            <div class="text-4xl font-bold text-accent">99.9%</div>
            <div class="text-sm text-base-content/60 mt-2">{{ $t('home.stats.uptime') }}</div>
          </div>
          <div class="text-center">
            <div class="text-4xl font-bold text-info">24/7</div>
            <div class="text-sm text-base-content/60 mt-2">{{ $t('home.stats.support') }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Features Section -->
    <div id="features" class="container mx-auto px-4 py-20 md:py-32">
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-4">{{ $t('home.features') }}</h2>
        <p class="text-xl text-base-content/60 max-w-2xl mx-auto">
          {{ $t('home.featuresSubtitle') }}
        </p>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
        <div class="card bg-base-100 shadow-lg hover:shadow-2xl transition-all duration-300 border border-base-300 group">
          <div class="card-body items-center text-center">
            <div class="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <h3 class="card-title text-2xl mb-2">{{ $t('home.feature1Title') }}</h3>
            <p class="text-base-content/70">{{ $t('home.feature1Description') }}</p>
          </div>
        </div>

        <div class="card bg-base-100 shadow-lg hover:shadow-2xl transition-all duration-300 border border-base-300 group">
          <div class="card-body items-center text-center">
            <div class="w-16 h-16 rounded-2xl bg-secondary/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <h3 class="card-title text-2xl mb-2">{{ $t('home.feature2Title') }}</h3>
            <p class="text-base-content/70">{{ $t('home.feature2Description') }}</p>
          </div>
        </div>

        <div class="card bg-base-100 shadow-lg hover:shadow-2xl transition-all duration-300 border border-base-300 group">
          <div class="card-body items-center text-center">
            <div class="w-16 h-16 rounded-2xl bg-accent/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="card-title text-2xl mb-2">{{ $t('home.feature3Title') }}</h3>
            <p class="text-base-content/70">{{ $t('home.feature3Description') }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- CTA Section -->
    <div class="bg-gradient-to-r from-primary to-secondary py-20">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-4xl md:text-5xl font-bold text-white mb-4">
          {{ $t('home.cta.title') }}
        </h2>
        <p class="text-xl text-white/90 mb-8 max-w-2xl mx-auto">
          {{ $t('home.cta.subtitle') }}
        </p>
        <router-link to="/plugins" class="btn btn-white btn-lg gap-2 shadow-xl">
          {{ $t('home.cta.button') }}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
          </svg>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.bg-grid-pattern {
  background-image: 
    linear-gradient(to right, currentColor 1px, transparent 1px),
    linear-gradient(to bottom, currentColor 1px, transparent 1px);
  background-size: 50px 50px;
}

.btn-white {
  @apply bg-white text-primary hover:bg-base-100;
}
</style>
