<script setup lang="ts">
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
</script>

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <header class="bg-white dark:bg-gray-800 shadow">
      <div class="max-w-7xl mx-auto px-4 py-6 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Car4Race</h1>
        <nav class="flex gap-4">
          <template v-if="userStore.isLoggedIn">
            <router-link to="/profile" class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
              {{ userStore.user?.nickname || '个人中心' }}
            </router-link>
            <button @click="userStore.logout()" class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
              退出
            </button>
          </template>
          <template v-else>
            <router-link to="/login" class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
              登录
            </router-link>
          </template>
        </nav>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 py-8">
      <div class="grid md:grid-cols-2 gap-8">
        <!-- 私域视频网站入口 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">私域视频网站</h2>
          <p class="text-gray-600 dark:text-gray-300 mb-4">
            个人知识分享平台，支持笔记展示、课程销售、会员体系
          </p>
          <router-link to="/courses" class="inline-block px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
            进入 →
          </router-link>
        </div>

        <!-- 圈速榜入口 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">圈速榜</h2>
          <p class="text-gray-600 dark:text-gray-300 mb-4">
            赛道日爱好者的圈速记录与社区榜单平台
          </p>
          <a href="#" class="inline-block px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700">
            进入 →
          </a>
        </div>
      </div>
    </main>
  </div>
</template>
