package rate

import (
	"net/http"
	"time"
)

type CacheBuilder struct {
	httpClient *http.Client
	interval   time.Duration
	client     Client
	config     Config
	synthetic  string
}

func NewCacheBuilder(appID string) *CacheBuilder {
	config := DefaultConfig.WithAppID(appID)
	return &CacheBuilder{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		interval: DefaultInterval,
		config:   config,
	}
}

func (b *CacheBuilder) Config(config Config) *CacheBuilder {
	b.config = config.Normalize()
	return b
}

func (b *CacheBuilder) Synthetic(base string) *CacheBuilder {
	b.synthetic = base
	b.config.WithSymbols(base)
	return b
}

func (b *CacheBuilder) Interval(interval time.Duration) *CacheBuilder {
	b.interval = interval
	return b
}

func (b *CacheBuilder) HTTPClient(httpClient *http.Client) *CacheBuilder {
	b.httpClient = httpClient
	return b
}

func (b *CacheBuilder) Client(client Client) *CacheBuilder {
	b.client = client
	return b
}

func (b *CacheBuilder) Build() *Cache {
	if b.client == nil {
		b.client = NewRestClient(b.httpClient, b.config)
	}

	if b.synthetic != "" && b.config.Base != b.synthetic {
		b.client = NewSyntheticClient(b.client, b.synthetic)
	}

	return NewCache(b.client, b.interval)
}
