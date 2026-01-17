<script setup lang="ts">
import { useUserStore } from '../../stores/user'
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'

const userStore = useUserStore()
const router = useRouter()

onMounted(async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  await userStore.fetchProfile()
  if (!userStore.isAdmin) {
    router.push('/')
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <nav class="bg-white dark:bg-gray-800 shadow">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex justify-between h-16">
          <div class="flex">
            <router-link to="/admin" class="flex items-center px-2 text-gray-900 dark:text-white font-semibold">
              管理后台
            </router-link>
            <div class="hidden sm:ml-6 sm:flex sm:space-x-4">
              <router-link
                to="/admin/courses"
                class="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
                active-class="text-blue-600 dark:text-blue-400"
              >
                课程管理
              </router-link>
              <router-link
                to="/admin/invite-codes"
                class="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
                active-class="text-blue-600 dark:text-blue-400"
              >
                邀请码管理
              </router-link>
            </div>
          </div>
          <div class="flex items-center gap-4">
            <router-link to="/" class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white text-sm">
              返回前台
            </router-link>
            <span class="text-gray-600 dark:text-gray-300 text-sm">
              {{ userStore.user?.nickname || userStore.user?.phone }}
            </span>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto py-6 px-4">
      <router-view v-if="userStore.isAdmin" />
      <div v-else class="text-center py-12 text-gray-500">
        权限验证中...
      </div>
    </main>
  </div>
</template>
