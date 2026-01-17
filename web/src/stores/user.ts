import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, userApi } from '../api'

interface User {
  id: number
  phone: string
  username: string
  nickname: string
  avatar: string
  role: string
}

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))

  const isLoggedIn = computed(() => !!token.value)
  const isVIP = computed(() => user.value?.role === 'vip')
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function sendCode(phone: string) {
    await authApi.sendCode(phone)
  }

  async function login(phone: string, code: string) {
    const res: any = await authApi.login(phone, code)
    token.value = res.data.token
    user.value = res.data.user
    localStorage.setItem('token', res.data.token)
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const res: any = await userApi.getProfile()
      user.value = res.data
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  return {
    user,
    token,
    isLoggedIn,
    isVIP,
    isAdmin,
    sendCode,
    login,
    fetchProfile,
    logout,
  }
})
