package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type LambdaConfiguration struct {
	General struct {
		Environment  string `env:"ENVIRONMENT" env-default:"local" env-description:"Lambda Environment"`
		Region       string `env:"AWS_REGION" env-default:"us-east-1" env-description:"AWS Region"`
		Country      string `env:"COUNTRY" env-default:"ARG" env-description:"ARG MEX COL (MULTI)" `
		FunctionName string `env:"AWS_LAMBDA_FUNCTION_NAME" env-default:"lambdaFunc" env-description:"Lambda Name" `
	}
}

func NewConfigLambda() *LambdaConfiguration {
	cfg := &LambdaConfiguration{}

	return cfg
}

func (cfg *LambdaConfiguration) LoadFromEnvs() {
	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(err)
	}
}

func (cfg *LambdaConfiguration) GetEnvsDescriptions() string {
	header := "Environment variables"
	help, _ := cleanenv.GetDescription(cfg, &header)

	return help
}
