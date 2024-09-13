/*
Copyright 2024 The corekit Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package redisdb

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v7"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis is the redis client.
type Redis struct {
	cluster     *redis.ClusterClient
	single      *redis.Client
	clusterMode bool
}

// NewRedis creates a new Redis instance.
func NewRedis(c Config) (*Redis, error) {
	r := &Redis{}

	if len(c.Address) == 1 {
		r.single = redis.NewClient(
			&redis.Options{
				Addr:         c.Address[0],
				Password:     c.Password,
				DB:           c.DB,
				DialTimeout:  3 * time.Second,
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
			})
		if err := r.single.Ping(context.Background()).Err(); err != nil {
			return nil, err
		}
		r.clusterMode = false
		return r, nil
	}

	r.cluster = redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:        c.Address,
			Password:     c.Password,
			DialTimeout:  3 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})
	if err := r.cluster.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	r.clusterMode = true
	return r, nil
}

func (r *Redis) Client() redis.UniversalClient {
	if r.clusterMode {
		return r.cluster
	}
	return r.single
}

func (r *Redis) RedSync() *redsync.Redsync {
	if r.clusterMode {
		return redsync.New(goredis.NewPool(r.cluster))
	}
	return redsync.New(goredis.NewPool(r.single))
}

func (r *Redis) Set(ctx context.Context, k, v string, t time.Duration) error {
	if r.clusterMode {
		return r.cluster.Set(ctx, k, v, t).Err()
	}
	return r.single.Set(ctx, k, v, t).Err()
}

func (r *Redis) Get(ctx context.Context, k string) interface{} {
	if r.clusterMode {
		return r.cluster.Get(ctx, k).Val()
	}
	return r.single.Get(ctx, k).Val()
}

func (r *Redis) HSet(ctx context.Context, k, field string, value interface{}) error {
	if r.clusterMode {
		return r.cluster.HSet(ctx, k, field, value).Err()
	}
	return r.single.HSet(ctx, k, field, value).Err()
}

func (r *Redis) HGet(ctx context.Context, k, field string) string {
	if r.clusterMode {
		return r.cluster.HGet(ctx, k, field).Val()
	}
	return r.single.HGet(ctx, k, field).Val()
}

func (r *Redis) HGetAll(ctx context.Context, k string) map[string]string {
	if r.clusterMode {
		return r.cluster.HGetAll(ctx, k).Val()
	}
	return r.single.HGetAll(ctx, k).Val()
}

func (r *Redis) HScan(ctx context.Context, k string, fn func(val string) error) error {
	if r.clusterMode {
		iter := r.cluster.HScan(ctx, k, 0, "", 50).Iterator()
		for iter.Next(ctx) {
			fn(iter.Val())
		}
		return iter.Err()
	}
	iter := r.single.HScan(ctx, k, 0, "", 50).Iterator()
	for iter.Next(ctx) {
		fn(iter.Val())
	}
	return iter.Err()
}

func (r *Redis) HDel(ctx context.Context, k, field string) error {
	if r.clusterMode {
		return r.cluster.HDel(ctx, k, field).Err()
	}
	return r.single.HDel(ctx, k, field).Err()
}

func (r *Redis) Expire(ctx context.Context, k string, t time.Duration) error {
	if r.clusterMode {
		return r.cluster.Expire(ctx, k, t).Err()
	}

	return r.single.Expire(ctx, k, t).Err()
}

func (r *Redis) HSetTTL(ctx context.Context, k, field string, value interface{}, t time.Duration) error {
	if r.clusterMode {
		if err := r.cluster.HSet(ctx, k, field, value).Err(); err != nil {
			return err
		}
		return r.cluster.Expire(ctx, k, t).Err()
	}
	if err := r.single.HSet(ctx, k, field, value).Err(); err != nil {
		return err
	}
	return r.single.Expire(ctx, k, t).Err()
}

func (r *Redis) Keys(ctx context.Context, k string) []string {
	if r.clusterMode {
		return r.cluster.Keys(ctx, k).Val()
	}
	return r.single.Keys(ctx, k).Val()
}

func (r *Redis) Del(ctx context.Context, k string) error {
	if r.clusterMode {
		return r.cluster.Del(ctx, k).Err()
	}
	return r.single.Del(ctx, k).Err()
}

func (r *Redis) Incr(ctx context.Context, k string) error {
	if r.clusterMode {
		return r.cluster.Incr(ctx, k).Err()
	}
	return r.single.Incr(ctx, k).Err()
}

func (r *Redis) BSet(ctx context.Context, k string, offset, value int64) error {
	if r.clusterMode {
		return r.cluster.SetBit(ctx, k, offset, int(value)).Err()
	}
	return r.single.SetBit(ctx, k, offset, int(value)).Err()
}

func (r *Redis) BGet(ctx context.Context, k string, offset int64) int64 {
	if r.clusterMode {
		return r.cluster.GetBit(ctx, k, offset).Val()
	}
	return r.single.GetBit(ctx, k, offset).Val()
}

func (r *Redis) BCount(ctx context.Context, k string, start, end int64) int64 {
	bc := &redis.BitCount{
		Start: start,
		End:   end,
	}
	if r.clusterMode {
		return r.cluster.BitCount(ctx, k, bc).Val()
	}
	return r.single.BitCount(ctx, k, bc).Val()
}

func (r *Redis) Add(ctx context.Context, b, k, v []byte) error {
	if r.clusterMode {
		return r.cluster.HSet(ctx, string(b), string(k), v).Err()
	}
	return r.single.HSet(ctx, string(b), string(k), v).Err()
}

func (r *Redis) Delete(ctx context.Context, b, k []byte) error {
	if r.clusterMode {
		return r.cluster.HDel(ctx, string(b), string(k)).Err()
	}
	return r.single.HDel(ctx, string(b), string(k)).Err()
}

func (r *Redis) All(ctx context.Context, k []byte) (interface{}, error) {
	if r.clusterMode {
		resp := r.cluster.HGetAll(ctx, string(k)).Val()
		if len(resp) == 0 {
			return nil, fmt.Errorf("%s key not found", string(k))
		}
		return resp, nil
	}
	resp := r.single.HGetAll(ctx, string(k)).Val()
	if len(resp) == 0 {
		return nil, fmt.Errorf("%s key not found", string(k))
	}
	return resp, nil
}

func (r *Redis) LPush(ctx context.Context, k, f string) error {
	if r.clusterMode {
		return r.cluster.LPush(ctx, k, f).Err()
	}
	return r.single.LPush(ctx, k, f).Err()
}

// Start starts the underlying database.use in go-zero
func (s *Redis) Start() {}

// Stop closes the underlying database.use in go-zero.
func (s *Redis) Stop() {
	if s.clusterMode {
		s.cluster.Close()
		return
	}
	s.single.Close()
}
