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
	"github.com/deepauto-io/filestype"
	"mime/multipart"
)

// UploadRequest is the upload request.
type UploadRequest struct {
	ConversationId string                `form:"conversation_id"`
	UploadType     string                `form:"upload_type"`
	File           *multipart.FileHeader `form:"files"`
}

func (r *UploadRequest) Validate() error {
	if r.File == nil {
		return fmt.Errorf("file is required")
	}

	if r.UploadType != string(filestype.MyFiles) && r.UploadType != string(filestype.Multimodal) {
		return fmt.Errorf("invalid upload type")
	}
	return nil
}

// UploadResponse is the upload response.
type UploadResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Size           int64  `json:"size"`
	DownloadUrl    string `json:"download_url"`
	ImageWidth     int    `json:"image_width"`
	ImageHeight    int    `json:"image_height"`
	UseCase        string `json:"use_case"`
	FileTokenSize  int    `json:"fileTokenSize,omitempty"`
	MimeType       string `json:"mimeType"`
	ConversationId string `json:"conversation_id"`
}
