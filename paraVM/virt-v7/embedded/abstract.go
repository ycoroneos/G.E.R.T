// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedded

import (
	"time"
)

type Pollfunc func() interface{}

func Poll(f Pollfunc, period time.Duration, sink chan interface{}) chan bool {
	kill := make(chan bool)
	go func(kill chan bool) {
		for {
			select {
			case <-kill:
				return
			default:
				if period > 0 {
					time.Sleep(period)
				}
				sink <- f()
			}
		}
	}(kill)
	return kill
}
