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

package corekit

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/image/webp"
	"image"
	"io"
	"net/http"
	"strings"
)

// Request sends a request to the given URL.
func Request(url string) (*http.Response, error) {
	return http.Get(url)
}

// DecodeUrlImage decodes an image from a URL.
func DecodeUrlImage(url string) (image.Config, string, error) {
	resp, err := Request(url)
	if err != nil {
		return image.Config{}, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return image.Config{}, "", fmt.Errorf("failed to fetch image: %s", resp.Status)
	}

	var readData []byte

	for _, limit := range []int64{1024 * 8, 1024 * 24, 1024 * 64} {
		additionalData := make([]byte, limit-int64(len(readData)))
		n, _ := io.ReadFull(resp.Body, additionalData)
		readData = append(readData, additionalData[:n]...)
		limitReader := io.MultiReader(bytes.NewReader(readData), resp.Body)

		var config image.Config
		var format string
		config, format, err = getImageConfig(limitReader)
		if err == nil {
			return config, format, nil
		}
	}
	return image.Config{}, "", err
}

// GetImageConfig gets the image type and base64 encoded data.
func GetImageConfig(imageBytes []byte) (image.Config, string, error) {
	reader := bytes.NewReader(imageBytes)
	return getImageConfig(reader)
}

// DecodeBase64Image decodes a base64 image string.
func DecodeBase64Image(base64String string) (image.Config, string, string, error) {
	if idx := strings.Index(base64String, ","); idx != -1 {
		base64String = base64String[idx+1:]
	}

	decodedData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return image.Config{}, "", "", err
	}
	reader := bytes.NewReader(decodedData)
	config, format, err := getImageConfig(reader)
	return config, format, base64String, err
}

// GetImageFromUrl gets the image type and base64 encoded data.
func getImageConfig(reader io.Reader) (image.Config, string, error) {
	config, format, err := image.DecodeConfig(reader)
	if err != nil {
		err = errors.New(fmt.Sprintf("fail to decode image config(gif, jpg, png): %s", err.Error()))
		config, err = webp.DecodeConfig(reader)
		if err != nil {
			err = errors.New(fmt.Sprintf("fail to decode image config(webp): %s", err.Error()))
		}
		format = "webp"
	}
	if err != nil {
		return image.Config{}, "", err
	}
	return config, format, nil
}
