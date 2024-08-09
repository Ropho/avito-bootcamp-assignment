package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	log "github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Init config from given file path
func Init(configPath, configName string) error {

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	return viper.ReadInConfig()
}

// ServiceCONF is the Main config struct
type ServiceCONF struct {
	*LoggerConfig
	*PostgresConfig
	JWTServiceConf

	ServerPort int
}

type PostgresConfig struct {
	URL string
}

type LoggerConfig struct {
	LogLevel  log.Level
	Directory string
	File      *os.File
}

// JWTServiceConf describes config for JWT generation service
type JWTServiceConf struct {
	AccessJWTSecret string `env:"JWTSECRET" envDefault:"JWTSECRET"`
	AccessJWTExp    string `env:"JWTEXP" envDefault:"10"`
}

func getLoggerCONF(rawLoggerConfig map[string]interface{}) (*LoggerConfig, error) {
	cfg := &LoggerConfig{}
	var (
		ok  bool
		err error
	)

	var logLevel string
	if logLevel, ok = rawLoggerConfig["log_level"].(string); !ok {
		return nil, errors.New("bad service log_level in config")
	}
	if cfg.LogLevel, err = log.ParseLevel(logLevel); err != nil {
		return nil, fmt.Errorf("bad log level: [%w]", err)
	}

	if cfg.Directory, ok = rawLoggerConfig["dir"].(string); !ok {
		return nil, fmt.Errorf("failed to parse directory")
	}
	err = os.MkdirAll(cfg.Directory, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger directory: [%w]", err)
	}

	var filePath string
	if filePath, ok = rawLoggerConfig["file"].(string); !ok {
		return nil, fmt.Errorf("failed to parse file")
	}
	cfg.File, err = os.OpenFile(fmt.Sprintf("%s/%s", cfg.Directory, filePath), os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for logging: [%w]", err)
	}

	return cfg, nil
}

func getPostgresCONF(rawPostgresConfig map[string]interface{}) (*PostgresConfig, error) {
	cfg := &PostgresConfig{}
	var (
		ok bool
	)

	if cfg.URL, ok = rawPostgresConfig["url"].(string); !ok {
		return nil, errors.New("bad postgres url in config")
	}

	return cfg, nil
}

// GetServiceCONF fills up the config struct
func GetServiceCONF() (*ServiceCONF, error) {
	serviceConfig := &ServiceCONF{}
	var err error

	serviceConfig.ServerPort = viper.GetInt("service.port")

	rawPostgresConfig := viper.GetViper().GetStringMap("service.postgres")
	if rawPostgresConfig == nil {
		return nil, errors.New("no postgres configuration specified")
	}
	serviceConfig.PostgresConfig, err = getPostgresCONF(rawPostgresConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: [%w]", err)
	}

	// Logger configuration
	rawLoggerConfig := viper.GetViper().GetStringMap("service.logger")
	if rawLoggerConfig == nil {
		return nil, errors.New("no logger configuration specified")
	}
	serviceConfig.LoggerConfig, err = getLoggerCONF(rawLoggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse logger config: [%w]", err)
	}

	err = serviceConfig.getEnvServiceConfig()
	if err != nil {
		return nil, err
	}

	return serviceConfig, nil
}

func (s *ServiceCONF) getEnvServiceConfig() error {

	err := env.Parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse env variables: [%w]", err)
	}

	if s.JWTServiceConf.AccessJWTSecret == "" || s.JWTServiceConf.AccessJWTExp == "" {
		return fmt.Errorf("empty secrets for JWT")
	}

	return nil
}
