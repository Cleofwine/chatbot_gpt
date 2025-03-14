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
		Host string
		Port int
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
		EnableContext     bool    `mapstructure:"enable_context"`
	}
	DependOnServices struct {
		ChatgptService struct {
			Address     string
			AccessToken string `mapstructure:"access_token"`
		} `mapstructure:"chatgpt_service"`
		ChatgptCrontab struct {
			Address     string
			AccessToken string `mapstructure:"access_token"`
		} `mapstructure:"chatgpt_crontab"`
	} `mapstructure:"depend_on_services"`
	Log struct {
		Level   string
		LogPath string `mapstructure:"log_path"`
	}
	Official struct {
		AppId string `mapstructure:"appid"`
		Token string `mapstructure:"token"`
	} `mapstructure:"official"`
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
