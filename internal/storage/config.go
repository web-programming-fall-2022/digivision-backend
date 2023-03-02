package storage

import (
	"fmt"
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name" yaml:"db_name"`
}

func (config *DBConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		config.Host, config.Port,
		config.Username, config.DBName, config.Password)
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultVal
	}
	return val
}

func getTestConfig() *DBConfig {
	port, err := strconv.Atoi(getEnv("TEST_DB_PORT", "5432"))
	if err != nil {
		panic(err)
	}
	return &DBConfig{
		Host:     getEnv("TEST_DB_HOST", "127.0.0.1"),
		Port:     port,
		Username: getEnv("TEST_DB_USER", "test_user"),
		Password: getEnv("TEST_DB_PASSWORD", "123"),
		DBName:   getEnv("TEST_DB_NAME", "postgres"),
	}
}
