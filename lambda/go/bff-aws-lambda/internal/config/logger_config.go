package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// LoggerConfiguration Struct for load envs.
type LoggerConfiguration struct {
	General struct {
		Environment string `env:"ENVIRONMENT" env-default:"local" env-description:"Lambda Environment"`
		Region      string `env:"AWS_REGION" env-default:"us-east-1" env-description:"AWS Region"`
		LogLevel    string `env:"LOG_LEVEL" env-default:"INFO" env-description:"Log Level (DEBUG,INFO)"`
		Flow        string `env:"FLOW" env-default:"flow" env-description:"Function Flow"`
		Country     string `env:"COUNTRY" env-default:"ARG" env-description:"ARG MEX COL (MULTI)"`
		Version     string `env:"AWS_LAMBDA_FUNCTION_VERSION" env-default:"1" env-description:"AWS Function Version"`
		Name        string `env:"AWS_LAMBDA_FUNCTION_NAME" env-default:"lambda" env-description:"AWS Function Name"`
	}
}

// NewLoggerConfiguration Return filled LoggerConfiguration.
func NewLoggerConfiguration() *LoggerConfiguration {
	cfg := &LoggerConfiguration{}
	cfg.LoadFromEnvs()

	return cfg
}

func (cfg *LoggerConfiguration) LoadFromEnvs() {
	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(err)
	}
}

// GetEnvsDescriptions get envs description.
func (cfg *LoggerConfiguration) GetEnvsDescriptions() string {
	header := "Environment variables"
	help, _ := cleanenv.GetDescription(cfg, &header) //nolint: errcheck

	return help
}
