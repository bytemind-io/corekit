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

	"github.com/asaskevich/govalidator"
	"github.com/sashabaranov/go-openai"
)

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

// ChatCompletionMessages is the messages for chat service.
type ChatCompletionMessages []*ChatCompletionMessage

// ChatCompletionMessage is the message for chat service.
type ChatCompletionMessage struct {
	Role        string      `json:"role"`
	Name        string      `json:"name"`
	Content     string      `json:"content"`
	Attachments Attachments `json:"attachments,omitempty"`
	Parts       Parts       `json:"parts,omitempty"`
}

// Parts is the parts for chat service.
type Parts []Part

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

// Attachments is the attachments for chat service.
type Attachments []Attachment

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
