package bitmap

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedis_CheckBits(t *testing.T) {
	type fields struct {
		key string
		m   uint64
	}
	type args struct {
		locs []uint64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		wantErr   bool
		doSetBits bool
	}{
		{
			name: "not exist",
			fields: fields{
				key: "not exist",
				m:   500,
			},
			args: args{
				locs: []uint64{
					10000,
					12345,
					45567,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "exist",
			fields: fields{
				key: "exist",
				m:   500,
			},
			args: args{
				locs: []uint64{
					10000,
					12345,
					45567,
				},
			},
			want:      true,
			wantErr:   false,
			doSetBits: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := miniredis.RunT(t)
			r := &Redis{
				client: redis.NewClient(&redis.Options{Addr: client.Addr()}),
				key:    tt.fields.key,
				m:      tt.fields.m,
			}
			if tt.doSetBits {
				r.SetBits(tt.args.locs)
			}
			got, err := r.CheckBits(tt.args.locs)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckBits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckBits() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRedisExpireTTL(t *testing.T) {
	m := miniredis.RunT(t)
	defer m.Close()

	key := "test-SetRedisExpireTTL"
	ttl := 10 * time.Second
	client := redis.NewClient(&redis.Options{Addr: m.Addr()})
	// do setbit prior to set expiry to avoid unsuccessful call with EXPIRE due to non-existing key
	client.SetBit(context.Background(), key, 0, 0)

	r := &Redis{
		client: client,
		key:    key,
		m:      10,
	}
	err := r.SetExpireTTL(ttl)
	assert.NoError(t, err)

	result := client.TTL(context.Background(), key)
	d, err := result.Result()
	assert.NoError(t, err)
	assert.Equal(t, ttl, d)
}

func TestNewRedis(t *testing.T) {
	m := miniredis.RunT(t)
	defer m.Close()

	key := "test-NewRedis"
	client := redis.NewClient(&redis.Options{Addr: m.Addr()})
	r := NewRedis(client, key, 10)
	// assert r.key is modified by NewRedis as `{key}_{NanoTime}`
	assert.Contains(t, r.key, fmt.Sprintf("%s_", key))

	result := client.Keys(context.Background(), fmt.Sprintf("%s_*", key))
	keys, err := result.Result()
	assert.NoError(t, err)
	// assert r.key was set to redis checked by Keys()
	assert.Contains(t, keys[0], fmt.Sprintf("%s_", key))
}
