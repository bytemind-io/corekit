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

package database

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"time"
)

// Config is the configuration for the database.
type Config struct {
	LogLevel    int           `json:",optional,env=LOG_LEVEL,default=2"                envconfig:"LOG_LEVEL"              default:"2"`        // 日志级别
	Driver      string        `json:",optional,env=DRIVER,default=postgres"            envconfig:"DRIVER"                 default:"postgres"` // 驱动
	Database    string        `json:",optional,env=DATABASE"                           envconfig:"DATABASE"`                                  // 数据库
	AutoMigrate bool          `json:",optional,env=DATABASE_AUTO_MIGRATE,default=true" envconfig:"DATABASE_AUTO_MIGRATE"  default:"true"`     // 是否自动迁移
	MaxLifetime time.Duration `json:",optional,env=DATABASE_MAX_LIFETIME,default=3s"   envconfig:"DATABASE_MAX_LIFETIME"  default:"3s"`       // 最大连接周期
	MaxIdleTime time.Duration `json:",optional,env=DATABASE_MAX_IDLETIME,default=5s"   envconfig:"DATABASE_MAX_IDLETIME"  default:"5s"`       // 最大空闲连接周期
	MaxOpen     int           `json:",optional,env=DATABASE_MAX_OPEN,default=50"       envconfig:"DATABASE_MAX_OPEN"      default:"50"`       // 最大连接数
	MaxIdle     int           `json:",optional,env=DATABASE_MAX_IDLE,default=20"       envconfig:"DATABASE_MAX_IDLE"      default:"20"`       // 最大空闲连接数
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if govalidator.IsNull(c.Driver) {
		return fmt.Errorf("driver is required")
	}

	if govalidator.IsNull(c.Database) {
		return fmt.Errorf("database is required")
	}
	return nil
}
