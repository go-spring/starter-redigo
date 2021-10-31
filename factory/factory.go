/*
 * Copyright 2012-2019 the original author or authors.
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

package factory

import (
	"fmt"

	"github.com/go-spring/spring-base/fastdev"
	"github.com/go-spring/spring-core/redis"
	"github.com/go-spring/spring-redigo"
	"github.com/go-spring/starter-core"
	g "github.com/gomodule/redigo/redis"
)

// NewClient 创建 Redis 客户端
func NewClient(config StarterCore.RedisConfig) (redis.Client, error) {

	if fastdev.ReplayMode() {
		return SpringRedigo.NewClient(nil), nil
	}

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := g.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	if _, err = conn.Do("PING"); err != nil {
		return nil, err
	}
	return SpringRedigo.NewClient(conn), nil
}
