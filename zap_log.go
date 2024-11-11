// #############################################################################
// # File: zap_log.go                                                          #
// # Project: zlog                                                             #
// # Created Date: 2024/10/08 15:18:55                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2024/11/11 11:22:11                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package zlog

import (
	"log"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/realjf/zlog/utils/fileutil"
)

const (
	logMaxSize         = 500
	logMaxAge          = 30
	logEncodingConsole = "console"
	logEncodingJson    = "json"
	logModeFile        = "file"
	logModeStdout      = "console"
)

type IZLog interface {
	Debug(msg string, fields ...zapcore.Field)
	Debugf(template string, args ...interface{})
	Info(msg string, fields ...zapcore.Field)
	Infof(template string, args ...interface{})
	Warn(msg string, fields ...zapcore.Field)
	Warnf(template string, args ...interface{})
	Error(msg string, fields ...zapcore.Field)
	Errorf(template string, args ...interface{})
	Fatal(msg string, fields ...zapcore.Field)
	Fatalf(template string, args ...interface{})
}

var localZLog *zLog

func init() {
	localZLog = newZLog(&ZLogConfig{
		Compress: true,
	})
}

func ZLog() IZLog {
	if localZLog == nil {
		panic("localZLog is nil")
	}
	return localZLog
}

func InitZLog(config *ZLogConfig, options ...zap.Option) {
	localZLog = newZLog(config, options...)
}

// =========================================================== 结构体 ===========================================================

type ZLogConfig struct {
	Level    LogLevel `yaml:"level"`    // 日志级别： debug|info|warn|error|fatal
	LogMode  string   `yaml:"log_mode"` // 日志模式 console|file
	MaxSize  int      `yaml:"max_size"` // 单日志文件最大字节/M
	MaxAge   int      `yaml:"max_age"`  // 日志文件最大存活天数
	Compress bool     `yaml:"compress"` // 是否启用压缩
	Encoding string   `yaml:"encoding"` // 日志编码 console|json
	LogFile  string   `yaml:"log_file"` // 日志文件路径
}

type zLog struct {
	logger *zap.Logger
}

// =========================================================== 构造方法 ===========================================================

func newZLog(config *ZLogConfig, options ...zap.Option) *zLog {
	var logger *zap.Logger

	if config.LogMode == logModeFile {
		dir := filepath.Dir(config.LogFile)
		if err := fileutil.MkdirIfNecessary(dir); err != nil {
			log.Panicf("创建日志目录[%s]失败：%+v\n", dir, errors.WithStack(err))
		}
		logger = newZLogWithFile(config, options...)
	} else {
		logger = newZLogWithConsole(config, options...)
	}

	return &zLog{
		logger: logger,
	}
}

func newZLogWithConsole(config *ZLogConfig, options ...zap.Option) (logger *zap.Logger) {
	conf := zap.Config{
		Level:            zap.NewAtomicLevelAt(config.Level.toZapLevel()),
		EncoderConfig:    newEncoderConfig(),
		Encoding:         logEncodingConsole,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if config.Encoding == logEncodingJson {
		conf.Encoding = logEncodingJson
	}
	logger, err := conf.Build(options...)
	if err != nil {
		log.Printf("Zap日志创建失败，使用NewExample创建\n")
		logger = zap.NewExample(options...)
	}

	return
}

func newZLogWithFile(config *ZLogConfig, options ...zap.Option) (logger *zap.Logger) {
	if config.MaxAge <= 0 {
		config.MaxAge = logMaxAge
	}
	if config.MaxSize <= 0 {
		config.MaxSize = logMaxSize
	}
	filename := filepath.Base(config.LogFile)
	hook := lumberjack.Logger{
		Filename:  filename,
		MaxSize:   config.MaxSize,
		MaxAge:    config.MaxAge,
		Compress:  config.Compress,
		LocalTime: true,
	}
	var encoder zapcore.Encoder
	if config.Encoding == logEncodingJson {
		encoder = zapcore.NewJSONEncoder(newEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(newEncoderConfig())
	}
	core := zapcore.NewCore(encoder, zapcore.AddSync(&hook), config.Level.toZapLevel())
	logger = zap.New(core)
	return
}

// =========================================================== 接口方法 ===========================================================

func (z *zLog) Debug(msg string, fields ...zapcore.Field) {
	z.logger.Debug(msg, fields...)
}

func (z *zLog) Debugf(template string, args ...interface{}) {
	z.logger.Sugar().Debugf(template, args...)
}

func (z *zLog) Info(msg string, fields ...zapcore.Field) {
	z.logger.Info(msg, fields...)
}

func (z *zLog) Infof(template string, args ...interface{}) {
	z.logger.Sugar().Infof(template, args...)
}

func (z *zLog) Warn(msg string, fields ...zapcore.Field) {
	z.logger.Warn(msg, fields...)
}

func (z *zLog) Warnf(template string, args ...interface{}) {
	z.logger.Sugar().Warnf(template, args...)
}

func (z *zLog) Error(msg string, fields ...zapcore.Field) {
	z.logger.Error(msg, fields...)
}

func (z *zLog) Errorf(template string, args ...interface{}) {
	z.logger.Sugar().Errorf(template, args...)
}

func (z *zLog) Fatal(msg string, fields ...zapcore.Field) {
	z.logger.Fatal(msg, fields...)
}

func (z *zLog) Fatalf(template string, args ...interface{}) {
	z.logger.Sugar().Fatalf(template, args...)
}

// =========================================================== 私有方法 ===========================================================

func newEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.CallerKey = ""
	return encoderConfig
}
