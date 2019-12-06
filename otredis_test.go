package otredis

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

func TestWrapClient(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		ctx    context.Context
		client redis.UniversalClient
		isErr  bool
	}{
		"fail": {
			ctx:    nil,
			client: nil,
			isErr:  true,
		},
		"client": {
			ctx:    context.TODO(),
			client: redis.NewClient(&redis.Options{}),
			isErr:  false,
		},
		"cluster": {
			ctx:    context.TODO(),
			client: redis.NewClusterClient(&redis.ClusterOptions{}),
			isErr:  false,
		},
		"ring": {
			ctx:    context.TODO(),
			client: redis.NewRing(&redis.RingOptions{}),
			isErr:  false,
		},
	}

	for k, tc := range tests {
		t.Run(k, func(t *testing.T) {
			_, err := WrapClient(tc.ctx, tc.client)
			assert.Equal(tc.isErr, err != nil)
		})
	}
}

func TestWithContext(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		ctx    context.Context
		client redis.UniversalClient
	}{
		"client": {
			ctx:    context.TODO(),
			client: redis.NewClient(&redis.Options{}),
		},
		"cluster": {
			ctx:    context.TODO(),
			client: redis.NewClusterClient(&redis.ClusterOptions{})},
		"ring": {
			ctx:    context.TODO(),
			client: redis.NewRing(&redis.RingOptions{}),
		},
	}

	for k, tc := range tests {
		t.Run(k, func(t *testing.T) {
			wrap, err := WrapClient(tc.ctx, tc.client)
			assert.NoError(err)
			equal := wrap.withContext(tc.ctx)
			assert.Equal(wrap, equal)
		})
	}
}
