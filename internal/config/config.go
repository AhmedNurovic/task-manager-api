package config

import (
	"github.com/spf13/viper"
	"fmt"
)

type Config struct {
	DBURL     string `mapstructure:"DB_URL"`
	Port      string `mapstructure:"PORT"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

func Load() (*Config, error) {
	viper.AddConfigPath(".")       
	viper.SetConfigName(".env")    
	viper.SetConfigType("env")    

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read the config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config into struct: %v", err)
	}

	if cfg.DBURL == "" || cfg.Port == "" || cfg.JWTSecret == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return &cfg, nil
}
