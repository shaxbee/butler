package rate

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Cache is a thread-safe cache of the latest rates.
type Cache struct {
	client Client
	config CacheConfig

	mtx    sync.Mutex
	latest *LatestResponse
}

func NewCache(client Client, config CacheConfig) *Cache {
	if client == nil {
		panic("rate: client is required")
	}

	config = config.Default()
	return &Cache{
		client: client,
		config: config,
	}
}

// Run fetches the latest rates at the specified interval.
//
// If handler is not nil, it is called when [Cache.Fetch] fails.
func (c *Cache) Run(ctx context.Context) error {
	next := func() time.Duration {
		interval := c.config.Interval
		return time.Now().
			Truncate(interval).
			Add(interval).
			Sub(time.Now())
	}

	timer := time.NewTimer(next())
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if err := c.Fetch(ctx); err != nil {
				return err
			}

			// schedule the next fetch
			timer.Reset(next())
		}
	}
}

// Fetch fetches the latest rates.
func (c *Cache) Fetch(ctx context.Context) error {
	req := LatestRequest{
		Base:            c.config.Base,
		Symbols:         c.config.Symbols,
		ShowAlternative: c.config.ShowAlternative,
	}

	latest := c.Latest()
	if latest != nil {
		req.ETag = c.Latest().ETag
		req.Modified = c.Latest().Timestamp
	}

	latest, err := c.client.Latest(ctx, req)
	switch {
	case StatusCode(err) == http.StatusNotModified:
		return nil
	case err != nil:
		return fmt.Errorf("rate: fetch latest: %w", err)
	}

	c.mtx.Lock()
	c.latest = latest
	c.mtx.Unlock()

	return nil
}

// Sync ensures the lates rates are available.
func (c *Cache) Sync(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.config.SyncTimeout)
	defer cancel()

	if c.Latest() == nil {
		return c.Fetch(ctx)
	}

	return nil
}

// Latest returns the latest rates.
//
// If the cache is not ready, it returns nil.
func (c *Cache) Latest() *LatestResponse {
	c.mtx.Lock()
	latest := c.latest
	c.mtx.Unlock()

	return latest
}

// Ready indicates if the latest rates are available.
func (c *Cache) Ready() bool {
	return c.Latest() != nil
}
