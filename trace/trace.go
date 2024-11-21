// #############################################################################
// # File: trace.go                                                            #
// # Project: trace                                                            #
// # Created Date: 2024/11/21 17:08:11                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2024/11/21 17:08:21                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package trace

import (
	"context"
	"crypto/rand"
	"fmt"
)

// =========================================================== 链路追踪方法 ===========================================================

type TraceContext struct {
	TraceID      string // 标识一个完整的分布式追踪链路，如一个请求的整个生命周期
	SpanID       string // 标识追踪链路中的一个单独的操作或阶段，如：http请求、数据库查询等
	ParentSpanID string // 标识上游服务的spanId
}

func NewTraceContext() *TraceContext {
	traceId := GenerateTraceID()
	spanId := GenerateSpanID()
	return &TraceContext{
		TraceID: traceId,
		SpanID:  spanId,
	}
}

func FromContext(ctx context.Context) (*TraceContext, bool) {
	tc, ok := ctx.Value(TraceContext{}).(*TraceContext)
	return tc, ok
}

func WithTraceContext(ctx context.Context, tc *TraceContext) context.Context {
	return context.WithValue(ctx, TraceContext{}, tc)
}

func StartSpan(ctx context.Context) context.Context {
	parentTC, _ := FromContext(ctx)
	childTC := &TraceContext{
		ParentSpanID: parentTC.SpanID,
		TraceID:      parentTC.TraceID,
		SpanID:       GenerateSpanID(),
	}
	return WithTraceContext(ctx, childTC)
}

// GenerateTraceID generates a 32-byte random traceID
func GenerateTraceID() string {
	var traceID [32]byte
	_, _ = rand.Read(traceID[:])
	return fmt.Sprintf("%x", traceID)
}

// GenerateSpanID generates an 16-byte random spanID
func GenerateSpanID() string {
	var spanID [16]byte
	_, _ = rand.Read(spanID[:])
	return fmt.Sprintf("%x", spanID)
}
