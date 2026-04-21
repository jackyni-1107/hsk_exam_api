import request from './request'

export interface LoginUser {
  id: number
  username: string
  nickname: string
  avatar: string
  roles?: string[]
  /** 与 sys_menu.permission 一致，用于按钮显隐（如 exam:paper:purge） */
  permissions?: string[]
}

export interface AuthPublicKey {
  public_key_hex: string
  algorithm: string
  cipher_mode: string
}

export function getLoginPublicKey() {
  return request.get<unknown, { data?: AuthPublicKey }>('/admin/auth/public-key')
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
