// #############################################################################
// # File: zap_log_test.go                                                     #
// # Project: zlog                                                             #
// # Created Date: 2024/10/08 18:04:40                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2025/06/05 07:47:37                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package zlog_test

import (
	"context"
	"testing"

	"github.com/realjf/zlog"
	"github.com/realjf/zlog/trace"
)

func TestZapLog(t *testing.T) {
	zlog.InitZLog([]*zlog.ZLogConfig{
		{
			Compress: true,
			LogMode:  "file|console",
			Encoding: "json",
			LogFile:  "./logs/zlog.log",
		},
	})
	zlog.ZLog().Infof("hello %s", "realjf")
}

func TestZapLogWithName(t *testing.T) {
	zlog.InitZLog([]*zlog.ZLogConfig{
		{
			Compress: true,
			LogMode:  "file|console",
			Encoding: "json",
			LogFile:  "./logs/zlog.log",
			Name:     "zlog",
			Default:  true,
		},
		{
			Compress: true,
			LogMode:  "file|console",
			Encoding: "json",
			LogFile:  "./logs/zlog2.log",
			Name:     "zlog2",
			Default:  false,
		},
	})
	zlog.ZLog().Infof("hello %s", "realjf")
	zlog.ZLog().WithName("zlog2").Infof("hello %s", "realjf2")
	go zlog.ZLog().WithName("zlog2").Infof("hello %s", "realjf3")
	zlog.ZLog().Infof("hello %s", "realjf4")
}

func TestZapLogWithTrace(t *testing.T) {
	tc := trace.NewTraceContext()
	ctx := trace.WithTraceContext(context.Background(), tc)
	zlog.InitZLog([]*zlog.ZLogConfig{
		{
			Compress: true,
			LogMode:  "file|console",
			Encoding: "json",
			LogFile:  "./logs/zlog.log",
		},
	})
	zlog.ZLog().InfofWithTrace(ctx, "hello %s", "realjf")

	childCtx := trace.StartSpan(ctx)
	zlog.ZLog().InfofWithTrace(childCtx, "hello %s", "child-realjf")

	child2Ctx := trace.StartSpan(childCtx)
	zlog.ZLog().InfofWithTrace(child2Ctx, "hello %s", "child2-realjf")
}

func TestWithPrefix(t *testing.T) {
	zlog.InitZLog([]*zlog.ZLogConfig{
		{
			Compress: true,
			LogMode:  "file|console",
			Encoding: "json",
			LogFile:  "./logs/zlog.log",
		},
	})
	zlog.ZLog().WithPrefix("[test]").Infof("hello %s", "realjf")
}

func TestWithPrefix2(t *testing.T) {
	zlog.InitZLog([]*zlog.ZLogConfig{
		{
			Compress: true,
			LogMode:  "file|console",
			Encoding: "json",
			LogFile:  "./logs/zlog.log",
			Name:     "zlog",
			Default:  true,
		},
		{
			Compress: true,
			LogMode:  "file|console",
			Encoding: "json",
			LogFile:  "./logs/zlog2.log",
			Name:     "zlog2",
			Default:  false,
		},
	})
	zlog.ZLog().WithPrefix("[test]").Infof("hello %s", "realjf")
	go zlog.ZLog().WithPrefix("[test2]").Infof("hello %s", "realjf2")
	zlog.ZLog().WithPrefix("[test3]").Infof("hello %s", "realjf3")
}
