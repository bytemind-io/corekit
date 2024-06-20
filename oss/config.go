package oss

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type (
	// Config represents the configuration for the oss service.
	Config struct {
		Driver          string `json:"driver" envconfig:"OSS_DRIVER"`
		Region          string `json:"region" envconfig:"OSS_REGION"`
		Bucket          string `json:"bucket" envconfig:"OSS_BUCKET"`
		EndPoint        string `json:"end_point" envconfig:"OSS_ENDPOINT"`
		AccessKeyId     string `json:"access_key_id" envconfig:"OSS_ACCESS_KEYID"`
		AccessKeySecret string `json:"access_key_secret" envconfig:"OSS_ACCESS_KEY_SECRET"`
		URL             string `json:"url" envconfig:"OSS_URL"`
		Secure          bool   `json:"secure" envconfig:"OSS_SECURE" default:"true"`
	}
)

func (c Config) Validate() error {
	if govalidator.IsNull(c.Region) {
		return errors.New("region is required")
	}

	if govalidator.IsNull(c.EndPoint) {
		return errors.New("endpoint is required")
	}

	if govalidator.IsNull(c.AccessKeyId) {
		return errors.New("minio access key id is required")
	}

	if govalidator.IsNull(c.AccessKeySecret) {
		return errors.New("secret access key is required")
	}

	if govalidator.IsNull(c.URL) {
		return errors.New("url is required")
	}

	if govalidator.IsNull(c.Bucket) {
		return errors.New("bucket name is required")
	}
	return nil
}

func MustLoadConfig(cfgPath string) Config {
	cfg := Config{}
	if cfgPath == "" {
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatal("load s3 config failed, error:", err)
		}
		return cfg
	}

	// todo load config from file
	return cfg
}
