import request from './request'

export interface LoginUser {
  id: number
  username: string
  nickname: string
  avatar: string
  roles?: string[]
}

export function login(payload: {
  username: string
  password: string
  captcha_id?: string
  captcha_answer?: string
}) {
  return request.post<
    unknown,
    { data?: { token?: string; user_info?: LoginUser } }
  >('/admin/auth/login', payload)
}

export function logout() {
  return request.post<unknown, unknown>('/admin/auth/logout')
}
