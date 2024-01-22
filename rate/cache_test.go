package rate

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bojanz/currency"
)

func ExampleNewCache() {
	ctx := context.Background()

	cache := NewCacheBuilder(*appID).Build()

	// run the cache in the background
	go cache.Run(ctx)

	// wait for the cache to be ready
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := cache.Sync(sctx); err != nil {
		panic(err)
	}

	// convert using the latest rates
	latest := cache.Latest()
	amount, err := currency.NewAmount("1234.56", "CAD")
	if err != nil {
		panic(err)
	}

	converted, err := latest.Convert(amount, "USD")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s = %s", amount, converted)
}

func TestCache(t *testing.T) {
	for _, base := range []string{"USD", "CAD"} {
		base := base
		t.Run(base, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			t.Cleanup(cancel)

			rec := setupRecorder(t)
			client := rec.GetDefaultClient()
			cache := NewCacheBuilder(*appID).HTTPClient(client).Synthetic(base).Build()

			go cache.Run(ctx)
			if err := cache.Sync(ctx); err != nil {
				t.Fatal("sync:", err)
			}

			if !cache.Ready() {
				t.Fatal("cache is not ready")
			}

			latest := cache.Latest()
			verifyLatest(t, base, symbols, latest)
		})
	}
}
