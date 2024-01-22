package rate

import (
	"context"
	"testing"
)

func TestSyntheticClient(t *testing.T) {
	ctx := context.Background()

	rec := setupRecorder(t)
	delegate := NewRestClient(rec.GetDefaultClient(), clientConfig().WithSymbols("CAD"))
	client := NewSyntheticClient(delegate, "CAD")

	rates, err := client.Latest(ctx)
	if err != nil {
		t.Fatal("client:", err)
	}

	verifyLatest(t, "CAD", symbols, rates)
}
