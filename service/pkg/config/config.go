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
	Chat struct {
		ProxyKey          string  `mapstructure:"proxy_key"`
		ProxyBaseURL      string  `mapstructure:"proxy_base_url"`
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
	Redis struct {
		Host string
		Port int
		PWD  string
	}
	DependOnServices struct {
		Tokenizer struct {
			Address     string
			AccessToken string `mapstructure:"access_token"`
		}
		ChatgptData struct {
			Address     string
			AccessToken string `mapstructure:"access_token"`
		} `mapstructure:"chatgpt_data"`
		ChatgptSensitive struct {
			Address     string
			AccessToken string `mapstructure:"access_token"`
		} `mapstructure:"chatgpt_sensitive"`
		ChatgptKeywords struct {
			Address     string
			AccessToken string `mapstructure:"access_token"`
		} `mapstructure:"chatgpt_keywords"`
	} `mapstructure:"depend_on_services"`
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
