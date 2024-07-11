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
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis is the redis client.
type Redis struct {
	cluster     *redis.ClusterClient
	single      *redis.Client
	mutex       *sync.Mutex
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
		r.mutex = new(sync.Mutex)
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

// Client returns the redis client.
func (r *Redis) Client() redis.UniversalClient {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.clusterMode {
		return r.cluster
	}
	return r.single
}

// Start starts the underlying database.use in go-zero
func (s *Redis) Start() {}

// Stop closes the underlying database.use in go-zero.
func (s *Redis) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.clusterMode {
		s.cluster.Close()
		return
	}
	s.single.Close()
}
