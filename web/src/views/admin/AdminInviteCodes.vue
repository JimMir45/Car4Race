<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '../../api'

interface InviteCode {
  id: number
  code: string
  course_id: number
  max_uses: number
  used_count: number
  expire_at: string | null
  is_active: boolean
  created_at: string
  course?: {
    id: number
    title: string
  }
}

interface Course {
  id: number
  title: string
}

const inviteCodes = ref<InviteCode[]>([])
const courses = ref<Course[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)

const showModal = ref(false)
const form = ref({
  course_id: 0,
  max_uses: 1,
  expire_at: '',
})

async function fetchInviteCodes() {
  loading.value = true
  try {
    const res: any = await adminApi.getInviteCodes({ page: page.value, page_size: 20 })
    inviteCodes.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error('获取邀请码失败', e)
  } finally {
    loading.value = false
  }
}

async function fetchCourses() {
  try {
    const res: any = await adminApi.getCourses({ page: 1, page_size: 100 })
    courses.value = res.data.list || []
  } catch (e) {
    console.error('获取课程列表失败', e)
  }
}

function openCreateModal() {
  form.value = {
    course_id: courses.value[0]?.id || 0,
    max_uses: 1,
    expire_at: '',
  }
  showModal.value = true
}

async function createInviteCode() {
  try {
    const data: any = {
      course_id: form.value.course_id,
      max_uses: form.value.max_uses,
    }
    if (form.value.expire_at) {
      data.expire_at = new Date(form.value.expire_at).toISOString()
    }
    await adminApi.createInviteCode(data)
    showModal.value = false
    fetchInviteCodes()
  } catch (e: any) {
    alert(e.message || '创建失败')
  }
}

function copyCode(code: string) {
  navigator.clipboard.writeText(code)
  alert('已复制到剪贴板')
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '永不过期'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  fetchInviteCodes()
  fetchCourses()
})
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">邀请码管理</h1>
      <button
        @click="openCreateModal"
        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        创建邀请码
      </button>
    </div>

    <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-700">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">邀请码</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">关联课程</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">使用情况</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">过期时间</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">状态</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-for="code in inviteCodes" :key="code.id" class="hover:bg-gray-50 dark:hover:bg-gray-700">
            <td class="px-6 py-4">
              <code class="text-sm bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded text-gray-900 dark:text-white">
                {{ code.code }}
              </code>
            </td>
            <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">
              {{ code.course?.title || '-' }}
            </td>
            <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">
              {{ code.used_count }} / {{ code.max_uses }}
            </td>
            <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">
              {{ formatDate(code.expire_at) }}
            </td>
            <td class="px-6 py-4">
              <span
                :class="code.is_active && code.used_count < code.max_uses
                  ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
                  : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'"
                class="px-2 py-1 text-xs rounded-full"
              >
                {{ code.is_active && code.used_count < code.max_uses ? '可用' : '已用完' }}
              </span>
            </td>
            <td class="px-6 py-4 text-right text-sm font-medium">
              <button @click="copyCode(code.code)" class="text-blue-600 hover:text-blue-900 dark:text-blue-400">
                复制
              </button>
            </td>
          </tr>
          <tr v-if="inviteCodes.length === 0 && !loading">
            <td colspan="6" class="px-6 py-12 text-center text-gray-500">暂无邀请码</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 创建邀请码模态框 -->
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full mx-4">
        <div class="px-6 py-4 border-b dark:border-gray-700">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">创建邀请码</h3>
        </div>
        <form @submit.prevent="createInviteCode" class="px-6 py-4 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">关联课程</label>
            <select
              v-model="form.course_id"
              required
              class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            >
              <option v-for="course in courses" :key="course.id" :value="course.id">
                {{ course.title }}
              </option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">最大使用次数</label>
            <input
              v-model.number="form.max_uses"
              type="number"
              min="1"
              required
              class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">过期时间（可选）</label>
            <input
              v-model="form.expire_at"
              type="datetime-local"
              class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <div class="flex justify-end gap-3 pt-4">
            <button
              type="button"
              @click="showModal = false"
              class="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded"
            >
              取消
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              创建
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
