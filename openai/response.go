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

import "encoding/json"

type ChatCodeResponseV4 struct {
	Created        int           `json:"created"`
	MessageID      string        `json:"message_id"`
	ConversationID string        `json:"conversation_id"`
	EndTurn        bool          `json:"end_turn"`
	Confirm        bool          `json:"confirm,omitempty"`
	Contents       []interface{} `json:"contents"`
	Downloads      []string      `json:"downloads,omitempty"`
	Usage          Usage         `json:"usage,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionResponse represents the response of the ChatCompletion API.
type ChatCompletionResponse struct {
	ConversationID string              `json:"conversation_id"`
	Error          interface{}         `json:"error"`
	Message        Message             `json:"message,omitempty"`
	Downloads      []map[string]string `json:"downloads,omitempty"`
	Model          string              `json:"model"`
	Stop           bool                `json:"stop"`
}

// Marshal returns the JSON encoding of ChatCompletionResponse.
func (c ChatCompletionResponse) Marshal() []byte {
	body, _ := json.Marshal(c)
	return body
}

// Message is the message for chat service.
type Message struct {
	Id         string      `json:"id"`
	Author     Author      `json:"author"`
	CreateTime float64     `json:"create_time"`
	UpdateTime float64     `json:"update_time"`
	Content    Content     `json:"content"`
	Status     string      `json:"status"`
	EndTurn    bool        `json:"end_turn"`
	Weight     float64     `json:"weight"`
	Metadata   interface{} `json:"metadata"`
	Recipient  string      `json:"recipient"`
}

// Content is the content for chat service.
type Content struct {
	ContentType string        `json:"content_type"`
	Parts       []interface{} `json:"parts"`
	Text        string        `json:"text"`
	Language    string        `json:"language"`
}

// Author is the author for chat service.
type Author struct {
	Role     string      `json:"role"`
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
}
