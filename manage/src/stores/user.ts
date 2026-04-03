import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { LoginUser } from '@/api/auth'

const TOKEN_KEY = 'admin_token'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem(TOKEN_KEY) || '')
  const userInfo = ref<LoginUser | null>(null)

  function setSession(t: string, info: LoginUser | null) {
    token.value = t
    userInfo.value = info
    if (t) {
      localStorage.setItem(TOKEN_KEY, t)
    } else {
      localStorage.removeItem(TOKEN_KEY)
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem(TOKEN_KEY)
  }

  return { token, userInfo, setSession, logout }
})
