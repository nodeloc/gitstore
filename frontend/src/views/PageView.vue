<template>
  <div class="bg-gray-50">
    <div class="container mx-auto px-4 py-8 max-w-4xl">
      <div v-if="loading" class="text-center py-12">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <div v-else-if="error" class="text-center py-12">
        <p class="text-red-600 text-lg">{{ error }}</p>
        <router-link to="/" class="text-blue-600 hover:underline mt-4 inline-block">
          {{ $t('common.backToHome') }}
        </router-link>
      </div>

      <div v-else-if="page" class="bg-white rounded-lg shadow-lg p-8">
        <h1 class="text-4xl font-bold mb-6 text-gray-900">{{ page.title }}</h1>
        <MdPreview :modelValue="page.content" class="prose prose-lg max-w-none" />
        <div class="mt-8 pt-6 border-t border-gray-200 text-sm text-gray-500">
          {{ $t('pages.lastUpdated') }}: {{ formatDate(page.updated_at) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import api from '@/utils/api'

const route = useRoute()
const { t } = useI18n()

const page = ref(null)
const loading = ref(true)
const error = ref(null)

const fetchPage = async () => {
  loading.value = true
  error.value = null
  
  try {
    const slug = route.params.slug
    const response = await api.get(`/pages/${slug}`)
    page.value = response.data
  } catch (err) {
    error.value = t('pages.notFound')
    console.error('Error fetching page:', err)
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

onMounted(() => {
  fetchPage()
})
</script>

<style scoped>
.prose {
  color: #374151;
}

.prose h1 {
  font-size: 2.25em;
  margin-top: 0;
  margin-bottom: 0.8888889em;
  line-height: 1.1111111;
  font-weight: 800;
}

.prose h2 {
  font-size: 1.5em;
  margin-top: 2em;
  margin-bottom: 1em;
  line-height: 1.3333333;
  font-weight: 700;
}

.prose h3 {
  font-size: 1.25em;
  margin-top: 1.6em;
  margin-bottom: 0.6em;
  line-height: 1.6;
  font-weight: 600;
}

.prose p {
  margin-top: 1.25em;
  margin-bottom: 1.25em;
  line-height: 1.75;
}

.prose a {
  color: #2563eb;
  text-decoration: underline;
  font-weight: 500;
}

.prose a:hover {
  color: #1d4ed8;
}

.prose ul {
  margin-top: 1.25em;
  margin-bottom: 1.25em;
  padding-left: 1.625em;
}

.prose li {
  margin-top: 0.5em;
  margin-bottom: 0.5em;
}

.prose strong {
  color: #111827;
  font-weight: 600;
}

.prose code {
  color: #111827;
  font-weight: 600;
  font-size: 0.875em;
  background-color: #f3f4f6;
  padding: 0.2em 0.4em;
  border-radius: 0.25rem;
}

.prose pre {
  background-color: #1f2937;
  color: #e5e7eb;
  overflow-x: auto;
  font-size: 0.875em;
  line-height: 1.7142857;
  margin-top: 1.7142857em;
  margin-bottom: 1.7142857em;
  border-radius: 0.375rem;
  padding: 0.8571429em 1.1428571em;
}

.prose pre code {
  background-color: transparent;
  border-width: 0;
  border-radius: 0;
  padding: 0;
  font-weight: inherit;
  color: inherit;
  font-size: inherit;
  font-family: inherit;
  line-height: inherit;
}

.prose blockquote {
  font-weight: 500;
  font-style: italic;
  color: #111827;
  border-left-width: 0.25rem;
  border-left-color: #e5e7eb;
  quotes: "\201C""\201D""\2018""\2019";
  margin-top: 1.6em;
  margin-bottom: 1.6em;
  padding-left: 1em;
}
</style>
