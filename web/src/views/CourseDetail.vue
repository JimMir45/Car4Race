<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { courseApi, orderApi, downloadApi } from '../api'
import { useUserStore } from '../stores/user'
import { marked } from 'marked'

interface CourseFile {
  id: number
  course_id: number
  file_type: string
  file_name: string
  file_size: number
  sort: number
}

interface Course {
  id: number
  title: string
  slug: string
  description: string
  cover_image: string
  price: number
  orig_price: number
  sales_count: number
  created_at: string
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const course = ref<Course | null>(null)
const purchased = ref(false)
const introContent = ref('')
const resourceFiles = ref<CourseFile[]>([])
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

const formatFileSize = (bytes: number): string => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

const renderedIntro = computed(() => {
  if (!introContent.value) return ''
  return marked(introContent.value)
})

const fetchCourse = async () => {
  loading.value = true
  try {
    const res: any = await courseApi.getDetail(slug.value)
    course.value = res.data.course
    purchased.value = res.data.purchased || false
    introContent.value = res.data.intro_content || ''
    resourceFiles.value = res.data.resource_files || []
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

const handleDownload = async (fileId?: number) => {
  if (!course.value) return

  errorMsg.value = ''
  successMsg.value = ''
  actionLoading.value = true

  try {
    const res: any = await downloadApi.createToken(course.value.id, fileId)
    const token = res.data.token
    const downloadUrl = res.data.download_url

    if (fileId) {
      // 直接下载指定文件
      window.location.href = downloadUrl
      successMsg.value = '下载已开始'
    } else {
      // 获取文件列表
      const downloadRes: any = await downloadApi.download(token)
      if (downloadRes.data.files && downloadRes.data.files.length > 0) {
        successMsg.value = '请选择要下载的文件'
      } else {
        errorMsg.value = '暂无可下载的资源文件'
      }
    }
  } catch (error: any) {
    errorMsg.value = error?.message || '获取下载链接失败'
  } finally {
    actionLoading.value = false
  }
}

const handleFileDownload = async (file: CourseFile) => {
  await handleDownload(file.id)
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
      <div v-else-if="course" class="space-y-6">
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
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
              <span>{{ course.sales_count }} 人已购</span>
              <span v-if="resourceFiles.length > 0">{{ resourceFiles.length }} 个资源文件</span>
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
                <span class="flex-1 px-6 py-3 bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400 rounded-lg text-center">
                  已购买，可下载资源
                </span>
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

        <!-- 课程介绍 Markdown -->
        <div v-if="renderedIntro" class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">课程介绍</h3>
          <div
            class="prose dark:prose-invert max-w-none"
            v-html="renderedIntro"
          ></div>
        </div>

        <!-- 资源文件列表 -->
        <div v-if="purchased && resourceFiles.length > 0" class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">课程资源</h3>
          <div class="space-y-3">
            <div
              v-for="file in resourceFiles"
              :key="file.id"
              class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg"
            >
              <div class="flex items-center gap-3">
                <span class="text-2xl">📦</span>
                <div>
                  <div class="text-gray-900 dark:text-white font-medium">{{ file.file_name }}</div>
                  <div class="text-sm text-gray-500 dark:text-gray-400">{{ formatFileSize(file.file_size) }}</div>
                </div>
              </div>
              <button
                @click="handleFileDownload(file)"
                :disabled="actionLoading"
                class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50"
              >
                下载
              </button>
            </div>
          </div>
        </div>

        <!-- 未购买时显示资源预览 -->
        <div v-if="!purchased && resourceFiles.length > 0" class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">课程资源预览</h3>
          <div class="space-y-3">
            <div
              v-for="file in resourceFiles"
              :key="file.id"
              class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg opacity-75"
            >
              <div class="flex items-center gap-3">
                <span class="text-2xl">📦</span>
                <div>
                  <div class="text-gray-900 dark:text-white font-medium">{{ file.file_name }}</div>
                  <div class="text-sm text-gray-500 dark:text-gray-400">{{ formatFileSize(file.file_size) }}</div>
                </div>
              </div>
              <span class="px-4 py-2 bg-gray-300 dark:bg-gray-600 text-gray-500 dark:text-gray-400 rounded cursor-not-allowed">
                购买后下载
              </span>
            </div>
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

<style>
/* Markdown 样式 */
.prose {
  color: inherit;
}

.prose h1, .prose h2, .prose h3, .prose h4, .prose h5, .prose h6 {
  color: inherit;
  font-weight: 600;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
}

.prose h1 { font-size: 1.5em; }
.prose h2 { font-size: 1.25em; }
.prose h3 { font-size: 1.1em; }

.prose p {
  margin-bottom: 1em;
  line-height: 1.7;
}

.prose ul, .prose ol {
  padding-left: 1.5em;
  margin-bottom: 1em;
}

.prose li {
  margin-bottom: 0.25em;
}

.prose code {
  background-color: rgba(0, 0, 0, 0.05);
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-size: 0.9em;
}

.dark .prose code {
  background-color: rgba(255, 255, 255, 0.1);
}

.prose pre {
  background-color: rgba(0, 0, 0, 0.05);
  padding: 1em;
  border-radius: 6px;
  overflow-x: auto;
  margin-bottom: 1em;
}

.dark .prose pre {
  background-color: rgba(255, 255, 255, 0.1);
}

.prose pre code {
  background-color: transparent;
  padding: 0;
}

.prose blockquote {
  border-left: 4px solid #e5e7eb;
  padding-left: 1em;
  margin-left: 0;
  margin-bottom: 1em;
  font-style: italic;
}

.dark .prose blockquote {
  border-left-color: #4b5563;
}

.prose a {
  color: #2563eb;
  text-decoration: underline;
}

.dark .prose a {
  color: #60a5fa;
}

.prose img {
  max-width: 100%;
  height: auto;
  border-radius: 6px;
  margin: 1em 0;
}

.prose table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 1em;
}

.prose th, .prose td {
  border: 1px solid #e5e7eb;
  padding: 0.5em;
  text-align: left;
}

.dark .prose th, .dark .prose td {
  border-color: #4b5563;
}

.prose th {
  background-color: rgba(0, 0, 0, 0.05);
}

.dark .prose th {
  background-color: rgba(255, 255, 255, 0.1);
}
</style>
