package rate

import (
	"context"
	"sync"
	"time"
)

// Cache is a thread-safe cache of the latest rates.
type Cache struct {
	client   Client
	interval time.Duration

	mtx    sync.Mutex
	latest *Latest
	ready  chan struct{}
}

func NewCache(client Client, interval time.Duration) *Cache {
	if client == nil {
		panic("rate: client is required")
	}

	if interval == 0 {
		interval = DefaultInterval
	}

	return &Cache{
		client:   client,
		interval: interval,
		ready:    make(chan struct{}, 1),
	}
}

// Latest returns the latest rates.
//
// If the cache is not ready, it returns nil.
func (c *Cache) Latest() *Latest {
	c.mtx.Lock()
	latest := c.latest
	c.mtx.Unlock()

	return latest
}

// Ready indicates if the latest rates are available.
func (c *Cache) Ready() bool {
	c.mtx.Lock()
	ready := c.ready == nil
	c.mtx.Unlock()

	return ready
}

// Fetch fetches the latest rates.
func (c *Cache) Fetch(ctx context.Context) error {
	latest, err := c.client.Latest(ctx)
	if err != nil {
		return err
	}

	c.mtx.Lock()
	ready := c.ready
	c.latest = latest
	c.ready = nil
	c.mtx.Unlock()

	// signal that the cache is ready
	if ready != nil {
		close(ready)
	}

	return nil
}

// Run fetches the latest rates at the specified interval.
//
// If handler is not nil, it is called when [Cache.Fetch] fails.
func (c *Cache) Run(ctx context.Context) error {
	// fetch the latest rates immediately
	timer := time.NewTimer(1)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			err := c.Fetch(ctx)
			if err != nil {
				return err
			}

			// schedule the next fetch
			now := time.Now()
			next := now.Truncate(c.interval).Add(c.interval)
			timer.Reset(next.Sub(now))
		}
	}
}

// Sync blocks until the cache is ready.
func (c *Cache) Sync(ctx context.Context) error {
	c.mtx.Lock()
	ready := c.ready
	c.mtx.Unlock()

	if ready == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ready:
		return nil
	}
}
