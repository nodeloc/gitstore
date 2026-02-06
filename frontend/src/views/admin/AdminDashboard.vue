<template>
  <div class="bg-base-200 pt-8 pb-8 min-h-full">
    <div class="container mx-auto px-4">
      <div class="flex gap-6">
        <!-- Left Sidebar -->
        <aside class="w-full md:w-64 flex-shrink-0">
          <div class="bg-base-100 rounded-xl border border-base-300 shadow-sm p-4 sticky top-4">
          <div class="mb-6 px-4">
            <h2 class="text-xs font-semibold text-base-content/50 uppercase tracking-wider">{{ $t('admin.menu') }}</h2>
          </div>
          
          <!-- Menu Section -->
          <nav class="space-y-1 mb-6">
            <button @click="switchTab('overview')" 
                    :class="['flex items-center gap-3 w-full px-4 py-3 rounded-lg transition-colors text-sm font-medium', 
                             activeTab === 'overview' ? 'bg-primary/10 text-primary' : 'text-base-content/70 hover:bg-base-200 hover:text-base-content']">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h4a1 1 0 011 1v7a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM14 5a1 1 0 011-1h4a1 1 0 011 1v7a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 16a1 1 0 011-1h4a1 1 0 011 1v3a1 1 0 01-1 1H5a1 1 0 01-1-1v-3zM14 16a1 1 0 011-1h4a1 1 0 011 1v3a1 1 0 01-1 1h-4a1 1 0 01-1-1v-3z" />
              </svg>
              <span>{{ $t('admin.overview') }}</span>
            </button>

            <button @click="switchTab('plugins')" 
                    :class="['flex items-center gap-3 w-full px-4 py-3 rounded-lg transition-colors text-sm font-medium', 
                             activeTab === 'plugins' ? 'bg-primary/10 text-primary' : 'text-base-content/70 hover:bg-base-200 hover:text-base-content']">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
              <span>{{ $t('admin.plugins') }}</span>
            </button>

            <button @click="switchTab('categories')" 
                    :class="['flex items-center gap-3 w-full px-4 py-3 rounded-lg transition-colors text-sm font-medium', 
                             activeTab === 'categories' ? 'bg-primary/10 text-primary' : 'text-base-content/70 hover:bg-base-200 hover:text-base-content']">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
              </svg>
              <span>{{ $t('admin.categories') }}</span>
            </button>

            <button @click="switchTab('pages')" 
                    :class="['flex items-center gap-3 w-full px-4 py-3 rounded-lg transition-colors text-sm font-medium', 
                             activeTab === 'pages' ? 'bg-primary/10 text-primary' : 'text-base-content/70 hover:bg-base-200 hover:text-base-content']">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
              <span>{{ $t('admin.pages') }}</span>
            </button>

            <button @click="switchTab('users')" 
                    :class="['flex items-center gap-3 w-full px-4 py-3 rounded-lg transition-colors text-sm font-medium', 
                             activeTab === 'users' ? 'bg-primary/10 text-primary' : 'text-base-content/70 hover:bg-base-200 hover:text-base-content']">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
              <span>{{ $t('admin.users') }}</span>
            </button>

            <button @click="switchTab('orders')" 
                    :class="['flex items-center gap-3 w-full px-4 py-3 rounded-lg transition-colors text-sm font-medium', 
                             activeTab === 'orders' ? 'bg-primary/10 text-primary' : 'text-base-content/70 hover:bg-base-200 hover:text-base-content']">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
              <span>{{ $t('admin.orders') }}</span>
            </button>

            <button @click="switchTab('licenses')" 
                    :class="['flex items-center gap-3 w-full px-4 py-3 rounded-lg transition-colors text-sm font-medium', 
                             activeTab === 'licenses' ? 'bg-primary/10 text-primary' : 'text-base-content/70 hover:bg-base-200 hover:text-base-content']">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
              </svg>
              <span>{{ $t('admin.licenses') }}</span>
            </button>
          </nav>

          <!-- System Section -->
          <div class="mt-6 pt-6 border-t border-base-300">
            <div class="mb-4 px-4">
              <h2 class="text-xs font-semibold text-base-content/50 uppercase tracking-wider">{{ $t('admin.system') }}</h2>
            </div>
            <nav class="space-y-1">
            <button @click="switchTab('settings')" 
                    :class="['flex items-center gap-3 w-full px-4 py-3 rounded-lg transition-colors text-sm font-medium', 
                             activeTab === 'settings' ? 'bg-primary/10 text-primary' : 'text-base-content/70 hover:bg-base-200 hover:text-base-content']">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span>{{ $t('admin.settings') }}</span>
            </button>
            </nav>
          </div>
          </div>
        </aside>

      <!-- Main Content -->
      <main class="flex-1 overflow-x-hidden">
        <!-- Overview Tab -->
        <div v-if="activeTab === 'overview'">
          <div class="mb-6">
            <h2 class="text-3xl font-bold mb-2">{{ $t('admin.overview') }}</h2>
            <p class="text-base-content/70">{{ $t('admin.welcome') }}</p>
          </div>

          <!-- Stats Cards -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div class="card bg-base-100 shadow">
              <div class="card-body">
                <div class="flex items-center gap-4">
                  <div class="p-3 rounded-lg bg-primary/10">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                  </div>
                  <div>
                    <p class="text-base-content/70 text-sm">{{ $t('admin.totalUsers') }}</p>
                    <p class="text-3xl font-bold">{{ adminStore.dashboardStats.total_users }}</p>
                  </div>
                </div>
                <div class="mt-2 text-sm text-success">
                  <span>{{ $t('admin.newUsersToday') }} +{{ adminStore.dashboardStats.new_users_today }}</span>
                </div>
              </div>
            </div>

            <div class="card bg-base-100 shadow">
              <div class="card-body">
                <div class="flex items-center gap-4">
                  <div class="p-3 rounded-lg bg-secondary/10">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                  </div>
                  <div>
                    <p class="text-base-content/70 text-sm">{{ $t('admin.totalPlugins') }}</p>
                    <p class="text-3xl font-bold">{{ adminStore.dashboardStats.total_plugins }}</p>
                  </div>
                </div>
                <div class="mt-2 text-sm text-base-content/70">
                  <span>{{ $t('admin.publishedPlugins') }}</span>
                </div>
              </div>
            </div>

            <div class="card bg-base-100 shadow">
              <div class="card-body">
                <div class="flex items-center gap-4">
                  <div class="p-3 rounded-lg bg-accent/10">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <div>
                    <p class="text-base-content/70 text-sm">{{ $t('admin.totalRevenue') }}</p>
                    <p class="text-3xl font-bold">${{ adminStore.dashboardStats.total_revenue?.toFixed(2) || '0.00' }}</p>
                  </div>
                </div>
                <div class="mt-2 text-sm text-success">
                  <span>{{ $t('admin.today') }} +${{ adminStore.dashboardStats.revenue_today?.toFixed(2) || '0.00' }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Additional Stats Row -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="card bg-base-100 shadow">
              <div class="card-body">
                <h3 class="card-title text-lg">{{ $t('admin.activeLicenses') }}</h3>
                <div class="flex items-baseline gap-2">
                  <p class="text-4xl font-bold">{{ adminStore.dashboardStats.active_licenses }}</p>
                  <p class="text-base-content/70">{{ $t('admin.licensesNotExpired') }}</p>
                </div>
              </div>
            </div>

            <div class="card bg-base-100 shadow">
              <div class="card-body">
                <h3 class="card-title text-lg">{{ $t('admin.quickActions') }}</h3>
                <div class="flex gap-2 flex-wrap mt-2">
                  <router-link to="/admin/plugins/create" class="btn btn-primary btn-sm">{{ $t('admin.createPlugin') }}</router-link>
                  <button @click="switchTab('orders')" class="btn btn-secondary btn-sm">{{ $t('admin.viewOrders') }}</button>
                  <button @click="switchTab('users')" class="btn btn-accent btn-sm">{{ $t('admin.manageUsers') }}</button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Other Tabs Content (wrapped in cards) -->
        <div v-if="activeTab !== 'overview'" class="card bg-base-100 shadow-xl">
          <div class="card-body">
          <!-- Plugins Tab -->
          <div v-if="activeTab === 'plugins'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="card-title">{{ $t('admin.plugins') }}</h2>
              <div class="flex gap-2">
                <input type="text" v-model="pluginSearch" :placeholder="$t('admin.search') + '...'"
                       class="input input-bordered input-sm" @input="debouncedSearchPlugins" />
                <select v-model="pluginStatusFilter" class="select select-bordered select-sm" @change="loadPlugins(1)">
                  <option value="">{{ $t('admin.filter') }}</option>
                  <option value="published">{{ $t('admin.published') }}</option>
                  <option value="draft">{{ $t('admin.draft') }}</option>
                  <option value="archived">{{ $t('admin.archived') }}</option>
                </select>
                <router-link to="/admin/plugins/create" class="btn btn-primary btn-sm">{{ $t('admin.createPlugin') }}</router-link>
              </div>
            </div>
            <div v-if="adminStore.loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else-if="adminStore.plugins.length === 0" class="alert alert-info">
              <span>{{ $t('admin.noPluginsClickCreate') }}</span>
            </div>
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>{{ $t('admin.name') }}</th>
                    <th>{{ $t('admin.slug') }}</th>
                    <th>{{ $t('admin.price') }}</th>
                    <th>{{ $t('admin.status') }}</th>
                    <th>{{ $t('admin.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="plugin in adminStore.plugins" :key="plugin.id">
                    <td>
                      <div class="flex items-center gap-2">
                        <div class="avatar placeholder" v-if="plugin.icon_url">
                          <div class="w-8 rounded">
                            <img :src="plugin.icon_url" :alt="plugin.name" />
                          </div>
                        </div>
                        <span class="font-medium">{{ plugin.name }}</span>
                      </div>
                    </td>
                    <td class="text-sm opacity-70">{{ plugin.slug }}</td>
                    <td>${{ plugin.price?.toFixed(2) || '0.00' }}</td>
                    <td>
                      <span :class="[
                        'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                        plugin.status === 'published' ? 'bg-green-100 text-green-800' : plugin.status === 'draft' ? 'bg-gray-100 text-gray-800' : 'bg-yellow-100 text-yellow-800'
                      ]">
                        {{ getStatusText(plugin.status) }}
                      </span>
                    </td>
                    <td>
                      <router-link :to="`/admin/plugins/${plugin.id}/edit`" class="btn btn-ghost btn-xs">{{ $t('admin.edit') }}</router-link>
                      <button class="btn btn-ghost btn-xs text-error" @click="confirmDeletePlugin(plugin)">{{ $t('admin.delete') }}</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <!-- Pagination -->
            <div class="flex justify-center mt-4" v-if="adminStore.pluginsPagination.total_pages > 1">
              <div class="join">
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.pluginsPagination.page <= 1"
                        @click="loadPlugins(adminStore.pluginsPagination.page - 1)">«</button>
                <button class="join-item btn btn-sm">
                  {{ adminStore.pluginsPagination.page }} / {{ adminStore.pluginsPagination.total_pages }}
                </button>
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.pluginsPagination.page >= adminStore.pluginsPagination.total_pages"
                        @click="loadPlugins(adminStore.pluginsPagination.page + 1)">»</button>
              </div>
            </div>
          </div>

          <!-- Categories Tab -->
          <div v-if="activeTab === 'categories'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="card-title">{{ $t('admin.categories') }}</h2>
              <button @click="openCreateCategoryModal" class="btn btn-primary btn-sm">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                {{ $t('admin.createCategory') }}
              </button>
            </div>
            
            <div v-if="categoriesLoading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>{{ $t('admin.name') }}</th>
                    <th>{{ $t('admin.slug') }}</th>
                    <th>{{ $t('admin.description') }}</th>
                    <th>{{ $t('admin.sortOrder') }}</th>
                    <th>{{ $t('admin.status') }}</th>
                    <th>{{ $t('admin.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="category in categoriesList" :key="category.id">
                    <td class="font-medium">{{ category.name }}</td>
                    <td><code class="text-xs bg-base-200 px-2 py-1 rounded">{{ category.slug }}</code></td>
                    <td class="max-w-xs truncate">{{ category.description || '-' }}</td>
                    <td>{{ category.sort_order }}</td>
                    <td>
                      <span :class="[
                        'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                        category.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                      ]">
                        {{ category.is_active ? $t('admin.enable') : $t('admin.disable') }}
                      </span>
                    </td>
                    <td>
                      <button @click="openEditCategoryModal(category)" class="btn btn-ghost btn-xs">{{ $t('admin.edit') }}</button>
                      <button @click="confirmDeleteCategory(category)" class="btn btn-ghost btn-xs text-error">{{ $t('admin.delete') }}</button>
                    </td>
                  </tr>
                </tbody>
              </table>

              <div v-if="categoriesList.length === 0" class="text-center py-8">
                <p class="text-base-content/60">{{ $t('admin.noCategoriesClickCreate') }}</p>
              </div>
            </div>
          </div>

          <!-- Pages Tab -->
          <div v-if="activeTab === 'pages'">
            <PageManagement />
          </div>

          <!-- Users Tab -->
          <div v-if="activeTab === 'users'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="card-title">{{ $t('admin.users') }}</h2>
              <div class="flex gap-2">
                <input type="text" v-model="userSearch" :placeholder="$t('admin.search') + '...'"
                       class="input input-bordered input-sm" @input="debouncedSearchUsers" />
                <select v-model="userRoleFilter" class="select select-bordered select-sm" @change="loadUsers(1)">
                  <option value="">{{ $t('admin.allRoles') }}</option>
                  <option value="admin">{{ $t('admin.admin') }}</option>
                  <option value="user">{{ $t('admin.user') }}</option>
                </select>
              </div>
            </div>
            <div v-if="adminStore.loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else-if="adminStore.users.length === 0" class="alert alert-info">
              <span>{{ $t('admin.noUsersData') }}</span>
            </div>
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>{{ $t('admin.user') }}</th>
                    <th>{{ $t('admin.email') }}</th>
                    <th>{{ $t('admin.role') }}</th>
                    <th>{{ $t('admin.status') }}</th>
                    <th>{{ $t('admin.registerTime') }}</th>
                    <th>{{ $t('admin.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="u in adminStore.users" :key="u.id">
                    <td>
                      <div class="flex items-center gap-3">
                        <div class="avatar placeholder">
                          <div class="w-10 rounded-full bg-neutral-focus text-neutral-content">
                            <span>{{ u.name?.[0] || u.email?.[0] || 'U' }}</span>
                          </div>
                        </div>
                        <div>{{ u.name || $t('admin.notSet') }}</div>
                      </div>
                    </td>
                    <td>{{ u.email }}</td>
                    <td>
                      <span :class="[
                        'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                        u.role === 'admin' ? 'bg-red-100 text-red-800' : 'bg-gray-100 text-gray-800'
                      ]">
                        {{ u.role === 'admin' ? $t('admin.admin') : $t('admin.user') }}
                      </span>
                    </td>
                    <td>
                      <span :class="[
                        'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                        u.is_active ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                      ]">
                        {{ u.is_active ? $t('admin.active') : $t('admin.inactive') }}
                      </span>
                    </td>
                    <td>{{ formatDate(u.created_at) }}</td>
                    <td>
                      <button class="btn btn-ghost btn-xs" @click="openUserModal(u)">{{ $t('admin.edit') }}</button>
                      <button class="btn btn-ghost btn-xs text-error" @click="confirmDeleteUser(u)"
                              v-if="u.id !== user?.id">{{ $t('admin.deactivate') }}</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <!-- Pagination -->
            <div class="flex justify-center mt-4" v-if="adminStore.usersPagination.total_pages > 1">
              <div class="join">
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.usersPagination.page <= 1"
                        @click="loadUsers(adminStore.usersPagination.page - 1)">«</button>
                <button class="join-item btn btn-sm">
                  {{ adminStore.usersPagination.page }} / {{ adminStore.usersPagination.total_pages }}
                </button>
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.usersPagination.page >= adminStore.usersPagination.total_pages"
                        @click="loadUsers(adminStore.usersPagination.page + 1)">»</button>
              </div>
            </div>
          </div>

          <!-- Orders Tab -->
          <div v-if="activeTab === 'orders'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="card-title">{{ $t('admin.orders') }}</h2>
              <div class="flex gap-2">
                <input type="text" v-model="orderSearch" :placeholder="$t('admin.search') + ' ' + $t('admin.orderNumber') + '...'"
                       class="input input-bordered input-sm" @input="debouncedSearchOrders" />
                <select v-model="orderStatusFilter" class="select select-bordered select-sm" @change="loadOrders(1)">
                  <option value="">{{ $t('admin.allStatuses') }}</option>
                  <option value="pending">{{ $t('admin.pending') }}</option>
                  <option value="paid">{{ $t('admin.paid') }}</option>
                  <option value="refunded">{{ $t('admin.refunded') }}</option>
                  <option value="failed">{{ $t('admin.failed') }}</option>
                </select>
              </div>
            </div>
            <div v-if="adminStore.loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else-if="adminStore.orders.length === 0" class="alert alert-info">
              <span>{{ $t('admin.noOrdersData') }}</span>
            </div>
            <div v-else class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>{{ $t('admin.orderNumber') }}</th>
                    <th>{{ $t('admin.user') }}</th>
                    <th>{{ $t('admin.plugin') }}</th>
                    <th>{{ $t('admin.amount') }}</th>
                    <th>{{ $t('admin.paymentMethod') }}</th>
                    <th>{{ $t('admin.status') }}</th>
                    <th>{{ $t('admin.createdAt') }}</th>
                    <th>{{ $t('admin.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="order in adminStore.orders" :key="order.id">
                    <td class="font-mono text-sm">{{ order.order_number }}</td>
                    <td>{{ order.user?.name || order.user?.email || '-' }}</td>
                    <td>{{ order.plugin?.name || '-' }}</td>
                    <td>${{ order.amount?.toFixed(2) }}</td>
                    <td>
                      <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
                        {{ getPaymentMethodText(order.payment_method) }}
                      </span>
                    </td>
                    <td>
                      <span :class="[
                        'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                        order.payment_status === 'paid' ? 'bg-green-100 text-green-800' : 
                        order.payment_status === 'failed' ? 'bg-red-100 text-red-800' : 
                        order.payment_status === 'refunded' ? 'bg-purple-100 text-purple-800' : 'bg-yellow-100 text-yellow-800'
                      ]">
                        {{ getOrderStatusText(order.payment_status) }}
                      </span>
                    </td>
                    <td>{{ formatDate(order.created_at) }}</td>
                    <td>
                      <div class="dropdown dropdown-end">
                        <label tabindex="0" class="btn btn-ghost btn-xs">{{ $t('admin.actions') }} ▼</label>
                        <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                          <li v-if="order.payment_status === 'pending'">
                            <a @click="updateOrderStatus(order, 'paid')">
                              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                              </svg>
                              {{ $t('admin.setAsPaid') }}
                            </a>
                          </li>
                          <li v-if="order.payment_status === 'pending'">
                            <a @click="updateOrderStatus(order, 'failed')">
                              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                              </svg>
                              {{ $t('admin.setAsFailed') }}
                            </a>
                          </li>
                          <li v-if="order.payment_status === 'paid' || order.payment_status === 'failed'">
                            <a @click="updateOrderStatus(order, 'pending')">
                              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                              </svg>
                              {{ $t('admin.setAsPending') }}
                            </a>
                          </li>
                          <li v-if="order.payment_status === 'paid'">
                            <a @click="confirmRefundOrder(order)" class="text-warning">
                              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                              </svg>
                              {{ $t('admin.refund') }}
                            </a>
                          </li>
                        </ul>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <!-- Pagination -->
            <div class="flex justify-center mt-4" v-if="adminStore.ordersPagination.total_pages > 1">
              <div class="join">
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.ordersPagination.page <= 1"
                        @click="loadOrders(adminStore.ordersPagination.page - 1)">«</button>
                <button class="join-item btn btn-sm">
                  {{ adminStore.ordersPagination.page }} / {{ adminStore.ordersPagination.total_pages }}
                </button>
                <button class="join-item btn btn-sm"
                        :disabled="adminStore.ordersPagination.page >= adminStore.ordersPagination.total_pages"
                        @click="loadOrders(adminStore.ordersPagination.page + 1)">»</button>
              </div>
            </div>
          </div>

          <!-- Licenses Tab -->
          <div v-if="activeTab === 'licenses'" class="space-y-6">
            <!-- Header with Filters -->
            <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
              <div>
                <h2 class="text-2xl font-bold text-base-content">{{ $t('admin.licenseManagement') }}</h2>
                <p class="text-sm text-base-content/60 mt-1">{{ $t('admin.manageAndMonitor') }}</p>
              </div>
              <div class="flex gap-3">
                <select v-model="licenseStatusFilter" 
                        class="select select-bordered select-sm bg-base-100 hover:bg-base-200 transition-colors" 
                        @change="loadLicenses(1)">
                  <option value="">📋 {{ $t('admin.allStatuses') }}</option>
                  <option value="active">✅ {{ $t('admin.active') }}</option>
                  <option value="expired">⏰ {{ $t('admin.expired') }}</option>
                  <option value="revoked">🚫 {{ $t('admin.revoked') }}</option>
                </select>
              </div>
            </div>

            <!-- Loading State -->
            <div v-if="adminStore.loading" class="flex flex-col items-center justify-center py-16">
              <span class="loading loading-spinner loading-lg text-primary"></span>
              <p class="mt-4 text-base-content/60">{{ $t('admin.loading') }}</p>
            </div>

            <!-- Empty State -->
            <div v-else-if="adminStore.licenses.length === 0" class="text-center py-16">
              <div class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-base-200 mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-base-content mb-2">{{ $t('admin.noLicensesData') }}</h3>
              <p class="text-base-content/60">{{ $t('admin.noLicenseRecords') }}</p>
            </div>

            <!-- Licenses Cards View -->
            <div v-else class="grid gap-4">
              <div v-for="license in adminStore.licenses" 
                   :key="license.id" 
                   class="card bg-base-100 shadow-sm hover:shadow-md transition-all duration-200 border border-base-300">
                <div class="card-body p-5">
                  <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                    <!-- User & Plugin Info -->
                    <div class="flex-1 space-y-3">
                      <div class="flex items-start gap-3">
                        <div class="avatar placeholder">
                          <div class="bg-primary text-primary-content rounded-full w-12 h-12">
                            <span class="text-lg">{{ (license.user?.name || license.user?.email || 'U')[0].toUpperCase() }}</span>
                          </div>
                        </div>
                        <div class="flex-1 min-w-0">
                          <div class="flex items-center gap-2 mb-1">
                            <h3 class="font-semibold text-base-content truncate">
                              {{ license.user?.name || license.user?.email || '-' }}
                            </h3>
                            <span :class="[
                              'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                              license.status === 'active' ? 'bg-green-100 text-green-800' : 
                              license.status === 'expired' ? 'bg-red-100 text-red-800' : 'bg-yellow-100 text-yellow-800'
                            ]">
                              {{ getLicenseStatusText(license.status) }}
                            </span>
                          </div>
                          <div class="flex flex-wrap items-center gap-2 text-sm text-base-content/60">
                            <span class="flex items-center gap-1">
                              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                              </svg>
                              {{ license.plugin?.name || '-' }}
                            </span>
                            <span class="text-base-content/40">•</span>
                            <span :class="['px-2 inline-flex text-xs leading-5 font-semibold rounded-full', license.license_type === 'permanent' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800']">
                              {{ license.license_type === 'permanent' ?  $t('admin.permanent') :  $t('admin.trial') }}
                            </span>
                          </div>
                        </div>
                      </div>
                      
                      <!-- License Details -->
                      <div class="grid grid-cols-2 gap-3 pl-15">
                        <div class="text-sm">
                          <span class="text-base-content/60">{{ $t('admin.maintenanceUntil') }}</span>
                          <p class="font-medium text-base-content mt-0.5">{{ formatDate(license.maintenance_until) }}</p>
                        </div>
                        <div class="text-sm">
                          <span class="text-base-content/60">{{ $t('admin.createdAt') }}</span>
                          <p class="font-medium text-base-content mt-0.5">{{ formatDate(license.created_at) }}</p>
                        </div>
                      </div>
                    </div>

                    <!-- Actions -->
                    <div class="flex lg:flex-col gap-2 lg:items-end">
                      <button v-if="license.status === 'active'"
                              class="btn btn-sm btn-outline btn-primary gap-2"
                              @click="openExtendModal(license)">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        {{ $t('admin.extend') }}
                      </button>
                      <button v-if="license.status === 'active'"
                              class="btn btn-sm btn-outline btn-error gap-2"
                              @click="confirmRevokeLicense(license)">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                        </svg>
                        {{ $t('admin.revokeAccess') }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Pagination -->
            <div class="flex justify-center mt-6" v-if="adminStore.licensesPagination.total_pages > 1">
              <div class="join shadow-sm">
                <button class="join-item btn btn-sm hover:btn-primary"
                        :disabled="adminStore.licensesPagination.page <= 1"
                        @click="loadLicenses(adminStore.licensesPagination.page - 1)">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                  </svg>
                </button>
                <button class="join-item btn btn-sm">
                  {{ $t('admin.pageOf', { current: adminStore.licensesPagination.page, total: adminStore.licensesPagination.total_pages }) }}
                </button>
                <button class="join-item btn btn-sm hover:btn-primary"
                        :disabled="adminStore.licensesPagination.page >= adminStore.licensesPagination.total_pages"
                        @click="loadLicenses(adminStore.licensesPagination.page + 1)">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <!-- Settings Tab -->
          <div v-if="activeTab === 'settings'">
            <h2 class="card-title mb-4">{{ $t('admin.systemSettings') }}</h2>
            <div v-if="adminStore.loading" class="flex justify-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
            <div v-else class="space-y-6 max-w-2xl">
              <div v-for="setting in sortedSettings" :key="setting.key" class="form-control">
                <label class="label">
                  <span class="label-text font-semibold">{{ $t(`settings.${setting.key}`) }}</span>
                  <span v-if="setting.description" class="label-text-alt text-base-content/60">{{ setting.description }}</span>
                </label>
                
                <!-- Boolean fields (true/false) - Toggle -->
                <input 
                  v-if="setting.value === 'true' || setting.value === 'false'"
                  type="checkbox" 
                  class="toggle toggle-primary" 
                  :checked="setting.value === 'true'"
                  @change="setting.value = $event.target.checked ? 'true' : 'false'"
                />
                
                <!-- Logo URL field - Image upload -->
                <div v-else-if="setting.key === 'logo_url'" class="space-y-2">
                  <div v-if="setting.value" class="flex items-center gap-4 mb-2">
                    <img :src="setting.value" alt="Logo" class="w-20 h-20 object-contain rounded-lg border border-base-300 bg-base-100 p-2" />
                    <button @click="setting.value = ''" class="btn btn-sm btn-ghost text-error">{{ $t('admin.removeLogo') }}</button>
                  </div>
                  <input 
                    type="file" 
                    accept="image/*" 
                    @change="(e) => uploadLogo(e, setting)"
                    class="file-input file-input-bordered w-full" 
                  />
                </div>
                
                <!-- Text fields -->
                <input 
                  v-else
                  type="text" 
                  v-model="setting.value" 
                  class="input input-bordered" 
                  :placeholder="$t(`settings.${setting.key}_placeholder`, setting.description)"
                />
              </div>
              <button class="btn btn-primary" @click="saveSettings">{{ $t('admin.saveSettings') }}</button>
            </div>
          </div>
          </div>
        </div>
      </main>
      </div>
    </div>

  <!-- Modals Section -->
    <!-- Plugin Modal -->
    <dialog ref="pluginModal" class="modal">
      <div class="modal-box w-11/12 max-w-3xl">
        <h3 class="font-bold text-lg">{{ editingPlugin ? $t('admin.editPlugin') : $t('admin.createPlugin') }}</h3>
        <form @submit.prevent="savePlugin" class="space-y-4 mt-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('admin.name') }} *</span></label>
              <input type="text" v-model="pluginForm.name" class="input input-bordered" required />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">Slug *</span></label>
              <input type="text" v-model="pluginForm.slug" class="input input-bordered" required />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('admin.price') }} (USD)</span></label>
              <input type="number" v-model="pluginForm.price" step="0.01" min="0" class="input input-bordered" />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('admin.status') }}</span></label>
              <select v-model="pluginForm.status" class="select select-bordered">
                <option value="draft">{{ $t('admin.draft') }}</option>
                <option value="published">{{ $t('admin.published') }}</option>
                <option value="archived">{{ $t('admin.archived') }}</option>
              </select>
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('admin.categories') }}</span></label>
              <input type="text" v-model="pluginForm.category" class="input input-bordered" />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('admin.version') }}</span></label>
              <input type="text" v-model="pluginForm.version" class="input input-bordered" placeholder="1.0.0" />
            </div>
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('admin.shortDescription') }}</span></label>
            <textarea v-model="pluginForm.description" class="textarea textarea-bordered" rows="2"></textarea>
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('admin.detailedDescription') }}</span>
              <span class="label-text-alt text-info">{{ $t('admin.supportsMarkdown') }}</span>
            </label>
            <MdEditor v-model="pluginForm.long_description" language="en-US" :preview="true" style="height: 400px;" />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="form-control col-span-2">
              <label class="label">
                <span class="label-text">{{ $t('admin.githubRepo') }}</span>
                <button type="button" class="btn btn-sm btn-ghost" @click="loadGitHubRepos">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                  {{ $t('admin.refreshList') }}
                </button>
              </label>
              <select v-if="githubRepos.length > 0" v-model="selectedGitHubRepo" class="select select-bordered" @change="onGitHubRepoSelected">
                <option :value="null">{{ $t('admin.manualInputOrSelect') }}</option>
                <option v-for="repo in githubRepos" :key="repo.id" :value="repo">
                  {{ repo.full_name }} {{ repo.private ? '(' + $t('admin.private') + ')' : '' }}
                </option>
              </select>
              <div v-if="loadingRepos" class="flex items-center justify-center py-2">
                <span class="loading loading-spinner loading-sm mr-2"></span>
                <span class="text-sm">{{ $t('admin.loadingRepoList') }}</span>
              </div>
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('admin.githubRepoURL') }}</span></label>
              <input type="text" v-model="pluginForm.github_repo_url" class="input input-bordered" />
            </div>
            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('admin.iconURL') }}</span></label>
              <input type="text" v-model="pluginForm.icon_url" class="input input-bordered" />
            </div>
          </div>
          <div class="modal-action">
            <button type="button" class="btn" @click="closePluginModal">{{ $t('admin.cancel') }}</button>
            <button type="submit" class="btn btn-primary">{{ editingPlugin ? $t('admin.save') : $t('admin.create') }}</button>
          </div>
        </form>
      </div>
    </dialog>

    <!-- Category Modal -->
    <dialog ref="categoryModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">{{ editingCategory ? $t('admin.editCategory') : $t('admin.createCategory') }}</h3>
        
        <form @submit.prevent="handleCategorySubmit" class="space-y-4">
          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('admin.categoryName') }} *</span></label>
            <input v-model="categoryForm.name" type="text" class="input input-bordered" required />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">Slug *</span></label>
            <input v-model="categoryForm.slug" type="text" class="input input-bordered" required placeholder="lowercase-with-dashes" />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('admin.description') }}</span></label>
            <textarea v-model="categoryForm.description" class="textarea textarea-bordered" rows="2"></textarea>
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('admin.icon') }}</span></label>
            <div v-if="categoryForm.icon_url" class="flex items-center gap-4 mb-2">
              <img :src="categoryForm.icon_url" alt="Icon" class="w-16 h-16 object-contain rounded-lg border border-base-300 bg-base-100 p-2" />
              <button type="button" @click="categoryForm.icon_url = ''" class="btn btn-sm btn-ghost text-error">{{ $t('admin.removeIcon') }}</button>
            </div>
            <input 
              type="file" 
              accept="image/*" 
              @change="uploadCategoryIcon"
              class="file-input file-input-bordered w-full" 
            />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('admin.sortOrderNum') }}</span></label>
            <input v-model.number="categoryForm.sort_order" type="number" class="input input-bordered" />
          </div>

          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-3">
              <input v-model="categoryForm.is_active" type="checkbox" class="checkbox" />
              <span class="label-text">{{ $t('admin.enable') }}</span>
            </label>
          </div>

          <div class="modal-action">
            <button type="button" @click="closeCategoryModal" class="btn btn-ghost">{{ $t('admin.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="categorySubmitting">
              <span v-if="categorySubmitting" class="loading loading-spinner loading-sm"></span>
              {{ editingCategory ? $t('admin.save') : $t('admin.create') }}
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <!-- User Modal -->
    <dialog ref="userModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">{{ $t('admin.editUser') }}</h3>
        <form @submit.prevent="saveUser" class="space-y-4 mt-4">
          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('admin.username') }}</span></label>
            <input type="text" v-model="userForm.name" class="input input-bordered" />
          </div>
          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('admin.role') }}</span></label>
            <select v-model="userForm.role" class="select select-bordered">
              <option value="user">{{ $t('admin.regularUser') }}</option>
              <option value="admin">{{ $t('admin.admin') }}</option>
            </select>
          </div>
          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">{{ $t('admin.accountStatus') }}</span>
              <input type="checkbox" v-model="userForm.is_active" class="toggle toggle-success" />
            </label>
          </div>
          <div class="modal-action">
            <button type="button" class="btn" @click="closeUserModal">{{ $t('admin.cancel') }}</button>
            <button type="submit" class="btn btn-primary">{{ $t('admin.save') }}</button>
          </div>
        </form>
      </div>
    </dialog>

    <!-- Extend License Modal -->
    <dialog ref="extendModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">{{ $t('admin.extendLicense') }}</h3>
        <form @submit.prevent="extendLicense" class="space-y-4 mt-4">
          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('admin.extendMonths') }}</span></label>
            <select v-model="extendMonths" class="select select-bordered">
              <option :value="1">1 {{ $t('admin.month') }}</option>
              <option :value="3">3 {{ $t('admin.month') }}</option>
              <option :value="6">6 {{ $t('admin.month') }}</option>
              <option :value="12">12 {{ $t('admin.month') }}</option>
              <option :value="24">24 {{ $t('admin.month') }}</option>
            </select>
          </div>
          <div class="modal-action">
            <button type="button" class="btn" @click="closeExtendModal">{{ $t('admin.cancel') }}</button>
            <button type="submit" class="btn btn-primary">{{ $t('admin.confirmExtend') }}</button>
          </div>
        </form>
      </div>
    </dialog>

    <!-- Confirm Dialog -->
    <dialog ref="confirmModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">{{ confirmTitle }}</h3>
        <p class="py-4">{{ confirmMessage }}</p>
        <div class="form-control" v-if="confirmNeedsReason">
          <label class="label"><span class="label-text">{{ $t('admin.reason') }}</span></label>
          <textarea v-model="confirmReason" class="textarea textarea-bordered" rows="2" required></textarea>
        </div>
        <div class="modal-action">
          <button class="btn" @click="closeConfirmModal">{{ $t('admin.cancel') }}</button>
          <button class="btn btn-error" @click="executeConfirm">{{ $t('admin.confirm') }}</button>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAdminStore } from '@/stores/admin'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import api from '@/utils/api'
import { toast } from '@/utils/toast'
import PageManagement from './PageManagement.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const authStore = useAuthStore()
const adminStore = useAdminStore()
const user = ref(authStore.user)

const activeTab = ref('overview')

// Sorted settings with priority order
const sortedSettings = computed(() => {
  const settings = adminStore.settings || []
  const order = ['site_name', 'site_subtitle', 'logo_url', 'support_email']
  
  return [...settings].sort((a, b) => {
    const indexA = order.indexOf(a.key)
    const indexB = order.indexOf(b.key)
    
    if (indexA !== -1 && indexB !== -1) return indexA - indexB
    if (indexA !== -1) return -1
    if (indexB !== -1) return 1
    return a.key.localeCompare(b.key)
  })
})

// Search and filters
const pluginSearch = ref('')
const pluginStatusFilter = ref('')
const userSearch = ref('')
const userRoleFilter = ref('')
const orderSearch = ref('')
const orderStatusFilter = ref('')
const licenseStatusFilter = ref('')

// Plugin form
const pluginModal = ref(null)
const editingPlugin = ref(null)
const pluginForm = ref({
  name: '',
  slug: '',
  description: '',
  long_description: '',
  price: 0,
  status: 'draft',
  category: '',
  version: '',
  github_repo_url: '',
  github_repo_id: 0,
  github_repo_name: '',
  icon_url: ''
})

// GitHub repos
const githubRepos = ref([])
const selectedGitHubRepo = ref(null)
const loadingRepos = ref(false)

// User form
const userModal = ref(null)
const editingUser = ref(null)
const userForm = ref({
  name: '',
  role: 'user',
  is_active: true
})

// Category management
const categoriesLoading = ref(false)
const categorySubmitting = ref(false)
const categoriesList = ref([])
const categoryModal = ref(null)
const editingCategory = ref(null)
const categoryForm = ref({
  name: '',
  slug: '',
  description: '',
  icon_url: '',
  sort_order: 0,
  is_active: true
})

// Extend modal
const extendModal = ref(null)
const extendingLicense = ref(null)
const extendMonths = ref(12)

// Confirm modal
const confirmModal = ref(null)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmNeedsReason = ref(false)
const confirmReason = ref('')
const confirmAction = ref(null)

// Debounce helper
function debounce(fn, delay) {
  let timer
  return (...args) => {
    clearTimeout(timer)
    timer = setTimeout(() => fn(...args), delay)
  }
}

const debouncedSearchPlugins = debounce(() => loadPlugins(1), 300)
const debouncedSearchUsers = debounce(() => loadUsers(1), 300)
const debouncedSearchOrders = debounce(() => loadOrders(1), 300)

onMounted(async () => {
  // Check if there's a tab parameter in the URL
  const route = useRoute()
  if (route.query.tab) {
    switchTab(route.query.tab)
  } else {
    await adminStore.fetchDashboardStats()
    await loadPlugins()
  }
})

function switchTab(tab) {
  activeTab.value = tab
  if (tab === 'overview') adminStore.fetchDashboardStats()
  if (tab === 'plugins') loadPlugins()
  else if (tab === 'categories') loadCategories()
  else if (tab === 'users') loadUsers()
  else if (tab === 'orders') loadOrders()
  else if (tab === 'licenses') loadLicenses()
  else if (tab === 'settings') adminStore.fetchSettings()
}

async function loadPlugins(page = 1) {
  await adminStore.fetchPlugins({
    page,
    page_size: 10,
    search: pluginSearch.value,
    status: pluginStatusFilter.value
  })
}

async function loadUsers(page = 1) {
  await adminStore.fetchUsers({
    page,
    page_size: 10,
    search: userSearch.value,
    role: userRoleFilter.value
  })
}

async function loadOrders(page = 1) {
  await adminStore.fetchOrders({
    page,
    page_size: 10,
    search: orderSearch.value,
    status: orderStatusFilter.value
  })
}

async function loadLicenses(page = 1) {
  await adminStore.fetchLicenses({
    page,
    page_size: 10,
    status: licenseStatusFilter.value
  })
}

// GitHub repos functions
async function loadGitHubRepos() {
  loadingRepos.value = true
  try {
    const response = await adminStore.fetchGitHubRepos()
    githubRepos.value = response.repositories || []
    if (githubRepos.value.length === 0) {
      toast.error(t('admin.githubRepoNotFound'))
    }
  } catch (err) {
    toast.error(t('admin.loadGithubRepoFailed') + ': ' + (err.response?.data?.error || err.message))
  } finally {
    loadingRepos.value = false
  }
}

function onGitHubRepoSelected() {
  if (selectedGitHubRepo.value) {
    const repo = selectedGitHubRepo.value
    pluginForm.value.github_repo_url = repo.html_url
    pluginForm.value.github_repo_id = repo.id
    pluginForm.value.github_repo_name = repo.full_name
    
    // 自动填充名称和描述（如果为空）
    if (!pluginForm.value.name) {
      pluginForm.value.name = repo.name
    }
    if (!pluginForm.value.slug) {
      pluginForm.value.slug = repo.name.toLowerCase().replace(/[^a-z0-9-]/g, '-')
    }
    if (!pluginForm.value.description && repo.description) {
      pluginForm.value.description = repo.description
    }
  }
}

// Plugin CRUD
function openPluginModal(plugin = null) {
  editingPlugin.value = plugin
  selectedGitHubRepo.value = null
  
  if (plugin) {
    pluginForm.value = { ...plugin }
  } else {
    pluginForm.value = {
      name: '',
      slug: '',
      description: '',
      long_description: '',
      price: 0,
      status: 'draft',
      category: '',
      version: '',
      github_repo_url: '',
      github_repo_id: 0,
      github_repo_name: '',
      icon_url: ''
    }
    // 自动加载 GitHub 仓库列表（仅在创建新插件时）
    if (githubRepos.value.length === 0) {
      loadGitHubRepos()
    }
  }
  pluginModal.value?.showModal()
}

function closePluginModal() {
  pluginModal.value?.close()
  editingPlugin.value = null
  selectedGitHubRepo.value = null
}

async function savePlugin() {
  try {
    if (editingPlugin.value) {
      await adminStore.updatePlugin(editingPlugin.value.id, pluginForm.value)
    } else {
      await adminStore.createPlugin(pluginForm.value)
    }
    closePluginModal()
    await loadPlugins()
  } catch (err) {
    toast.error(t('admin.saveFailed') + ': ' + (err.response?.data?.error || err.message))
  }
}

function confirmDeletePlugin(plugin) {
  confirmTitle.value = t('admin.confirmDelete')
  confirmMessage.value = t('admin.confirmDeletePlugin', { name: plugin.name })
  confirmNeedsReason.value = false
  confirmAction.value = async () => {
    try {
      await adminStore.deletePlugin(plugin.id)
      await loadPlugins()
    } catch (err) {
      toast.error(t('admin.deleteFailed') + ': ' + (err.response?.data?.error || err.message))
    }
  }
  confirmModal.value?.showModal()
}

// User CRUD
function openUserModal(u) {
  editingUser.value = u
  userForm.value = {
    name: u.name || '',
    role: u.role || 'user',
    is_active: u.is_active !== false
  }
  userModal.value?.showModal()
}

function closeUserModal() {
  userModal.value?.close()
  editingUser.value = null
}

async function saveUser() {
  try {
    await adminStore.updateUser(editingUser.value.id, userForm.value)
    closeUserModal()
    await loadUsers()
  } catch (err) {
    toast.error(t('admin.saveFailed') + ': ' + (err.response?.data?.error || err.message))
  }
}

function confirmDeleteUser(u) {
  confirmTitle.value = t('admin.confirmDeactivate')
  confirmMessage.value = t('admin.confirmDeleteUser', { name: u.name || u.email })
  confirmNeedsReason.value = false
  confirmAction.value = async () => {
    try {
      await adminStore.deleteUser(u.id)
      await loadUsers()
    } catch (err) {
      toast.error(t('admin.operationFailed') + ': ' + (err.response?.data?.error || err.message))
    }
  }
  confirmModal.value?.showModal()
}

// Order actions
async function updateOrderStatus(order, newStatus) {
  const statusTexts = {
    'pending': t('admin.pending'),
    'paid': t('admin.paid'),
    'failed': t('admin.failed'),
    'refunded': t('admin.refunded')
  }
  
  const confirm = window.confirm(t('admin.confirmUpdateOrderStatus', {
    orderNumber: order.order_number,
    status: statusTexts[newStatus]
  }))
  if (!confirm) return
  
  try {
    await adminStore.updateOrderStatus(order.id, newStatus)
    await loadOrders()
  } catch (err) {
    toast.error(t('admin.updateOrderStatusFailed') + ': ' + (err.response?.data?.error || err.message))
  }
}

function confirmRefundOrder(order) {
  confirmTitle.value = t('admin.confirmRefund')
  confirmMessage.value = t('admin.confirmRefundMessage', { 
    orderNumber: order.order_number, 
    amount: order.amount?.toFixed(2) 
  })
  confirmNeedsReason.value = false
  confirmAction.value = async () => {
    try {
      await adminStore.refundOrder(order.id)
      await loadOrders()
    } catch (err) {
      toast.error(t('admin.refundFailed') + ': ' + (err.response?.data?.error || err.message))
    }
  }
  confirmModal.value?.showModal()
}

// License actions
function openExtendModal(license) {
  extendingLicense.value = license
  extendMonths.value = 12
  extendModal.value?.showModal()
}

function closeExtendModal() {
  extendModal.value?.close()
  extendingLicense.value = null
}

async function extendLicense() {
  try {
    await adminStore.extendLicense(extendingLicense.value.id, extendMonths.value)
    closeExtendModal()
    await loadLicenses()
  } catch (err) {
    toast.error(t('admin.extendFailed') + ': ' + (err.response?.data?.error || err.message))
  }
}

function confirmRevokeLicense(license) {
  confirmTitle.value = t('admin.confirmRevokeTitle')
  confirmMessage.value = t('admin.confirmRevokeMessage', {
    user: license.User?.name || license.User?.email,
    plugin: license.Plugin?.name
  })
  confirmNeedsReason.value = true
  confirmReason.value = ''
  confirmAction.value = async () => {
    if (!confirmReason.value) {
      toast.error(t('admin.pleaseProvideReason'))
      return
    }
    try {
      await adminStore.revokeLicense(license.id, confirmReason.value)
      await loadLicenses()
    } catch (err) {
      toast.error(t('admin.revokeFailed') + ': ' + (err.response?.data?.error || err.message))
    }
  }
  confirmModal.value?.showModal()
}

// Settings
async function uploadLogo(event, setting) {
  const file = event.target.files[0]
  if (!file) return

  // Validate file type
  if (!file.type.startsWith('image/')) {
    toast.error(t('admin.uploadFailed') + ': Only image files are allowed')
    event.target.value = ''
    return
  }

  // Validate file size (max 5MB)
  if (file.size > 5 * 1024 * 1024) {
    toast.error(t('admin.uploadFailed') + ': File size exceeds 5MB limit')
    event.target.value = ''
    return
  }

  // Convert to base64
  const reader = new FileReader()
  reader.onload = (e) => {
    setting.value = e.target.result
    toast.success(t('admin.logoUploaded'))
  }
  reader.onerror = () => {
    toast.error(t('admin.uploadFailed'))
    event.target.value = ''
  }
  reader.readAsDataURL(file)
}

async function saveSettings() {
  try {
    await adminStore.updateSettings(adminStore.settings)
    toast.success(t('admin.settingsSaved'))
  } catch (err) {
    toast.error(t('admin.saveFailed') + ': ' + (err.response?.data?.error || err.message))
  }
}

// Category management functions
async function loadCategories() {
  try {
    categoriesLoading.value = true
    const response = await api.get('/admin/categories')
    categoriesList.value = response.data.categories || []
  } catch (error) {
    console.error('Failed to load categories:', error)
    toast.error(t('admin.loadCategoriesFailed'))
  } finally {
    categoriesLoading.value = false
  }
}

function openCreateCategoryModal() {
  editingCategory.value = null
  categoryForm.value = {
    name: '',
    slug: '',
    description: '',
    icon_url: '',
    sort_order: 0,
    is_active: true
  }
  categoryModal.value?.showModal()
}

function openEditCategoryModal(category) {
  editingCategory.value = category
  categoryForm.value = {
    name: category.name,
    slug: category.slug,
    description: category.description || '',
    icon_url: category.icon_url || '',
    sort_order: category.sort_order || 0,
    is_active: category.is_active
  }
  categoryModal.value?.showModal()
}

function closeCategoryModal() {
  categoryModal.value?.close()
}

async function uploadCategoryIcon(event) {
  const file = event.target.files[0]
  if (!file) return

  // Validate file type
  if (!file.type.startsWith('image/')) {
    toast.error(t('admin.uploadFailed') + ': Only image files are allowed')
    event.target.value = ''
    return
  }

  // Validate file size (max 5MB)
  if (file.size > 5 * 1024 * 1024) {
    toast.error(t('admin.uploadFailed') + ': File size exceeds 5MB limit')
    event.target.value = ''
    return
  }

  // Convert to base64
  const reader = new FileReader()
  reader.onload = (e) => {
    categoryForm.value.icon_url = e.target.result
    toast.success(t('admin.iconUploaded'))
  }
  reader.onerror = () => {
    toast.error(t('admin.uploadFailed'))
    categoryForm.value.icon_url = ''
    event.target.value = ''
  }
  reader.readAsDataURL(file)
}

async function handleCategorySubmit() {
  try {
    categorySubmitting.value = true
    
    if (editingCategory.value) {
      await api.put(`/admin/categories/${editingCategory.value.id}`, categoryForm.value)
      toast.success(t('admin.categoryUpdated'))
    } else {
      await api.post('/admin/categories', categoryForm.value)
      toast.success(t('admin.categoryCreated'))
    }
    
    closeCategoryModal()
    await loadCategories()
  } catch (error) {
    console.error('Failed to save category:', error)
    toast.error(error.response?.data?.error || t('admin.saveFailed'))
  } finally {
    categorySubmitting.value = false
  }
}

function confirmDeleteCategory(category) {
  confirmTitle.value = t('admin.confirmDelete')
  confirmMessage.value = t('admin.confirmDeleteCategory', { name: category.name })
  confirmNeedsReason.value = false
  confirmAction.value = async () => {
    try {
      await api.delete(`/admin/categories/${category.id}`)
      toast.success(t('admin.categoryDeleteSuccess'))
      await loadCategories()
    } catch (error) {
      console.error('Failed to delete category:', error)
      toast.error(error.response?.data?.error || t('admin.categoryDeleteFailed'))
    }
  }
  confirmModal.value?.showModal()
}

// Confirm modal
function closeConfirmModal() {
  confirmModal.value?.close()
  confirmAction.value = null
  confirmReason.value = ''
}

async function executeConfirm() {
  if (confirmAction.value) {
    await confirmAction.value()
  }
  closeConfirmModal()
}

// Helpers
function formatDate(date) {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

function getStatusBadgeClass(status) {
  switch (status) {
    case 'published': return 'badge-success'
    case 'draft': return 'badge-warning'
    case 'archived': return 'badge-neutral'
    default: return 'badge-neutral'
  }
}

function getStatusText(status) {
  switch (status) {
    case 'published': return t('admin.published')
    case 'draft': return t('admin.draft')
    case 'archived': return t('admin.archived')
    default: return status
  }
}

function getOrderStatusBadgeClass(status) {
  switch (status) {
    case 'paid': return 'badge-success'
    case 'pending': return 'badge-warning'
    case 'refunded': return 'badge-info'
    case 'failed': return 'badge-error'
    default: return 'badge-neutral'
  }
}

function getOrderStatusText(status) {
  switch (status) {
    case 'paid': return t('admin.paid')
    case 'pending': return t('admin.pending')
    case 'refunded': return t('admin.refunded')
    case 'failed': return t('admin.failed')
    default: return status
  }
}

function getPaymentMethodText(method) {
  switch (method) {
    case 'stripe': return 'Stripe'
    case 'paypal': return 'PayPal'
    case 'alipay': return t('admin.alipay')
    default: return method
  }
}

function getLicenseStatusBadgeClass(status) {
  switch (status) {
    case 'active': return 'badge-success'
    case 'expired': return 'badge-warning'
    case 'revoked': return 'badge-error'
    default: return 'badge-neutral'
  }
}

function getLicenseStatusText(status) {
  switch (status) {
    case 'active': return t('admin.active')
    case 'expired': return t('admin.expired')
    case 'revoked': return t('admin.revoked')
    default: return status
  }
}

function logout() {
  authStore.logout()
  router.push('/')
}
</script>

