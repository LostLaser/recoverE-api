package utils

import (
	"log"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Get should be called to retrieve any value from the yaml.
func Get(propertyName string) interface{} {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to retrive calling environment for property retrieval")
	}
	basepath := path.Clean(filepath.Dir(b))

	viper.AddConfigPath(basepath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.BindEnv("port")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Config file get error: #%v ", err)
	}
	viper.AutomaticEnv()

	return viper.Get(propertyName)
}
