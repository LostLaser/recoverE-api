package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	initialized = false
)

// Get should be called to retrieve any value from the yaml.
func Get(propertyName string) interface{} {
	setup()
	if !viper.IsSet(propertyName) {
		panic("Can't find property: " + propertyName)
	}
	return viper.Get(propertyName)
}

func setup() {
	if initialized {
		return
	}
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.BindEnv("port")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Config file get error: #%v ", err)
	}
	viper.AutomaticEnv()

	initialized = true
}
