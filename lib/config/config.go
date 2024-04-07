package config

import "github.com/spf13/viper"

type Config struct {
	// Postgres
	Engine   string `mapstructure:"engine"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`

	AlchemyUrl           string `mapstructure:"alchemy_url"`
	SmartContractAddress string `mapstructure:"smart_contract_address"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
