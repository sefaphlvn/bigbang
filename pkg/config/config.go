package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func Read(cfgFile string) *AppConfig {
	var appConfig AppConfig
	viper.SetConfigType("yaml")

	if _, found := os.LookupEnv("isProd"); found {
		// viper.AutomaticEnv()
		BindEnvs(&appConfig, "")
	} else {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		panic(err)
	}

	return &appConfig
}

func BindEnvs(iface interface{}, prefix string) {
	val := reflect.ValueOf(iface).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if mapstructureTag, ok := field.Tag.Lookup("mapstructure"); ok {
			envVar := mapstructureTag
			if prefix != "" {
				envVar = prefix + "_" + envVar
			}
			envVar = strings.ToUpper(envVar)
			viper.BindEnv(field.Name, envVar)
		}
	}
}
