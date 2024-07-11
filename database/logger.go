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
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type slog struct {
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewSlog(config logger.Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "[%.3fms] [rows:%v] %s"
		traceWarnStr = "[%s][%.3fms] [rows:%v] %s"
		traceErrStr  = "[%s] [%.3fms] [rows:%v] %s"
	)
	return &slog{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceErrStr:  traceErrStr,
		traceWarnStr: traceWarnStr,
	}
}

func (l *slog) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *slog) Info(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Info("")
	if l.LogLevel >= logger.Info {
		logx.WithContext(ctx).Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (l *slog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logx.WithContext(ctx).Infof(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (l *slog) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logx.WithContext(ctx).Errorf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (l *slog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logx.WithContext(ctx).WithFields(logx.Field("db_caller", utils.FileWithLineNum())).Errorf(l.traceErrStr, err.Error(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logx.WithContext(ctx).WithFields(logx.Field("db_caller", utils.FileWithLineNum())).Errorf(l.traceErrStr, err.Error(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logx.WithContext(ctx).WithFields(logx.Field("db_caller", utils.FileWithLineNum())).Infof(l.traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logx.WithContext(ctx).WithFields(logx.Field("db_caller", utils.FileWithLineNum())).Infof(l.traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			logx.WithContext(ctx).WithFields(logx.Field("db_caller", utils.FileWithLineNum())).Infof(l.traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logx.WithContext(ctx).WithFields(logx.Field("db_caller", utils.FileWithLineNum())).Infof(l.traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
