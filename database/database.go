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
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
)

// SqlStore represents a SQL store.
type SqlStore struct {
	db          *gorm.DB
	autoMigrate bool
}

// New creates a new SQL store.
func New(opt Config) (*SqlStore, error) {
	var (
		db  *gorm.DB
		err error
	)

	gormCfg := &gorm.Config{
		Logger: NewSlog(logger.Config{
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.LogLevel(opt.LogLevel),
		}),
	}

	switch opt.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(opt.Database), gormCfg)
	case "postgres":
		db, err = gorm.Open(postgres.Open(opt.Database), gormCfg)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(opt.Database), gormCfg)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(opt.Database), gormCfg)
	default:
		return nil, fmt.Errorf("can not find the driver type")
	}

	if err != nil {
		return nil, err
	}

	s := &SqlStore{
		db:          db,
		autoMigrate: opt.AutoMigrate,
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(opt.MaxIdle)
	sqlDB.SetMaxOpenConns(opt.MaxOpen)
	sqlDB.SetConnMaxLifetime(opt.MaxLifetime)
	sqlDB.SetConnMaxIdleTime(opt.MaxIdleTime)

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return s, nil
}

// JSONSerializer is a serializer for JSON.
type JSONSerializer struct{}

// Scan scans the value from the database.
func (JSONSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	fieldValue := reflect.New(field.FieldType)
	if dbValue != nil {
		var bytes []byte
		switch v := dbValue.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		default:
			return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
		}
		err = json.Unmarshal(bytes, fieldValue.Interface())
	}
	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
	return
}

// Value returns the value to be stored in the database.
func (JSONSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}

// CreateTable creates the table.
func (s *SqlStore) CreateTable(tables func() []schema.Tabler) error {
	schema.RegisterSerializer("json", JSONSerializer{})
	if s.autoMigrate {
		var errs []string
		for _, t := range tables() {
			if !s.db.Migrator().HasTable(t) {
				// Create the table.
				if err := s.db.Migrator().CreateTable(t); err != nil {
					errs = append(errs, err.Error())
				}
			} else {
				// Migrate the table.
				if err := s.db.AutoMigrate(t); err != nil {
					errs = append(errs, err.Error())
				}
			}
			if len(errs) > 0 {
				return fmt.Errorf("failed to shut down server: [%s]", strings.Join(errs, ","))
			}
		}
	}

	schema.RegisterSerializer("json", JSONSerializer{})
	return nil
}

// DB returns the underlying database.
func (s *SqlStore) DB() *gorm.DB {
	return s.db
}

// Start starts the SQL store.
func (s *SqlStore) Start() {}

// Stop closes the underlying database.
func (s *SqlStore) Stop() {
	db, err := s.db.DB()
	if err != nil {
		panic(err)
	}
	if err := db.Close(); err != nil {
		panic(err)
	}
}
