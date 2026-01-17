package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env       string
	Port      string
	DBPath    string
	JWTSecret string

	// MinIO 配置
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool

	// 短信服务配置
	SMSProvider   string // aliyun | tencent
	SMSAccessKey  string
	SMSSecretKey  string
	SMSSignName   string
	SMSTemplateID string
}

func Load() (*Config, error) {
	cfg := &Config{
		Env:       getEnv("ENV", "development"),
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "./data/car4race.db"),
		JWTSecret: getEnv("JWT_SECRET", "car4race-dev-secret-key"),

		// MinIO 配置
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "car4race"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "car4race123"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "car4race"),
		MinIOUseSSL:    getEnvBool("MINIO_USE_SSL", false),

		SMSProvider:   getEnv("SMS_PROVIDER", "aliyun"),
		SMSAccessKey:  getEnv("SMS_ACCESS_KEY", ""),
		SMSSecretKey:  getEnv("SMS_SECRET_KEY", ""),
		SMSSignName:   getEnv("SMS_SIGN_NAME", "Car4Race"),
		SMSTemplateID: getEnv("SMS_TEMPLATE_ID", ""),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return b
	}
	return defaultValue
}
