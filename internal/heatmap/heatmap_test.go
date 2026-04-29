package heatmap

import (
	"testing"

	"github.com/simrantanwani226/compete-finder/internal/provider"
)

func TestGroupByBatch(t *testing.T) {
	startups := []provider.Startup{
		{Name: "A", Batch: "W24"},
		{Name: "B", Batch: "S23"},
		{Name: "C", Batch: "W24"},
		{Name: "D", Batch: "S23"},
		{Name: "E", Batch: "W24"},
	}

	groups := groupByBatch(startups)

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if len(groups["W24"]) != 3 {
		t.Errorf("expected 3 in W24, got %d", len(groups["W24"]))
	}
	if len(groups["S23"]) != 2 {
		t.Errorf("expected 2 in S23, got %d", len(groups["S23"]))
	}
}

func TestSortBatches(t *testing.T) {
	batches := []string{"W24", "S23", "S24", "W23"}
	got := sortBatches(batches)

	want := []string{"S23", "W23", "S24", "W24"}
	for i, b := range want {
		if got[i] != b {
			t.Errorf("position %d: expected %q, got %q", i, b, got[i])
		}
	}
}

func TestBuildHeatmap(t *testing.T) {
	startups := []provider.Startup{
		{Name: "A", Batch: "S23"},
		{Name: "B", Batch: "S23"},
		{Name: "C", Batch: "W24"},
		{Name: "D", Batch: "W24"},
		{Name: "E", Batch: "W24"},
	}

	result := BuildHeatmap(startups)

	if len(result.Batches) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(result.Batches))
	}

	// First batch should be S23 (chronologically first)
	if result.Batches[0].Batch != "S23" {
		t.Errorf("expected S23 first, got %q", result.Batches[0].Batch)
	}
	if result.Batches[0].Count != 2 {
		t.Errorf("expected count 2 for S23, got %d", result.Batches[0].Count)
	}
	if result.Batches[0].Trend != "stable" {
		t.Errorf("expected stable for first batch, got %q", result.Batches[0].Trend)
	}

	// Second batch should be W24
	if result.Batches[1].Batch != "W24" {
		t.Errorf("expected W24 second, got %q", result.Batches[1].Batch)
	}
	if result.Batches[1].Count != 3 {
		t.Errorf("expected count 3 for W24, got %d", result.Batches[1].Count)
	}
	if result.Batches[1].Trend != "growing" {
		t.Errorf("expected growing for W24, got %q", result.Batches[1].Trend)
	}

	// Growth = 3/2 = 1.5, so market status = "growing" (not "hot" since not > 1.5)
	if result.GrowthFactor != 1.5 {
		t.Errorf("expected growth 1.5, got %f", result.GrowthFactor)
	}
	if result.MarketStatus != "hot" {
		t.Errorf("expected hot, got %q", result.MarketStatus)
	}
}

func TestBuildHeatmapEmpty(t *testing.T) {
	result := BuildHeatmap([]provider.Startup{})

	if len(result.Batches) != 0 {
		t.Errorf("expected 0 batches, got %d", len(result.Batches))
	}
	if result.MarketStatus != "stable" {
		t.Errorf("expected stable, got %q", result.MarketStatus)
	}
	if result.GrowthFactor != 1.0 {
		t.Errorf("expected growth 1.0, got %f", result.GrowthFactor)
	}
}
