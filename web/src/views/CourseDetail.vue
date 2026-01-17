<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { courseApi, orderApi, downloadApi } from '../api'
import { useUserStore } from '../stores/user'

interface Course {
  id: number
  title: string
  slug: string
  description: string
  cover_image: string
  price: number
  orig_price: number
  video_url: string
  duration: number
  sales_count: number
  created_at: string
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const course = ref<Course | null>(null)
const purchased = ref(false)
const loading = ref(true)
const redeemCode = ref('')
const showRedeemModal = ref(false)
const actionLoading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

const slug = computed(() => route.params.slug as string)

const formatPrice = (price: number) => {
  return `¥${price.toFixed(0)}`
}

const formatDuration = (seconds: number) => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  if (hours > 0) {
    return `${hours}小时${minutes}分钟`
  }
  return `${minutes}分钟`
}

const fetchCourse = async () => {
  loading.value = true
  try {
    const res = await courseApi.getDetail(slug.value)
    course.value = res.data.course
    purchased.value = res.data.purchased || false
  } catch (error: any) {
    if (error?.code === 40403) {
      router.push('/courses')
    }
    console.error('Failed to fetch course:', error)
  } finally {
    loading.value = false
  }
}

const handleRedeem = async () => {
  if (!redeemCode.value.trim()) {
    errorMsg.value = '请输入邀请码'
    return
  }

  errorMsg.value = ''
  successMsg.value = ''
  actionLoading.value = true

  try {
    await orderApi.redeemCode(redeemCode.value.trim())
    successMsg.value = '兑换成功！'
    purchased.value = true
    showRedeemModal.value = false
    redeemCode.value = ''
  } catch (error: any) {
    errorMsg.value = error?.message || '兑换失败，请检查邀请码是否正确'
  } finally {
    actionLoading.value = false
  }
}

const handleDownload = async () => {
  if (!course.value) return

  errorMsg.value = ''
  successMsg.value = ''
  actionLoading.value = true

  try {
    const res = await downloadApi.createToken(course.value.id)
    const token = res.data.token

    // 获取下载链接
    const downloadRes = await downloadApi.download(token)
    const videoUrl = downloadRes.data.video_url

    // 在新窗口打开视频
    window.open(videoUrl, '_blank')
    successMsg.value = '下载链接已生成'
  } catch (error: any) {
    errorMsg.value = error?.message || '获取下载链接失败'
  } finally {
    actionLoading.value = false
  }
}

const goToLogin = () => {
  router.push({ path: '/login', query: { redirect: route.fullPath } })
}

onMounted(() => {
  fetchCourse()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <header class="bg-white dark:bg-gray-800 shadow">
      <div class="max-w-7xl mx-auto px-4 py-6 flex justify-between items-center">
        <div class="flex items-center gap-4">
          <router-link to="/courses" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            ← 返回课程
          </router-link>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">课程详情</h1>
        </div>
        <nav class="flex gap-4">
          <template v-if="userStore.isLoggedIn">
            <router-link to="/orders" class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
              我的订单
            </router-link>
            <router-link to="/profile" class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
              {{ userStore.user?.nickname || '个人中心' }}
            </router-link>
          </template>
          <template v-else>
            <router-link to="/login" class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
              登录
            </router-link>
          </template>
        </nav>
      </div>
    </header>

    <main class="max-w-4xl mx-auto px-4 py-8">
      <!-- 加载中 -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <!-- 课程详情 -->
      <div v-else-if="course" class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        <!-- 封面 -->
        <div class="aspect-video bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
          <img
            v-if="course.cover_image"
            :src="course.cover_image"
            :alt="course.title"
            class="w-full h-full object-cover"
          />
          <span v-else class="text-gray-400 dark:text-gray-500 text-6xl">📚</span>
        </div>

        <!-- 内容 -->
        <div class="p-6">
          <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">
            {{ course.title }}
          </h2>

          <p class="text-gray-600 dark:text-gray-400 mb-6">
            {{ course.description }}
          </p>

          <!-- 课程信息 -->
          <div class="flex flex-wrap gap-4 mb-6 text-sm text-gray-500 dark:text-gray-400">
            <span v-if="course.duration">时长: {{ formatDuration(course.duration) }}</span>
            <span>{{ course.sales_count }} 人已购</span>
          </div>

          <!-- 价格 -->
          <div class="flex items-center gap-4 mb-6">
            <span class="text-3xl font-bold text-red-600">{{ formatPrice(course.price) }}</span>
            <span v-if="course.orig_price > course.price" class="text-lg text-gray-400 line-through">
              {{ formatPrice(course.orig_price) }}
            </span>
          </div>

          <!-- 提示信息 -->
          <div v-if="errorMsg" class="mb-4 p-3 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded">
            {{ errorMsg }}
          </div>
          <div v-if="successMsg" class="mb-4 p-3 bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400 rounded">
            {{ successMsg }}
          </div>

          <!-- 操作按钮 -->
          <div class="flex gap-4">
            <template v-if="purchased">
              <!-- 已购买 -->
              <button
                @click="handleDownload"
                :disabled="actionLoading"
                class="flex-1 px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50"
              >
                {{ actionLoading ? '处理中...' : '下载课程' }}
              </button>
            </template>
            <template v-else-if="userStore.isLoggedIn">
              <!-- 已登录未购买 -->
              <button
                @click="showRedeemModal = true"
                class="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                使用邀请码兑换
              </button>
            </template>
            <template v-else>
              <!-- 未登录 -->
              <button
                @click="goToLogin"
                class="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                登录后购买
              </button>
            </template>
          </div>
        </div>
      </div>

      <!-- 课程不存在 -->
      <div v-else class="text-center py-12">
        <p class="text-gray-500 dark:text-gray-400">课程不存在</p>
        <router-link to="/courses" class="mt-4 inline-block text-blue-600 hover:underline">
          返回课程列表
        </router-link>
      </div>
    </main>

    <!-- 邀请码弹窗 -->
    <div v-if="showRedeemModal" class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full p-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          输入邀请码
        </h3>

        <div v-if="errorMsg" class="mb-4 p-3 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded text-sm">
          {{ errorMsg }}
        </div>

        <input
          v-model="redeemCode"
          type="text"
          placeholder="请输入邀请码"
          class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white mb-4"
          @keyup.enter="handleRedeem"
        />

        <div class="flex gap-4">
          <button
            @click="showRedeemModal = false; errorMsg = ''"
            class="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            取消
          </button>
          <button
            @click="handleRedeem"
            :disabled="actionLoading"
            class="flex-1 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
          >
            {{ actionLoading ? '兑换中...' : '确认兑换' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
