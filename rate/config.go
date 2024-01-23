package rate

import (
	"slices"
	"time"
)

var (
	DefaultEndpoint    = "https://openexchangerates.org"
	DefaultTimeout     = 10 * time.Second
	DefaultSyncTimeout = 30 * time.Second
	DefaultBase        = "USD"
	DefaultInterval    = 24 * time.Hour
)

type CacheConfig struct {
	Interval        time.Duration `koanf:"interval"`
	SyncTimeout     time.Duration `koanf:"sync_timeout"`
	Base            string        `koanf:"base"`
	Symbols         []string      `koanf:"symbols"`
	ShowAlternative bool          `koanf:"show_alternative"`
}

type ClientConfig struct {
	Endpoint string
	AppID    string
	Timeout  time.Duration
}

func (c CacheConfig) Default() CacheConfig {
	if c.SyncTimeout == 0 {
		c.SyncTimeout = DefaultSyncTimeout
	}

	if c.Base == "" {
		c.Base = DefaultBase
	}

	symbols := slices.Clone(c.Symbols)
	slices.Sort(symbols)
	slices.Compact(symbols)

	c.Symbols = symbols

	return c
}

func (c ClientConfig) Default() ClientConfig {
	if c.Endpoint == "" {
		c.Endpoint = DefaultEndpoint
	}

	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}

	return c
}
