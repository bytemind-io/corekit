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
	"github.com/bytemind-io/corekit/openai"
	"github.com/spf13/cast"
)

// ClaudeRequest is the request for chat. https://docs.anthropic.com/claude/reference/messages_post
// claude3.5 [官网](https://docs.anthropic.com/claude/reference/messages_post)
type ClaudeRequest struct {
	Model             string      `json:"model"`
	MaxTokens         int         `json:"max_tokens"`
	Messages          Messages    `json:"messages"`
	Stream            bool        `json:"stream"`
	Metadata          interface{} `json:"metadata,omitempty"`
	MaxTokensToSample float64     `json:"max_tokens_to_sample,omitempty"`
	Temperature       float64     `json:"temperature"`
	TopP              float64     `json:"top_p"`
	TopK              float64     `json:"top_k"`
}

// OpenaiConvertClaude is the request for chat. https://docs.anthropic.com/claude/reference/messages_post
func OpenaiConvertClaude(r openai.ChatCompletionRequest) *ClaudeRequest {
	req := &ClaudeRequest{
		Model:     r.Model,
		Stream:    r.Stream,
		MaxTokens: r.MaxTokens,
	}

	for _, message := range r.Messages {
		if len(message.Parts) != 0 {
			var contents Contents
			for _, part := range message.Parts {
				contents = append(contents, Content{
					Type: "image",
					Source: &Source{
						Type:      "base64",
						MediaType: part.MimeType,
						Data:      part.ImageData,
					},
				})
			}

			contents = append(contents, Content{
				Type: "text",
				Text: cast.ToString(message.Content),
			})

			req.Messages = append(req.Messages, Message{
				Role:    message.Role,
				Content: contents,
			})
			return req
		}
	}

	for _, message := range r.Messages {
		req.Messages = append(req.Messages, Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return req
}
