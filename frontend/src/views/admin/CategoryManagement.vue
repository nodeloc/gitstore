<template>
  <div class="container mx-auto px-4 py-8">
    <div class="max-w-6xl mx-auto">
      <div class="mb-6">
        <button @click="$router.back()" class="btn btn-ghost btn-sm">
          ← 返回
        </button>
      </div>

      <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold">分类管理</h1>
        <button @click="openCreateModal" class="btn btn-primary">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          创建分类
        </button>
      </div>

      <div v-if="loading" class="flex justify-center py-12">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <div v-else class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="overflow-x-auto">
            <table class="table">
              <thead>
                <tr>
                  <th>名称</th>
                  <th>Slug</th>
                  <th>描述</th>
                  <th>排序</th>
                  <th>状态</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="category in categories" :key="category.id">
                  <td>{{ category.name }}</td>
                  <td><code class="text-xs">{{ category.slug }}</code></td>
                  <td>{{ category.description || '-' }}</td>
                  <td>{{ category.sort_order }}</td>
                  <td>
                    <span class="badge" :class="category.is_active ? 'badge-success' : 'badge-ghost'">
                      {{ category.is_active ? '启用' : '禁用' }}
                    </span>
                  </td>
                  <td>
                    <button @click="openEditModal(category)" class="btn btn-ghost btn-xs">编辑</button>
                    <button @click="deleteCategory(category)" class="btn btn-ghost btn-xs text-error">删除</button>
                  </td>
                </tr>
              </tbody>
            </table>

            <div v-if="categories.length === 0" class="text-center py-8">
              <p class="text-base-content/60">暂无分类，点击"创建分类"按钮添加</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <dialog ref="categoryModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">{{ editingCategory ? '编辑分类' : '创建分类' }}</h3>
        
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div class="form-control">
            <label class="label"><span class="label-text">分类名称 *</span></label>
            <input v-model="form.name" type="text" class="input input-bordered" required />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">Slug *</span></label>
            <input v-model="form.slug" type="text" class="input input-bordered" required placeholder="lowercase-with-dashes" />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">描述</span></label>
            <textarea v-model="form.description" class="textarea textarea-bordered" rows="2"></textarea>
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">图标 URL</span></label>
            <input v-model="form.icon_url" type="text" class="input input-bordered" placeholder="https://..." />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">排序顺序</span></label>
            <input v-model.number="form.sort_order" type="number" class="input input-bordered" />
          </div>

          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-3">
              <input v-model="form.is_active" type="checkbox" class="checkbox" />
              <span class="label-text">启用</span>
            </label>
          </div>

          <div class="modal-action">
            <button type="button" @click="closeModal" class="btn btn-ghost">取消</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              {{ editingCategory ? '保存' : '创建' }}
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'
import { toast } from '@/utils/toast'

const loading = ref(false)
const submitting = ref(false)
const categories = ref([])
const categoryModal = ref(null)
const editingCategory = ref(null)

const form = ref({
  name: '',
  slug: '',
  description: '',
  icon_url: '',
  sort_order: 0,
  is_active: true
})

onMounted(() => {
  loadCategories()
})

const loadCategories = async () => {
  try {
    loading.value = true
    const response = await api.get('/admin/categories')
    categories.value = response.data.categories || []
  } catch (error) {
    console.error('Failed to load categories:', error)
    toast.error('加载分类失败')
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  editingCategory.value = null
  form.value = {
    name: '',
    slug: '',
    description: '',
    icon_url: '',
    sort_order: 0,
    is_active: true
  }
  categoryModal.value?.showModal()
}

const openEditModal = (category) => {
  editingCategory.value = category
  form.value = {
    name: category.name,
    slug: category.slug,
    description: category.description || '',
    icon_url: category.icon_url || '',
    sort_order: category.sort_order || 0,
    is_active: category.is_active
  }
  categoryModal.value?.showModal()
}

const closeModal = () => {
  categoryModal.value?.close()
}

const handleSubmit = async () => {
  try {
    submitting.value = true
    
    if (editingCategory.value) {
      await api.put(`/admin/categories/${editingCategory.value.id}`, form.value)
      toast.success('分类更新成功')
    } else {
      await api.post('/admin/categories', form.value)
      toast.success('分类创建成功')
    }
    
    closeModal()
    await loadCategories()
  } catch (error) {
    console.error('Failed to save category:', error)
    toast.error(error.response?.data?.error || '保存失败')
  } finally {
    submitting.value = false
  }
}

const deleteCategory = async (category) => {
  if (!confirm(`确定要删除分类"${category.name}"吗？`)) {
    return
  }
  
  try {
    await api.delete(`/admin/categories/${category.id}`)
    toast.success('分类删除成功')
    await loadCategories()
  } catch (error) {
    console.error('Failed to delete category:', error)
    toast.error(error.response?.data?.error || '删除失败')
  }
}
</script>
