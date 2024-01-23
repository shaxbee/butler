package rate

import "net/http"

type CacheBuilder struct {
	config       CacheConfig
	clientConfig ClientConfig
	synthetic    string

	httpClient *http.Client
	client     Client
}

func NewCacheBuilder(appID string) *CacheBuilder {
	config := CacheConfig{}
	clientConfig := ClientConfig{
		AppID: appID,
	}
	return &CacheBuilder{
		config:       config.Default(),
		clientConfig: clientConfig.Default(),
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

func (b *CacheBuilder) ClientConfig(config ClientConfig) *CacheBuilder {
	b.clientConfig = config.Default()
	return b
}

func (b *CacheBuilder) Symbols(symbols ...string) *CacheBuilder {
	b.config.Symbols = append(b.config.Symbols, symbols...)
	return b
}

func (b *CacheBuilder) Synthetic(base string) *CacheBuilder {
	b.synthetic = base
	return b
}

func (b *CacheBuilder) HTTPClient(httpClient *http.Client) *CacheBuilder {
	b.httpClient = httpClient
	return b
}

func (b *CacheBuilder) Build() *Cache {
	if b.client == nil {
		b.client = NewRestClient(b.httpClient, b.clientConfig)
	}

	if b.synthetic != "" && b.config.Base != b.synthetic {
		b.client = NewSyntheticClient(b.client, b.synthetic)
	}

	return NewCache(b.client, b.config)
}
