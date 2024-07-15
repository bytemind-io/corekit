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

package redisdb

import (
	"fmt"
	"github.com/asaskevich/govalidator"
)

// Config is the configuration for the redisdb.
type Config struct {
	Address  []string `json:",optional,env=REDIS_ADDRESS,default=127.0.0.1:6379"                envconfig:"REDIS_ADDRESS"             default:"127.0.0.1:6379"` // redis 集群地址
	Password string   `json:",optional,env=REDIS_ADDRESS"                                       envconfig:"REDIS_PASSWORD"`                                     // openssl rand -base64 12
	DB       int      `json:",optional,env=REDIS_DB,default=0"                                  envconfig:"REDIS_DB"                  default:"0"`
}

func NewConfig() Config {
	return Config{
		Address:  []string{"127.0.0.1:36379"},
		Password: "test-1",
	}
}

// Validate validates the Config.
func (c Config) Validate() error {
	if len(c.Address) == 0 {
		return fmt.Errorf("redis cluster is required")
	}
	if govalidator.IsNull(c.Password) {
		return fmt.Errorf("redis password is required")
	}
	return nil
}
