package provider

import (
	"context"
	"testing"
)

func TestPersonZeroValue(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	var p Person

	if p.Name != "" {
		t.Errorf("expected empty string, got %q", p.Name)
	}
	if p.Age != 0 {
		t.Errorf("expected 0, got %d", p.Age)
	}
}

func TestSliceBasics(t *testing.T) {
	tags := NewTags("fintech", "b2b")

	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "fintech" {
		t.Errorf("expected 'fintech', got %q", tags[0])
	}
	if tags[1] != "b2b" {
		t.Errorf("expected 'b2b', got %q", tags[1])
	}
}

func TestSliceNilZeroValue(t *testing.T) {
	var s []string

	if s != nil {
		t.Error("expected nil slice")
	}
	if len(s) != 0 {
		t.Errorf("expected len 0, got %d", len(s))
	}
}

func TestStartupStruct(t *testing.T) {
	s := Startup{
		Name:        "Acme Corp",
		Description: "We build rockets",
		Industries:  []string{"Aerospace", "B2B"},
		Batch:       "W24",
		TeamSize:    5,
		Status:      "Active",
		URL:         "https://acme.com",
	}

	if s.Name != "Acme Corp" {
		t.Errorf("expected 'Acme Corp', got %q", s.Name)
	}
	if s.Description != "We build rockets" {
		t.Errorf("expected 'We build rockets', got %q", s.Description)
	}
	if len(s.Industries) != 2 {
		t.Fatalf("expected 2 industries, got %d", len(s.Industries))
	}
	if s.Industries[0] != "Aerospace" {
		t.Errorf("expected 'Aerospace', got %q", s.Industries[0])
	}
	if s.Batch != "W24" {
		t.Errorf("expected 'W24', got %q", s.Batch)
	}
	if s.TeamSize != 5 {
		t.Errorf("expected 5, got %d", s.TeamSize)
	}
	if s.Status != "Active" {
		t.Errorf("expected 'Active', got %q", s.Status)
	}
	if s.URL != "https://acme.com" {
		t.Errorf("expected 'https://acme.com', got %q", s.URL)
	}
}

func TestStartupZeroValues(t *testing.T) {
	var s Startup

	if s.Name != "" {
		t.Errorf("expected empty Name, got %q", s.Name)
	}
	if s.TeamSize != 0 {
		t.Errorf("expected 0 TeamSize, got %d", s.TeamSize)
	}
	if s.Industries != nil {
		t.Error("expected nil Industries slice")
	}
}

// fakeProvider is a test type that should satisfy the Provider interface.
type fakeProvider struct {
	name     string
	startups []Startup
	err      error
}

func (f *fakeProvider) Name() string {
	return f.name
}

func (f *fakeProvider) Fetch(ctx context.Context) ([]Startup, error) {
	return f.startups, f.err
}

func TestProviderInterface(t *testing.T) {
	fake := &fakeProvider{
		name: "test",
		startups: []Startup{
			{Name: "TestCo", Description: "A test company"},
		},
	}

	// This line will only compile if Provider interface exists
	// and fakeProvider satisfies it.
	var p Provider = fake

	if p.Name() != "test" {
		t.Errorf("expected 'test', got %q", p.Name())
	}

	startups, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(startups) != 1 {
		t.Fatalf("expected 1 startup, got %d", len(startups))
	}
	if startups[0].Name != "TestCo" {
		t.Errorf("expected 'TestCo', got %q", startups[0].Name)
	}
}
