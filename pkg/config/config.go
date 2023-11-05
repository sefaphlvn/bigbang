package config

import (
	"github.com/spf13/viper"
)

func Read(cfgFile string) *AppConfig {

	var (
		err       error
		appConfig *AppConfig
	)

	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgFile)

	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err = viper.Unmarshal(&appConfig); err != nil {
		panic(err)
	}

	return appConfig
}
