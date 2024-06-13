package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

var cfg config

func Config() config {
	return cfg
}

type config struct {
	EnvType string        `yaml:"env_type"`
	Storage storageConfig `yaml:"storage"`
	Server  serverConfig  `yaml:"server"`
	JWT     jwtConfig     `yaml:"jwt"`
}

type storageConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type serverConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type jwtConfig struct {
	TokenTTL   time.Duration   `yaml:"token_ttl"`
	PrivateKey *rsa.PrivateKey `yaml:"-"`
	PublicKey  *rsa.PublicKey  `yaml:"-"`
}

func LoadConfig() {
	envType := getEnvType()
	path := getConfigFilePath(envType)
	cleanenv.ReadConfig(path, &cfg)
	readKeys()
}

func getConfigFilePath(envType string) string {
	path := fmt.Sprintf("./config/%s.yaml", envType)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("%s file not found", path)
	}
	return path
}

func getEnvType() string {
	envType := os.Getenv("ENV_TYPE")
	if envType == "" {
		log.Fatal("Empty ENV_TYPE variable")
	}
	if envType != EnvProd {
		log.Printf("!!! Using %s env type. Not for production !!!", envType)
		log.Printf("!!! Using %s env type. Not for production !!!", envType)
		log.Printf("!!! Using %s env type. Not for production !!!", envType)
	}
	return envType
}

func readKeys() {
	// Read private
	data, err := os.ReadFile("keys/private.pem")
	if err != nil {
		log.Fatal(err)
	}
	cfg.JWT.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		log.Fatal(err)
	}

	// Read public
	data, err = os.ReadFile("keys/public.pem")
	if err != nil {
		log.Fatal(err)
	}
	cfg.JWT.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		log.Fatal(err)
	}
}
