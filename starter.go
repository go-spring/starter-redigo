/*
 * Copyright 2025 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package StarterRedigo

import (
	"github.com/go-spring/spring-core/gs"
	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Addr     string `value:"${addr}"`
	Password string `value:"${password:=}"`
}

func init() {
	gs.Group("${spring.redigo}",
		func(c Config) (*redis.Pool, error) { // init
			return &redis.Pool{
				Dial: func() (redis.Conn, error) {
					return redis.Dial("tcp", c.Addr, redis.DialPassword(c.Password))
				},
			}, nil
		},
		func(pool *redis.Pool) error { // destroy
			return pool.Close()
		})
}
