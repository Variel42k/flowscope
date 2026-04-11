import axios from 'axios'

import { AuthSession } from '../lib/types'

const baseURL = import.meta.env.VITE_API_BASE ?? '/api'

export const apiClient = axios.create({
  baseURL,
  timeout: 30000,
})

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('flowscope_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error?.response?.status
    const requestUrl = String(error?.config?.url ?? '')
    if (status === 401 && !requestUrl.includes('/api/auth/login')) {
      clearAuthSession()
    }
    return Promise.reject(error)
  },
)

export async function login(username: string, password: string) {
  const { data } = await apiClient.post<AuthSession>('/api/auth/login', { username, password })
  setAuthSession(data)
  return data
}

export function logout() {
  clearAuthSession()
}

export function isLoggedIn() {
  return Boolean(localStorage.getItem('flowscope_token'))
}

export function getCurrentUser() {
  return localStorage.getItem('flowscope_user') ?? ''
}

export function getCurrentRole() {
  return localStorage.getItem('flowscope_role') ?? 'viewer'
}

export async function startOIDCLogin() {
  const { data } = await apiClient.get<{ authorize_url: string }>('/api/auth/oidc/start')
  if (data?.authorize_url) {
    window.location.href = data.authorize_url
  }
}

export function applyAuthFromQuery(search: string) {
  const params = new URLSearchParams(search)
  const token = params.get('token')
  const user = params.get('user')
  const role = params.get('role') || 'viewer'
  if (!token || !user) return null
  setAuthSession({ token, user, role, auth_mode: 'oidc' })
  return { token, user, role }
}

export async function fetchMe() {
  const { data } = await apiClient.get<{ user: string; role: string }>('/api/auth/me')
  if (data?.user) {
    localStorage.setItem('flowscope_user', data.user)
  }
  if (data?.role) {
    localStorage.setItem('flowscope_role', data.role)
  }
  return data
}

function setAuthSession(data: AuthSession) {
  localStorage.setItem('flowscope_token', data.token)
  localStorage.setItem('flowscope_user', data.user)
  localStorage.setItem('flowscope_role', data.role || 'viewer')
}

function clearAuthSession() {
  localStorage.removeItem('flowscope_token')
  localStorage.removeItem('flowscope_user')
  localStorage.removeItem('flowscope_role')
}
