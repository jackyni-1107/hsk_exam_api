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
	// S3ForcePathStyle 为 true 时预签名为 path-style（路径含 /bucket/）；为 false 时为 virtual-hosted（路径仅为对象 key，桶在主机名）。
	S3ForcePathStyle bool
}
