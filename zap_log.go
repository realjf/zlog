// #############################################################################
// # File: zap_log.go                                                          #
// # Project: zlog                                                             #
// # Created Date: 2024/10/08 15:18:55                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2025/06/10 10:46:44                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package zlog

import (
	"context"
	"log"
	"path/filepath"
	"strings"
	"sync"
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

	DebugWithTrace(ctx context.Context, msg string, fields ...zapcore.Field)
	DebugfWithTrace(ctx context.Context, template string, args ...interface{})
	InfoWithTrace(ctx context.Context, msg string, fields ...zapcore.Field)
	InfofWithTrace(ctx context.Context, template string, args ...interface{})
	WarnWithTrace(ctx context.Context, msg string, fields ...zapcore.Field)
	WarnfWithTrace(ctx context.Context, template string, args ...interface{})
	ErrorWithTrace(ctx context.Context, msg string, fields ...zapcore.Field)
	ErrorfWithTrace(ctx context.Context, template string, args ...interface{})
	FatalWithTrace(ctx context.Context, msg string, fields ...zapcore.Field)
	FatalfWithTrace(ctx context.Context, template string, args ...interface{})

	WithPrefix(prefix string) IZLog
	WithName(name ...string) IZLog

	GetZCore(name string) *zap.Logger
}

var localZLog *zLog

func init() {
	localZLog = newZLog([]*ZLogConfig{
		{
			Compress: true,
		},
	})
}

func ZLog() IZLog {
	if localZLog == nil {
		panic("localZLog is nil")
	}
	return localZLog
}

func InitZLog(configs []*ZLogConfig, options ...zap.Option) {
	localZLog = newZLog(configs, options...)
}

// =========================================================== 结构体 ===========================================================

type ZLogConfig struct {
	Level      LogLevel `yaml:"level"`       // 日志级别： debug|info|warn|error|fatal
	LogMode    string   `yaml:"log_mode"`    // 日志模式 console|file
	MaxSize    int      `yaml:"max_size"`    // 单日志文件最大字节/M
	MaxAge     int      `yaml:"max_age"`     // 日志文件最大存活天数
	MaxBackups int      `yaml:"max_backups"` // 日志文件最大数
	Compress   bool     `yaml:"compress"`    // 是否启用压缩
	Encoding   string   `yaml:"encoding"`    // 日志编码 console|json
	LogFile    string   `yaml:"log_file"`    // 日志文件路径
	Name       string   `yaml:"name"`        // 日志名称
	Default    bool     `yaml:"default"`     // 默认日志记录器
}

type zLog struct {
	loggers map[string]*zap.Logger
	cfgs    map[string]*ZLogConfig
	options []zap.Option

	lock        sync.Mutex
	usedLoggers map[string]*zap.Logger

	prefix string
}

// =========================================================== 构造方法 ===========================================================

func NewZLog(configs []*ZLogConfig, options ...zap.Option) IZLog {
	return newZLog(configs, options...)
}

func newZLog(configs []*ZLogConfig, options ...zap.Option) *zLog {
	cfgs := make(map[string]*ZLogConfig)
	loggers := make(map[string]*zap.Logger)
	usedLoggers := make(map[string]*zap.Logger, 0)
	for _, config := range configs {
		var logger *zap.Logger

		var err error
		config.LogFile, err = filepath.Abs(config.LogFile)
		if err != nil {
			log.Panicf("获取日志文件绝对路径失败：%v\n", err.Error())
		}

		if strings.Contains(config.LogMode, logModeFile) && strings.Contains(config.LogMode, logModeStdout) {
			logger = newZLogWithFileAndConsole(config, options...)
		} else if config.LogMode == logModeFile {
			logger = newZLogWithFile(config, options...)
		} else {
			logger = newZLogWithConsole(config, options...)
		}
		cfgs[config.Name] = config
		loggers[config.Name] = logger
		if config.Default {
			usedLoggers[config.Name] = logger
		}
	}

	if len(usedLoggers) == 0 {
		// use first logger as default
		cfgs[configs[0].Name].Default = true
		usedLoggers[configs[0].Name] = loggers[configs[0].Name]
	}

	return &zLog{
		loggers:     loggers,
		cfgs:        cfgs,
		usedLoggers: usedLoggers,
		options:     options,
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
	core := newFileCore(config, options...)
	logger = zap.New(core)
	return
}

func newFileCore(config *ZLogConfig, options ...zap.Option) zapcore.Core {
	dir := filepath.Dir(config.LogFile)
	if err := fileutil.MkdirIfNecessary(dir); err != nil {
		log.Panicf("创建日志目录[%s]失败：%+v\n", dir, errors.WithStack(err))
	}
	if config.MaxAge <= 0 {
		config.MaxAge = logMaxAge
	}
	if config.MaxSize <= 0 {
		config.MaxSize = logMaxSize
	}
	// filename := filepath.Base(config.LogFile)
	hook := lumberjack.Logger{
		Filename:   config.LogFile,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
		LocalTime:  true,
	}
	var encoder zapcore.Encoder
	if config.Encoding == logEncodingJson {
		encoder = zapcore.NewJSONEncoder(newEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(newEncoderConfig())
	}
	core := zapcore.NewCore(encoder, zapcore.AddSync(&hook), config.Level.toZapLevel())
	return core
}

func newZLogWithFileAndConsole(config *ZLogConfig, options ...zap.Option) (logger *zap.Logger) {
	consoleCore := newZLogWithConsole(config, options...)
	fileCore := newFileCore(config, options...)

	core := zapcore.NewTee(consoleCore.Core(), fileCore)
	logger = zap.New(core)

	return
}

// =========================================================== 接口方法 ===========================================================

func (z *zLog) Debug(msg string, fields ...zapcore.Field) {
	z.withName(func(logger *zap.Logger) {
		logger.Debug(z.withPrefix(msg), fields...)
	})
}

func (z *zLog) Debugf(template string, args ...interface{}) {
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Debugf(z.withPrefix(template), args...)
	})
}

func (z *zLog) Info(msg string, fields ...zapcore.Field) {
	z.withName(func(logger *zap.Logger) {
		logger.Info(z.withPrefix(msg), fields...)
	})
}

func (z *zLog) Infof(template string, args ...interface{}) {
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Infof(z.withPrefix(template), args...)
	})
}

func (z *zLog) Warn(msg string, fields ...zapcore.Field) {
	z.withName(func(logger *zap.Logger) {
		logger.Warn(z.withPrefix(msg), fields...)
	})
}

func (z *zLog) Warnf(template string, args ...interface{}) {
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Warnf(z.withPrefix(template), args...)
	})
}

func (z *zLog) Error(msg string, fields ...zapcore.Field) {
	z.withName(func(logger *zap.Logger) {
		logger.Error(z.withPrefix(msg), fields...)
	})
}

func (z *zLog) Errorf(template string, args ...interface{}) {
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Errorf(z.withPrefix(template), args...)
	})
}

func (z *zLog) Fatal(msg string, fields ...zapcore.Field) {
	z.withName(func(logger *zap.Logger) {
		logger.Fatal(z.withPrefix(msg), fields...)
	})
}

func (z *zLog) Fatalf(template string, args ...interface{}) {
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Fatalf(z.withPrefix(template), args...)
	})
}

// =========================================================== 带链路追踪的接口方法 ===========================================================

func (z *zLog) DebugWithTrace(ctx context.Context, msg string, fields ...zapcore.Field) {
	msg = z.withPrefix(msg)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Debug(msg, fields...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Debug(msg, fields...)
	})
}

func (z *zLog) DebugfWithTrace(ctx context.Context, template string, args ...interface{}) {
	template = z.withPrefix(template)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Sugar().Debugf(template, args...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Debugf(template, args...)
	})
}

func (z *zLog) InfoWithTrace(ctx context.Context, msg string, fields ...zapcore.Field) {
	msg = z.withPrefix(msg)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Info(msg, fields...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Info(msg, fields...)
	})
}

func (z *zLog) InfofWithTrace(ctx context.Context, template string, args ...interface{}) {
	template = z.withPrefix(template)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Sugar().Infof(template, args...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Infof(template, args...)
	})
}

func (z *zLog) WarnWithTrace(ctx context.Context, msg string, fields ...zapcore.Field) {
	msg = z.withPrefix(msg)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Warn(msg, fields...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Warn(msg, fields...)
	})
}

func (z *zLog) WarnfWithTrace(ctx context.Context, template string, args ...interface{}) {
	template = z.withPrefix(template)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Sugar().Warnf(template, args...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Warnf(template, args...)
	})
}

func (z *zLog) ErrorWithTrace(ctx context.Context, msg string, fields ...zapcore.Field) {
	msg = z.withPrefix(msg)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Error(msg, fields...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Error(msg, fields...)
	})
}

func (z *zLog) ErrorfWithTrace(ctx context.Context, template string, args ...interface{}) {
	template = z.withPrefix(template)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Sugar().Errorf(template, args...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Errorf(template, args...)
	})
}

func (z *zLog) FatalWithTrace(ctx context.Context, msg string, fields ...zapcore.Field) {
	msg = z.withPrefix(msg)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Fatal(msg, fields...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Fatal(msg, fields...)
	})
}

func (z *zLog) FatalfWithTrace(ctx context.Context, template string, args ...interface{}) {
	template = z.withPrefix(template)
	opt := WithTrace(ctx)
	if nz, err := opt(z); err == nil {
		nz.withName(func(logger *zap.Logger) {
			logger.Sugar().Fatalf(template, args...)
		})
		return
	}
	z.withName(func(logger *zap.Logger) {
		logger.Sugar().Fatalf(template, args...)
	})
}

func (z *zLog) GetZCore(name string) *zap.Logger {
	return z.loggers[name]
}

// =========================================================== 带前缀打印的接口方法 ===========================================================

func (z *zLog) WithPrefix(prefix string) IZLog {
	z.lock.Lock()
	defer z.lock.Unlock()

	cfgs := make([]*ZLogConfig, 0)
	for _, cfg := range z.cfgs {
		cfgs = append(cfgs, cfg)
	}

	newZlog := newZLog(cfgs, z.options...)
	newZlog.prefix = prefix
	return newZlog
}

// 使用指定日志记录器
func (z *zLog) WithName(names ...string) IZLog {
	z.lock.Lock()
	defer z.lock.Unlock()

	cfgs := make([]*ZLogConfig, 0)
	for _, cfg := range z.cfgs {
		cfgs = append(cfgs, cfg)
	}

	newZlog := newZLog(cfgs, z.options...)
	usedLoggers := make(map[string]*zap.Logger, 0)
	for _, name := range names {
		if logger, ok := z.loggers[name]; ok {
			usedLoggers[name] = logger
		}
	}
	newZlog.usedLoggers = usedLoggers

	return newZlog
}

// =========================================================== 私有方法 ===========================================================

func (z *zLog) withName(f func(logger *zap.Logger)) {
	z.lock.Lock()
	defer z.lock.Unlock()

	for _, logger := range z.usedLoggers {
		f(logger)
	}

	z.resetUsedLogger()
}

func (z *zLog) resetUsedLogger() {
	z.usedLoggers = make(map[string]*zap.Logger, 0)
	for name, config := range z.cfgs {
		if config.Default {
			z.usedLoggers[name] = z.loggers[name]
		}
	}
}

func (z *zLog) withPrefix(original string) string {
	if z.prefix != "" {
		original = z.prefix + " " + original
	}
	return original
}

func newEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.CallerKey = ""
	return encoderConfig
}
