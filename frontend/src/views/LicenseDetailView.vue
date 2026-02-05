<template>
  <div class="container mx-auto px-4 py-8">
    <div class="mb-6">
      <router-link to="/licenses" class="btn btn-ghost btn-sm">
        ← {{ $t('common.back') }}
      </router-link>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="license" class="space-y-6">
      <!-- License Info Card -->
      <div class="card bg-base-200 shadow-xl">
        <div class="card-body">
          <h1 class="card-title text-3xl mb-4">{{ license.plugin?.name || 'License' }}</h1>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p class="text-sm opacity-70">{{ $t('licenses.licenseId') }}</p>
              <p class="font-mono text-sm">{{ license.id }}</p>
            </div>
            
            <div>
              <p class="text-sm opacity-70">{{ $t('licenses.status') }}</p>
              <span :class="getStatusClass(license.status)">
                {{ license.status }}
              </span>
            </div>
            
            <div>
              <p class="text-sm opacity-70">{{ $t('licenses.type') }}</p>
              <p>{{ license.license_type }}</p>
            </div>
            
            <div>
              <p class="text-sm opacity-70">{{ $t('licenses.maintenanceUntil') }}</p>
              <p>{{ formatDate(license.maintenance_until) }}</p>
            </div>
            
            <div>
              <p class="text-sm opacity-70">{{ $t('licenses.createdAt') }}</p>
              <p>{{ formatDate(license.created_at) }}</p>
            </div>
            
            <div>
              <p class="text-sm opacity-70">GitHub Account</p>
              <p>{{ license.github_account?.login || 'N/A' }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Usage Instructions Card -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title text-2xl mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
            {{ $t('licenseDetail.howToUse') }}
          </h2>
          
          <div class="prose max-w-none">
            <!-- Success Alert -->
            <div class="alert alert-success mb-6">
              <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div>
                <strong>{{ $t('licenseDetail.activated') }}</strong>
                {{ $t('licenseDetail.activatedDesc') }}
              </div>
            </div>

            <!-- Step 1: Accept Invitation -->
            <div class="mb-8">
              <h3 class="flex items-center gap-2 text-xl font-bold mb-4">
                <span class="badge badge-primary badge-lg">1</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                {{ $t('licenseDetail.step1Title') }}
              </h3>
              
              <div class="card bg-primary/10 p-6 mb-4 border-2 border-primary">
                <p class="mb-4">{{ $t('licenseDetail.step1Desc') }}</p>
                
                <div class="space-y-4">
                  <div class="flex items-start gap-3">
                    <span class="badge badge-primary">①</span>
                    <div class="flex-1">
                      <p class="font-semibold mb-2">{{ $t('licenseDetail.step1Check') }}</p>
                      <p class="text-sm opacity-80">{{ $t('licenseDetail.step1CheckDesc') }}</p>
                      <code class="block mt-2 bg-base-300 px-3 py-2 rounded">{{ license.github_account?.login }}</code>
                    </div>
                  </div>
                  
                  <div class="flex items-start gap-3">
                    <span class="badge badge-primary">②</span>
                    <div class="flex-1">
                      <p class="font-semibold mb-2">{{ $t('licenseDetail.step1Open') }}</p>
                      <a href="https://github.com/notifications" target="_blank" class="btn btn-primary btn-sm mt-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                        </svg>
                        {{ $t('licenseDetail.step1OpenBtn') }}
                      </a>
                    </div>
                  </div>
                  
                  <div class="flex items-start gap-3">
                    <span class="badge badge-primary">③</span>
                    <div class="flex-1">
                      <p class="font-semibold mb-2">{{ $t('licenseDetail.step1Find') }}</p>
                      <p class="text-sm opacity-80">{{ $t('licenseDetail.step1FindDesc') }}</p>
                      <div class="mt-2 p-3 bg-base-200 rounded-lg">
                        <p class="text-sm font-mono">
                          <strong>{{ license.plugin?.github_repo_name }}</strong>
                        </p>
                      </div>
                    </div>
                  </div>
                  
                  <div class="flex items-start gap-3">
                    <span class="badge badge-primary">④</span>
                    <div class="flex-1">
                      <p class="font-semibold mb-2">{{ $t('licenseDetail.step1Accept') }}</p>
                      <p class="text-sm opacity-80">{{ $t('licenseDetail.step1AcceptDesc') }}</p>
                    </div>
                  </div>
                </div>
              </div>

              <div class="alert alert-info">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <div>
                  <p class="font-bold">{{ $t('licenseDetail.invitationNote') }}</p>
                  <p class="text-sm">{{ $t('licenseDetail.invitationNoteDesc') }}</p>
                </div>
              </div>
            </div>

            <!-- Step 2: Clone Repository -->
            <!-- Step 2: Clone Repository -->
            <div class="mb-8">
              <h3 class="flex items-center gap-2 text-xl font-bold mb-4">
                <span class="badge badge-primary badge-lg">2</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                {{ $t('licenseDetail.step2Title') }}
              </h3>
              
              <p class="mb-4">{{ $t('licenseDetail.step2Desc') }}</p>
              
              <div class="card bg-base-200 p-4 mb-4">
                <p class="font-semibold mb-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 inline mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                  </svg>
                  {{ $t('licenseDetail.repositoryUrl') }}
                </p>
                <a :href="`https://github.com/${license.plugin?.github_repo_name}`" 
                   target="_blank" 
                   class="link link-primary font-mono text-sm break-all">
                  https://github.com/{{ license.plugin?.github_repo_name }}
                </a>
              </div>
              
              <div class="space-y-4">
                <div>
                  <p class="font-semibold mb-2">{{ $t('licenseDetail.cloneCommand') }}</p>
                  <div class="mockup-code">
                    <pre><code>git clone https://github.com/{{ license.plugin?.github_repo_name }}.git</code></pre>
                  </div>
                </div>
                
                <div class="alert alert-warning">
                  <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                  <div>
                    <p class="font-bold">{{ $t('licenseDetail.cloneTroubleshoot') }}</p>
                    <ul class="text-sm list-disc ml-4 mt-1">
                      <li>{{ $t('licenseDetail.troubleCheck1') }} <code>{{ license.github_account?.login }}</code></li>
                      <li>{{ $t('licenseDetail.troubleCheck2') }}</li>
                      <li>{{ $t('licenseDetail.troubleCheck3') }}</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <!-- Step 3: Installation -->
            <div class="mb-8">
              <h3 class="flex items-center gap-2 text-xl font-bold mb-4">
                <span class="badge badge-primary badge-lg">3</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                </svg>
                {{ $t('licenseDetail.step3Title') }}
              </h3>
              
              <p class="mb-4">{{ $t('licenseDetail.step3Desc') }}</p>
              
              <div class="card bg-base-200 p-4">
                <p class="text-sm opacity-80">{{ $t('licenseDetail.step3Note') }}</p>
              </div>
            </div>

            <div class="divider"></div>

            <!-- Maintenance Period Info -->
            <h3 class="flex items-center gap-2 text-xl font-bold mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ $t('licenseDetail.maintenanceTitle') }}
            </h3>
            
            <p class="mb-4">
              {{ $t('licenseDetail.maintenanceDesc1') }}
              <strong class="text-primary">{{ formatDate(license.maintenance_until) }}</strong>
              {{ $t('licenseDetail.maintenanceDesc2') }}
            </p>
            
            <div class="overflow-x-auto mb-6">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>{{ $t('licenseDetail.tableFeature') }}</th>
                    <th>{{ $t('licenseDetail.tableDuring') }}</th>
                    <th>{{ $t('licenseDetail.tableAfter') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>{{ $t('licenseDetail.featureAccess') }}</td>
                    <td><span class="badge badge-success">✓ {{ $t('licenseDetail.available') }}</span></td>
                    <td><span class="badge badge-error">✗ {{ $t('licenseDetail.removed') }}</span></td>
                  </tr>
                  <tr>
                    <td>{{ $t('licenseDetail.featureUpdates') }}</td>
                    <td><span class="badge badge-success">✓ {{ $t('licenseDetail.yes') }}</span></td>
                    <td><span class="badge badge-ghost">✗ {{ $t('licenseDetail.no') }}</span></td>
                  </tr>
                  <tr>
                    <td>{{ $t('licenseDetail.featureSupport') }}</td>
                    <td><span class="badge badge-success">✓ {{ $t('licenseDetail.provided') }}</span></td>
                    <td><span class="badge badge-ghost">✗ {{ $t('licenseDetail.notProvided') }}</span></td>
                  </tr>
                  <tr>
                    <td>{{ $t('licenseDetail.featureDownloaded') }}</td>
                    <td><span class="badge badge-success">✓ {{ $t('licenseDetail.usable') }}</span></td>
                    <td><span class="badge badge-success">✓ {{ $t('licenseDetail.permanent') }}</span></td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="alert alert-warning">
              <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <div>
                <p class="font-bold">{{ $t('licenseDetail.importantReminder') }}</p>
                <p class="text-sm">{{ $t('licenseDetail.expiryWarning') }}</p>
                <ul class="text-sm list-disc ml-4 mt-1">
                  <li>{{ $t('licenseDetail.suggestionBackup') }}</li>
                  <li>{{ $t('licenseDetail.suggestionRenew') }}</li>
                </ul>
              </div>
            </div>

            <div class="divider"></div>

            <!-- Need Help Section -->
            <h3 class="flex items-center gap-2 text-xl font-bold mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ $t('licenseDetail.needHelp') }}
            </h3>
            
            <p class="mb-4">{{ $t('licenseDetail.helpIntro') }}</p>
            <ul class="space-y-2">
              <li>{{ $t('licenseDetail.helpDocs') }} <router-link to="/plugins" class="link link-primary">{{ $t('licenseDetail.helpDocsLink') }}</router-link></li>
              <li>{{ $t('licenseDetail.helpSupport') }} <a href="mailto:support@nodeloc.com" class="link link-primary">support@nodeloc.com</a></li>
              <li>{{ $t('licenseDetail.helpCommunity') }} <a href="https://github.com/nodeloc/discussions" target="_blank" class="link link-primary">{{ $t('licenseDetail.helpCommunityLink') }}</a></li>
            </ul>

            <div class="divider"></div>

            <!-- FAQ Section -->
            <h3 class="flex items-center gap-2 text-xl font-bold mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ $t('licenseDetail.faq') }}
            </h3>
            
            <details class="collapse collapse-arrow bg-base-200 mb-2">
              <summary class="collapse-title font-medium">{{ $t('licenseDetail.faq1Question') }}</summary>
              <div class="collapse-content">
                <p class="mb-2">{{ $t('licenseDetail.faq1Intro') }}</p>
                <ul class="list-decimal ml-6 space-y-1">
                  <li>{{ $t('licenseDetail.faq1Check1') }}</li>
                  <li>{{ $t('licenseDetail.faq1Check2') }} <code>{{ license.github_account?.login }}</code></li>
                  <li>{{ $t('licenseDetail.faq1Check3') }} <code>{{ formatDate(license.maintenance_until) }}</code></li>
                </ul>
              </div>
            </details>

            <details class="collapse collapse-arrow bg-base-200 mb-2">
              <summary class="collapse-title font-medium">{{ $t('licenseDetail.faq2Question') }}</summary>
              <div class="collapse-content">
                <p class="mb-2">{{ $t('licenseDetail.faq2Intro') }}</p>
                <ul class="space-y-1">
                  <li>✗ {{ $t('licenseDetail.faq2Item1') }}</li>
                  <li>✗ {{ $t('licenseDetail.faq2Item2') }}</li>
                  <li>✓ {{ $t('licenseDetail.faq2Item3') }}</li>
                  <li>💡 {{ $t('licenseDetail.faq2Item4') }}</li>
                </ul>
              </div>
            </details>

            <details class="collapse collapse-arrow bg-base-200 mb-2">
              <summary class="collapse-title font-medium">{{ $t('licenseDetail.faq3Question') }}</summary>
              <div class="collapse-content">
                <p>{{ $t('licenseDetail.faq3Answer') }} <code>{{ license.github_account?.login }}</code>{{ $t('licenseDetail.faq3Contact') }}</p>
              </div>
            </details>

            <details class="collapse collapse-arrow bg-base-200 mb-2">
              <summary class="collapse-title font-medium">{{ $t('licenseDetail.faq4Question') }}</summary>
              <div class="collapse-content">
                <p>{{ $t('licenseDetail.faq4Answer') }}</p>
              </div>
            </details>

          </div>
        </div>
      </div>

      <!-- Actions Card -->
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title mb-4">{{ $t('licenses.actions') }}</h2>
          <div class="flex gap-4">
            <button class="btn btn-primary" @click="renewLicense" v-if="needsRenewal">
              {{ $t('licenses.renewMaintenance') }}
            </button>
            <button class="btn btn-outline" @click="copyLicenseId">
              {{ $t('licenses.copyLicenseId') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-12">
      <p class="text-lg">{{ $t('licenses.notFound') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const license = ref(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const response = await api.get(`/licenses/${route.params.id}`)
    license.value = response.data.license
  } catch (error) {
    console.error('Failed to load license:', error)
  } finally {
    loading.value = false
  }
})

const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString()
}

const getStatusClass = (status) => {
  const classes = {
    'active': 'badge badge-success',
    'expired': 'badge badge-error',
    'revoked': 'badge badge-warning'
  }
  return classes[status] || 'badge'
}

const needsRenewal = computed(() => {
  if (!license.value?.maintenance_until) return false
  const maintenanceDate = new Date(license.value.maintenance_until)
  const now = new Date()
  const daysUntilExpiry = (maintenanceDate - now) / (1000 * 60 * 60 * 24)
  return daysUntilExpiry < 30
})

const renewLicense = () => {
  router.push(`/purchase/${license.value.plugin_id}`)
}

const copyLicenseId = async () => {
  try {
    await navigator.clipboard.writeText(license.value.id)
    alert('License ID copied to clipboard!')
  } catch (error) {
    console.error('Failed to copy:', error)
  }
}
</script>
