import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error)
  }
)

// Auth API
export const authApi = {
  sendCode: (phone: string) => api.post('/auth/send-code', { phone }),
  login: (phone: string, code: string) => api.post('/auth/login', { phone, code }),
}

// User API
export const userApi = {
  getProfile: () => api.get('/user/profile'),
  updateProfile: (data: { nickname?: string; avatar?: string }) =>
    api.put('/user/profile', data),
}

// Course API
export const courseApi = {
  getList: (params?: { page?: number; page_size?: number; sort?: string }) =>
    api.get('/hpa/courses', { params }),
  getDetail: (slug: string) => api.get(`/hpa/courses/${slug}`),
}

// Order API
export const orderApi = {
  getList: (params?: { page?: number; page_size?: number }) =>
    api.get('/hpa/orders', { params }),
  create: (courseId: number) => api.post('/hpa/orders', { course_id: courseId }),
  redeemCode: (code: string) => api.post('/hpa/redeem', { code }),
}

// Download API
export const downloadApi = {
  createToken: (courseId: number, fileId?: number) =>
    api.post('/hpa/download', { course_id: courseId, file_id: fileId }),
  download: (token: string) => api.get(`/hpa/download/${token}`),
}

// Admin API
export const adminApi = {
  // 课程管理
  getCourses: (params?: { page?: number; page_size?: number }) =>
    api.get('/admin/courses', { params }),
  createCourse: (data: {
    title: string
    slug: string
    description?: string
    cover_image?: string
    price: number
    orig_price?: number
    is_public?: boolean
    sort?: number
  }) => api.post('/admin/courses', data),
  updateCourse: (id: number, data: {
    title: string
    slug: string
    description?: string
    cover_image?: string
    price: number
    orig_price?: number
    is_public?: boolean
    sort?: number
  }) => api.put(`/admin/courses/${id}`, data),
  deleteCourse: (id: number) => api.delete(`/admin/courses/${id}`),

  // 课程文件管理
  getCourseFiles: (courseId: number) => api.get(`/admin/courses/${courseId}/files`),
  uploadCourseFile: (courseId: number, file: File, fileType: 'intro' | 'resource') => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('file_type', fileType)
    return api.post(`/admin/courses/${courseId}/files`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  deleteCourseFile: (courseId: number, fileId: number) =>
    api.delete(`/admin/courses/${courseId}/files/${fileId}`),

  // 邀请码管理
  getInviteCodes: (params?: { page?: number; page_size?: number }) =>
    api.get('/admin/invite-codes', { params }),
  createInviteCode: (data: {
    course_id: number
    max_uses?: number
    expire_at?: string
  }) => api.post('/admin/invite-codes', data),
}

export default api
