<template>
  <footer class="footer footer-center p-10 bg-base-300 text-base-content mt-auto">
    <div class="grid grid-flow-col gap-4">
      <a class="link link-hover">{{ $t('footer.about') }}</a>
      <a class="link link-hover">{{ $t('footer.contact') }}</a>
      <a class="link link-hover">{{ $t('footer.privacy') }}</a>
      <a class="link link-hover">{{ $t('footer.terms') }}</a>
    </div>
    <div>
      <p class="font-bold">
        GitStore
      </p>
      <p>{{ $t('footer.poweredBy') }}</p>
    </div>
    <div>
      <p>{{ $t('footer.copyright') }}</p>
    </div>
  </footer>
</template>

<script setup>
</script>
