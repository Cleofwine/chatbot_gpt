package server

import (
	"chatgpt-service/pkg/config"
	"chatgpt-service/pkg/log"
	"chatgpt-service/proto"
	chatcontext "chatgpt-service/server/chat-context"
	"chatgpt-service/services/datas"
	datas_proto "chatgpt-service/services/datas/proto"
	defaultclient "chatgpt-service/services/default_client"
	"chatgpt-service/services/keywords"
	keywords_proto "chatgpt-service/services/keywords/proto"
	"chatgpt-service/services/sensitive-words"
	sensitive_proto "chatgpt-service/services/sensitive-words/proto"
	"chatgpt-service/services/tokenizer"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

const ChatPrimedTokens = 2

type chatgptServiceServer struct {
	proto.UnimplementedChatGPTServiceServerServer
	config *config.Config
	log    *log.Logger
}

func NewChatGPTServiceServer(config *config.Config, log *log.Logger) proto.ChatGPTServiceServerServer {
	return &chatgptServiceServer{
		config: config,
		log:    log,
	}
}

type chatGPTAPP struct {
	config config.Config
	log    *log.Logger
}

func (s *chatgptServiceServer) getChatGPTAPP(in *proto.ChatCompletionReq) *chatGPTAPP {
	conf := *s.config
	if in.ChatParam != nil {
		if in.ChatParam.Model != "" {
			conf.Chat.Model = in.ChatParam.Model
		}
		conf.Chat.TopP = in.ChatParam.TopP
		conf.Chat.FrequencyPenalty = in.ChatParam.FrequencyPenalty
		conf.Chat.PresencePenalty = in.ChatParam.PresencePenalty
		conf.Chat.Temperature = in.ChatParam.Temperature
		conf.Chat.BotDesc = in.ChatParam.BotDesc
		if conf.Chat.MaxTokens != 0 {
			conf.Chat.MaxTokens = int(in.ChatParam.MaxTokens)
		}
		if conf.Chat.ContextLen != 0 {
			conf.Chat.ContextLen = int(in.ChatParam.ContextLen)
		}
		if conf.Chat.ContextTTL != 0 {
			conf.Chat.ContextTTL = int(in.ChatParam.ContextTTL)
		}
		if conf.Chat.MinResponseTokens != 0 {
			conf.Chat.MinResponseTokens = int(in.ChatParam.MinResponseTokens)
		}
	}
	app := &chatGPTAPP{
		log:    s.log,
		config: conf,
	}
	return app
}

func (s *chatgptServiceServer) ChatCompletion(ctx context.Context, in *proto.ChatCompletionReq) (*proto.ChatCompletionRes, error) {
	app := s.getChatGPTAPP(in)
	// 敏感词过滤
	ok, msg, err := app.sensitiveWords(in)
	if err != nil {
		app.log.Error(err)
		return nil, err
	}
	if !ok {
		res := app.buildChatCompletionResponse(msg)
		return res, nil
	}
	// 关键词查找
	keywords := app.keywords(in)
	// 正式与gpt通信
	client := app.getChatGPTClient()
	contextList, tokensNumAll, currTokensNum, currMessage, req, err := app.buildChatCompletionRequest(in, false)
	if err != nil {
		app.log.Error(err)
		return nil, err
	}
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		app.log.Error(err)
		return nil, err
	}
	res := &proto.ChatCompletionRes{}
	res.Id = resp.ID
	res.Object = resp.Object
	res.Created = resp.Created
	res.Model = resp.Model
	usage := proto.Usage{
		PromptTokens:     int32(resp.Usage.PromptTokens),
		CompletionTokens: int32(resp.Usage.CompletionTokens),
		TotalTokens:      int32(resp.Usage.TotalTokens),
	}
	res.Usage = &usage
	res.Choices = []*proto.ChatCompletionChoice{}
	for i := 0; i < len(resp.Choices); i++ {
		index := resp.Choices[i].Index
		message := resp.Choices[i].Message
		finishReason := resp.Choices[i].FinishReason
		msg := proto.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
			Name:    message.Name,
		}
		choices := proto.ChatCompletionChoice{
			Index:        int32(index),
			Message:      &msg,
			FinishReason: string(finishReason),
		}
		res.Choices = append(res.Choices, &choices)
	}
	// 保存上下文，开启协程提升效率，上下文即使保存失败也不会影响正常流程
	go func() {
		reqContext := &chatcontext.ChatMessage{
			Message:   currMessage,
			TokensNum: currTokensNum,
		}
		resContext := &chatcontext.ChatMessage{
			ID:        uuid.New().String(),
			Message:   resp.Choices[0].Message,
			TokensNum: resp.Usage.CompletionTokens,
		}
		if in.Endpoint == proto.ChatEndpoint_WEB {
			reqContext.ID = in.Id
			reqContext.PID = in.Pid
			if resp.ID != "" {
				resContext.ID = resp.ID
			}
			resContext.PID = reqContext.ID
		}
		err := app.saveContext(in, reqContext, resContext, contextList)
		if err != nil {
			app.log.Error(err)
			return
		}
	}()
	// 调用数据服务
	go func() {
		err := app.saveData(in, keywords, currTokensNum, resp.Choices[0].Message.Content, resp.Usage.CompletionTokens, tokensNumAll)
		if err != nil {
			app.log.Error(err)
		}
	}()
	return res, err
}

func (s *chatgptServiceServer) ChatCompletionStream(in *proto.ChatCompletionReq, stream proto.ChatGPTServiceServer_ChatCompletionStreamServer) error {
	app := s.getChatGPTAPP(in)
	// 敏感词过滤
	ok, msg, err := app.sensitiveWords(in)
	if err != nil {
		app.log.Error(err)
		return err
	}
	if !ok {
		resId := uuid.New().String()
		startRes := app.buildChatCompletionStreamResponse(resId, "", "")
		endRes := app.buildChatCompletionStreamResponse(resId, "", "stop")
		err = stream.Send(startRes)
		if err != nil {
			app.log.Error(err)
			return err
		}
		resList := app.buildChatCompletionStreamResponseList(resId, msg)
		for _, res := range resList {
			err = stream.Send(res)
			if err != nil {
				app.log.Error(err)
				return err
			}
		}
		err = stream.Send(endRes)
		if err != nil {
			app.log.Error(err)
			return err
		}
		return nil
	}
	// 关键词查找
	keywords := app.keywords(in)
	// 正式与gpt通信
	client := app.getChatGPTClient()
	contextList, tokensNumAll, currTokensNum, currMessage, req, err := app.buildChatCompletionRequest(in, true)
	if err != nil {
		app.log.Error(err)
		return err
	}
	chatStream, err := client.CreateChatCompletionStream(stream.Context(), req)
	if err != nil {
		app.log.Error(err)
		return err
	}
	defer chatStream.Close()
	completionContent := ""
	resultID := ""
	for {
		resp, err := chatStream.Recv()
		if err != nil && err != io.EOF {
			app.log.Error(err)
			return err
		}
		if err == io.EOF {
			break
		}
		if resultID == "" {
			resultID = resp.ID
		}
		completionContent += resp.Choices[0].Delta.Content
		res := &proto.ChatCompletionStreamRes{}
		res.Id = resp.ID
		res.Object = resp.Object
		res.Created = resp.Created
		res.Model = resp.Model
		res.Choices = []*proto.ChatCompletionStreamChoice{}
		for i := 0; i < len(resp.Choices); i++ {
			index := resp.Choices[i].Index
			delta := resp.Choices[i].Delta
			finishReason := resp.Choices[i].FinishReason
			del := proto.ChatCompletionStreamDelta{
				Role:    delta.Role,
				Content: delta.Content,
			}
			choices := proto.ChatCompletionStreamChoice{
				Index:        int32(index),
				Delta:        &del,
				FinishReason: string(finishReason),
			}
			res.Choices = append(res.Choices, &choices)
		}
		err = stream.Send(res)
		if err != nil {
			app.log.Error(err)
			return err
		}
	}
	resultCompletion := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: completionContent,
	}
	resultTokensNum, err := tokenizer.GetTokenCount(resultCompletion, s.config.Chat.Model)
	if err != nil {
		app.log.Error(err)
	}
	// 保存上下文
	go func() {
		reqContext := &chatcontext.ChatMessage{
			Message:   currMessage,
			TokensNum: currTokensNum,
		}
		resContext := &chatcontext.ChatMessage{
			ID: uuid.New().String(),
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completionContent,
			},
			TokensNum: resultTokensNum,
		}
		if in.Endpoint == proto.ChatEndpoint_WEB {
			reqContext.ID = in.Id
			reqContext.PID = in.Pid
			if resultID != "" {
				resContext.ID = resultID
			}
			resContext.PID = reqContext.ID
		}
		err := app.saveContext(in, reqContext, resContext, contextList)
		if err != nil {
			app.log.Error(err)
		}
	}()

	// 调用数据服务
	go func() {
		err := app.saveData(in, keywords, currTokensNum, completionContent, resultTokensNum, tokensNumAll)
		if err != nil {
			app.log.Error(err)
		}
	}()
	return nil
}

func (s *chatGPTAPP) getChatGPTClient() *openai.Client {
	conf := s.config
	accessToken := conf.Chat.ProxyKey
	config := openai.DefaultConfig(accessToken)
	config.BaseURL = conf.Chat.ProxyBaseURL
	client := openai.NewClientWithConfig(config)
	return client
}

func (s *chatGPTAPP) buildChatCompletionRequest(in *proto.ChatCompletionReq, stream bool) (contextList []*chatcontext.ChatMessage, tokenNum int, currTokensNum int, currMessage openai.ChatCompletionMessage, req openai.ChatCompletionRequest, err error) {
	conf := s.config
	currMessage = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: in.Message,
	}
	req = openai.ChatCompletionRequest{
		Model: conf.Chat.Model,
		Messages: []openai.ChatCompletionMessage{
			currMessage,
		},
		// 表示当前请求的回复，最大tokens数
		MaxTokens:        conf.Chat.MinResponseTokens,
		Temperature:      conf.Chat.Temperature,
		TopP:             conf.Chat.TopP,
		PresencePenalty:  conf.Chat.PresencePenalty,
		FrequencyPenalty: conf.Chat.FrequencyPenalty,
		Stream:           stream,
	}
	contextList = make([]*chatcontext.ChatMessage, 0)
	var value interface{}
	if in.EnableContext {
		cache := chatcontext.GetCacheContext(in.Endpoint)
		cacheID := in.Id
		if in.Endpoint == proto.ChatEndpoint_WEB {
			cacheID = in.Pid
		}
		value, err = cache.Get(cacheID, in.GroupId, in.Endpoint)
		if err != nil {
			s.log.Error(err)
			return
		}
		contextList = value.([]*chatcontext.ChatMessage)
	}
	tokenNum, currTokensNum, req.Messages, err = s.buildMessage(contextList, currMessage)
	if err != nil {
		s.log.Error(err)
		return
	}
	req.MaxTokens = conf.Chat.MaxTokens - tokenNum
	return
}

func (s *chatGPTAPP) buildMessage(contextList []*chatcontext.ChatMessage, currMessage openai.ChatCompletionMessage) (tokensNum int, currTokensNum int, messages []openai.ChatCompletionMessage, err error) {
	conf := s.config
	var sysMessage openai.ChatCompletionMessage
	messages = []openai.ChatCompletionMessage{currMessage}
	currTokensNum, err = tokenizer.GetTokenCount(currMessage, conf.Chat.Model)
	if err != nil {
		s.log.Error(err)
		return
	}
	botTokens := 0
	if conf.Chat.BotDesc != "" {
		sysMessage = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: conf.Chat.BotDesc,
		}
		botTokens, err = tokenizer.GetTokenCount(sysMessage, conf.Chat.Model)
		if err != nil {
			s.log.Error(err)
			return
		}
	}
	if currTokensNum > conf.Chat.MaxTokens-conf.Chat.MinResponseTokens-botTokens {
		err = fmt.Errorf("上下文合计tokens最大%d，回复保留tokens数%d，ai特征使用tokens %d，剩余可用tokens %d，当前消息tokens %d", conf.Chat.MaxTokens, conf.Chat.MinResponseTokens, botTokens, conf.Chat.MaxTokens-conf.Chat.MinResponseTokens-botTokens, currTokensNum)
		s.log.Error(err)
		return
	}
	tokensNum = currTokensNum + botTokens + ChatPrimedTokens
	for _, item := range contextList {
		if tokensNum+item.TokensNum > conf.Chat.MaxTokens-conf.Chat.MinResponseTokens {
			break
		}
		messages = append(messages, item.Message)
		tokensNum += item.TokensNum + ChatPrimedTokens
	}
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	if botTokens > 0 {
		messages = append([]openai.ChatCompletionMessage{sysMessage}, messages...)
	}
	return
}

func (s *chatGPTAPP) saveContext(in *proto.ChatCompletionReq, reqContext, resContext *chatcontext.ChatMessage, contextList []*chatcontext.ChatMessage) (err error) {
	cache := chatcontext.GetCacheContext(in.Endpoint)
	if in.Endpoint == proto.ChatEndpoint_WEB {
		err = cache.Set(reqContext.ID, in.GroupId, in.Endpoint, reqContext, s.config.Chat.ContextTTL)
		if err != nil {
			s.log.Error(err)
			return
		}
		err = cache.Set(resContext.ID, in.GroupId, in.Endpoint, resContext, s.config.Chat.ContextTTL)
		if err != nil {
			s.log.Error(err)
			return
		}
		return nil
	}
	contextList = append([]*chatcontext.ChatMessage{resContext, reqContext}, contextList...)
	if len(contextList) > s.config.Chat.ContextLen {
		contextList = contextList[:s.config.Chat.ContextLen]
	}
	err = cache.Set(in.Id, in.GroupId, in.Endpoint, contextList, s.config.Chat.ContextTTL)
	if err != nil {
		s.log.Error(err)
		return
	}
	return nil
}

func (s *chatGPTAPP) sensitiveWords(in *proto.ChatCompletionReq) (ok bool, msg string, err error) {
	sensitiveClientPool := sensitive.GetSensitiveClientPool()
	conn := sensitiveClientPool.Get()
	defer sensitiveClientPool.Put(conn)
	client := sensitive_proto.NewChatGPTSensitiveClient(conn)
	ctx := context.Background()
	ctx = defaultclient.AppendBearerTokenToContext(ctx, s.config.DependOnServices.ChatgptSensitive.AccessToken)
	req := &sensitive_proto.ValidateReq{
		Text: in.Message,
	}
	res, err := client.Validate(ctx, req)
	if err != nil {
		s.log.Error(err)
		return false, "", err
	}
	ok = res.Ok
	if !ok {
		msg = "触及到了知识盲区，请换个问题再问吧"
	}
	return
}

func (s *chatGPTAPP) keywords(in *proto.ChatCompletionReq) []string {
	keywordsClientPool := keywords.GetKeywordsClientPool()
	conn := keywordsClientPool.Get()
	defer keywordsClientPool.Put(conn)
	client := keywords_proto.NewChatGPTKeywordsClient(conn)
	ctx := context.Background()
	ctx = defaultclient.AppendBearerTokenToContext(ctx, s.config.DependOnServices.ChatgptKeywords.AccessToken)
	req := &keywords_proto.FindAllReq{
		Text: in.Message,
	}
	res, err := client.FindAll(ctx, req)
	if err != nil {
		s.log.Error(err)
		return []string{}
	}
	return res.Words
}

func (s *chatGPTAPP) saveData(in *proto.ChatCompletionReq, keywords []string, userMsgTokens int, aiMsg string, aiMsgTokens int, reqTokens int) error {
	dataClientPool := datas.GetChatGPTDataClientPool()
	conn := dataClientPool.Get()
	defer dataClientPool.Put(conn)
	client := datas_proto.NewChatGPTDataClient(conn)
	ctx := context.Background()
	ctx = defaultclient.AppendBearerTokenToContext(ctx, s.config.DependOnServices.ChatgptData.AccessToken)
	req := &datas_proto.Record{
		UserMsg:         in.Message,
		UserMsgTokens:   int32(userMsgTokens),
		AiMsg:           aiMsg,
		AiMsgTokens:     int32(aiMsgTokens),
		UserMsgKeywords: keywords,
		ReqTokens:       int32(reqTokens),
		CreateAt:        time.Now().Unix(),
		Endpoint:        int32(in.Endpoint),
		EnterpriseId:    in.EnterpriseId,
		EndpointAccount: in.EndpointAccount,
	}
	if in.Endpoint != proto.ChatEndpoint_WEB {
		req.Account = in.Id
		req.GroupId = in.GroupId
	}
	_, err := client.AddRecord(ctx, req)
	return err
}

func (s *chatGPTAPP) buildChatCompletionResponse(msg string) *proto.ChatCompletionRes {
	res := &proto.ChatCompletionRes{
		Id:      uuid.New().String(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   s.config.Chat.Model,
		Choices: []*proto.ChatCompletionChoice{
			{
				Message: &proto.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: msg,
				},
				FinishReason: "stop",
			},
		},
		Usage: &proto.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
	}
	return res
}

func (s *chatGPTAPP) buildChatCompletionStreamResponseList(id, msg string) []*proto.ChatCompletionStreamRes {
	list := make([]*proto.ChatCompletionStreamRes, 0)
	for _, delta := range msg {
		list = append(list, s.buildChatCompletionStreamResponse(id, string(delta), ""))
	}
	return list
}

func (s *chatGPTAPP) buildChatCompletionStreamResponse(id, delta, finishReason string) *proto.ChatCompletionStreamRes {
	res := &proto.ChatCompletionStreamRes{
		Id:      id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   config.GetConf().Chat.Model,
		Choices: []*proto.ChatCompletionStreamChoice{
			{
				Index: 0,
				Delta: &proto.ChatCompletionStreamDelta{
					Content: delta,
					Role:    openai.ChatMessageRoleAssistant,
				},
				FinishReason: finishReason,
			},
		},
	}
	return res
}
