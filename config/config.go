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

package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/zeromicro/go-zero/core/conf"
)

// Load loads configuration from file.
func Load(file string, v any, df func(v any), opts ...conf.Option) error {
	if err := conf.Load(file, &v, opts...); err != nil {
		return err
	}
	// Load configuration from environment variables.
	if err := envconfig.Process("", &v); err != nil {
		return err
	}

	df(&v)
	return nil
}
