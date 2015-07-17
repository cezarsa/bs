// Copyright 2015 bs authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric

import (
	"os"
	"time"

	"github.com/tsuru/bs/container"
)

type runner struct {
	dockerEndpoint string
	interval       time.Duration
	finish         chan bool
}

func NewRunner(dockerEndpoint string, interval time.Duration) *runner {
	return &runner{
		finish:         make(chan bool),
		dockerEndpoint: dockerEndpoint,
		interval:       interval,
	}
}

func (r *runner) Start() error {
	client, err := container.NewClient(r.dockerEndpoint)
	if err != nil {
		return err
	}
	reporter := &Reporter{
		backend:    os.Getenv("METRICS_BACKEND"),
		infoClient: client,
	}
	go func() {
		for {
			reporter.Do()
			select {
			case <-r.finish:
				return
			case <-time.After(r.interval):
			}

		}
	}()
	return nil
}

func (r *runner) Stop() {
	r.finish <- true
}
