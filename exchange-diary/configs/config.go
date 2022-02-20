package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DBConfig    DBConfig    `mapstructure:"db-config"`
	KakaoClient KakaoClient `mapstructure:"kakao-client"`
}

type DBConfig struct {
	Host string `mapstructure:"host"`
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type OAuthConfig struct {
	ClientId     string `mapstructure:"client-id"`
	ClientSecret string `mapstructure:"client-secret"`
	RedirectUrl  string `mapstructure:"redirect-url"`
}

type KakaoClient struct {
	Oauth OAuthConfig `mapstructure:"oauth"`
}

const (
	typeExtension = "yaml"
)

func Load(path string, name string) (Config, error) {
	config := Config{}
	fmt.Println("Load config file - profile:", name)

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(typeExtension)

	err := viper.ReadInConfig()

	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
