// #############################################################################
// # File: zap_log_test.go                                                     #
// # Project: zlog                                                             #
// # Created Date: 2024/10/08 18:04:40                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2024/11/21 17:51:10                                        #
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
	zlog.InitZLog(&zlog.ZLogConfig{
		Compress: true,
		LogMode:  "file|console",
		Encoding: "json",
		LogFile:  "./logs/zlog.log",
	})
	zlog.ZLog().Infof("hello %s", "realjf")
}

func TestZapLogWithTrace(t *testing.T) {
	tc := trace.NewTraceContext()
	ctx := trace.WithTraceContext(context.Background(), tc)
	zlog.InitZLog(&zlog.ZLogConfig{
		Compress: true,
		LogMode:  "file|console",
		Encoding: "json",
		LogFile:  "./logs/zlog.log",
	})
	zlog.ZLog().InfofWithTrace(ctx, "hello %s", "realjf")

	childCtx := trace.StartSpan(ctx)
	zlog.ZLog().InfofWithTrace(childCtx, "hello %s", "child-realjf")

	child2Ctx := trace.StartSpan(childCtx)
	zlog.ZLog().InfofWithTrace(child2Ctx, "hello %s", "child2-realjf")
}
