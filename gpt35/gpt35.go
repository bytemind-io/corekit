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

package gpt35

import (
	"github.com/bytemind-io/corekit/openai"
	"github.com/google/uuid"
)

// GPT35 is the GPT35.
func GPT35(in openai.ChatCompletionRequest) ChatCompletionRequest {
	req := ChatCompletionRequest{
		ConversationMode: map[string]interface{}{
			"kind": "primary_assistant",
		},
		ConversationId:             in.ConversationID,
		ParentMessageId:            in.ParentMessageID,
		Model:                      in.Model,
		Stream:                     in.Stream,
		ForceParagen:               false,
		ForceRateLimit:             false,
		TimezoneOffsetMin:          -480,
		HistoryAndTrainingDisabled: false,
		WebsocketRequestId:         uuid.NewString(),
	}

	if in.Action == "" {
		req.Action = "next"
	}

	if in.ParentMessageID == "" {
		req.ParentMessageId = uuid.NewString()
	}

	if in.Model == "" {
		req.Model = "text-davinci-002-render-sha"
	}

	if len(in.Messages) != 0 {
		for _, message := range in.Messages {
			msg := Message{
				ID: uuid.NewString(),
				Author: Author{
					Role: message.Role,
					Name: message.Name,
				},
			}

			if len(message.Parts) == 0 {
				msg.Content = Content{
					ContentType: "text",
					Parts: []interface{}{
						message.Content,
					},
				}
			}
			req.Messages = append(req.Messages, msg)
		}
	}
	return req
}
