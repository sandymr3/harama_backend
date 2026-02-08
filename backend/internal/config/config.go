package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	GeminiAPIKey      string
	MinioEndpoint     string
	MinioAccessKey    string
	MinioSecretKey    string
	MinioBucket       string
	MinioUseSSL       bool
	SupabaseURL       string
	SupabaseJWTSecret string
	CORSOrigin        string
}

func Load() *Config {
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	return &Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://harama:pass@localhost:5432/harama?sslmode=disable"),
		GeminiAPIKey:      getEnv("GEMINI_API_KEY", ""),
		MinioEndpoint:     getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:    getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:    getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioBucket:       getEnv("MINIO_BUCKET", "harama"),
		MinioUseSSL:       getEnv("MINIO_USE_SSL", "false") == "true",
		SupabaseURL:       getEnv("SUPABASE_URL", ""),
		SupabaseJWTSecret: getEnv("SUPABASE_JWT_SECRET", ""),
		CORSOrigin:        getEnv("CORS_ORIGIN", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
