package oss

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestVultrPutObject(t *testing.T) {
	cfg := Config{
		Bucket:          "testtesttest",
		EndPoint:        "Your EndPoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your EndPoint",
		Secure:          true,
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	body, err := os.ReadFile("image.png")
	if err != nil {
		t.Error(err)
	}

	fileInfo, err := os.Stat("image.png")
	if err != nil {
		t.Error(err)
	}

	//上传文件
	obj := Object{
		UserId:      "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
		Bucket:      cfg.Bucket,
		FileName:    fileInfo.Name(),
		FileBytes:   body,
		FileSize:    fileInfo.Size(),
		ContentType: "image/png",
		VersionId:   "",
		CreatedAt:   time.Time{},
	}
	if meta, err := client.PutObject(context.Background(), &obj); err != nil {
		t.Log(err)
	} else {
		t.Log(meta.Url)
	}
}

func TestVultrListObject(t *testing.T) {
	cfg := Config{
		Bucket:          "testtesttest",
		EndPoint:        "Your EndPoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your EndPoint",
		Secure:          true,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if meta, err := client.ListObject(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(meta)
	} else {
		t.Log(meta)
	}
}

func TestVultrDeleteObject(t *testing.T) {
	cfg := Config{
		Bucket:          "testtesttest",
		EndPoint:        "Your EndPoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your EndPoint",
		Secure:          true,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if err := client.DeleteObject(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		t.Log("delete object success")
	}
}

func TestVultrSize(t *testing.T) {
	cfg := Config{
		Bucket:          "testtesttest",
		EndPoint:        "Your EndPoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your EndPoint",
		Secure:          true,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if total, err := client.Size(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(total)
	}
}

func TestVultrClear(t *testing.T) {
	cfg := Config{
		Bucket:          "testtesttest",
		EndPoint:        "Your EndPoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your EndPoint",
		Secure:          true,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if err := client.Clear(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		t.Log("clear object success")
	}
}

func TestVultrDownload(t *testing.T) {
	cfg := Config{
		Bucket:          "testtesttest",
		EndPoint:        "Your EndPoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your EndPoint",
		Secure:          true,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if bytes, err := client.Download(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		ioutil.WriteFile("image.png", bytes, 0666)
		t.Log("download object success")
	}
}

func TestVultrCreateBucket(t *testing.T) {
	cfg := Config{
		Bucket:          "testtesttest",
		EndPoint:        "Your EndPoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your EndPoint",
		Secure:          true,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if err := client.CreateBucket(context.Background(), cfg.Bucket, cfg.Region); err != nil {
		t.Error(err)
	} else {
		t.Log("create bucket success")
	}
}

func TestVultrSetBucketPolicy(t *testing.T) {
	cfg := Config{
		Bucket:          "testtesttest",
		EndPoint:        "Your EndPoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		Secure:          true,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if err := client.CreateBucket(context.Background(), cfg.Bucket, cfg.Region); err != nil {
		t.Error(err)
	} else {
		t.Log("create bucket success")
	}

	policy := `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": "*",
                "Action": [
                    "s3:GetObject"
                ],
                "Resource": [
                    "arn:aws:s3:::` + cfg.Bucket + `/*"
                ]
            }
        ]
    }`
	err = client.SetBucketPolicy(context.Background(), cfg.Bucket, policy)
	if err != nil {
		t.Fatal(err)
	} else {
		log.Printf("Successfully set policy for bucket %s\n", cfg.Bucket)
	}
}

func TestVultrURL(t *testing.T) {
	cfg := Config{
		Bucket:          "testtesttest",
		EndPoint:        "Your EndPoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		Secure:          true,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if str, err := client.URL(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(str)
	}
}

func TestMinIOPutObject(t *testing.T) {
	cfg := Config{
		Region:          "us-east-1",
		Bucket:          "test-test",
		EndPoint:        "Your Endpoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your Endpoint",
		Secure:          false,
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	body, err := os.ReadFile("image.png")
	if err != nil {
		t.Error(err)
	}

	fileInfo, err := os.Stat("image.png")
	if err != nil {
		t.Error(err)
	}

	//上传文件
	obj := Object{
		UserId:      "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
		Bucket:      cfg.Bucket,
		FileName:    fileInfo.Name(),
		FileBytes:   body,
		FileSize:    fileInfo.Size(),
		ContentType: "image/png",
		VersionId:   "",
		CreatedAt:   time.Time{},
	}
	if meta, err := client.PutObject(context.Background(), &obj); err != nil {
		t.Log(err)
	} else {
		t.Log(meta.Url)
	}
}

func TestMinIOListObject(t *testing.T) {
	cfg := Config{
		Region:          "us-east-1",
		Bucket:          "test-test",
		EndPoint:        "Your Endpoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your Endpoint",
		Secure:          false,
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if meta, err := client.ListObject(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(meta)
	} else {
		t.Log(meta)
	}
}

func TestMinIODeleteObject(t *testing.T) {
	cfg := Config{
		Region:          "us-east-1",
		Bucket:          "test-test",
		EndPoint:        "Your Endpoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your Endpoint",
		Secure:          false,
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if err := client.DeleteObject(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		t.Log("delete object success")
	}
}

func TestMinIOCreateBucket(t *testing.T) {
	cfg := Config{
		Region:          "us-east-1",
		Bucket:          "test-test",
		EndPoint:        "Your Endpoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your Endpoint",
		Secure:          false,
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if err := client.CreateBucket(context.Background(), cfg.Bucket, cfg.Region); err != nil {
		t.Error(err)
	} else {
		t.Log("create bucket success")
	}
}

func TestMinIOSetBucketPolicy(t *testing.T) {
	cfg := Config{
		Region:          "us-east-1",
		Bucket:          "test-test",
		EndPoint:        "Your Endpoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your Endpoint",
		Secure:          false,
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if err := client.CreateBucket(context.Background(), cfg.Bucket, cfg.Region); err != nil {
		t.Error(err)
	} else {
		t.Log("create bucket success")
	}

	policy := `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": "*",
                "Action": [
                    "s3:GetObject"
                ],
                "Resource": [
                    "arn:aws:s3:::` + cfg.Bucket + `/*"
                ]
            }
        ]
    }`
	err = client.SetBucketPolicy(context.Background(), cfg.Bucket, policy)
	if err != nil {
		t.Fatal(err)
	} else {
		log.Printf("Successfully set policy for bucket %s\n", cfg.Bucket)
	}
}

func TestMinIOSize(t *testing.T) {
	cfg := Config{
		Region:          "us-east-1",
		Bucket:          "test-test",
		EndPoint:        "Your Endpoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your Endpoint",
		Secure:          false,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if total, err := client.Size(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(total)
	}
}

func TestMinIOClear(t *testing.T) {
	cfg := Config{
		Region:          "us-east-1",
		Bucket:          "test-test",
		EndPoint:        "Your Endpoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your Endpoint",
		Secure:          false,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if err := client.Clear(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		t.Log("clear object success")
	}
}

func TestMinIODownload(t *testing.T) {
	cfg := Config{
		Region:          "us-east-1",
		Bucket:          "test-test",
		EndPoint:        "Your Endpoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your Endpoint",
		Secure:          false,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if bytes, err := client.Download(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		ioutil.WriteFile("image.png", bytes, 0666)
		t.Log("download object success")
	}
}

func TestMinIOURL(t *testing.T) {
	cfg := Config{
		Region:          "us-east-1",
		Bucket:          "test-test",
		EndPoint:        "Your Endpoint",
		AccessKeyId:     "Your Access Key",
		AccessKeySecret: "Your Access Secret",
		URL:             "http://Your Endpoint",
		Secure:          false,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}

	if str, err := client.URL(context.Background(), Metadata{
		BucketName: cfg.Bucket,
		ObjectName: "image.png",
		UserID:     "e335cc2a-f237-43f1-84c1-2b4b803a2f56",
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(str)
	}
}
