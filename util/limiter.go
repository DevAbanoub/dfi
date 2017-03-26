// This is free and unencumbered software released into the public domain.
// 
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
// 
// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
// 
// For more information, please refer to <http://unlicense.org/>
package util

import "time"

type Limiter struct {
	Throttle chan time.Time
	Ticker   *time.Ticker
	quit     chan bool
}

// Return a new rate limiter. This is used to make sure that something like a
// requrest for instance does not run too many times. However, it does allow
// bursting. For example, it may refill at a rate of 3 tokens per minute, and
// have a burst of three. This means that if it has been running for more than
// a minute without being used then, it will be able to be used 3 times in
// rapid succession - no limiting will apply.
func NewLimiter(rate time.Duration, burst int, fill bool) *Limiter {
	tick := time.NewTicker(rate)
	throttle := make(chan time.Time, burst)
	quit := make(chan bool)

	if fill {
		for i := 0; i < burst; i++ {
			throttle <- time.Now()
		}
	}

	go func() {
		for t := range tick.C {
			select {
			case _ = <-quit:
				return
			case throttle <- t:
			default:
			}
		}
	}()

	return &Limiter{throttle, tick, quit}
}

// Block until the given time has elapsed. Or just use a token from the bucket.
func (l *Limiter) Wait() {
	_, _ = <-l.Throttle
}

// Finish running.
func (l *Limiter) Stop() {
	l.Ticker.Stop()
	l.quit <- true
	close(l.Throttle)
}

// Limits requests from peers
type PeerLimiter struct {
	queryLimiter    *Limiter
	announceLimiter *Limiter
}

func (pl *PeerLimiter) Setup() {
	// Allow an announce every 10 minutes, bursting to allow three.
	// The burst is there as people may make "mistakes" with titles or descriptions
	pl.announceLimiter = NewLimiter(time.Minute*10, 3, true)

	pl.queryLimiter = NewLimiter(time.Second/3, 3, true)
}
