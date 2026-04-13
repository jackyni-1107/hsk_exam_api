package storage

import (
	"context"
	"encoding/json"

	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"
)

// jsonConfig 对应表字段 config_json。
type jsonConfig struct {
	BasePath                string `json:"base_path"`
	Endpoint                string `json:"endpoint"`
	Bucket                  string `json:"bucket"`
	AccessKey               string `json:"access_key"`
	SecretKey               string `json:"secret_key"`
	Region                  string `json:"region"`
	PublicBaseURL           string `json:"public_base_url"`
	PresignSignatureVersion string `json:"presign_signature_version"`
	// S3ForcePathStyle 显式为 true 时使用 path-style；缺省或 false 使用 virtual-hosted（预签名路径不含桶前缀）。
	S3ForcePathStyle *bool `json:"s3_force_path_style"`
}

// GetActiveConfig 读取当前启用的存储配置，供 NewAdapter 使用。
func GetActiveConfig(ctx context.Context) (cfg Config, cleanupDays int) {
	cfg = Config{Type: "local", BasePath: "./storage", Bucket: "default"}
	cleanupDays = 30
	var e sysentity.SysFileStorageConfig
	err := sysdao.SysFileStorageConfig.Ctx(ctx).
		Where("is_active", 1).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&e)
	if err != nil || e.Id == 0 {
		return cfg, cleanupDays
	}
	cfg.Type = e.StorageType
	if e.CleanupBeforeDays > 0 {
		cleanupDays = e.CleanupBeforeDays
	}
	if e.ConfigJson != "" {
		var jc jsonConfig
		if json.Unmarshal([]byte(e.ConfigJson), &jc) == nil {
			cfg.BasePath = jc.BasePath
			if cfg.BasePath == "" {
				cfg.BasePath = "./storage"
			}
			cfg.Endpoint = jc.Endpoint
			cfg.Bucket = jc.Bucket
			if cfg.Bucket == "" {
				cfg.Bucket = "default"
			}
			cfg.AccessKey = jc.AccessKey
			cfg.SecretKey = jc.SecretKey
			cfg.Region = jc.Region
			cfg.PublicBaseURL = jc.PublicBaseURL
			cfg.PresignSignatureVersion = normalizePresignSigVersion(ctx, jc.PresignSignatureVersion)
			// 缺省 true：与历史行为一致（path-style，URL 路径含 /bucket/）；显式 false 为 virtual-hosted。
			cfg.S3ForcePathStyle = true
			if jc.S3ForcePathStyle != nil {
				cfg.S3ForcePathStyle = *jc.S3ForcePathStyle
			}
		}
	}
	return cfg, cleanupDays
}
