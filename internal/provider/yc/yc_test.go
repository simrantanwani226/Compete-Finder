package yc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// WHY a fake HTTP server?
// We don't want tests hitting the real YC API because:
// 1. Tests would be slow (network call)
// 2. Tests would be flaky (API could be down)
// 3. Tests wouldn't be deterministic (data changes)
//
// httptest.NewServer creates a local HTTP server that returns our test data.
// The provider doesn't know it's fake — it just sees a URL.

const testJSON = `[
	{
		"name": "TestCo",
		"one_liner": "A test company",
		"long_description": "A longer description",
		"industries": ["Fintech", "B2B"],
		"batch": "W24",
		"team_size": 10,
		"status": "Active",
		"website": "https://testco.com"
	},
	{
		"name": "NullCo",
		"one_liner": "",
		"long_description": "Only has long desc",
		"industries": [],
		"batch": "S23",
		"team_size": 0,
		"status": "Dead",
		"website": ""
	}
]`

func TestFetch(t *testing.T) {
	// Create fake HTTP server that returns our test JSON
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testJSON))
	}))
	defer srv.Close() // Clean up after test

	// Create provider pointing at our fake server (not the real YC API)
	p := New(srv.URL)
	startups, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(startups) != 2 {
		t.Fatalf("expected 2 startups, got %d", len(startups))
	}

	// Verify field mapping works correctly
	if startups[0].Name != "TestCo" {
		t.Errorf("expected name TestCo, got %s", startups[0].Name)
	}
	if startups[0].Description != "A test company" {
		t.Errorf("expected one_liner as description, got %s", startups[0].Description)
	}
	if len(startups[0].Industries) != 2 {
		t.Errorf("expected 2 industries, got %d", len(startups[0].Industries))
	}

	// KEY TEST: When one_liner is empty, fall back to long_description
	if startups[1].Description != "Only has long desc" {
		t.Errorf("expected long_description fallback, got %s", startups[1].Description)
	}
}

func TestFetchServerError(t *testing.T) {
	// Server that always returns 500
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	p := New(srv.URL)
	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}
