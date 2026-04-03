package storage

// Config 存储配置（来自 DB 活动项 + config_json）。
type Config struct {
	Type                    string
	BasePath                string
	Endpoint                string
	Bucket                  string
	AccessKey               string
	SecretKey               string
	Region                  string
	PublicBaseURL           string
	PresignSignatureVersion string
}
