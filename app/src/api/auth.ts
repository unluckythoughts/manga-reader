import { apiClient } from './client'
import type { LoginResponse, User } from '../types'

export const authApi = {
  login: (email: string, password: string) =>
    apiClient.post<LoginResponse>('/auth/login', { email, password }),
  register: (name: string, email: string, password: string) =>
    apiClient.post<User>('/auth/register', { name, email, password }),
  getUser: () => apiClient.get<User>('/user')
}
