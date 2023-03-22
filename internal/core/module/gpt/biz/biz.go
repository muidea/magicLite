package biz

import (
	"context"
	log "github.com/cihub/seelog"
	"github.com/muidea/magicLite/internal/config"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	openai "github.com/sashabaranov/go-openai"

	"github.com/muidea/magicLite/pkg/common"

	"github.com/muidea/magicLite/internal/core/base/biz"
)

type GPT struct {
	biz.Base
}

func New(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *GPT {
	return &GPT{
		Base: biz.New(common.GPTModule, eventHub, backgroundRoutine),
	}
}

func (s *GPT) pickAuthToken() string {
	tokens := config.GetAuthToken()
	return tokens[0]
}

func (s *GPT) QueryMessage(message string) (ret string, err error) {
	client := openai.NewClient(s.pickAuthToken())
	respVal, respErr := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if respErr != nil {
		err = respErr
		log.Errorf("ChatCompletion error: %v\n", respErr)
		return
	}

	ret = respVal.Choices[0].Message.Content
	return
}
