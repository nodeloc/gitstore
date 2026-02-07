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
            <div class="text-4xl font-bold text-primary">20+</div>
            <div class="text-sm text-base-content/60 mt-2">{{ $t('home.stats.plugins') }}</div>
          </div>
          <div class="text-center">
            <div class="text-4xl font-bold text-secondary">100+</div>
            <div class="text-sm text-base-content/60 mt-2">{{ $t('home.stats.users') }}</div>
          </div>
          <div class="text-center">
            <div class="text-4xl font-bold text-accent">+100%</div>
            <div class="text-sm text-base-content/60 mt-2">{{ $t('home.stats.engagement') }}</div>
          </div>
          <div class="text-center">
            <div class="text-4xl font-bold text-info">24/7</div>
            <div class="text-sm text-base-content/60 mt-2">{{ $t('home.stats.support') }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Value Proposition Section -->
    <div id="features" class="container mx-auto px-4 py-20 md:py-32">
      <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-bold mb-4">{{ $t('home.forumValueTitle') }}</h2>
        <p class="text-xl text-base-content/60 max-w-2xl mx-auto">
          {{ $t('home.forumValueSubtitle') }}
        </p>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-6xl mx-auto">
        <!-- Enhanced Features -->
        <div class="card bg-base-100 shadow-lg hover:shadow-2xl transition-all duration-300 border border-primary/20 group">
          <div class="card-body">
            <div class="flex items-start gap-4">
              <div class="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z" />
                </svg>
              </div>
              <div>
                <h3 class="card-title text-2xl mb-2">{{ $t('home.enhanceFeatures') }}</h3>
                <p class="text-base-content/70">{{ $t('home.enhanceFeaturesDesc') }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Better User Experience -->
        <div class="card bg-base-100 shadow-lg hover:shadow-2xl transition-all duration-300 border border-secondary/20 group">
          <div class="card-body">
            <div class="flex items-start gap-4">
              <div class="w-14 h-14 rounded-xl bg-secondary/10 flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <h3 class="card-title text-2xl mb-2">{{ $t('home.improveExperience') }}</h3>
                <p class="text-base-content/70">{{ $t('home.improveExperienceDesc') }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Increase Engagement -->
        <div class="card bg-base-100 shadow-lg hover:shadow-2xl transition-all duration-300 border border-success/20 group">
          <div class="card-body">
            <div class="flex items-start gap-4">
              <div class="w-14 h-14 rounded-xl bg-success/10 flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
              </div>
              <div>
                <h3 class="card-title text-2xl mb-2">{{ $t('home.increaseEngagement') }}</h3>
                <p class="text-base-content/70">{{ $t('home.increaseEngagementDesc') }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Professional Support -->
        <div class="card bg-base-100 shadow-lg hover:shadow-2xl transition-all duration-300 border border-accent/20 group">
          <div class="card-body">
            <div class="flex items-start gap-4">
              <div class="w-14 h-14 rounded-xl bg-accent/10 flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192l-3.536 3.536M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
              </div>
              <div>
                <h3 class="card-title text-2xl mb-2">{{ $t('home.professionalSupport') }}</h3>
                <p class="text-base-content/70">{{ $t('home.professionalSupportDesc') }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- How It Works Section -->
    <div class="bg-base-200 py-20 md:py-32">
      <div class="container mx-auto px-4">
        <div class="text-center mb-16">
          <h2 class="text-4xl md:text-5xl font-bold mb-4">{{ $t('howItWorks.title') }}</h2>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8 max-w-6xl mx-auto">
          <div class="text-center">
            <div class="w-16 h-16 rounded-2xl bg-primary text-primary-content flex items-center justify-center text-2xl font-bold mx-auto mb-4">1</div>
            <h3 class="text-xl font-bold mb-2">{{ $t('howItWorks.step1') }}</h3>
            <p class="text-base-content/70">{{ $t('howItWorks.step1Desc') }}</p>
          </div>
          <div class="text-center">
            <div class="w-16 h-16 rounded-2xl bg-secondary text-secondary-content flex items-center justify-center text-2xl font-bold mx-auto mb-4">2</div>
            <h3 class="text-xl font-bold mb-2">{{ $t('howItWorks.step2') }}</h3>
            <p class="text-base-content/70">{{ $t('howItWorks.step2Desc') }}</p>
          </div>
          <div class="text-center">
            <div class="w-16 h-16 rounded-2xl bg-accent text-accent-content flex items-center justify-center text-2xl font-bold mx-auto mb-4">3</div>
            <h3 class="text-xl font-bold mb-2">{{ $t('howItWorks.step3') }}</h3>
            <p class="text-base-content/70">{{ $t('howItWorks.step3Desc') }}</p>
          </div>
          <div class="text-center">
            <div class="w-16 h-16 rounded-2xl bg-success text-success-content flex items-center justify-center text-2xl font-bold mx-auto mb-4">4</div>
            <h3 class="text-xl font-bold mb-2">{{ $t('howItWorks.step4') }}</h3>
            <p class="text-base-content/70">{{ $t('howItWorks.step4Desc') }}</p>
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
