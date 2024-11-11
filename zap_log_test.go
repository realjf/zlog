// #############################################################################
// # File: zap_log_test.go                                                     #
// # Project: zlog                                                             #
// # Created Date: 2024/10/08 18:04:40                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2024/11/11 11:19:27                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package zlog_test

import (
	"testing"

	"github.com/realjf/zlog"
)

func TestZapLog(t *testing.T) {
	zlog.InitZLog(&zlog.ZLogConfig{
		Compress: true,
		LogMode:  "file",
		Encoding: "json",
		LogFile:  "./zlog.log",
	})
	zlog.ZLog().Infof("hello %s", "realjf")
}
