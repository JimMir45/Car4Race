<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '../../api'

interface Course {
  id: number
  title: string
  slug: string
  description: string
  cover_image: string
  price: number
  orig_price: number
  sales_count: number
  is_public: boolean
  sort: number
  created_at: string
}

interface CourseFile {
  id: number
  course_id: number
  file_type: string
  file_name: string
  file_size: number
  sort: number
  created_at: string
}

const courses = ref<Course[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)

// 模态框状态
const showModal = ref(false)
const showFilesModal = ref(false)
const editingCourse = ref<Course | null>(null)
const currentCourseFiles = ref<CourseFile[]>([])
const currentCourseId = ref<number | null>(null)
const uploading = ref(false)

// 表单数据
const form = ref({
  title: '',
  slug: '',
  description: '',
  cover_image: '',
  price: 0,
  orig_price: 0,
  is_public: true,
  sort: 0,
})

async function fetchCourses() {
  loading.value = true
  try {
    const res: any = await adminApi.getCourses({ page: page.value, page_size: 20 })
    courses.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error('获取课程失败', e)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingCourse.value = null
  form.value = {
    title: '',
    slug: '',
    description: '',
    cover_image: '',
    price: 0,
    orig_price: 0,
    is_public: true,
    sort: 0,
  }
  showModal.value = true
}

function openEditModal(course: Course) {
  editingCourse.value = course
  form.value = {
    title: course.title,
    slug: course.slug,
    description: course.description,
    cover_image: course.cover_image,
    price: course.price,
    orig_price: course.orig_price,
    is_public: course.is_public,
    sort: course.sort,
  }
  showModal.value = true
}

async function saveCourse() {
  try {
    if (editingCourse.value) {
      await adminApi.updateCourse(editingCourse.value.id, form.value)
    } else {
      await adminApi.createCourse(form.value)
    }
    showModal.value = false
    fetchCourses()
  } catch (e: any) {
    alert(e.message || '保存失败')
  }
}

async function deleteCourse(course: Course) {
  if (!confirm(`确定要删除课程 "${course.title}" 吗？`)) return
  try {
    await adminApi.deleteCourse(course.id)
    fetchCourses()
  } catch (e: any) {
    alert(e.message || '删除失败')
  }
}

async function openFilesModal(course: Course) {
  currentCourseId.value = course.id
  showFilesModal.value = true
  await fetchCourseFiles(course.id)
}

async function fetchCourseFiles(courseId: number) {
  try {
    const res: any = await adminApi.getCourseFiles(courseId)
    currentCourseFiles.value = res.data || []
  } catch (e) {
    console.error('获取文件列表失败', e)
    currentCourseFiles.value = []
  }
}

async function handleFileUpload(event: Event, fileType: 'intro' | 'resource') {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file || !currentCourseId.value) return

  uploading.value = true
  try {
    await adminApi.uploadCourseFile(currentCourseId.value, file, fileType)
    await fetchCourseFiles(currentCourseId.value)
  } catch (e: any) {
    alert(e.message || '上传失败')
  } finally {
    uploading.value = false
    target.value = ''
  }
}

async function deleteFile(fileId: number) {
  if (!currentCourseId.value || !confirm('确定要删除该文件吗？')) return
  try {
    await adminApi.deleteCourseFile(currentCourseId.value, fileId)
    await fetchCourseFiles(currentCourseId.value)
  } catch (e: any) {
    alert(e.message || '删除失败')
  }
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

onMounted(() => {
  fetchCourses()
})
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">课程管理</h1>
      <button
        @click="openCreateModal"
        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        创建课程
      </button>
    </div>

    <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-700">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">标题</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">价格</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">销量</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">状态</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-for="course in courses" :key="course.id" class="hover:bg-gray-50 dark:hover:bg-gray-700">
            <td class="px-6 py-4">
              <div class="text-sm font-medium text-gray-900 dark:text-white">{{ course.title }}</div>
              <div class="text-sm text-gray-500 dark:text-gray-400">{{ course.slug }}</div>
            </td>
            <td class="px-6 py-4">
              <span class="text-sm text-gray-900 dark:text-white">¥{{ course.price }}</span>
              <span v-if="course.orig_price > course.price" class="ml-2 text-xs text-gray-400 line-through">
                ¥{{ course.orig_price }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">{{ course.sales_count }}</td>
            <td class="px-6 py-4">
              <span
                :class="course.is_public
                  ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
                  : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'"
                class="px-2 py-1 text-xs rounded-full"
              >
                {{ course.is_public ? '公开' : '隐藏' }}
              </span>
            </td>
            <td class="px-6 py-4 text-right text-sm font-medium space-x-2">
              <button @click="openFilesModal(course)" class="text-green-600 hover:text-green-900 dark:text-green-400">
                文件
              </button>
              <button @click="openEditModal(course)" class="text-blue-600 hover:text-blue-900 dark:text-blue-400">
                编辑
              </button>
              <button @click="deleteCourse(course)" class="text-red-600 hover:text-red-900 dark:text-red-400">
                删除
              </button>
            </td>
          </tr>
          <tr v-if="courses.length === 0 && !loading">
            <td colspan="5" class="px-6 py-12 text-center text-gray-500">暂无课程</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 创建/编辑课程模态框 -->
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div class="px-6 py-4 border-b dark:border-gray-700">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">
            {{ editingCourse ? '编辑课程' : '创建课程' }}
          </h3>
        </div>
        <form @submit.prevent="saveCourse" class="px-6 py-4 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">标题</label>
            <input
              v-model="form.title"
              type="text"
              required
              class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Slug</label>
            <input
              v-model="form.slug"
              type="text"
              required
              class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">描述</label>
            <textarea
              v-model="form.description"
              rows="3"
              class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            ></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">封面图片URL</label>
            <input
              v-model="form.cover_image"
              type="text"
              class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">价格</label>
              <input
                v-model.number="form.price"
                type="number"
                step="0.01"
                required
                class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">原价</label>
              <input
                v-model.number="form.orig_price"
                type="number"
                step="0.01"
                class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">排序</label>
              <input
                v-model.number="form.sort"
                type="number"
                class="w-full px-3 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              />
            </div>
            <div class="flex items-center pt-6">
              <label class="flex items-center">
                <input v-model="form.is_public" type="checkbox" class="mr-2" />
                <span class="text-sm text-gray-700 dark:text-gray-300">公开</span>
              </label>
            </div>
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
              保存
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 文件管理模态框 -->
    <div v-if="showFilesModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div class="px-6 py-4 border-b dark:border-gray-700 flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">课程文件管理</h3>
          <button @click="showFilesModal = false" class="text-gray-500 hover:text-gray-700">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="px-6 py-4">
          <!-- 上传区域 -->
          <div class="mb-6 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                上传课程介绍 (Markdown)
              </label>
              <input
                type="file"
                accept=".md,.markdown"
                @change="(e) => handleFileUpload(e, 'intro')"
                :disabled="uploading"
                class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                上传课程资源 (ZIP)
              </label>
              <input
                type="file"
                accept=".zip,.rar,.7z"
                @change="(e) => handleFileUpload(e, 'resource')"
                :disabled="uploading"
                class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-green-50 file:text-green-700 hover:file:bg-green-100"
              />
            </div>
            <div v-if="uploading" class="text-sm text-blue-600">上传中...</div>
          </div>

          <!-- 文件列表 -->
          <div class="space-y-2">
            <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300">已上传文件</h4>
            <div v-if="currentCourseFiles.length === 0" class="text-sm text-gray-500 py-4 text-center">
              暂无文件
            </div>
            <div
              v-for="file in currentCourseFiles"
              :key="file.id"
              class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded"
            >
              <div class="flex items-center gap-3">
                <span
                  :class="file.file_type === 'intro'
                    ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300'
                    : 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'"
                  class="px-2 py-1 text-xs rounded"
                >
                  {{ file.file_type === 'intro' ? '介绍' : '资源' }}
                </span>
                <span class="text-sm text-gray-900 dark:text-white">{{ file.file_name }}</span>
                <span class="text-xs text-gray-500">{{ formatFileSize(file.file_size) }}</span>
              </div>
              <button
                @click="deleteFile(file.id)"
                class="text-red-600 hover:text-red-800 text-sm"
              >
                删除
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
