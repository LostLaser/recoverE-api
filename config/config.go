package config

import (
	"log"

	"github.com/spf13/viper"
)

// Get should be called to retrieve any value from the yaml.
func Get(propertyName string) interface{} {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.BindEnv("port")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Config file get error: #%v ", err)
	}
	viper.AutomaticEnv()

	v := viper.Get(propertyName)
	if v == nil {
		panic("Can't find property: " + propertyName)
	}
	return v
}
