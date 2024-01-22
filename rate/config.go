package rate

import (
	"slices"
	"time"
)

var (
	DefaultEndpoint = "https://openexchangerates.org"
	DefaultTimeout  = 10 * time.Second
	DefaultBase     = "USD"
	DefaultInterval = 24 * time.Hour
)

type Config struct {
	Endpoint string
	AppID    string
	Base     string
	Symbols  []string
	Timeout  time.Duration
}

var DefaultConfig = Config{
	Endpoint: DefaultEndpoint,
	Timeout:  DefaultTimeout,
	Base:     DefaultBase,
}

func (c Config) WithAppID(appID string) Config {
	c.AppID = appID
	return c
}

func (c Config) WithBase(base string) Config {
	c.Base = base
	return c
}

func (c Config) WithSymbols(symbols ...string) Config {
	c.Symbols = append(c.Symbols, symbols...)
	return c
}

func (c Config) Normalize() Config {
	if c.Endpoint == "" {
		c.Endpoint = DefaultEndpoint
	}

	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
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
