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
	"fmt"

	"github.com/go-spring/spring-core/gs"
	"github.com/gomodule/redigo/redis"
)

// Config defines Redis connection configuration.
type Config struct {

	// Addr is the Redis server address.
	Addr string `value:"${addr}"`

	// Password is the Redis server password, default is empty.
	Password string `value:"${password:=}"`

	// Driver specifies which Redis driver to use, defaults to DefaultDriver.
	Driver string `value:"${driver:=DefaultDriver}"`
}

func init() {
	// Register a group of beans under the key "${spring.redigo}".
	// This group manages the lifecycle of Redis clients.
	gs.Group("${spring.redigo}",
		// create function creates a new Redis client
		func(c Config) (*redis.Pool, error) {
			d, ok := driverRegistry[c.Driver]
			if !ok {
				return nil, fmt.Errorf("redis driver not found: %s", c.Driver)
			}
			pool, err := d.CreateClient(c)
			if err != nil {
				return nil, fmt.Errorf("failed to create redis client: %w", err)
			}
			return pool, nil
		},
		// destroy function closes the Redis client
		func(pool *redis.Pool) error {
			return pool.Close()
		})
}

var driverRegistry = map[string]Driver{}

func init() {
	RegisterDriver("DefaultDriver", DefaultDriver{})
}

// Driver interface defines how to create a Redis client.
type Driver interface {
	CreateClient(c Config) (*redis.Pool, error)
}

// RegisterDriver registers a Redis driver with the given name.
// It panics if the driver name has already been registered.
func RegisterDriver(name string, driver Driver) {
	if _, ok := driverRegistry[name]; ok {
		panic("redis driver already registered: " + name)
	}
	driverRegistry[name] = driver
}

// DefaultDriver is the default implementation of the Driver interface.
type DefaultDriver struct{}

// CreateClient creates a new Redis client based on the provided configuration.
func (DefaultDriver) CreateClient(c Config) (*redis.Pool, error) {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", c.Addr, redis.DialPassword(c.Password))
		},
	}, nil
}
