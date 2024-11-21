// #############################################################################
// # File: trace_test.go                                                       #
// # Project: trace                                                            #
// # Created Date: 2024/11/21 17:08:36                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2024/11/21 17:14:06                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package trace_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/realjf/zlog/trace"
)

func TestTrace(t *testing.T) {
	tc := trace.NewTraceContext()
	ctx := trace.WithTraceContext(context.Background(), tc)

	if tc, ok := trace.FromContext(ctx); ok {
		fmt.Printf("initial traceID[%s], spanID[%s]\n", tc.TraceID, tc.SpanID)
	}

	childCtx := trace.StartSpan(ctx)
	if childTC, ok := trace.FromContext(childCtx); ok {
		fmt.Printf("child traceID[%s], child spanID[%s], parent spanID[%s]\n", childTC.TraceID, childTC.SpanID, childTC.ParentSpanID)
	}

	child2Ctx := trace.StartSpan(ctx)
	if child2TC, ok := trace.FromContext(child2Ctx); ok {
		fmt.Printf("child2 traceID[%s], child2 spanID[%s], parent spanID[%s]\n", child2TC.TraceID, child2TC.SpanID, child2TC.ParentSpanID)
	}

	child3Ctx := trace.StartSpan(childCtx)
	if child3TC, ok := trace.FromContext(child3Ctx); ok {
		fmt.Printf("child2 traceID[%s], child2 spanID[%s], parent spanID[%s]\n", child3TC.TraceID, child3TC.SpanID, child3TC.ParentSpanID)
	}
}
