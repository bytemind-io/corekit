package oss

import (
	"errors"
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/kelseyhightower/envconfig"
)

type (
	// Config represents the configuration for the oss service.
	Config struct {
		Driver          string `json:",optional,env=OSS_DRIVER" envconfig:"OSS_DRIVER"`
		Region          string `json:",optional,env=OSS_REGION" envconfig:"OSS_REGION"`
		Bucket          string `json:",optional,env=OSS_BUCKET" envconfig:"OSS_BUCKET"`
		EndPoint        string `json:",optional,env=OSS_ENDPOINT" envconfig:"OSS_ENDPOINT"`
		AccessKeyId     string `json:",optional,env=OSS_ACCESS_KEYID" envconfig:"OSS_ACCESS_KEYID"`
		AccessKeySecret string `json:",optional,env=OSS_ACCESS_KEY_SECRET" envconfig:"OSS_ACCESS_KEY_SECRET"`
		URL             string `json:",optional,env=OSS_URL" envconfig:"OSS_URL"`
		Secure          bool   `json:",optional,env=OSS_SECURE,default=true" envconfig:"OSS_SECURE" default:"true"`
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
