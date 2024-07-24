package configs

import (
	"github.com/spf13/viper"
)

type conf struct {
	WebServerPort string     `mapstructure:"WEB_SERVER_PORT"`
	REDIS_HOST    string     `mapstructure:"REDIS_HOST"`
	REDIS_PORT    string     `mapstructure:"REDIS_PORT"`
	RTLIP         int        `mapstructure:"RTL_IP"`
	RTLBlockTime  int        `mapstructure:"RTL_BLOCK_TIME"`
	RTLTokens     []RTLToken `mapstructure:"RTL_TOKENS"`
}

type RTLToken struct {
	Token          string `mapstructure:"token"`
	ExpirationTime int    `mapstructure:"expiration_time"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	viper.SetConfigType("json")
	viper.SetConfigFile("tokens.json")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.UnmarshalKey("RTL_TOKENS", &cfg.RTLTokens)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
