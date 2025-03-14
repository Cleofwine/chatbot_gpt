package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host        string
		Port        int
		AccessToken string `mapstructure:"access_token"`
	}
	Mysql struct {
		DSN         string `mapstructure:"dsn"`
		MaxLeftTime int    `mapstructure:"max_life_time"`
		MaxOpenConn int    `mapstructure:"max_open_conn"`
		MaxIdleConn int    `mapstructure:"max_idle_conn"`
	}
	Log struct {
		Level   string
		LogPath string `mapstructure:"log_path"`
	}
}

var cfg *Config

func InitConf(configPath string) {
	if configPath == "" {
		panic("请指定应用程序配置文件")
	}
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		panic("配置文件不存在")
	}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(configPath)
	v.ReadInConfig()
	cfg = &Config{}
	err = v.Unmarshal(cfg)
	if err != nil {
		panic(err.Error())
	}
}

func GetConf() *Config {
	return cfg
}
