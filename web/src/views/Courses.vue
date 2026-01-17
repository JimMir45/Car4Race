<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { courseApi } from '../api'
import { useUserStore } from '../stores/user'

interface Course {
  id: number
  title: string
  slug: string
  description: string
  cover_image: string
  price: number
  orig_price: number
  sales_count: number
}

const router = useRouter()
const userStore = useUserStore()

const courses = ref<Course[]>([])
const loading = ref(true)
const sortBy = ref('newest')
const total = ref(0)

const sortOptions = [
  { value: 'newest', label: '最新' },
  { value: 'price_asc', label: '价格升序' },
  { value: 'price_desc', label: '价格降序' },
  { value: 'sales', label: '销量' },
]

const fetchCourses = async () => {
  loading.value = true
  try {
    const res = await courseApi.getList({ sort: sortBy.value })
    courses.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('Failed to fetch courses:', error)
  } finally {
    loading.value = false
  }
}

const goToDetail = (slug: string) => {
  router.push(`/courses/${slug}`)
}

const formatPrice = (price: number) => {
  return `¥${price.toFixed(0)}`
}

watch(sortBy, () => {
  fetchCourses()
})

onMounted(() => {
  fetchCourses()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <header class="bg-white dark:bg-gray-800 shadow">
      <div class="max-w-7xl mx-auto px-4 py-6 flex justify-between items-center">
        <div class="flex items-center gap-4">
          <router-link to="/" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            ← 返回首页
          </router-link>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">课程中心</h1>
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

    <main class="max-w-7xl mx-auto px-4 py-8">
      <!-- 排序控制 -->
      <div class="mb-6 flex justify-between items-center">
        <p class="text-gray-600 dark:text-gray-400">共 {{ total }} 门课程</p>
        <div class="flex items-center gap-2">
          <span class="text-gray-600 dark:text-gray-400">排序:</span>
          <select
            v-model="sortBy"
            class="px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
          >
            <option v-for="opt in sortOptions" :key="opt.value" :value="opt.value">
              {{ opt.label }}
            </option>
          </select>
        </div>
      </div>

      <!-- 加载中 -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <!-- 课程列表 -->
      <div v-else-if="courses.length > 0" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="course in courses"
          :key="course.id"
          @click="goToDetail(course.slug)"
          class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden cursor-pointer hover:shadow-lg transition-shadow"
        >
          <div class="aspect-video bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
            <img
              v-if="course.cover_image"
              :src="course.cover_image"
              :alt="course.title"
              class="w-full h-full object-cover"
            />
            <span v-else class="text-gray-400 dark:text-gray-500 text-4xl">📚</span>
          </div>
          <div class="p-4">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
              {{ course.title }}
            </h3>
            <p class="text-gray-600 dark:text-gray-400 text-sm mb-4 line-clamp-2">
              {{ course.description }}
            </p>
            <div class="flex justify-between items-center">
              <div>
                <span class="text-xl font-bold text-red-600">{{ formatPrice(course.price) }}</span>
                <span v-if="course.orig_price > course.price" class="ml-2 text-sm text-gray-400 line-through">
                  {{ formatPrice(course.orig_price) }}
                </span>
              </div>
              <span class="text-sm text-gray-500">{{ course.sales_count }} 人已购</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="text-center py-12">
        <p class="text-gray-500 dark:text-gray-400">暂无课程</p>
      </div>
    </main>
  </div>
</template>
