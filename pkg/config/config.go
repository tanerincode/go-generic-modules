package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	AppConfig     *Configuration
	viperInstance *viper.Viper
)

func Init(configName string) error {
	var once sync.Once
	var initErr error
	once.Do(func() {
		viperInstance = viper.New()
		viperInstance.SetConfigName(configName)
		viperInstance.SetConfigType("yaml")

		// Set the prefix that Viper looks for with environment variables
		// and tell Viper to read environment variables that match the config keys.
		viperInstance.SetEnvPrefix("APP") // It will look for env variables with prefix "APP"
		viperInstance.AutomaticEnv()
		viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %s", err)
		}

		configDir := filepath.Join(cwd, "./configs")
		viperInstance.AddConfigPath(configDir)

		if err := viperInstance.ReadInConfig(); err == nil {
			log.Printf("Using config file: %s", viperInstance.ConfigFileUsed())
		} else {
			initErr = fmt.Errorf("fatal error config file: %s", err)
			return
		}

		AppConfig = &Configuration{}
		if err := viperInstance.Unmarshal(AppConfig); err != nil {
			initErr = fmt.Errorf("unable to unmarshal into struct, %v", err)
			return
		}
	})
	return initErr
}

func ResyncEnv() error {
	// Refresh the environment variables
	viperInstance.AutomaticEnv()

	// Re-unmarshal the environment variables
	if err := viperInstance.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("unable to re-unmarshal env vars into config: %w", err)
	}

	return nil
}

func GetConfig(key string) interface{} {
	return viperInstance.Get(key)
}

func SetConfig(key string, value interface{}) {
	viperInstance.Set(key, value)
}

type Configuration struct {
	Server struct {
		Host    string `mapstructure:"host"`
		Port    int    `mapstructure:"port"`
		Timeout struct {
			Read  int `mapstructure:"read"`
			Write int `mapstructure:"write"`
			Idle  int `mapstructure:"idle"`
		} `mapstructure:"timeout"`
	} `mapstructure:"server"`
	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"logging"`
	Database struct {
		Type            string `mapstructure:"type"`
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		Name            string `mapstructure:"name"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"database"`
	Auth struct {
		JwtSecret   string `mapstructure:"jwt_secret"`
		TokenExpiry int    `mapstructure:"token_expiry"`
	} `mapstructure:"auth"`
	CORS struct {
		AllowedOrigins   []string `mapstructure:"allowed_origins"`
		AllowedMethods   []string `mapstructure:"allowed_methods"`
		AllowedHeaders   []string `mapstructure:"allowed_headers"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
		MaxAge           int      `mapstructure:"max_age"`
	} `mapstructure:"cors"`
	RateLimiter struct {
		Requests int `mapstructure:"requests"`
		Duration int `mapstructure:"duration"`
	} `mapstructure:"rate_limiter"`
	Features struct {
		NewFeature  bool `mapstructure:"new_feature"`
		BetaFeature bool `mapstructure:"beta_feature"`
	} `mapstructure:"features"`
	ExternalServices struct {
		EmailService struct {
			APIKey   string `mapstructure:"api_key"`
			Endpoint string `mapstructure:"endpoint"`
		} `mapstructure:"email_service"`
	} `mapstructure:"external_services"`
}
