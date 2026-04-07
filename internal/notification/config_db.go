package notification

import (
	"context"
	"encoding/json"

	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"
)

// SMTPConfig 闁喕娆?SMTP 闁板秶鐤
type SMTPConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	From string `json:"from"`
}

// AliyunSMSConfig 闂冨潡鍣锋禍鎴犵叚娣囷繝鍘ら敓
type AliyunSMSConfig struct {
	AccessKey    string `json:"access_key"`
	SecretKey    string `json:"secret_key"`
	SignName     string `json:"sign_name"`
	TemplateCode string `json:"template_code"`
}

// TencentSMSConfig 閼垫崘顔嗘禍鎴犵叚娣囷繝鍘ら敓
type TencentSMSConfig struct {
	SecretId   string `json:"secret_id"`
	SecretKey  string `json:"secret_key"`
	SdkAppId   string `json:"sdk_app_id"`
	SignName   string `json:"sign_name"`
	TemplateId string `json:"template_id"`
}

// GetActiveEmailConfig 閼惧嘲褰囪ぐ鎾冲閸氼垳鏁ら惃鍕仏娴犲爼鍘ら敓
func GetActiveEmailConfig(ctx context.Context) (provider string, cfg interface{}, ok bool) {
	var e sysentity.SysNotificationChannelConfig
	err := sysdao.SysNotificationChannelConfig.Ctx(ctx).
		Where("channel", "email").
		Where("is_active", 1).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&e)
	if err != nil || e.Id == 0 {
		return "", nil, false
	}
	if e.Provider == "smtp" && e.ConfigJson != "" {
		var c SMTPConfig
		if json.Unmarshal([]byte(e.ConfigJson), &c) == nil {
			return "smtp", &c, true
		}
	}
	return e.Provider, nil, false
}

// GetActiveSMSConfig 閼惧嘲褰囪ぐ鎾冲閸氼垳鏁ら惃鍕叚娣囷繝鍘ら敓
func GetActiveSMSConfig(ctx context.Context) (provider string, cfg interface{}, ok bool) {
	var e sysentity.SysNotificationChannelConfig
	err := sysdao.SysNotificationChannelConfig.Ctx(ctx).
		Where("channel", "sms").
		Where("is_active", 1).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&e)
	if err != nil || e.Id == 0 {
		return "", nil, false
	}
	if e.ConfigJson == "" {
		return e.Provider, nil, false
	}
	switch e.Provider {
	case "aliyun":
		var c AliyunSMSConfig
		if json.Unmarshal([]byte(e.ConfigJson), &c) == nil {
			return "aliyun", &c, true
		}
	case "tencent":
		var c TencentSMSConfig
		if json.Unmarshal([]byte(e.ConfigJson), &c) == nil {
			return "tencent", &c, true
		}
	}
	return e.Provider, nil, false
}
