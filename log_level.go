// #############################################################################
// # File: log_level.go                                                        #
// # Project: zlog                                                             #
// # Created Date: 2024/10/08 15:32:56                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2024/10/09 13:49:48                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel string

const (
	logLevelDebug LogLevel = "debug"
	logLevelInfo  LogLevel = "info"
	logLevelWarn  LogLevel = "warn"
	logLevelError LogLevel = "error"
	logLevelFatal LogLevel = "fatal"
)

func (l LogLevel) String() string {
	return string(l)
}

func (l LogLevel) toZapLevel() (level zapcore.Level) {
	switch l {
	case logLevelInfo:
		level = zap.InfoLevel
	case logLevelWarn:
		level = zap.WarnLevel
	case logLevelError:
		level = zap.ErrorLevel
	case logLevelFatal:
		level = zap.FatalLevel
	default:
		level = zap.DebugLevel
	}
	return level
}
