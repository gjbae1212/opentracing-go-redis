package otredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
)

type UniversalClient interface {
	redis.UniversalClient

	// withContext is to inject context and to add hook.
	// it is an internal method.
	withContext(ctx context.Context) UniversalClient
}

type redisClient struct {
	*redis.Client
}

// WithContext is to inject context and to add hook.
func (rc *redisClient) withContext(ctx context.Context) UniversalClient {
	opts := rc.Client.Options()
	rc.Client = rc.Client.WithContext(ctx)
	rc.AddHook(hook{addrs: []string{opts.Addr}, database: opts.DB})
	return rc
}

type redisClusterClient struct {
	*redis.ClusterClient
}

// WithContext is to inject context and to add hook.
func (rc *redisClusterClient) withContext(ctx context.Context) UniversalClient {
	rc.ClusterClient = rc.ClusterClient.WithContext(ctx)
	rc.AddHook(hook{addrs: rc.ClusterClient.Options().Addrs, database: 0})
	return rc
}

type redisRing struct {
	*redis.Ring
}

// WithContext is to inject context and to add hook.
func (rc *redisRing) withContext(ctx context.Context) UniversalClient {
	opts := rc.Ring.Options()
	rc.Ring = rc.Ring.WithContext(ctx)
	addrs := make([]string, len(opts.Addrs))
	i := 0
	for _, v := range opts.Addrs {
		addrs[i] = v
		i += 1
	}
	rc.AddHook(hook{addrs: addrs, database: opts.DB})
	return rc
}

// WrapClient is to wrap context and to add hooks for opentracing.
func WrapClient(ctx context.Context, client redis.UniversalClient) (UniversalClient, error) {
	if ctx == nil || client == nil {
		return nil, fmt.Errorf("[err] WrapClient invalid params")
	}
	var wrapClient UniversalClient
	switch client.(type) {
	case *redis.Client:
		wrapClient = &redisClient{Client: client.(*redis.Client)}
	case *redis.ClusterClient:
		wrapClient = &redisClusterClient{ClusterClient: client.(*redis.ClusterClient)}
	case *redis.Ring:
		wrapClient = &redisRing{Ring: client.(*redis.Ring)}
	default:
		return nil, fmt.Errorf("[err] WrapClient not support client")
	}

	wrapClient = wrapClient.withContext(ctx)
	return wrapClient, nil
}
