package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type UploadConfig struct {
	Path         string   `mapstructure:"path"`
	MaxSizeMB    int      `mapstructure:"max_size_mb"`
	AllowedTypes []string `mapstructure:"allowed_types"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Docker 环境变量覆盖（大写 + 下划线）
	bindEnvOverrides()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func bindEnvOverrides() {
	overrides := map[string]string{
		"database.host":     "DATABASE_HOST",
		"database.port":     "DATABASE_PORT",
		"database.user":     "DATABASE_USER",
		"database.password": "DATABASE_PASSWORD",
		"database.dbname":   "DATABASE_DBNAME",
		"redis.host":        "REDIS_HOST",
		"redis.port":        "REDIS_PORT",
		"redis.password":    "REDIS_PASSWORD",
		"redis.db":          "REDIS_DB",
		"server.port":       "SERVER_PORT",
		"jwt.secret":        "JWT_SECRET",
	}
	for key, envKey := range overrides {
		if v := os.Getenv(envKey); v != "" {
			viper.Set(key, v)
		}
	}
}
