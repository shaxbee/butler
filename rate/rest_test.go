package rate

import (
	"context"
	"flag"
	"path"
	"regexp"
	"strings"
	"testing"

	_ "embed"

	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var mockAppID = "deadbeef"
var appID = flag.String("app-id", mockAppID, "Open Exchange Rates App ID")
var update = flag.Bool("update", false, "update VCR cassettes")
var symbols = []string{"THB", "KRW", "JPY", "EUR"}

func TestClient(t *testing.T) {
	rec := setupRecorder(t)
	client := NewRestClient(rec.GetDefaultClient(), clientConfig())

	ctx := context.Background()
	latest, err := client.Latest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	verifyLatest(t, client.Base(), symbols, latest)

	usd := amount("1234", "USD")
	if actual := mustConvert(t, latest, usd, "USD"); !actual.Equal(usd) {
		t.Errorf("expected %s, got %s", usd, actual)
	}

	converted := make([]string, len(symbols))
	for i, symbol := range symbols {
		amount := mustConvert(t, latest, usd, symbol)
		if amount.Equal(usd) {
			t.Errorf("expected %s to be different from USD", symbol)
		}

		converted[i] = amount.String()
	}

	if actual := mustConvert(t, latest, usd, "THB"); actual.Equal(usd) {
		t.Errorf("expected THB to be different from USD")
	}

	ts := latest.Timestamp.Format("2006-01-02 15:04")
	t.Logf("%s %s: %s", ts, usd.String(), strings.Join(converted, ", "))
}

func setupRecorder(t testing.TB) *recorder.Recorder {
	t.Helper()

	options := &recorder.Options{
		CassetteName:       path.Join("testdata", t.Name(), "vcr"),
		Mode:               recorder.ModeReplayOnly,
		SkipRequestLatency: true,
	}
	if *update {
		options.Mode = recorder.ModeRecordOnly
	}

	rec, err := recorder.NewWithOptions(options)
	if err != nil {
		t.Fatal("recorder:", err)
	}
	t.Cleanup(func() {
		if err := rec.Stop(); err != nil {
			t.Error("recorder:", err)
		}
	})

	// replace real app ID
	re := regexp.MustCompile(`app_id=[a-z0-9a-f]+`)
	hook := func(i *cassette.Interaction) error {
		i.Request.URL = re.ReplaceAllString(i.Request.URL, "app_id="+mockAppID)
		return nil
	}
	rec.AddHook(hook, recorder.AfterCaptureHook)

	return rec
}

func clientConfig() Config {
	return Config{
		AppID:   *appID,
		Symbols: symbols,
	}
}
