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

package token

import (
	"context"
	"fmt"
	"image"
	"math"
	"strings"

	"github.com/bytemind-io/corekit"
	"github.com/pkoukk/tiktoken-go"
	openai "github.com/sashabaranov/go-openai"
)

type gpt4 struct{}

// Name returns the name of the model.
func (s *gpt4) Name() string {
	return "gpt-4"
}

// Encode encodes the tiktoken.
func (s *gpt4) Encode() (*tiktoken.Tiktoken, error) {
	encoder, err := tiktoken.EncodingForModel(s.Name())
	if err != nil {
		return nil, err
	}
	return encoder, nil
}

// Calculate calculates the completion.
func (s *gpt4) Calculate(ctx context.Context, msg openai.ChatCompletionRequest) (int, error) {
	// Encode the message.tiktoken
	encoder, err := s.Encode()
	if err != nil {
		return 0, err
	}

	// Calculate the number of tokens per message and name.
	var tokensPerMessage, tokensPerName, tokenNum, tokenFun = 3, 1, 0, func(text string) int {
		return len(encoder.Encode(text, nil, nil))
	}

	for _, message := range msg.Messages {
		tokenNum += tokensPerMessage
		tokenNum += tokenFun(message.Role)
		if len(message.Content) > 0 {
			stringContent := message.Content
			tokenNum += tokenFun(stringContent)
			if message.Name != "" {
				tokenNum += tokensPerName
				tokenNum += tokenFun(message.Name)
			}
		} else {
			for _, part := range message.MultiContent {
				if part.Type == openai.ChatMessagePartTypeImageURL {

				}
			}
		}
	}

	return 0, nil
}

// Image returns the number of tokens for an image.
func (s *gpt4) Image(ctx context.Context, img *openai.ChatMessageImageURL) (int, error) {
	if img.Detail == "low" {
		return 85, nil
	}

	if img.Detail == "auto" || img.Detail == "" {
		img.Detail = "high"
	}

	var (
		config image.Config
		err    error
	)
	if strings.HasPrefix(img.URL, "http") {
		config, _, err = corekit.DecodeUrlImage(img.URL)
	} else {
		config, _, _, err = corekit.DecodeBase64Image(img.URL)
	}
	if err != nil {
		return 0, err
	}

	if config.Width == 0 || config.Height == 0 {
		return 0, fmt.Errorf("fail to decode image config: %s", img.URL)
	}

	shortSide := config.Width
	otherSide := config.Height
	scale := 1.0
	if config.Height < shortSide {
		shortSide = config.Height
		otherSide = config.Width
	}
	if shortSide > 768 {
		scale = float64(shortSide) / 768
		shortSide = 768
	}

	otherSide = int(math.Ceil(float64(otherSide) / scale))
	tiles := (shortSide + 511) / 512 * ((otherSide + 511) / 512)
	return tiles*170 + 85, nil
}
