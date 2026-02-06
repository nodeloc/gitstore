<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold">{{ $t('pages.title') }}</h1>
      <button 
        @click="showCreateModal = true" 
        class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
      >
        {{ $t('pages.createPage') }}
      </button>
    </div>

    <!-- Pages List -->
    <div class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
              {{ $t('pages.title') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
              {{ $t('pages.slug') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
              {{ $t('pages.status') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
              {{ $t('pages.updatedAt') }}
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
              {{ $t('common.actions') }}
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="page in pages" :key="page.id">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ page.title }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-500">{{ page.slug }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span 
                :class="page.status === 'published' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'" 
                class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
              >
                {{ $t(`pages.status_${page.status}`) }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatDate(page.updated_at) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button 
                @click="editPage(page)" 
                class="text-blue-600 hover:text-blue-900 mr-4"
              >
                {{ $t('common.edit') }}
              </button>
              <button 
                @click="deletePage(page)" 
                class="text-red-600 hover:text-red-900"
              >
                {{ $t('common.delete') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal || showEditModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-20 mx-auto p-5 border w-11/12 max-w-4xl shadow-lg rounded-md bg-white">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-medium">
            {{ showEditModal ? $t('pages.editPage') : $t('pages.createPage') }}
          </h3>
          <button @click="closeModal" class="text-gray-400 hover:text-gray-500">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="submitForm">
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              {{ $t('pages.title') }}
            </label>
            <input 
              v-model="formData.title" 
              type="text" 
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              {{ $t('pages.slug') }}
            </label>
            <input 
              v-model="formData.slug" 
              type="text" 
              required
              :disabled="showEditModal"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100"
              :placeholder="$t('pages.slugPlaceholder')"
            />
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              {{ $t('pages.content') }}
            </label>
            <textarea 
              v-model="formData.content" 
              rows="15"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 font-mono text-sm"
              :placeholder="$t('pages.contentPlaceholder')"
            ></textarea>
            <p class="mt-1 text-sm text-gray-500">{{ $t('pages.markdownSupport') }}</p>
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              {{ $t('pages.status') }}
            </label>
            <select 
              v-model="formData.status"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="draft">{{ $t('pages.status_draft') }}</option>
              <option value="published">{{ $t('pages.status_published') }}</option>
            </select>
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              {{ $t('pages.sortOrder') }}
            </label>
            <input 
              v-model.number="formData.sort_order" 
              type="number" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div class="flex justify-end gap-3">
            <button 
              type="button" 
              @click="closeModal" 
              class="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
            >
              {{ $t('common.cancel') }}
            </button>
            <button 
              type="submit" 
              class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
              :disabled="loading"
            >
              {{ loading ? $t('common.saving') : $t('common.save') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/utils/api'
import { toast } from '@/utils/toast'

const { t } = useI18n()

const pages = ref([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const loading = ref(false)
const currentPage = ref(null)

const formData = ref({
  title: '',
  slug: '',
  content: '',
  status: 'draft',
  sort_order: 0
})

const fetchPages = async () => {
  try {
    const response = await api.get('/admin/pages')
    pages.value = response.data
  } catch (error) {
    toast.error(t('pages.fetchError'))
    console.error('Error fetching pages:', error)
  }
}

const editPage = async (page) => {
  try {
    const response = await api.get(`/admin/pages/${page.id}`)
    currentPage.value = response.data
    formData.value = {
      title: response.data.title,
      slug: response.data.slug,
      content: response.data.content,
      status: response.data.status,
      sort_order: response.data.sort_order
    }
    showEditModal.value = true
  } catch (error) {
    toast.error(t('pages.fetchError'))
    console.error('Error fetching page:', error)
  }
}

const deletePage = async (page) => {
  if (!confirm(t('pages.deleteConfirm'))) {
    return
  }

  try {
    await api.delete(`/admin/pages/${page.id}`)
    toast.success(t('pages.deleteSuccess'))
    fetchPages()
  } catch (error) {
    toast.error(t('pages.deleteError'))
    console.error('Error deleting page:', error)
  }
}

const submitForm = async () => {
  loading.value = true
  try {
    if (showEditModal.value) {
      await api.put(`/admin/pages/${currentPage.value.id}`, formData.value)
      toast.success(t('pages.updateSuccess'))
    } else {
      await api.post('/admin/pages', formData.value)
      toast.success(t('pages.createSuccess'))
    }
    closeModal()
    fetchPages()
  } catch (error) {
    toast.error(
      showEditModal.value ? t('pages.updateError') : t('pages.createError')
    )
    console.error('Error saving page:', error)
  } finally {
    loading.value = false
  }
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  currentPage.value = null
  formData.value = {
    title: '',
    slug: '',
    content: '',
    status: 'draft',
    sort_order: 0
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

onMounted(() => {
  fetchPages()
})
</script>
