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
	"fmt"

	"github.com/spf13/cast"

	weboai "github.com/bytemind-io/corekit/openai"
	"github.com/sashabaranov/go-openai"
)

// BedrockRequest is the bedrock request. https://us-east-1.console.aws.amazon.com/bedrock/home?region=us-east-1#/providers?model=anthropic.claude-3-haiku-20240307-v1:0
// AWS Bedrock [官网](https://us-east-1.console.aws.amazon.com/bedrock/home?region=us-east-1#/providers?model=anthropic.claude-3-haiku-20240307-v1:0)
type BedrockRequest struct {
	AnthropicVersion  string   `json:"anthropic_version"`
	MaxTokens         int      `json:"max_tokens,omitempty"`
	System            string   `json:"system,omitempty"`
	Messages          Messages `json:"messages"`
	MaxTokensToSample float32  `json:"max_tokens_to_sample,omitempty"`
	Temperature       float32  `json:"temperature,omitempty"`
	TopP              float32  `json:"top_p,omitempty"`
}

// OpenaiConvertSonnet is the sonnet.TODO only use in Charlie W. Johnson.
func OpenaiConvertSonnet(in openai.ChatCompletionRequest) (*BedrockRequest, error) {
	version, ok := AnthropicVersion[in.Model]
	if !ok {
		return nil, fmt.Errorf("model %s not found.(Aws Anthropic Version)", in.Model)
	}

	req := &BedrockRequest{
		AnthropicVersion: version,
		MaxTokens:        in.MaxTokens,
		Temperature:      in.Temperature,
		TopP:             in.TopP,
	}

	if in.MaxTokens <= 0 {
		req.MaxTokens = 1024
	}

	if in.Temperature <= 0 {
		req.Temperature = 1.0
	}

	var (
		msgs   Messages
		sysCnt int
	)

	for idx, message := range in.Messages {
		if message.Role == openai.ChatMessageRoleSystem {
			if in.Messages[idx+1].Role != openai.ChatMessageRoleUser {
				return nil, fmt.Errorf("system must be followed by user message")
			}

			sysCnt++
			req.System = message.Content
		} else {
			if message.Role == openai.ChatMessageRoleUser {
				if idx == 0 {
					msgs = append(msgs, Message{
						Role: openai.ChatMessageRoleUser,
						Content: []Content{
							{
								Type: "text",
								Text: message.Content,
							},
						},
					})
					continue
				}

				if in.Messages[idx-1].Role == openai.ChatMessageRoleSystem || in.Messages[idx-1].Role == openai.ChatMessageRoleAssistant {
					msgs = append(msgs, Message{
						Role: openai.ChatMessageRoleUser,
						Content: []Content{
							{
								Type: "text",
								Text: message.Content,
							},
						},
					})
					continue
				}

				return nil, fmt.Errorf("user must be followed by system and assistant message")
			}

			if message.Role == openai.ChatMessageRoleAssistant {
				if in.Messages[idx-1].Role != openai.ChatMessageRoleUser {
					return nil, fmt.Errorf("assistant must be followed by user message")
				}

				msgs = append(msgs, Message{
					Role: openai.ChatMessageRoleAssistant,
					Content: []Content{
						{
							Type: "text",
							Text: message.Content,
						},
					},
				})
				continue
			}
			return nil, fmt.Errorf("role must be system, user or assistant")
		}
	}

	if sysCnt > 1 {
		return nil, fmt.Errorf("system message must be only one.(%d)", sysCnt)
	}

	req.Messages = msgs
	return req, nil
}

// OpenaiWebConvertSonnet is the sonnet.TODO only use in Charlie W. Johnson.[将Openai网页数据抓换成Bedrock]
func OpenaiWebConvertSonnet(r weboai.ChatCompletionRequest) (*BedrockRequest, error) {
	version, ok := AnthropicVersion[r.Model]
	if !ok {
		return nil, fmt.Errorf("model %s not found.(Aws Anthropic Version)", r.Model)
	}

	req := &BedrockRequest{
		AnthropicVersion: version,
		MaxTokens:        r.MaxTokens,
		Temperature:      r.Temperature,
		TopP:             r.TopP,
	}

	if r.MaxTokens <= 0 {
		req.MaxTokens = 1024
	}

	if r.Temperature <= 0 {
		req.Temperature = 1.0
	}

	for _, message := range r.Messages {
		// system has two maybe has bug?
		if message.Role == openai.ChatMessageRoleSystem {
			req.System = message.Content
		}

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
			return req, nil
		}
	}

	for _, message := range r.Messages {
		req.Messages = append(req.Messages, Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return req, nil
}

/*
	example 1:
{
	  "modelId": "anthropic.claude-3-haiku-20240307-v1:0",
	  "contentType": "application/json",
	  "accept": "application/json",
	  "body": {
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens": 1000,
		"messages": [
		  {
			"role": "user",
			"content": [
			  {
				"type": "image",
				"source": {
				  "type": "base64",
				  "media_type": "image/jpeg",
				  "data": "iVBORw..."
				}
			  },
			  {
				"type": "text",
				"text": "What's in this image?"
			  }
			]
		  }
		]
	  }
	}

	example 2:

	{"role": "user", "content": [
	  {
		"type": "image",
		"source": {
		  "type": "base64",
		  "media_type": "image/jpeg",
		  "data": "/9j/4AAQSkZJRg...",
		}
	  },
	  {"type": "text", "text": "What is in this image?"}
	]}
*/

// Messages is the messages.
type Messages []Message

// Message is the message.
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// Contents is the contents.
type Contents []Content

// Content is the content.
type Content struct {
	Type   string  `json:"type"`
	Source *Source `json:"source,omitempty"`
	Text   string  `json:"text,omitempty"`
}

// Source is the source.
type Source struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}
