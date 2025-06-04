// #############################################################################
// # File: options.go                                                          #
// # Project: zlog                                                             #
// # Created Date: 2024/11/21 17:15:14                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2025/06/05 00:08:10                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package zlog

import (
	"context"

	"github.com/realjf/zlog/trace"
	"go.uber.org/zap"
)

type Option func(*zLog) (*zLog, error)

func WithTrace(ctx context.Context) Option {
	return func(z *zLog) (*zLog, error) {
		if tc, ok := trace.FromContext(ctx); ok {
			cfgs := make([]*ZLogConfig, 0)
			for _, cfg := range z.cfgs {
				cfgs = append(cfgs, cfg)
			}

			newZlog := newZLog(cfgs, z.options...)
			newZlog.usedLoggers = make(map[string]*zap.Logger, 0)
			for name, _ := range z.loggers {
				logger := z.loggers[name].With(
					zap.String("traceID", tc.TraceID),
					zap.String("spanID", tc.SpanID),
					zap.String("parentSpanID", tc.ParentSpanID),
				)
				newZlog.loggers[name] = logger
				if newZlog.cfgs[name].Default {
					newZlog.usedLoggers[name] = logger
				}
			}

			if len(newZlog.usedLoggers) == 0 {
				newZlog.cfgs[cfgs[0].Name].Default = true
				newZlog.usedLoggers[cfgs[0].Name] = newZlog.loggers[cfgs[0].Name]
			}

			return newZlog, nil
		}
		return z, nil
	}
}
