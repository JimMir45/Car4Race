<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()

onMounted(() => {
  userStore.fetchProfile()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <header class="bg-white dark:bg-gray-800 shadow">
      <div class="max-w-7xl mx-auto px-4 py-6 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">个人中心</h1>
        <router-link to="/" class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
          返回首页
        </router-link>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 py-8">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div v-if="userStore.user" class="space-y-4">
          <div class="flex items-center gap-4">
            <div class="w-16 h-16 bg-gray-300 dark:bg-gray-600 rounded-full flex items-center justify-center text-2xl">
              {{ userStore.user.nickname?.charAt(0) || '?' }}
            </div>
            <div>
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                {{ userStore.user.nickname }}
              </h2>
              <p class="text-gray-500 dark:text-gray-400">
                @{{ userStore.user.username }}
              </p>
            </div>
          </div>

          <div class="border-t dark:border-gray-700 pt-4 space-y-2">
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">手机号</span>
              <span class="text-gray-900 dark:text-white">{{ userStore.user.phone }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">账号类型</span>
              <span class="text-gray-900 dark:text-white">
                {{ userStore.user.role === 'vip' ? 'VIP会员' : userStore.user.role === 'admin' ? '管理员' : '普通用户' }}
              </span>
            </div>
          </div>

          <div class="border-t dark:border-gray-700 pt-4">
            <button
              @click="userStore.logout()"
              class="w-full py-2 px-4 bg-red-600 text-white rounded-lg hover:bg-red-700"
            >
              退出登录
            </button>
          </div>
        </div>

        <div v-else class="text-center py-8 text-gray-500 dark:text-gray-400">
          加载中...
        </div>
      </div>
    </main>
  </div>
</template>
