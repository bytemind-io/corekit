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

package corekit

import "github.com/zeromicro/go-zero/core/logx"

// Config defines the logx configuration.
type Config struct {
	EnableStat  bool   `json:",optional,env=ENABLE_STAT,default=false"       envconfig:"ENABLE_STAT"              default:"false"`
	ServiceName string `json:",optional,env=ZEROLOG_SERVICE_NAME"            envconfig:"ZEROLOG_SERVICE_NAME"`
	LogLevel    string `json:",optional,env=ZEROLOG_LEVEL,default=info"      envconfig:"ZEROLOG_LEVEL"            default:"info"`
	LogEncoding string `json:",optional,env=ZEROLOG_ENCODING,default=json"   envconfig:"ZEROLOG_ENCODING"         default:"json"`
}

// SetZeroLogx sets up logx with the given configuration.
func SetZeroLogx(stat bool, name, level, encoding string) {
	logx.SetUp(logx.LogConf{
		ServiceName: name,
		Encoding:    encoding,
		Level:       level,
	})
	if !stat {
		logx.DisableStat()
	}
	logx.AddGlobalFields(
		logx.Field("service_name", name),
	)
}
