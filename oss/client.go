package oss

import (
	"bytes"
	"context"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/url"
	"sync"
	"time"
)

var (
	OssClient     Oss
	OssClientOnce sync.Once
)

func NewClient(cfg Config) (Oss, error) {
	var err error
	OssClientOnce.Do(func() {
		OssClient, err = newMinio(cfg)
	})

	return OssClient, err
}

type (
	Client struct {
		config Config // OSS client configuration
		client *minio.Client
	}
)

func newMinio(cfg Config) (Oss, error) {
	client, err := minio.New(cfg.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyId, cfg.AccessKeySecret, ""),
		Region: cfg.Region,
		Secure: cfg.Secure,
	})

	return &Client{
		config: cfg,
		client: client,
	}, err
}

func (c *Client) CreateBucket(ctx context.Context, bucket, region string) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}
	if region == "" {
		region = c.config.Region
	}
	return c.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
		Region: region,
	})
}

func (c *Client) DeleteBucket(ctx context.Context, bucket string) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}
	return c.client.RemoveBucket(ctx, bucket)
}

func (c *Client) SetBucketPolicy(ctx context.Context, bucket, policy string) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}
	return c.client.SetBucketPolicy(ctx, bucket, policy)
}

func (c *Client) ListBuckets(ctx context.Context) ([]Bucket, error) {
	buckets := make([]Bucket, 0)
	list, err := c.client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		buckets = append(buckets, Bucket{
			Name:         v.Name,
			CreationDate: v.CreationDate,
		})
	}
	return buckets, nil
}

func (c *Client) URL(ctx context.Context, metadata Metadata) (string, error) {
	if metadata.BucketName == "" {
		metadata.BucketName = c.config.Bucket
	}
	query := make(url.Values)
	query.Set("response-content-disposition", "attachment; filename=\""+metadata.ObjectName+"\"")
	ul, err := c.client.PresignedGetObject(ctx, metadata.BucketName, metadata.RelativeFilePath(), 60*5*time.Second, query)
	if err != nil {
		return "", err
	}
	return ul.String(), nil
}

func (c *Client) PutObject(ctx context.Context, obj *Object) (*Object, error) {
	if len(obj.FileBytes) == 0 {
		return nil, nil
	}
	if obj.Bucket == "" {
		obj.Bucket = c.config.Bucket
	}
	if obj.ContentType == "" {
		obj.ContentType = mimetype.Detect(obj.FileBytes).String()
	}
	objectId := obj.UUID()

	opts := minio.PutObjectOptions{
		ContentType:        obj.ContentType,
		ContentDisposition: "inline",
	}
	if c.config.Driver == DriverAliyun {
		opts.Header().Set("x-oss-object-acl", "public-read")
		opts.Header().Set("x-oss-tagging", "UserID="+obj.UserId)
	} else {
		opts.UserTags = map[string]string{
			"UserID": obj.UserId,
		}
		opts.UserMetadata = map[string]string{"x-amz-acl": "public-read"}
	}
	_, err := c.client.PutObject(ctx, obj.Bucket, objectId, bytes.NewReader(obj.FileBytes), obj.FileSize, opts)
	if err != nil {
		return obj, err
	}
	obj.Url = obj.GetFileUrl(c.config)
	return obj, err
}

func (c *Client) PutObjectReader(ctx context.Context, obj *ObjectReader) (*ObjectReader, error) {
	if obj == nil {
		return nil, nil
	}
	if obj.Bucket == "" {
		obj.Bucket = c.config.Bucket
	}
	if obj.ContentType == "" {
		if mime, err := mimetype.DetectReader(obj.Reader); err != nil {
			return nil, err
		} else {
			obj.ContentType = mime.String()
		}
	}
	objectId := obj.UUID()

	opts := minio.PutObjectOptions{
		ContentType:        obj.ContentType,
		ContentDisposition: "inline",
	}
	if c.config.Driver == DriverAliyun {
		opts.Header().Set("x-oss-object-acl", "public-read")
		opts.Header().Set("x-oss-tagging", "UserID="+obj.UserId)
	} else {
		opts.UserTags = map[string]string{
			"UserID": obj.UserId,
		}
		opts.UserMetadata = map[string]string{"x-amz-acl": "public-read"}
	}
	_, err := c.client.PutObject(ctx, obj.Bucket, objectId, obj.Reader, obj.FileSize, opts)
	if err != nil {
		return obj, err
	}
	obj.Url = obj.GetFileUrl(c.config)
	return obj, err
}

func (c *Client) GetObject(ctx context.Context, metadata Metadata) (*Object, error) {
	if metadata.BucketName == "" {
		metadata.BucketName = c.config.Bucket
	}
	var getObjOpts minio.GetObjectOptions
	if metadata.GetObjOpts == nil {
		getObjOpts = minio.GetObjectOptions{}
	} else {
		getObjOpts = *metadata.GetObjOpts
	}

	reader, err := c.client.GetObject(ctx, metadata.BucketName, metadata.RelativeFilePath(), getObjOpts)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("file not found")
	}

	obj := &Object{
		UserId:      metadata.UserID,
		Bucket:      metadata.BucketName,
		FileName:    metadata.ObjectName,
		FileBytes:   data,
		FileSize:    int64(len(data)),
		ContentType: mimetype.Detect(data).String(),
	}
	obj.Url = obj.GetFileUrl(c.config)
	return obj, err
}

func (c *Client) GetObjectReader(ctx context.Context, metadata Metadata) (*ObjectReader, error) {
	if metadata.BucketName == "" {
		metadata.BucketName = c.config.Bucket
	}
	var getObjOpts minio.GetObjectOptions
	if metadata.GetObjOpts == nil {
		getObjOpts = minio.GetObjectOptions{}
	} else {
		getObjOpts = *metadata.GetObjOpts
	}

	reader, err := c.client.GetObject(ctx, metadata.BucketName, metadata.RelativeFilePath(), getObjOpts)
	if err != nil {
		return nil, err
	}

	var fileSize int64
	if _, err := reader.Seek(0, io.SeekEnd); err != nil {
		return nil, err
	}
	if stat, err := reader.Stat(); err != nil {
		return nil, err
	} else {
		fileSize = stat.Size
	}

	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var detectReaderMime *mimetype.MIME
	if detectReaderMime, err = mimetype.DetectReader(reader); err != nil {
		return nil, err
	}

	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	obj := &ObjectReader{
		UserId:      metadata.UserID,
		Bucket:      metadata.BucketName,
		FileName:    metadata.ObjectName,
		Reader:      reader,
		FileSize:    fileSize,
		ContentType: detectReaderMime.String(),
	}
	obj.Url = obj.GetFileUrl(c.config)
	return obj, err
}

func (c *Client) DeleteObject(ctx context.Context, metadata Metadata) error {
	if metadata.BucketName == "" {
		metadata.BucketName = c.config.Bucket
	}
	return c.client.RemoveObject(ctx, metadata.BucketName, metadata.RelativeFilePath(), minio.RemoveObjectOptions{})
}

func (c *Client) ListObject(ctx context.Context, metadata Metadata) ([]Object, error) {
	if metadata.BucketName == "" {
		metadata.BucketName = c.config.Bucket
	}
	var listObjOpts minio.ListObjectsOptions
	if metadata.ListObjOpts == nil {
		listObjOpts = minio.ListObjectsOptions{
			Prefix:    metadata.UserID + "/",
			Recursive: true,
		}
	} else {
		listObjOpts = *metadata.ListObjOpts
	}
	objectCh := c.client.ListObjects(ctx, metadata.BucketName, listObjOpts)
	objectList := make([]Object, 0)
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objectList = append(objectList, Object{
			UserId:      object.UserTags["UserID"],
			Bucket:      metadata.BucketName,
			FileName:    object.Key,
			FileSize:    object.Size,
			ContentType: object.ContentType,
			VersionId:   object.VersionID,
			CreatedAt:   object.LastModified,
		})
	}
	return objectList, nil
}

func (c *Client) Clear(ctx context.Context, metadata Metadata) error {
	if metadata.BucketName == "" {
		metadata.BucketName = c.config.Bucket
	}
	var listObjOpts minio.ListObjectsOptions
	if metadata.ListObjOpts == nil {
		listObjOpts = minio.ListObjectsOptions{
			Prefix:    metadata.UserID + "/",
			Recursive: true,
		}
	} else {
		listObjOpts = *metadata.ListObjOpts
	}
	objectCh := c.client.ListObjects(ctx, metadata.BucketName, listObjOpts)

	var removeObjOpts minio.RemoveObjectOptions
	if metadata.RemoveObjOpts == nil {
		removeObjOpts = minio.RemoveObjectOptions{}
	} else {
		removeObjOpts = *metadata.RemoveObjOpts
	}

	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		err := c.client.RemoveObject(ctx, metadata.BucketName, object.Key, removeObjOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Download(ctx context.Context, metadata Metadata) ([]byte, error) {
	if metadata.BucketName == "" {
		metadata.BucketName = c.config.Bucket
	}
	var getObjOpts minio.GetObjectOptions
	if metadata.GetObjOpts == nil {
		getObjOpts = minio.GetObjectOptions{}
	} else {
		getObjOpts = *metadata.GetObjOpts
	}
	reader, err := c.client.GetObject(ctx, metadata.BucketName, metadata.RelativeFilePath(), getObjOpts)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) Size(ctx context.Context, metadata Metadata) (int64, error) {
	if metadata.BucketName == "" {
		metadata.BucketName = c.config.Bucket
	}
	var listObjOpts minio.ListObjectsOptions
	if metadata.ListObjOpts == nil {
		listObjOpts = minio.ListObjectsOptions{
			Prefix:    metadata.UserID + "/",
			Recursive: true,
		}
	} else {
		listObjOpts = *metadata.ListObjOpts
	}
	objectCh := c.client.ListObjects(ctx, metadata.BucketName, listObjOpts)

	var totalSize int64
	for object := range objectCh {
		if object.Err != nil {
			return 0, object.Err
		}
		totalSize += object.Size
	}
	return totalSize, nil
}
