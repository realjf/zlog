// #############################################################################
// # File: std_log.go                                                          #
// # Project: zlog                                                             #
// # Created Date: 2024/10/08 15:19:03                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2024/10/08 17:58:32                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package zlog

import "gopkg.in/natefinch/lumberjack.v2"

func NewLumberjackLogger(config *ZLogConfig) *lumberjack.Logger {
	if config.MaxSize <= 0 {
		config.MaxSize = logMaxSize
	}
	if config.MaxSize <= 0 {
		config.MaxAge = logMaxAge
	}
	loghook := lumberjack.Logger{ //定义日志分割器
		Filename:  config.LogFile,  // 日志文件路径
		MaxSize:   config.MaxSize,  // 文件最大M字节
		MaxAge:    config.MaxAge,   // 最多保留几天
		Compress:  config.Compress, // 是否压缩
		LocalTime: true,
	}
	return &loghook
}
