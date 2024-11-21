// #############################################################################
// # File: options.go                                                          #
// # Project: zlog                                                             #
// # Created Date: 2024/11/21 17:15:14                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2024/11/21 17:47:48                                        #
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
			logger := z.logger.With(
				zap.String("traceID", tc.TraceID),
				zap.String("spanID", tc.SpanID),
				zap.String("parentSpanID", tc.ParentSpanID),
			)

			newZlog := newZLog(z.cfg, z.options...)
			newZlog.logger = logger
			return newZlog, nil
		}
		return z, nil
	}
}
