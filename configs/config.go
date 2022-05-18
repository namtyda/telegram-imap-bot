package configs

import "github.com/spf13/viper"

func InitConfig() error {
	viper.AddConfigPath("../../configs")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
