package yc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestToStartup(t *testing.T) {
	c := ycCompany{
		Name:            "TestCo",
		OneLiner:        "A test company",
		LongDescription: "A longer description",
		Industries:      []string{"Fintech", "B2B"},
		Batch:           "W24",
		TeamSize:        10,
		Status:          "Active",
		Website:         "https://testco.com",
	}

	s := c.toStartup()

	if s.Name != "TestCo" {
		t.Errorf("expected 'TestCo', got %q", s.Name)
	}
	// one_liner should be used as Description when present
	if s.Description != "A test company" {
		t.Errorf("expected 'A test company', got %q", s.Description)
	}
	if len(s.Industries) != 2 {
		t.Errorf("expected 2 industries, got %d", len(s.Industries))
	}
	if s.Batch != "W24" {
		t.Errorf("expected 'W24', got %q", s.Batch)
	}
	if s.TeamSize != 10 {
		t.Errorf("expected 10, got %d", s.TeamSize)
	}
	if s.Status != "Active" {
		t.Errorf("expected 'Active', got %q", s.Status)
	}
	if s.URL != "https://testco.com" {
		t.Errorf("expected 'https://testco.com', got %q", s.URL)
	}
}

func TestToStartupFallback(t *testing.T) {
	c := ycCompany{
		Name:            "NullCo",
		OneLiner:        "",
		LongDescription: "Only has long desc",
	}

	s := c.toStartup()

	// When one_liner is empty, should fall back to long_description
	if s.Description != "Only has long desc" {
		t.Errorf("expected fallback to long_description, got %q", s.Description)
	}
}

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

func TestName(t *testing.T) {
	p := New("http://example.com")
	if p.Name() != "yc" {
		t.Errorf("expected 'yc', got %q", p.Name())
	}
}

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testJSON))
	}))
	defer srv.Close()

	p := New(srv.URL)
	startups, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(startups) != 2 {
		t.Fatalf("expected 2 startups, got %d", len(startups))
	}
	if startups[0].Name != "TestCo" {
		t.Errorf("expected 'TestCo', got %q", startups[0].Name)
	}
	if startups[0].Description != "A test company" {
		t.Errorf("expected 'A test company', got %q", startups[0].Description)
	}
	// Fallback test
	if startups[1].Description != "Only has long desc" {
		t.Errorf("expected fallback description, got %q", startups[1].Description)
	}
}

func TestFetchServerError(t *testing.T) {
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

func TestFetchBadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not valid json`))
	}))
	defer srv.Close()

	p := New(srv.URL)
	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error on bad JSON")
	}
}
