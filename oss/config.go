package oss

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type (
	// Config represents the configuration for the oss service.
	Config struct {
		Driver          string `json:"driver" envconfig:"OSS_DRIVER"`
		Region          string `json:"region" envconfig:"OSS_REGION"`
		Bucket          string `json:"bucket" envconfig:"OSS_BUCKET"`
		EndPoints       string `json:"endPoints" envconfig:"OSS_ENDPOINTS"`
		AccessKeyId     string `json:"accessKeyId" envconfig:"OSS_ACCESS_KEYID"`
		AccessKeySecret string `json:"access_key_secret" envconfig:"OSS_ACCESS_KEY_SECRET"`
		URL             string `json:"url" envconfig:"OSS_URL"`
		Secure          bool   `json:"secure" envconfig:"OSS_SECURE" default:"true"`
	}
)

func MustLoadConfig(cfgPath string) *Config {
	if cfgPath == "" {
		cfg := &Config{}
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatal("load s3 config failed, error:", err)
		}
		return cfg
	}

	return nil
}
