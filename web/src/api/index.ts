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
  createToken: (courseId: number) =>
    api.post('/hpa/download', { course_id: courseId }),
  download: (token: string) => api.get(`/hpa/download/${token}`),
}

export default api
