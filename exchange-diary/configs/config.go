package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	db           DBConfig
	kakaoClient  KakaoClient
	googleClient GoogleClient
}

type DBConfig struct {
	host string
	name string
	port int
}

type KakaoClient struct {
	apiKey string
}

type GoogleClient struct {
	apiKey string
}

const (
	typeExtension = "yaml"
)

func Load(path string, name string) (Config, error) {
	config := Config{}
	fmt.Println("hello")

	viper.SetConfigName(name)
	viper.SetConfigType(typeExtension)
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
