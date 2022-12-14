# go-bloomfilter

![build workflow](https://github.com/x0rworld/go-bloomfilter/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/x0rworld/go-bloomfilter.svg)](https://pkg.go.dev/github.com/x0rworld/go-bloomfilter)

[go-bloomfilter] is implemented by Golang which supports in-memory and Redis. Moreover, it’s available for a
duration-based rotation.

## Resources

- [Documentation]
- [Examples]

## Features

- [Supporting bitmap]: in-memory and Redis
- [Rotation]

## Installation

```shell
go get github.com/x0rworld/go-bloomfilter
```

## Quickstart

```go
package main

import (
	"context"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/factory"
	"log"
)

func main() {
	m, k := bloom.EstimateParameters(100, 0.01)
	// configure factory config
	cfg := config.NewDefaultFactoryConfig()
	// modify config
	cfg.FilterConfig.M = uint64(m)
	cfg.FilterConfig.K = uint64(k)
	// create factory by config
	ff, err := factory.NewFilterFactory(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	// create filter by factory
	f, err := ff.NewFilter(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	// manipulate filter: Exist & Add
	data := "hello world"
	exist, err := f.Exist(data)
	if err != nil {
		log.Println(err)
		return
	}
	// data: hello world, exist: false
	log.Printf("data: %v, exist: %v\n", data, exist)
	err = f.Add(data)
	if err != nil {
		log.Println(err)
		return
	}
	// add data: hello world
	log.Printf("add data: %s\n", data)
	exist, err = f.Exist(data)
	if err != nil {
		log.Println(err)
		return
	}
	// data: hello world, exist: true
	log.Printf("data: %v, exist: %v\n", data, exist)
}
```

More examples such as rotation could be found in [Examples].

[go-bloomfilter]: https://github.com/x0rworld/go-bloomfilter

[Examples]: ./example

[Documentation]: https://pkg.go.dev/github.com/x0rworld/go-bloomfilter

[Supporting bitmap]: ./bitmap

[Rotation]: ./filter/rotator
