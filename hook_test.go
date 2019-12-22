package otredis

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/stretchr/testify/assert"
)

func TestHook_BeforeProcess(t *testing.T) {
	assert := assert.New(t)
	opentracing.SetGlobalTracer(mocktracer.New())

	tests := map[string]struct {
		ctx context.Context
		hk  hook
		cmd redis.Cmder
	}{
		"success": {
			ctx: context.TODO(),
			hk:  hook{addrs: []string{"127.0.0.1:6379", "127.0.0.1:6378"}, database: 10},
			cmd: redis.NewStringCmd("GET", "ALLAN"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			newCtx, err := tc.hk.BeforeProcess(tc.ctx, tc.cmd)
			assert.NoError(err)
			span := opentracing.SpanFromContext(newCtx)
			tags := span.(*mocktracer.MockSpan).Tags()
			for k, v := range tags {
				switch k {
				case string(ext.DBType):
					assert.Equal("redis", v.(string))
				case string(ext.DBInstance):
					assert.Equal("10", v.(string))
				case string(ext.PeerAddress):
					assert.Equal("127.0.0.1:6379, 127.0.0.1:6378", v)
				case string(ext.PeerService):
					assert.Equal("redis", v.(string))
				case string(ext.SpanKind):
					assert.Equal("client", string(v.(ext.SpanKindEnum)))
				case string(ext.DBStatement):
					assert.Equal("GET", v)
				default:
					panic("unknown tag")
				}
			}
		})
	}
}

func TestHook_AfterProcess(t *testing.T) {
	assert := assert.New(t)
	opentracing.SetGlobalTracer(mocktracer.New())

	failCtx, failCancel := context.WithCancel(context.Background())
	failCancel()

	tests := map[string]struct {
		ctx    context.Context
		hk     hook
		cmd    redis.Cmder
		ctxErr error
	}{
		"fail": {
			ctx:    failCtx,
			hk:     hook{addrs: []string{"127.0.0.1:6379", "127.0.0.1:6378"}, database: 10},
			cmd:    redis.NewStringCmd("GET", "ALLAN"),
			ctxErr: context.Canceled,
		},
		"success": {
			ctx: context.TODO(),
			hk:  hook{addrs: []string{"127.0.0.1:6379", "127.0.0.1:6378"}, database: 10},
			cmd: redis.NewStringCmd("GET", "ALLAN"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// before
			newCtx, err := tc.hk.BeforeProcess(tc.ctx, tc.cmd)
			assert.NoError(err)
			now := time.Now()
			time.Sleep(1 * time.Millisecond)

			// after
			err = tc.hk.AfterProcess(newCtx, tc.cmd)
			assert.NoError(err)
			span := opentracing.SpanFromContext(newCtx).(*mocktracer.MockSpan)
			assert.True(span.FinishTime.UnixNano() > now.UnixNano())
			if tc.ctxErr != nil {
				records := span.Logs()
				assert.Len(records, 1)
				assert.Len(records[0].Fields, 1)
				field := records[0].Fields[0]
				assert.Equal("error", field.Key)
				assert.Equal(tc.ctxErr.Error(), field.ValueString)
			}
		})
	}
}

func TestHook_BeforeProcessPipeline(t *testing.T) {
	assert := assert.New(t)
	opentracing.SetGlobalTracer(mocktracer.New())

	tests := map[string]struct {
		ctx  context.Context
		hk   hook
		cmds []redis.Cmder
	}{
		"success": {
			ctx:  context.TODO(),
			hk:   hook{addrs: []string{"127.0.0.1:6379", "127.0.0.1:6378"}, database: 10},
			cmds: []redis.Cmder{redis.NewStringCmd("GET", "ALLAN"), redis.NewStringCmd("SET", "ALLAN")},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			newCtx, err := tc.hk.BeforeProcessPipeline(tc.ctx, tc.cmds)
			assert.NoError(err)
			span := opentracing.SpanFromContext(newCtx)
			tags := span.(*mocktracer.MockSpan).Tags()
			for k, v := range tags {
				switch k {
				case string(ext.DBType):
					assert.Equal("redis", v.(string))
				case string(ext.DBInstance):
					assert.Equal("10", v.(string))
				case string(ext.PeerAddress):
					assert.Equal("127.0.0.1:6379, 127.0.0.1:6378", v)
				case string(ext.PeerService):
					assert.Equal("redis", v.(string))
				case string(ext.SpanKind):
					assert.Equal("client", string(v.(ext.SpanKindEnum)))
				case string(ext.DBStatement):
					assert.Equal("GET --> SET", v)
				default:
					panic("unknown tag")
				}
			}
		})
	}
}

func TestHook_AfterProcessPipeline(t *testing.T) {
	assert := assert.New(t)
	opentracing.SetGlobalTracer(mocktracer.New())

	failCtx, failCancel := context.WithCancel(context.Background())
	failCancel()

	tests := map[string]struct {
		ctx    context.Context
		hk     hook
		cmds   []redis.Cmder
		ctxErr error
	}{
		"fail": {
			ctx:    failCtx,
			hk:     hook{addrs: []string{"127.0.0.1:6379", "127.0.0.1:6378"}, database: 10},
			cmds:   []redis.Cmder{redis.NewStringCmd("GET", "ALLAN"), redis.NewStringCmd("SET", "ALLAN")},
			ctxErr: context.Canceled,
		},
		"success": {
			ctx:  context.TODO(),
			hk:   hook{addrs: []string{"127.0.0.1:6379", "127.0.0.1:6378"}, database: 10},
			cmds: []redis.Cmder{redis.NewStringCmd("GET", "ALLAN"), redis.NewStringCmd("SET", "ALLAN")},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// before
			newCtx, err := tc.hk.BeforeProcessPipeline(tc.ctx, tc.cmds)
			assert.NoError(err)
			assert.NoError(err)
			now := time.Now()
			time.Sleep(1 * time.Millisecond)

			// after
			err = tc.hk.AfterProcessPipeline(newCtx, tc.cmds)
			assert.NoError(err)
			span := opentracing.SpanFromContext(newCtx).(*mocktracer.MockSpan)
			assert.True(span.FinishTime.UnixNano() > now.UnixNano())
			if tc.ctxErr != nil {
				records := span.Logs()
				assert.Len(records, 1)
				assert.Len(records[0].Fields, 1)
				field := records[0].Fields[0]
				assert.Equal("error", field.Key)
				assert.Equal(tc.ctxErr.Error(), field.ValueString)
			}
		})
	}

}
