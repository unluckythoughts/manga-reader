import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api/auth'
import type { User } from '../types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('auth_token'))

  const isLoggedIn = computed(() => !!token.value)

  async function login(email: string, password: string) {
    const res = await authApi.login(email, password)
    token.value = res.token
    // The login response wraps the user under `data.user` via the envelope, but
    // the auth handler returns { token } only — fetch the full user separately.
    localStorage.setItem('auth_token', res.token)
    await fetchUser()
  }

  async function fetchUser() {
    if (!token.value) return
    try {
      user.value = await authApi.getUser()
    } catch {
      // Token is invalid or expired — clear it
      token.value = null
      user.value = null
      localStorage.removeItem('auth_token')
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('auth_token')
  }

  return { user, token, isLoggedIn, login, fetchUser, logout }
})
