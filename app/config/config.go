package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

const (
	TestModeValue = "test"

	appExternalPort  = "APP_EXTERNAL_PORT"
	handlerMode      = "HANDLER_MODE"
	postgresHost     = "POSTGRES_HOST"
	postgresPort     = "POSTGRES_PORT"
	postgresUser     = "POSTGRES_USER"
	postgresPassword = "POSTGRES_PASSWORD"
	postgresDatabase = "POSTGRES_DATABASE"

	tokenSigningKey = "TOKEN_SIGNING_KEY"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SslMode  string
}

type AuthConfig struct {
	TokenSigningKey string
}

type Config struct {
	AppPort     string
	HandlerMode string
	DbConfig
	AuthConfig
}

func (c *Config) Init(envFilePath string) (*Config, error) {
	reader, err := getReader(envFilePath)
	if err != nil {
		return nil, err
	}
	defer func(reader *os.File) {
		_ = reader.Close()
	}(reader)
	env, err := readEnvFile(reader)
	if err != nil {
		return nil, err
	}

	appPort, err := envValue(env, appExternalPort)
	if err != nil {
		return nil, err
	}

	apiMode, err := envValue(env, handlerMode)
	if err != nil {
		return nil, err
	}

	dbHost, err := envValue(env, postgresHost)
	if err != nil {
		return nil, err
	}

	dbPort, err := envValue(env, postgresPort)
	if err != nil {
		return nil, err
	}

	dbUser, err := envValue(env, postgresUser)
	if err != nil {
		return nil, err
	}

	dbPassword, err := envValue(env, postgresPassword)
	if err != nil {
		return nil, err
	}

	dbName, err := envValue(env, postgresDatabase)
	if err != nil {
		return nil, err
	}

	authTokenSigningKey, err := envValue(env, tokenSigningKey)
	if err != nil {
		return nil, err
	}

	return &Config{
		AppPort:     appPort,
		HandlerMode: apiMode,
		DbConfig: DbConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DbName:   dbName,
			SslMode:  "disable",
		},
		AuthConfig: AuthConfig{
			TokenSigningKey: authTokenSigningKey,
		},
	}, nil
}

func getReader(envFilePath string) (*os.File, error) {
	return os.Open(envFilePath)
}

func readEnvFile(reader *os.File) (map[string]string, error) {
	env, err := godotenv.Parse(reader)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func envValue(env map[string]string, envKey string) (string, error) {
	/*
		value, found := os.LookupEnv(envKey)
		if found && value != "" {
			return value, nil
		}
	*/
	value, found := env[envKey]
	if !found {
		return "", fmt.Errorf("%s not found", envKey)
	}

	if value == "" {
		return "", fmt.Errorf("%s is empty", envKey)
	}

	return value, nil
}
