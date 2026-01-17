<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { orderApi, downloadApi } from '../api'
import { useUserStore } from '../stores/user'

interface Course {
  id: number
  title: string
  slug: string
  price: number
}

interface Order {
  id: number
  order_no: string
  course_id: number
  amount: number
  status: string
  pay_method: string
  pay_time: string | null
  created_at: string
  course: Course
}

const router = useRouter()
const userStore = useUserStore()

const orders = ref<Order[]>([])
const loading = ref(true)
const total = ref(0)
const page = ref(1)
const pageSize = 10
const downloadingId = ref<number | null>(null)
const errorMsg = ref('')

const statusMap: Record<string, { label: string; color: string }> = {
  pending: { label: '待支付', color: 'text-yellow-600' },
  paid: { label: '已完成', color: 'text-green-600' },
  cancelled: { label: '已取消', color: 'text-gray-500' },
  refunded: { label: '已退款', color: 'text-red-600' },
}

const payMethodMap: Record<string, string> = {
  invite_code: '邀请码',
  wechat: '微信支付',
  alipay: '支付宝',
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const formatPrice = (price: number) => {
  return price === 0 ? '免费' : `¥${price.toFixed(0)}`
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const res = await orderApi.getList({ page: page.value, page_size: pageSize })
    orders.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  } finally {
    loading.value = false
  }
}

const goToCourse = (slug: string) => {
  router.push(`/courses/${slug}`)
}

const handleDownload = async (order: Order) => {
  if (order.status !== 'paid') return

  downloadingId.value = order.id
  errorMsg.value = ''

  try {
    const res = await downloadApi.createToken(order.course_id)
    const token = res.data.token

    const downloadRes = await downloadApi.download(token)
    const videoUrl = downloadRes.data.video_url

    window.open(videoUrl, '_blank')
  } catch (error: any) {
    errorMsg.value = error?.message || '获取下载链接失败'
  } finally {
    downloadingId.value = null
  }
}

onMounted(() => {
  fetchOrders()
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
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">我的订单</h1>
        </div>
        <nav class="flex gap-4">
          <router-link to="/profile" class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white">
            {{ userStore.user?.nickname || '个人中心' }}
          </router-link>
        </nav>
      </div>
    </header>

    <main class="max-w-4xl mx-auto px-4 py-8">
      <!-- 错误提示 -->
      <div v-if="errorMsg" class="mb-4 p-3 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded">
        {{ errorMsg }}
      </div>

      <!-- 加载中 -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <!-- 订单列表 -->
      <div v-else-if="orders.length > 0" class="space-y-4">
        <div
          v-for="order in orders"
          :key="order.id"
          class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden"
        >
          <!-- 订单头部 -->
          <div class="px-4 py-3 bg-gray-50 dark:bg-gray-700/50 flex justify-between items-center text-sm">
            <span class="text-gray-500 dark:text-gray-400">
              订单号: {{ order.order_no }}
            </span>
            <span :class="statusMap[order.status]?.color || 'text-gray-500'">
              {{ statusMap[order.status]?.label || order.status }}
            </span>
          </div>

          <!-- 订单内容 -->
          <div class="p-4 flex justify-between items-center">
            <div class="flex-1">
              <h3
                @click="goToCourse(order.course.slug)"
                class="text-lg font-medium text-gray-900 dark:text-white hover:text-blue-600 dark:hover:text-blue-400 cursor-pointer"
              >
                {{ order.course.title }}
              </h3>
              <div class="mt-2 text-sm text-gray-500 dark:text-gray-400 space-y-1">
                <p>支付方式: {{ payMethodMap[order.pay_method] || order.pay_method }}</p>
                <p>下单时间: {{ formatDate(order.created_at) }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-xl font-bold text-gray-900 dark:text-white">
                {{ formatPrice(order.amount) }}
              </p>
              <button
                v-if="order.status === 'paid'"
                @click="handleDownload(order)"
                :disabled="downloadingId === order.id"
                class="mt-2 px-4 py-1.5 text-sm bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50"
              >
                {{ downloadingId === order.id ? '获取中...' : '下载课程' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="text-center py-12">
        <p class="text-gray-500 dark:text-gray-400 mb-4">暂无订单</p>
        <router-link to="/courses" class="text-blue-600 hover:underline">
          去看看课程
        </router-link>
      </div>
    </main>
  </div>
</template>
