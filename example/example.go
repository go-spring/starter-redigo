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

package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-spring/spring-core/gs"
	"github.com/gomodule/redigo/redis"

	_ "github.com/go-spring/starter-redigo"
)

type Service struct {
	Redis *redis.Pool `autowire:""`
}

func main() {

	// Here `s` is not referenced by any other object,
	// so we need to register it as a root object.
	s := &Service{}
	gs.Root(gs.Object(s))

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		c := s.Redis.Get()
		defer func() { _ = c.Close() }()
		str, err := redis.String(c.Do("GET", "key"))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte(str))
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		c := s.Redis.Get()
		defer func() { _ = c.Close() }()
		str, err := redis.String(c.Do("SET", "key", "value"))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte(str))
	})

	gs.Run()

	// Example usage:
	//
	// ~ curl http://127.0.0.1:9090/get
	// redigo: nil returned%
	// ~ curl http://127.0.0.1:9090/set
	// OK%
	// ~ curl http://127.0.0.1:9090/get
	// value%
}

// ----------------------------------------------------------------------------
// Change working directory
// ----------------------------------------------------------------------------

// init sets the working directory of the application to the directory
// where this source file resides.
// This ensures that any relative file operations are based on the source file location,
// not the process launch path.
func init() {
	var execDir string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		execDir = filepath.Dir(filename)
	}
	err := os.Chdir(execDir)
	if err != nil {
		panic(err)
	}
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(workDir)
}
