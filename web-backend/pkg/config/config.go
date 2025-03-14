package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Enterprise struct {
		Id string
	}
	Http struct {
		Host              string
		Port              int
		BasicAuthUser     string `mapstructure:"basic_auth_user"`
		BasicAuthPassword string `mapstructure:"basic_auth_password"`
	}
	DependOnServices struct {
		ChatgptServices struct {
			Address     string
			AccessToken string `mapstructure:"access_token"`
		} `mapstructure:"chatgpt_services"`
	} `mapstructure:"depend_on_services"`
	Log struct {
		Level   string
		LogPath string `mapstructure:"log_path"`
	}
	Chat struct {
		Model             string  `mapstructure:"model"`
		MaxTokens         int     `mapstructure:"max_tokens"`
		Temperature       float32 `mapstructure:"temperature"`
		TopP              float32 `mapstructure:"top_p"`
		PresencePenalty   float32 `mapstructure:"presence_penalty"`
		FrequencyPenalty  float32 `mapstructure:"frequency_penalty"`
		BotDesc           string  `mapstructure:"bot_desc"`
		ContextTTL        int     `mapstructure:"context_ttl"`
		ContextLen        int     `mapstructure:"context_len"`
		MinResponseTokens int     `mapstructure:"min_response_tokens"`
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
