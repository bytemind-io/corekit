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
	"encoding/json"
	"fmt"
	customOpenai "github.com/bytemind-io/corekit/openai"
	"github.com/spf13/cast"
	"image"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/bytemind-io/corekit"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

var (
	tokenMap      = map[string]*tiktoken.Tiktoken{}
	defaultToken  *tiktoken.Tiktoken
	cl100kToken   *tiktoken.Tiktoken
	o200kToken    *tiktoken.Tiktoken
	p50kToken     *tiktoken.Tiktoken
	p50kEditToken *tiktoken.Tiktoken
	r50kToken     *tiktoken.Tiktoken
)

func init() {
	token, err := tiktoken.GetEncoding(tiktoken.MODEL_CL100K_BASE)
	if err != nil {
		panic(err)
	}
	cl100kToken = token
	defaultToken = cl100kToken

	token, err = tiktoken.GetEncoding(tiktoken.MODEL_O200K_BASE)
	if err != nil {
		panic(err)
	}
	o200kToken = token

	token, err = tiktoken.GetEncoding(tiktoken.MODEL_P50K_BASE)
	if err != nil {
		panic(err)
	}
	p50kToken = token

	token, err = tiktoken.GetEncoding(tiktoken.MODEL_P50K_EDIT)
	if err != nil {
		panic(err)
	}
	p50kEditToken = token

	token, err = tiktoken.GetEncoding(tiktoken.MODEL_R50K_BASE)
	if err != nil {
		panic(err)
	}
	r50kToken = token

	//预设的模型编码
	for model, encoding := range tiktoken.MODEL_TO_ENCODING {
		switch encoding {
		case tiktoken.MODEL_CL100K_BASE:
			tokenMap[model] = cl100kToken
		case tiktoken.MODEL_O200K_BASE:
			tokenMap[model] = o200kToken
		case tiktoken.MODEL_P50K_BASE:
			tokenMap[model] = p50kToken
		case tiktoken.MODEL_P50K_EDIT:
			tokenMap[model] = p50kEditToken
		case tiktoken.MODEL_R50K_BASE:
			tokenMap[model] = r50kToken
		default:
			tokenMap[model] = defaultToken
		}
	}

	//自定义的模型编码
	for model, encoding := range ModelToEncoding {
		switch encoding {
		case tiktoken.MODEL_CL100K_BASE:
			tokenMap[model] = cl100kToken
		case tiktoken.MODEL_O200K_BASE:
			tokenMap[model] = o200kToken
		case tiktoken.MODEL_P50K_BASE:
			tokenMap[model] = p50kToken
		case tiktoken.MODEL_P50K_EDIT:
			tokenMap[model] = p50kEditToken
		case tiktoken.MODEL_R50K_BASE:
			tokenMap[model] = r50kToken
		default:
			tokenMap[model] = cl100kToken
		}
	}
}

// getModelDefaultTokenEncoder returns the default token encoder for the given model.
func getModelDefaultTokenEncoder(model string) *tiktoken.Tiktoken {
	if strings.HasPrefix(model, "gpt-4o") {
		return o200kToken
	}
	return defaultToken
}

// getTokenEncoder returns the token encoder for the given model.
func getTokenEncoder(model string) *tiktoken.Tiktoken {
	tokenEncoder, ok := tokenMap[model]
	if !ok {
		tokenEncoder, err := tiktoken.EncodingForModel(model)
		if err != nil {
			tokenEncoder = getModelDefaultTokenEncoder(model)
			return tokenEncoder
		}
		tokenMap[model] = tokenEncoder
		return tokenEncoder
	}

	return tokenEncoder
}

// CalculateRequestToken calculates the chat token for the given model.[请求相关]
func CalculateRequestToken(in *openai.ChatCompletionRequest, model string) (int, error) {
	if model == "" {
		model = in.Model
	}

	tkm := 0
	msgTokens, err := CalculateMessage(in.Messages, model)
	if err != nil {
		return 0, err
	}
	tkm += msgTokens

	if in.Tools != nil {
		toolsData, _ := json.Marshal(in.Tools)
		var openaiTools []openai.Tool
		err := json.Unmarshal(toolsData, &openaiTools)
		if err != nil {
			return 0, fmt.Errorf("count_tools_token_fail: %s", err)
		}
		countStr := ""
		for _, tool := range openaiTools {
			countStr = tool.Function.Name
			if tool.Function.Description != "" {
				countStr += tool.Function.Description
			}
			if tool.Function.Parameters != nil {
				countStr += fmt.Sprintf("%v", tool.Function.Parameters)
			}
		}
		toolTokens, err := CalculateTextToken(countStr, model)
		if err != nil {
			return 0, err
		}
		tkm += 8
		tkm += toolTokens
	}
	return tkm, nil
}

// CalculateStreamResponseToken calculates the chat token for the given model.
func CalculateStreamResponseToken(out *openai.ChatCompletionStreamResponse, model string) (int, error) {
	if model == "" {
		model = out.Model
	}

	messages := make([]openai.ChatCompletionStreamChoiceDelta, 0)
	for _, v := range out.Choices {
		messages = append(messages, v.Delta)
	}
	return CalculateStreamMessage(messages, model)
}

func CalculateStreamMessage(messages []openai.ChatCompletionStreamChoiceDelta, model string) (int, error) {
	tokenEncoder := getTokenEncoder(model)

	tokenNum := 0
	for _, message := range messages {
		if len(message.Content) > 0 {
			if message.Content != "" {
				stringContent := message.Content
				tokenNum += getTokenNum(tokenEncoder, stringContent)
			}
		}
	}

	return tokenNum, nil
}

// CalculateResponseToken calculates the chat token for the given model.
func CalculateResponseToken(out *openai.ChatCompletionResponse, model string) (int, error) {
	if model == "" {
		model = out.Model
	}

	messages := make([]openai.ChatCompletionMessage, 0)
	for _, v := range out.Choices {
		messages = append(messages, v.Message)
	}
	return CalculateMessage(messages, model)
}

// CalculateCustomResponseToken calculates custom response token for the given model.
func CalculateCustomResponseToken(model string, out ...customOpenai.ChatCompletionResponse) (int, error) {
	if len(out) == 0 {
		return 0, nil
	}

	if model == "" {
		model = out[0].Model
	}

	text := ""
	for _, v := range out {
		if v.Message.Content.ContentType == customOpenai.ContentTypeText {
			text += v.Message.Content.Text
			for _, part := range v.Message.Content.Parts {
				text += cast.ToString(part)
			}
		}
	}

	imageList := make([]string, 0)
	for _, v := range out {
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
		num, err := CalculateTextToken(text, model)
		if err != nil {
			return 0, err
		}
		tokenNum += num
	}

	for _, url := range imageList {
		num, err := CalculateImageToken(&openai.ChatMessageImageURL{URL: url}, model)
		if err != nil {
			return 0, err
		}
		tokenNum += num
	}

	return tokenNum, nil

}

// CalculateteMessage calculates the chat token for the given model.[消息Token计算]
// https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
// https://github.com/pkoukk/tiktoken-go/issues/6
func CalculateMessage(messages []openai.ChatCompletionMessage, model string) (int, error) {
	tokenEncoder := getTokenEncoder(model)
	var (
		tokensPerMessage int
		tokensPerName    int
	)

	if model == "gpt-3.5-turbo-0301" {
		tokensPerMessage = 4
		tokensPerName = -1
	} else {
		tokensPerMessage = 3
		tokensPerName = 1
	}

	tokenNum := 0
	for _, message := range messages {
		tokenNum += tokensPerMessage
		tokenNum += getTokenNum(tokenEncoder, message.Role)
		if len(message.Content) > 0 {
			if message.Content != "" {
				stringContent := message.Content
				tokenNum += getTokenNum(tokenEncoder, stringContent)
				if message.Name != "" {
					tokenNum += tokensPerName
					tokenNum += getTokenNum(tokenEncoder, message.Name)
				}
			} else {
				arrayContent := message.MultiContent
				for _, m := range arrayContent {
					if m.Type == openai.ChatMessagePartTypeImageURL {
						imageTokenNum, err := CalculateImageToken(m.ImageURL, model)
						if err != nil {
							return 0, err
						}
						tokenNum += imageTokenNum
					} else {
						tokenNum += getTokenNum(tokenEncoder, m.Text)
					}
				}
			}
		}
	}

	tokenNum += 3
	return tokenNum, nil
}

// CalculateImageBytesToken calculates the chat token for the given image bytes.[计算图片Token]
func CalculateImageBytesToken(imageBytes []byte, model string) (int, error) {
	if model == "glm-4v" {
		return 1047, nil
	}

	config, _, err := corekit.GetImageConfig(imageBytes)
	if err != nil {
		return 0, err
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

// CalculateImageToken gets the token number for the given image URL.[计算图片Token]
func CalculateImageToken(imageUrl *openai.ChatMessageImageURL, model string) (int, error) {
	if model == "glm-4v" {
		return 1047, nil
	}
	if imageUrl.Detail == "low" {
		return 85, nil
	}
	if imageUrl.Detail == "auto" || imageUrl.Detail == "" {
		imageUrl.Detail = "high"
	}

	var (
		config image.Config
		err    error
	)
	if strings.HasPrefix(imageUrl.URL, "http") {
		config, _, err = corekit.DecodeUrlImage(imageUrl.URL)
	} else {
		config, _, _, err = corekit.DecodeBase64Image(imageUrl.URL)
	}
	if err != nil {
		return 0, err
	}

	if config.Width == 0 || config.Height == 0 {
		return 0, fmt.Errorf("fail to decode image config: %s", imageUrl.URL)
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

// CalculateTextToken calculates the chat token for the given model.[文本Token计算]
func CalculateTextToken(in interface{}, model string) (int, error) {
	switch v := in.(type) {
	case string:
		return CounterText(v, model)
	case []string:
		text := ""
		for _, s := range v {
			text += s
		}
		return CounterText(text, model)
	}
	return CalculateTextToken(fmt.Sprintf("%v", in), model)
}

// CalculateAudioToken calculates the chat token for the given model.[音频Token计算]
func CalculateAudioToken(text, model string) (int, error) {
	if strings.HasPrefix(model, "tts") {
		return utf8.RuneCountInString(text), nil
	} else {
		return CounterText(text, model)
	}
}

// CounterText counts the number of tokens in the given text.[计数相关]
func CounterText(text string, model string) (int, error) {
	tokenEncoder := getTokenEncoder(model)
	return getTokenNum(tokenEncoder, text), nil
}

// getTokenNum returns the number of tokens in the given text.
func getTokenNum(tokenEncoder *tiktoken.Tiktoken, text string) int {
	return len(tokenEncoder.Encode(text, nil, nil))
}
