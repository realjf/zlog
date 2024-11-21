# zlog
log tool


### Quick Start

```go
func main() {
    	zlog.InitZLog(&zlog.ZLogConfig{
		Compress: true,
		LogMode:  "file|console",
		Encoding: "json",
		LogFile:  "./logs/zlog.log",
	})
	zlog.ZLog().Infof("hello %s", "realjf")
}

```

### Start With Trace Context

```go
func main() {
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

```
