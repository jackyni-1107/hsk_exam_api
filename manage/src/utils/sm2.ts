import { sm2 } from 'sm-crypto'

export function encryptPasswordWithSM2(password: string, publicKeyHex: string): string {
  const key = publicKeyHex.trim()
  if (!key) {
    throw new Error('SM2 公钥为空')
  }
  const noPrefix = key.replace(/^04/i, '')
  try {
    return sm2.doEncrypt(password, noPrefix, 1)
  } catch {
    // 部分实现要求带 04 前缀，兼容重试一次。
    return sm2.doEncrypt(password, key, 1)
  }
}
