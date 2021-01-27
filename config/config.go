package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	initialized = false
)

// Get should be called to retrieve any value from the yaml.
func Get(propertyName string) interface{} {
	setup()
	if !viper.IsSet(propertyName) {
		log.Fatal("Can't find property: " + propertyName)
	}
	return viper.Get(propertyName)
}

func setup() {
	if initialized {
		return
	}
	viper.AddConfigPath("./")
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.BindEnv("port")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	viper.AutomaticEnv()

	initialized = true
}
