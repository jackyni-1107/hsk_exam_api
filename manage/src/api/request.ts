import axios, { type AxiosInstance, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

const request: AxiosInstance = axios.create({
  baseURL: '/api',
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
  (res: AxiosResponse) => {
    const body = res.data as { code?: number; message?: string; data?: unknown }
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
