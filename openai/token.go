/*
Copyright 2024 The corekit Authors.

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
	"github.com/bytemind-io/corekit/token"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cast"
)

// CalculateCustomResponseToken calculates custom response token for the given model.
func CalculateCustomResponseToken(model string, out ...ChatCompletionResponse) (int, error) {
	if len(out) == 0 {
		return 0, nil
	}

	if model == "" {
		model = out[0].Model
	}

	text := ""
	imageList := make([]string, 0)
	for _, v := range out {
		// If the response is text or code, attach the response content.
		switch v.Message.Content.ContentType {
		case ContentTypeText:
			if len(v.Message.Content.Parts) > 0 && v.Message.Author.Role != ContentTypeText {
				for _, part := range v.Message.Content.Parts {
					text += cast.ToString(part)
				}
			}
		case ContentTypeCode, ContentTypeExecutionOutput:
			text += v.Message.Content.Text
		}

		if len(v.Downloads) != 0 {
			for _, list := range v.Downloads {
				for _, url := range list {
					imageList = append(imageList, url)
				}
			}
		}
	}

	tokenNum := 0
	if text != "" {
		num, err := token.CalculateTextToken(text, model)
		if err != nil {
			return 0, err
		}
		tokenNum += num
	}

	for _, url := range imageList {
		num, err := token.CalculateImageToken(&openai.ChatMessageImageURL{URL: url}, model)
		if err != nil {
			return 0, err
		}
		tokenNum += num
	}

	return tokenNum, nil
}
