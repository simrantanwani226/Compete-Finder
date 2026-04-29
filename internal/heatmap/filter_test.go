package heatmap

import (
	"testing"

	"github.com/simrantanwani226/compete-finder/internal/provider"
)

func TestFilterBySector(t *testing.T) {
	startups := []provider.Startup{
		{Name: "PayCo", Industries: []string{"Fintech", "B2B"}},
		{Name: "HealthBot", Industries: []string{"Healthcare", "AI"}},
		{Name: "BankAPI", Industries: []string{"Fintech", "API"}},
	}

	got := FilterBySector(startups, "fintech")

	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0].Name != "PayCo" {
		t.Errorf("expected PayCo, got %q", got[0].Name)
	}
	if got[1].Name != "BankAPI" {
		t.Errorf("expected BankAPI, got %q", got[1].Name)
	}
}

func TestFilterBySectorCaseInsensitive(t *testing.T) {
	startups := []provider.Startup{
		{Name: "PayCo", Industries: []string{"Fintech"}},
	}

	got := FilterBySector(startups, "FINTECH")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestFilterBySectorEmpty(t *testing.T) {
	startups := []provider.Startup{
		{Name: "PayCo", Industries: []string{"Fintech"}},
		{Name: "HealthBot", Industries: []string{"Healthcare"}},
	}

	got := FilterBySector(startups, "")
	if len(got) != 2 {
		t.Fatalf("expected all 2 results when sector is empty, got %d", len(got))
	}
}

func TestFilterBySectorNoMatch(t *testing.T) {
	startups := []provider.Startup{
		{Name: "PayCo", Industries: []string{"Fintech"}},
	}

	got := FilterBySector(startups, "healthcare")
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}
