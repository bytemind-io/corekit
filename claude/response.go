/*
Copyright 2024 The corego Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package claude

import (
	"time"

	"github.com/bytemind-io/corekit"
	"github.com/bytemind-io/corekit/openai"
	"github.com/google/uuid"
	sysopenai "github.com/sashabaranov/go-openai"
)

// ClaudeResponse is the response from the Claude service.
type ClaudeResponse struct {
	Type         string      `json:"type"`
	Index        int         `json:"index"`
	Delta        Content     `json:"delta"`
	Id           string      `json:"id"`
	Role         string      `json:"role"`
	Model        string      `json:"model"`
	Content      []Content   `json:"content"`
	StopReason   string      `json:"stop_reason"`
	StopSequence interface{} `json:"stop_sequence"`
}

// OpenAIWeb opens the ClaudeResponse.
func (r *ClaudeResponse) OpenAIWeb(in *openai.ChatCompletionRequest) (list []openai.ChatCompletionResponse) {
	for _, content := range r.Content {
		resp := openai.ChatCompletionResponse{
			ConversationID: in.ConversationID,
			Message: openai.Message{
				Id: r.Id,
				Author: openai.Author{
					Role: r.Role,
				},
				CreateTime: corekit.UnixNano(),
				Content: openai.Content{
					ContentType: content.Type,
					Parts: []interface{}{
						content.Text,
					},
				},
				Status:    corekit.FinishedSuccessfully,
				EndTurn:   false, //TODO maybe has bug ?
				Weight:    1.0,
				Recipient: corekit.All,
			},
		}
		list = append(list, resp)
	}
	return
}

func (r *ClaudeResponse) Openai() sysopenai.ChatCompletionResponse {
	req := sysopenai.ChatCompletionResponse{
		ID:      r.Id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   r.Model,
		Choices: []sysopenai.ChatCompletionChoice{
			{
				Index: 0,
				Message: sysopenai.ChatCompletionMessage{
					Role:    sysopenai.ChatMessageRoleAssistant,
					Content: r.Content[0].Text,
				},
				FinishReason: sysopenai.FinishReason(r.StopReason),
			},
		},
		Usage: sysopenai.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
		SystemFingerprint: "fp_" + uuid.NewString(),
	}

	req.Usage.PromptTokens = 0
	req.Usage.CompletionTokens = 0
	req.Usage.TotalTokens = 0
	return req
}
