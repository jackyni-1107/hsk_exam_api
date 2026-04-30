import axios, { type AxiosInstance, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

type ApiEnvelope<T = unknown> = {
  code?: number
  message?: string
  data?: T
}

const request: AxiosInstance = axios.create({
  baseURL: (import.meta.env.VITE_API_BASE_URL || '/api') as string,
  timeout: 60000,
})

request.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('admin_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (res: AxiosResponse): any => {
    const body = res.data as ApiEnvelope
    if (body && typeof body.code === 'number' && body.code !== 0) {
      const msg = body.message || '请求失败'
      ElMessage.error(msg)
      return Promise.reject(new Error(msg))
    }
    return body
  },
  (err) => {
    const msg = err?.response?.data?.message || err?.message || '网络错误'
    ElMessage.error(msg)
    return Promise.reject(err)
  }
)

export default request
