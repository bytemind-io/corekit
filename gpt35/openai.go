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

// ChatCompletionRequest is the request for chat. https://beta.openai.com/docs/api-reference/messages/post
type ChatCompletionRequest struct {
	Action                     string      `json:"action"`
	Messages                   []Message   `json:"messages"`
	ConversationMode           interface{} `json:"conversation_mode"`
	Model                      string      `json:"model"`
	Stream                     bool        `json:"stream"`
	ConversationId             string      `json:"conversation_id,omitempty"`
	ParentMessageId            string      `json:"parent_message_id"`
	ForceParagen               bool        `json:"force_paragen,omitempty"`
	ForceRateLimit             bool        `json:"force_rate_limit,omitempty"`
	TimezoneOffsetMin          int         `json:"timezone_offset_min,omitempty"`
	HistoryAndTrainingDisabled bool        `json:"history_and_training_disabled,omitempty"`
	WebsocketRequestId         string      `json:"websocket_request_id,omitempty"`
}

// Message is the messages.
type Message struct {
	ID        string      `json:"id"`
	Author    Author      `json:"author"`
	Content   Content     `json:"content"`
	Metadata  interface{} `json:"metadata,omitempty"`
	Recipient string      `json:"recipient,omitempty"`
}

// Author is the author.
type Author struct {
	Role string `json:"role"`
	Name string `json:"name,omitempty"`
}

// Content is the content.
type Content struct {
	ContentType string        `json:"content_type"`
	Parts       []interface{} `json:"parts"`
}

type Parts []Part

type Part struct {
	AssetPointer string `json:"asset_pointer"`
	SizeBytes    int    `json:"size_bytes"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

type Attachments []Attachment

type Attachment struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Size          int    `json:"size"`
	FileTokenSize int    `json:"fileTokenSize,omitempty"`
	MimeType      string `json:"mimeType"`
	Width         int    `json:"width,omitempty"`
	Height        int    `json:"height,omitempty"`
}
