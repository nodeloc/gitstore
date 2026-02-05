<template>
  <div class="container mx-auto px-4 py-8">
    <div class="max-w-4xl mx-auto">
      <div class="mb-6">
        <button @click="$router.back()" class="btn btn-ghost btn-sm">
          ← 返回
        </button>
      </div>

      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title text-3xl mb-6">{{ isEdit ? '编辑插件' : '创建插件' }}</h2>

          <form @submit.prevent="handleSubmit" class="space-y-6">
            <div class="form-control">
              <label class="label"><span class="label-text">插件名称 *</span></label>
              <input v-model="form.name" type="text" class="input input-bordered" required />
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text">Slug (URL标识) *</span></label>
              <input v-model="form.slug" type="text" class="input input-bordered" required placeholder="lowercase-with-dashes" />
              <label class="label"><span class="label-text-alt">将用于 URL，如：/plugins/your-slug</span></label>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text">分类</span></label>
              <select v-model="form.category" class="select select-bordered">
                <option value="">选择分类</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.slug">
                  {{ cat.name }}
                </option>
              </select>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text">简短描述</span></label>
              <textarea v-model="form.description" class="textarea textarea-bordered" rows="2"></textarea>
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text">详细描述</span>
                <span class="label-text-alt text-info">支持 Markdown 格式，可粘贴或拖拽上传图片</span>
              </label>
              <MdEditor 
                v-model="form.long_description" 
                language="zh-CN" 
                :preview="true" 
                style="height: 500px;"
                @onUploadImg="handleUploadImg"
              />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="form-control col-span-2">
                <label class="label">
                  <span class="label-text">GitHub 仓库</span>
                  <button type="button" class="btn btn-sm btn-ghost" @click="loadGitHubRepos">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                    刷新列表
                  </button>
                </label>
                
                <div v-if="githubRepos.length > 0" class="space-y-2">
                  <input 
                    v-model="repoSearch" 
                    type="text" 
                    class="input input-bordered w-full input-sm" 
                    placeholder="搜索仓库名称..."
                  />
                  <select v-model="selectedGitHubRepo" class="select select-bordered" @change="onGitHubRepoSelected">
                    <option value="">选择仓库 (共 {{ filteredRepos.length }} 个)</option>
                    <option v-for="repo in filteredRepos" :key="repo.id" :value="JSON.stringify(repo)">
                      {{ repo.full_name }}
                    </option>
                  </select>
                </div>
                <div v-else class="alert alert-info">
                  <span>点击刷新按钮加载您的 GitHub 仓库列表</span>
                </div>
              </div>

              <div class="form-control">
                <label class="label"><span class="label-text">GitHub 仓库 ID</span></label>
                <input v-model="form.github_repo_id" type="text" class="input input-bordered" placeholder="123456789" />
              </div>

              <div class="form-control">
                <label class="label"><span class="label-text">GitHub 仓库名称</span></label>
                <input v-model="form.github_repo_name" type="text" class="input input-bordered" placeholder="owner/repo" />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="form-control">
                <label class="label"><span class="label-text">价格 ($) *</span></label>
                <input v-model.number="form.price" type="number" step="0.01" class="input input-bordered" required />
              </div>

              <div class="form-control">
                <label class="label"><span class="label-text">货币</span></label>
                <select v-model="form.currency" class="select select-bordered">
                  <option value="USD">USD</option>
                  <option value="CNY">CNY</option>
                  <option value="EUR">EUR</option>
                </select>
              </div>

              <div class="form-control">
                <label class="label"><span class="label-text">版本</span></label>
                <input v-model="form.version" type="text" class="input input-bordered" placeholder="1.0.0" />
              </div>

              <div class="form-control">
                <label class="label"><span class="label-text">默认维护月数</span></label>
                <input v-model.number="form.default_maintenance_months" type="number" class="input input-bordered" />
              </div>
            </div>

            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-3">
                <input v-model="form.is_active" type="checkbox" class="checkbox" />
                <span class="label-text">启用（用户可见并可购买）</span>
              </label>
            </div>

            <div class="divider"></div>

            <div class="flex gap-4 justify-end">
              <button type="button" @click="$router.back()" class="btn btn-ghost">取消</button>
              <button type="submit" class="btn btn-primary" :disabled="loading">
                <span v-if="loading" class="loading loading-spinner loading-sm"></span>
                {{ isEdit ? '保存更改' : '创建插件' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import api from '@/utils/api'
import { toast } from '@/utils/toast'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const githubRepos = ref([])
const categories = ref([])
const selectedGitHubRepo = ref('')
const repoSearch = ref('')

const filteredRepos = computed(() => {
  // 暂时显示所有仓库以便调试
  let repos = githubRepos.value
  
  // 如果有搜索条件，进一步过滤
  if (repoSearch.value) {
    const search = repoSearch.value.toLowerCase()
    repos = repos.filter(repo => 
      repo.full_name.toLowerCase().includes(search) ||
      (repo.description && repo.description.toLowerCase().includes(search))
    )
  }
  
  return repos
})

const form = ref({
  name: '',
  slug: '',
  category: '',
  description: '',
  long_description: '',
  price: 0,
  currency: 'USD',
  version: '1.0.0',
  default_maintenance_months: 12,
  github_repo_id: '',
  github_repo_name: '',
  is_active: true
})

onMounted(async () => {
  await loadCategories()
  if (isEdit.value) {
    await loadPlugin()
  }
})

const loadCategories = async () => {
  try {
    const response = await api.get('/categories')
    categories.value = response.data.categories || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

const loadPlugin = async () => {
  try {
    loading.value = true
    const response = await api.get(`/admin/plugins/${route.params.id}`)
    form.value = response.data.plugin
  } catch (error) {
    console.error('Failed to load plugin:', error)
    toast.error('加载插件失败')
  } finally {
    loading.value = false
  }
}

const loadGitHubRepos = async () => {
  try {
    const response = await api.get('/admin/github/repositories')
    githubRepos.value = response.data.repositories || []
    console.log('✅ 已加载仓库总数:', githubRepos.value.length)
    const privateCount = githubRepos.value.filter(r => r.private).length
    console.log('🔒 私有仓库数量:', privateCount)
    console.log('📂 公开仓库数量:', githubRepos.value.length - privateCount)
    
    // 调试：显示前5个仓库的详细信息
    console.log('前5个仓库详情:', githubRepos.value.slice(0, 5).map(r => ({
      name: r.full_name,
      private: r.private,
      fork: r.fork
    })))
  } catch (error) {
    console.error('Failed to load GitHub repos:', error)
    toast.error('加载 GitHub 仓库列表失败')
  }
}

const onGitHubRepoSelected = () => {
  if (selectedGitHubRepo.value) {
    const repo = JSON.parse(selectedGitHubRepo.value)
    form.value.github_repo_id = String(repo.id)
    form.value.github_repo_name = repo.full_name
    if (!form.value.name) {
      form.value.name = repo.name
    }
    if (!form.value.slug) {
      form.value.slug = repo.name.toLowerCase()
    }
  }
}

const handleSubmit = async () => {
  try {
    loading.value = true
    
    if (isEdit.value) {
      await api.put(`/admin/plugins/${route.params.id}`, form.value)
      toast.success('插件更新成功')
    } else {
      await api.post('/admin/plugins', form.value)
      toast.success('插件创建成功')
    }
    
    router.push('/admin?tab=plugins')
  } catch (error) {
    console.error('Failed to save plugin:', error)
    toast.error(error.response?.data?.error || '保存失败')
  } finally {
    loading.value = false
  }
}

const handleUploadImg = async (files, callback) => {
  try {
    const res = await Promise.all(
      files.map(async (file) => {
        const formData = new FormData()
        formData.append('file', file)
        
        const response = await api.post('/admin/upload/image', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        
        // Return the full URL including the backend host
        const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
        return response.data.url.startsWith('http') 
          ? response.data.url 
          : `${baseURL}${response.data.url}`
      })
    )
    
    callback(res)
  } catch (error) {
    console.error('Failed to upload image:', error)
    toast.error('图片上传失败')
  }
}
</script>
