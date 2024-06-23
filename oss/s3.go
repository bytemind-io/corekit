/*
Copyright 2024 The modhub Authors.

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

package oss

import (
	"context"
	"fmt"
	"path/filepath"
	"time"
)

// Metadata is the metadata for the s3.
type Metadata struct {
	BucketName string
	ObjectName string
	UserID     string
}

func (m Metadata) RelativeFilePath() string {
	return filepath.Join(m.UserID, m.ObjectName)
}

// Object is the object for the s3.
type Object struct {
	UserId      string    `json:"user_id"`
	Bucket      string    `json:"bucket"`
	FileName    string    `json:"file_name"`
	FileBytes   []byte    `json:"file_bytes"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
	VersionId   string    `json:"version_id"`
	CreatedAt   time.Time `json:"created_at"`
	Url         string    `json:"url"`
}

// UUID returns the uuid.
func (obj Object) UUID() string {
	s := filepath.Join(obj.UserId, obj.FileName)
	return s
}

func (obj Object) GetFileUrl(cfg Config) string {
	switch cfg.Driver {
	case DriverAliyun:
		return fmt.Sprintf("%s/%s", cfg.URL, obj.UUID())
	default:
		return fmt.Sprintf("%s/%s/%s", cfg.URL, obj.Bucket, obj.UUID())
	}
}

// Bucket container for bucket metadata.
type Bucket struct {
	// The name of the bucket.
	Name string `json:"name"`
	// Date the bucket was created.
	CreationDate time.Time `json:"creationDate"`
}

// Oss is the service for the s3.
type Oss interface {
	// CreateBucket creates a bucket.
	CreateBucket(ctx context.Context, bucket, region string) error
	// DeleteBucket delete a bucket.
	DeleteBucket(ctx context.Context, bucket string) error
	// ListBuckets list all buckets owned by this authenticated user.
	ListBuckets(ctx context.Context) ([]Bucket, error)
	// SetBucketPolicy sets the bucket policy.
	SetBucketPolicy(ctx context.Context, bucket, policy string) error
	// URL returns the url.
	URL(ctx context.Context, metadata Metadata) (string, error)
	// PutObject puts the object to the oss.
	PutObject(ctx context.Context, obj *Object) (*Object, error)
	// GetObject get the object from the oss.
	GetObject(ctx context.Context, metadata Metadata) (*Object, error)
	// DeleteObject deletes the object.
	DeleteObject(ctx context.Context, metadata Metadata) error
	// ListObject list all object owned by this authenticated user.
	ListObject(ctx context.Context, metadata Metadata) ([]Object, error)
	// Clear clears the object.
	Clear(ctx context.Context, metadata Metadata) error
	// Download downloads the object.
	Download(ctx context.Context, metadata Metadata) ([]byte, error)
	// Size returns the size.
	Size(ctx context.Context, metadata Metadata) (int64, error)
}
