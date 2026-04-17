package v1

import "github.com/gogf/gf/v2/frame/g"

type PublicKeyReq struct {
	g.Meta `path:"/auth/public-key" method:"get" tags:"客户端认证" summary:"获取登录 SM2 公钥"`
}

type PublicKeyRes struct {
	PublicKeyHex string `json:"public_key_hex" dc:"SM2 公钥（hex）"`
	Algorithm    string `json:"algorithm" dc:"加密算法（固定 sm2）"`
	CipherMode   string `json:"cipher_mode" dc:"密文模式（固定 c1c3c2）"`
}
