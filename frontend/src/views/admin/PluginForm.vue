<template>
  <div class="bg-base-200 py-8">
  <div class="container mx-auto px-4">
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
            <!-- GitHub 仓库选择 - 放在最前面 -->
            <div class="form-control">
              <label class="label">
                <span class="label-text">GitHub 仓库</span>
                <span class="label-text-alt text-info">选择要关联的 GitHub 仓库</span>
              </label>
              
              <!-- 未加载状态：显示可点击的输入框 -->
              <div v-if="!reposLoaded && githubRepos.length === 0">
                <button
                  type="button"
                  class="input input-bordered w-full text-left flex items-center justify-between"
                  @click="loadGitHubRepos"
                  :disabled="loading"
                >
                  <span v-if="!loading" class="text-base-content/60">点击加载 GitHub 仓库列表</span>
                  <span v-else class="flex items-center gap-2">
                    <span class="loading loading-spinner loading-sm"></span>
                    加载中...
                  </span>
                </button>
              </div>
              
              <!-- 已加载状态：自定义下拉选择 -->
              <div v-else class="relative">
                <!-- 搜索输入框 / 显示选中项 -->
                <div class="relative">
                  <input 
                    v-model="repoSearch" 
                    type="text" 
                    class="input input-bordered w-full pr-24" 
                    :placeholder="form.github_repo_name ? `已选: ${form.github_repo_name}` : `搜索仓库... (共 ${githubRepos.length} 个)`"
                    @focus="showRepoDropdown = true"
                  />
                  <div class="absolute right-2 top-1/2 -translate-y-1/2 flex gap-1">
                    <button 
                      v-if="repoSearch"
                      type="button"
                      class="btn btn-ghost btn-xs btn-circle"
                      @click="repoSearch = ''"
                    >
                      ✕
                    </button>
                    <button 
                      type="button"
                      class="btn btn-ghost btn-xs"
                      @click="showRepoDropdown = !showRepoDropdown"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
                
                <!-- 下拉列表 -->
                <div 
                  v-show="showRepoDropdown && filteredRepos.length > 0"
                  class="absolute z-50 w-full mt-2 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-96 overflow-auto"
                >
                  <div class="p-2 text-xs text-base-content/60 border-b border-base-300">
                    找到 {{ filteredRepos.length }} 个仓库
                  </div>
                  <ul class="menu">
                    <li v-for="repo in filteredRepos" :key="repo.id">
                      <a 
                        @click="selectRepo(repo)"
                        class="flex flex-col items-start gap-1 py-3"
                        :class="{ 'active': form.github_repo_id === repo.id }"
                      >
                        <div class="font-semibold">{{ repo.full_name }}</div>
                        <div v-if="repo.description" class="text-xs text-base-content/60 line-clamp-1">
                          {{ repo.description }}
                        </div>
                        <div class="flex gap-2 text-xs">
                          <span v-if="repo.private" class="badge badge-sm badge-warning">Private</span>
                          <span v-else class="badge badge-sm badge-ghost">Public</span>
                          <span v-if="repo.language" class="badge badge-sm badge-ghost">{{ repo.language }}</span>
                        </div>
                      </a>
                    </li>
                  </ul>
                </div>
                
                <!-- 点击外部关闭下拉 -->
                <div 
                  v-show="showRepoDropdown"
                  class="fixed inset-0 z-40"
                  @click="showRepoDropdown = false"
                ></div>
                
                <!-- 显示已选择的仓库信息 -->
                <div v-if="form.github_repo_name" class="mt-2 text-sm text-success flex items-center gap-1">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  已关联: {{ form.github_repo_name }} (ID: {{ form.github_repo_id }})
                </div>
              </div>
            </div>

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

            <div class="form-c
const showRepoDropdown = ref(false)ontrol">
              <label class="label"><span class="label-text">状态</span></label>
              <select v-model="form.status" class="select select-bordered">
                <option value="draft">草稿（Draft）- 仅管理员可见</option>
                <option value="published">已发布（Published）- 用户可见并可购买</option>
                <option value="archived">已归档（Archived）- 已下架</option>
              </select>
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
const reposLoaded = ref(false)
const showRepoDropdown = ref(false)

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
  status: 'draft',
  github_repo_id: 0,
  github_repo_name: ''
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
    const plugin = response.data.plugin
    // Ensure status field is set correctly
    form.value = {
      ...plugin,
      status: plugin.status || 'draft'
    }
  } catch (error) {
    console.error('Failed to load plugin:', error)
    toast.error('加载插件失败')
  } finally {
    loading.value = false
  }
}

const loadGitHubRepos = async () => {
  if (reposLoaded.value) return // 已加载过则不重复加载
  
  try {
    loading.value = true
    const response = await api.get('/admin/github/repositories')
    githubRepos.value = response.data.repositories || []
    reposLoaded.value = true
    
    console.log('✅ 已加载仓库总数:', githubRepos.value.length)
    const privateCount = githubRepos.value.filter(r => r.private).length
    console.log('🔒 私有仓库数量:', privateCount)
    console.log('📂 公开仓库数量:', githubRepos.value.length - privateCount)
    
    if (githubRepos.value.length > 0) {
      toast.success(`已加载 ${githubRepos.value.length} 个仓库`)
    } else {
      toast.warning('未找到可用的 GitHub 仓库')
    }
  } catch (error) {
    console.error('Failed to load GitHub repos:', error)
    toast.error('加载 GitHub 仓库列表失败')
  } finally {
    loading.value = false
  }
}

const onGitHubRepoSelected = () => {
  if (selectedGitHubRepo.value) {
    const repo = JSON.parse(selectedGitHubRepo.value)
    form.value.github_repo_id = repo.id
    form.value.github_repo_name = repo.full_name
    if (!form.value.name) {
      form.value.name = repo.name
    }
    if (!form.value.slug) {
      form.value.slug = repo.name.toLowerCase()
    }
  }
}

const selectRepo = (repo) => {
  form.value.github_repo_id = repo.id
  form.value.github_repo_name = repo.full_name
  showRepoDropdown.value = false
  repoSearch.value = ''
  
  // 总是更新插件名称和 slug（基于新选择的仓库）
  form.value.name = repo.name
  form.value.slug = repo.name.toLowerCase().replace(/[^a-z0-9]+/g, '-')
  
  // 如果新仓库有描述，则更新描述
  if (repo.description) {
    form.value.description = repo.description
  }
  
  toast.success(`已选择仓库: ${repo.full_name}`)
}

const handleSubmit = async () => {
  try {
    loading.value = true
    
    // Ensure github_repo_id is a number
    const payload = {
      ...form.value,
      github_repo_id: form.value.github_repo_id ? Number(form.value.github_repo_id) : 0
    }
    
    if (isEdit.value) {
      await api.put(`/admin/plugins/${route.params.id}`, payload)
      toast.success('插件更新成功')
    } else {
      await api.post('/admin/plugins', payload)
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
        
        // Return the URL directly (it's already a valid path like /uploads/xxx)
        return response.data.url
      })
    )
    
    callback(res)
  } catch (error) {
    console.error('Failed to upload image:', error)
    toast.error('图片上传失败')
  }
}
</script>
