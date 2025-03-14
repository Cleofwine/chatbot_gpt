// 单元测试：功能测试
package test_case

import (
	"context"
	"flag"
	"io"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

// 单元测试中，不能使用config对象，因为它读取不到我们的命令行参数
var configFile = flag.String("config", "../config.yaml", "测试用例配置文件")
var addr = flag.String("addr", "http://localhost:4002/v1", "chatgpt_proxy 地址")
var conf *viper.Viper

// 要记得先起服务，再测试
// 功能测试
// go test ./test-case/proxy_test.go --config=../dev.config.yaml --addr=http://localhost:4002/v1 -v
// 性能测试
// go test -bench . -benchmem test-case/*.go --config=../dev.config.yaml --addr=http://localhost:4002/v1 -v

// flag的使用和数据导出
// go test -bench . -benchtime 2x -count 2 -cpu 1,2 -cover -benchmem \
// -blockprofile block.out -cpuprofile cpu.out -memprofile mem.out \
// -mutexprofile mutex.out -outputdir ./testout \
// test-case/*.go --config=../dev.config.yaml --addr=http://localhost:4002/v1

// 只运行并发测试，不运行功能测试
// go test -run ^$ -bench ^BenchmarkProxyChatCompletionParallel$ -benchtime 2x -cpu 1,2 -benchmem  test-case/*.go --config=../dev.config.yaml --addr=http://localhost:4002/v1 -v
// go test -run ^$ -bench ^BenchmarkProxyChatCompletionParallel$ -benchtime 100x -cpu 100 -benchmem  test-case/*.go --config=../dev.config.yaml --addr=http://localhost:4002/v1 -v

func TestMain(m *testing.M) {
	flag.Parse()
	conf = viper.New()
	conf.SetConfigType("yaml")
	conf.SetConfigFile(*configFile)
	conf.ReadInConfig()
	m.Run()
}

func TestProxyChatCompletion(t *testing.T) {
	accessToken := conf.GetString("http.access_token")
	config := openai.DefaultConfig(accessToken)
	config.BaseURL = *addr
	client := openai.NewClientWithConfig(config)
	req := openai.ChatCompletionRequest{
		Model: conf.GetString("chat.model"),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
		},
		MaxTokens:        512,
		Temperature:      float32(conf.GetFloat64("chat.Temperature")),
		TopP:             float32(conf.GetFloat64("chat.top_p")),
		PresencePenalty:  float32(conf.GetFloat64("chat.presence_penalty")),
		FrequencyPenalty: float32(conf.GetFloat64("chat.frequency_penalty")),
	}
	req = setBotDest(req, "你是一个AI助手，我需要你模拟一名资深的软件工程师来回答我的问题")
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.Choices[0].Message.Content)
}

func TestProxyChatCompletionStream(t *testing.T) {
	accessToken := conf.GetString("http.access_token")
	config := openai.DefaultConfig(accessToken)
	config.BaseURL = *addr
	client := openai.NewClientWithConfig(config)
	req := openai.ChatCompletionRequest{
		Model: conf.GetString("chat.model"),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "你好",
			},
		},
		MaxTokens:        512,
		Temperature:      float32(conf.GetFloat64("chat.Temperature")),
		TopP:             float32(conf.GetFloat64("chat.top_p")),
		PresencePenalty:  float32(conf.GetFloat64("chat.presence_penalty")),
		FrequencyPenalty: float32(conf.GetFloat64("chat.frequency_penalty")),
		Stream:           true,
	}
	req = setBotDest(req, "你是一个AI助手，我需要你模拟一名资深的软件工程师来回答我的问题")
	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	defer stream.Close()
	for {
		resp, err := stream.Recv()
		if err != nil && err != io.EOF {
			t.Error(err)
			return
		}
		if err == io.EOF {
			break
		}
		t.Log(resp.Choices[0].Delta.Content)
	}
}

func setBotDest(request openai.ChatCompletionRequest, botDesc string) openai.ChatCompletionRequest {
	if request.Messages[0].Role != openai.ChatMessageRoleSystem {
		systemMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: botDesc,
		}
		newMessages := append([]openai.ChatCompletionMessage{systemMessage}, request.Messages...)
		request.Messages = newMessages
	}
	return request
}
