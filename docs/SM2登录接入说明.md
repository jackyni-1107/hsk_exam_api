# SM2 登录接入说明

## 目标

- 登录接口统一使用 `password` 字段传 SM2 密文。
- 服务端使用 SM2 私钥解密后，再复用原有密码哈希校验流程。

## 接口

- 管理端公钥：`GET /api/admin/auth/public-key`
- 客户端公钥：`GET /api/client/auth/public-key`
- 管理端登录：`POST /api/admin/auth/login`
- 客户端登录：`POST /api/client/auth/login`

## 公钥响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "public_key_hex": "04....",
    "algorithm": "sm2",
    "cipher_mode": "c1c3c2"
  }
}
```

## 登录请求参数

- `username`: 用户名
- `password`: SM2 密文（必填，hex 或 base64）

## 前端加密约定

- 算法：SM2
- 模式：`c1c3c2`
- 编码：十六进制字符串（推荐），服务端兼容 Base64
- JS 示例：`sm-crypto` 的 `sm2.doEncrypt(plain, publicKey, 1)`

## 注意事项

- SM2 仅作为应用层二次保护，生产环境仍必须强制 HTTPS。
- 私钥仅保存在服务端配置，不得写入前端或代码仓库。
- 密文格式错误或解密失败时，服务端返回统一登录失败语义错误。
