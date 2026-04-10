import axios from 'axios'

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
      localStorage.removeItem('flowscope_token')
      localStorage.removeItem('flowscope_user')
    }
    return Promise.reject(error)
  },
)

export async function login(username: string, password: string) {
  const { data } = await apiClient.post('/api/auth/login', { username, password })
  localStorage.setItem('flowscope_token', data.token)
  localStorage.setItem('flowscope_user', data.user)
  return data
}

export function logout() {
  localStorage.removeItem('flowscope_token')
  localStorage.removeItem('flowscope_user')
}

export function isLoggedIn() {
  return Boolean(localStorage.getItem('flowscope_token'))
}
