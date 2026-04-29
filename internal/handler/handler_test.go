package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	competev1 "github.com/simrantanwani226/compete-finder/gen/compete/v1"
	"github.com/simrantanwani226/compete-finder/gen/compete/v1/competev1connect"
	"github.com/simrantanwani226/compete-finder/internal/provider"
)

// fakeProvider returns hardcoded data for testing.
type fakeProvider struct{}

func (f *fakeProvider) Name() string { return "fake" }

func (f *fakeProvider) Fetch(ctx context.Context) ([]provider.Startup, error) {
	return []provider.Startup{
		{
			Name:        "PayCo",
			Description: "fintech payments processing platform",
			Industries:  []string{"Fintech", "B2B"},
			Batch:       "W24",
			TeamSize:    10,
			Status:      "Active",
			URL:         "https://payco.com",
		},
		{
			Name:        "HealthBot",
			Description: "healthcare data analytics tools",
			Industries:  []string{"Healthcare", "AI"},
			Batch:       "S23",
			TeamSize:    5,
			Status:      "Active",
			URL:         "https://healthbot.com",
		},
		{
			Name:        "BankAPI",
			Description: "fintech banking api infrastructure",
			Industries:  []string{"Fintech", "API"},
			Batch:       "W24",
			TeamSize:    8,
			Status:      "Active",
			URL:         "https://bankapi.com",
		},
	}, nil
}

func setupTestServer(t *testing.T) competev1connect.CompeteServiceClient {
	t.Helper()
	h := NewHandler(&fakeProvider{})
	path, handler := competev1connect.NewCompeteServiceHandler(h)
	mux := http.NewServeMux()
	mux.Handle(path, handler)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	client := competev1connect.NewCompeteServiceClient(
		http.DefaultClient,
		srv.URL,
	)
	return client
}

func TestFindCompetitors(t *testing.T) {
	client := setupTestServer(t)

	resp, err := client.FindCompetitors(context.Background(), connect.NewRequest(&competev1.FindCompetitorsRequest{
		Name:        "MyFintech",
		Description: "fintech payments",
		Sector:      "fintech",
		Limit:       2,
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Msg.Competitors) != 2 {
		t.Fatalf("expected 2 competitors, got %d", len(resp.Msg.Competitors))
	}

	// PayCo should rank first (shares "fintech" and "payments")
	if resp.Msg.Competitors[0].Name != "PayCo" {
		t.Errorf("expected PayCo first, got %q", resp.Msg.Competitors[0].Name)
	}
}

func TestGetMarketHeatmap(t *testing.T) {
	client := setupTestServer(t)

	resp, err := client.GetMarketHeatmap(context.Background(), connect.NewRequest(&competev1.GetMarketHeatmapRequest{
		Sector: "fintech",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Msg.BatchTrends) == 0 {
		t.Fatal("expected at least one batch trend")
	}

	if resp.Msg.MarketStatus == "" {
		t.Error("expected a market status")
	}

	if resp.Msg.GrowthFactor == 0 {
		t.Error("expected a non-zero growth factor")
	}
}
