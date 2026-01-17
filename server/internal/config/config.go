package config

import (
	"os"
)

type Config struct {
	Env       string
	Port      string
	DBPath    string
	JWTSecret string

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
