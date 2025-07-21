package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Db       DbCfg
	Kafka    Kafka
	External External
}

type Kafka struct {
	KafkaServer          string `env-default:"172.30.1.81:9092"`
	KafkaGroupId         string ``
	KafkaOffsetResetType string ``
}

type DbCfg struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Option   string
}

type External struct {
	BaseUrl1C string
}

func MustInit(configPath string) *Config {
	godotenv.Load(configPath)

	return &Config{
		Kafka: Kafka{
			KafkaServer:          MustGetEnv("KAFKA_SERVER"),
			KafkaGroupId:         MustGetEnv("KAFKA_GROUP_ID"),
			KafkaOffsetResetType: MustGetEnv("KAFKA_OFFSET_RESET_TYPE"),
		},
		Db: DbCfg{
			Driver:   MustGetEnv("PG_DRIVER"),
			Host:     MustGetEnv("PG_HOST"),
			Port:     MustGetEnv("PG_PORT"),
			Name:     MustGetEnv("PG_NAME"),
			User:     MustGetEnv("PG_USER"),
			Password: MustGetEnv("PG_PASSWORD"),
			Option:   MustGetEnv("PG_OPTION"),
		},
		External: External{
			BaseUrl1C: MustGetEnv("BASE_URL_1C"),
		},
	}
}

func PathDefault(workDir string, filename *string) string {
	if filename == nil {
		return filepath.Join(workDir, ".env")
	}

	return filepath.Join(workDir, *filename)
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("no variable in env: %s", key)
	}
	return value
}

func MustGetEnvAsInt(name string) int {
	valueStr := MustGetEnv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return -1
}

func ParseConfigPathFromCl(currentDir string) string {
	var configPath string
	flag.StringVar(&configPath, "config", PathDefault(currentDir, nil), "path to config file")
	flag.Parse()

	return configPath
}
