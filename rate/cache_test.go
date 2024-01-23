package rate

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func ExampleNewCache() {
	ctx := context.Background()

	cache := NewCacheBuilder(*appID).Build()

	// run the cache in the background
	go cache.Run(ctx)

	// wait for the cache to be ready
	{
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := cache.Sync(ctx); err != nil {
			panic(err)
		}
	}

	// convert using the latest rates
	latest := cache.Latest()
	amount, err := decimal.NewFromString("100.00")
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
		t.Run(base, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			t.Cleanup(cancel)

			rec := setupRecorder(t)
			cache := NewCacheBuilder(*appID).
				Symbols(symbols...).
				Synthetic(base).
				HTTPClient(rec.GetDefaultClient()).
				Build()

			// sync cache with the latest rates
			if err := cache.Sync(ctx); err != nil {
				t.Fatal(err)
			}

			if !cache.Ready() {
				t.Fatal("cache is not ready")
			}

			latest := cache.Latest()
			verifyLatest(t, base, symbols, latest)

			// fetch 2nd time to verify http caching
			prev := latest
			if err := cache.Fetch(ctx); err != nil {
				t.Fatal(err)
			}

			if !latest.Equal(prev) {
				t.Error("expected cache to return the same response")
			}
		})
	}
}

func TestSyntheticCache(t *testing.T) {
	ctx := context.Background()

	rec := setupRecorder(t)
	cache := NewCacheBuilder(*appID).
		HTTPClient(rec.GetDefaultClient()).
		Symbols(symbols...).
		Synthetic("CAD").
		Build()

	if err := cache.Sync(ctx); err != nil {
		t.Fatal("sync:", err)
	}

	latest := cache.Latest()

	verifyLatest(t, "CAD", symbols, latest)
}
