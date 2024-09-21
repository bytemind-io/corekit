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

package openai

import (
	"fmt"
	"github.com/bytemind-io/corekit/token"

	"github.com/asaskevich/govalidator"
	"github.com/sashabaranov/go-openai"
)

// ConversationV4Request represents the request for the conversation endpoint
type ConversationV4Request struct {
	GizmoId                    string      `json:"gizmo_id,omitempty"`
	Message                    string      `json:"message"`
	ParentMessageID            string      `json:"parent_message_id,omitempty"`
	ConversationID             string      `json:"conversation_id,omitempty"`
	Stream                     bool        `json:"stream,omitempty"`
	Model                      string      `json:"model"`
	Attachments                Attachments `json:"attachments,omitempty"`
	Parts                      Parts       `json:"parts,omitempty"`
	HistoryAndTrainingDisabled bool        `json:"history_and_training_disabled,omitempty"`
}

// ChatCompletionRequest represents the request for the conversation endpoint
type ChatCompletionRequest struct {
	GizmoId                    string                 `json:"gizmo_id"`
	Action                     string                 `json:"action"`
	Messages                   ChatCompletionMessages `json:"messages"`
	ParentMessageID            string                 `json:"parent_message_id,omitempty"`
	ConversationID             string                 `json:"conversation_id,omitempty"`
	Stream                     bool                   `json:"stream,omitempty"`
	Model                      string                 `json:"model"`
	HistoryAndTrainingDisabled bool                   `json:"history_and_training_disabled,omitempty"`
	Namespace                  string                 `json:"namespace"`
	Temperature                float32                `json:"temperature,omitempty"`
	RepetitionPenalty          float32                `json:"repetition_penalty,omitempty"`
	TopP                       float32                `json:"top_p,omitempty"`
	TopK                       float32                `json:"top_k,omitempty"`
	MaxTokensToSample          float32                `json:"max_tokens_to_sample,omitempty"`
	MaxTokens                  int                    `json:"max_tokens,omitempty"`
	MaxNewTokens               int                    `json:"max_new_tokens,omitempty"`
}

// Validate validates the request for the conversation endpoint.
func (r *ChatCompletionRequest) Validate() error {
	if r.Messages == nil || len(r.Messages) == 0 {
		return fmt.Errorf("message is required")
	}

	if govalidator.IsNull(r.Model) {
		return fmt.Errorf("model is required")
	}

	for idx, message := range r.Messages {
		if len(message.Parts) == 0 && len(message.Attachments) == 0 {
			if govalidator.IsNull(message.Content) {
				return fmt.Errorf("%d message content is empty", idx)
			}
		}

		if govalidator.IsNull(message.Role) {
			message.Role = openai.ChatMessageRoleUser
		}
	}

	if govalidator.IsNull(r.Action) {
		r.Action = "next"
	}

	if r.MaxTokens < 0 {
		r.MaxTokens = 0
	}
	return nil
}

// OpenAI marshals to openai.ChatChatCompletionRequest
func (r *ChatCompletionRequest) OpenAI() *openai.ChatCompletionRequest {
	req := &openai.ChatCompletionRequest{
		Model:            r.Model,
		MaxTokens:        r.MaxTokens,
		Temperature:      r.Temperature,
		TopP:             r.TopP,
		N:                0,
		Stream:           r.Stream,
		Stop:             nil,
		PresencePenalty:  0,
		ResponseFormat:   nil,
		Seed:             nil,
		FrequencyPenalty: 0,
		LogitBias:        nil,
		LogProbs:         false,
		TopLogProbs:      0,
		User:             "",
		Functions:        nil,
		FunctionCall:     nil,
		Tools:            nil,
		ToolChoice:       nil,
		StreamOptions:    nil,
	}

	if r.Messages != nil {
		req.Messages = r.Messages.Marshal()
	}
	return req
}

func (r *ChatCompletionRequest) CalculateRequestToken() (int, error) {
	return token.CalculateRequestToken(r.OpenAI(), "")
}

// ChatCompletionMessages is the messages for chat service.
type ChatCompletionMessages []*ChatCompletionMessage

func (m ChatCompletionMessages) Marshal() []openai.ChatCompletionMessage {
	res := make([]openai.ChatCompletionMessage, 0, len(m))
	for _, v := range m {
		res = append(res, v.Marshal())
	}
	return res
}

// ChatCompletionMessage is the message for chat service.
type ChatCompletionMessage struct {
	Role        string      `json:"role"`
	Name        string      `json:"name"`
	Content     string      `json:"content"`
	Attachments Attachments `json:"attachments,omitempty"`
	Parts       Parts       `json:"parts,omitempty"`
}

func (m ChatCompletionMessage) Marshal() openai.ChatCompletionMessage {
	if len(m.Parts) != 0 {
		var contents []openai.ChatMessagePart

		if len(m.Content) != 0 {
			contents = append(contents, openai.ChatMessagePart{
				Type: openai.ChatMessagePartTypeText,
				Text: m.Content,
			})
		}

		contents = append(contents, m.Parts.Marshal()...)
		if len(m.Attachments) != 0 {
			// TODO 需要移除！ 目前文件上传|几乎所有厂家都不支持
			//contents = append(contents, m.Attachments.Marshal()...)
		}
		return openai.ChatCompletionMessage{
			Role:         m.Role,
			MultiContent: contents,
			Name:         m.Name,
		}
	}
	return openai.ChatCompletionMessage{
		Role:    m.Role,
		Content: m.Content,
		Name:    m.Name,
	}
}

// Parts is the parts for chat service.
type Parts []Part

func (p Parts) Marshal() []openai.ChatMessagePart {
	res := make([]openai.ChatMessagePart, 0, len(p))
	for _, v := range p {
		res = append(res, v.MarshalToOpenaiPart())
	}
	return res
}

// Part is the part for chat service.
type Part struct {
	Name         string `json:"name"`
	AssetPointer string `json:"asset_pointer"`
	SizeBytes    int    `json:"size_bytes"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	MimeType     string `json:"mimeType"`
	ImageData    string `json:"image_data,omitempty"`
	OssUrl       string `json:"oss_url,omitempty"`
}

func (p Part) MarshalToOpenaiPart() openai.ChatMessagePart {
	return openai.ChatMessagePart{
		Type: openai.ChatMessagePartTypeImageURL,
		ImageURL: &openai.ChatMessageImageURL{
			URL: p.OssUrl,
		},
	}
}

// Attachments is the attachments for chat service.
type Attachments []Attachment

func (a Attachments) Marshal() []openai.ChatMessagePart {
	res := make([]openai.ChatMessagePart, 0, len(a))
	for _, v := range a {
		res = append(res, v.MarshalToOpenaiPart())
	}
	return res
}

// Attachment is the attachment for chat service.
type Attachment struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	FileTokenSize int    `json:"fileTokenSize,omitempty"`
	MimeType      string `json:"mimeType"`
	Width         int    `json:"width,omitempty"`
	Height        int    `json:"height,omitempty"`
	OssUrl        string `json:"oss_url,omitempty"`
}

func (a Attachment) MarshalToOpenaiPart() openai.ChatMessagePart {
	return openai.ChatMessagePart{
		Type: openai.ChatMessagePartTypeText,
		Text: "",
	}
}
