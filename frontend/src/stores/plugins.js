import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/api'

export const usePluginsStore = defineStore('plugins', () => {
  const plugins = ref([])
  const currentPlugin = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const pagination = ref({
    total: 0,
    page: 1,
    pageSize: 12
  })

  async function fetchPlugins(params = {}) {
    try {
      loading.value = true
      error.value = null

      const response = await api.get('/plugins', { params })

      plugins.value = response.data.plugins || []
      if (response.data.pagination) {
        pagination.value = response.data.pagination
      }
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchPlugin(slug) {
    try {
      loading.value = true
      error.value = null

      const response = await api.get(`/plugins/${slug}`)
      currentPlugin.value = response.data.plugin

      return response.data.plugin
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  function clearCurrentPlugin() {
    currentPlugin.value = null
  }

  return {
    plugins,
    currentPlugin,
    loading,
    error,
    pagination,
    fetchPlugins,
    fetchPlugin,
    clearCurrentPlugin
  }
})
