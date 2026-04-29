package heatmap

import (
	"sort"

	"github.com/simrantanwani226/compete-finder/internal/provider"
)

type BatchTrend struct {
	Batch string
	Count int
	Trend string
}
type HeatmapResult struct {
	Batches      []BatchTrend
	MarketStatus string
	GrowthFactor float64
}

func groupByBatch(startups []provider.Startup) map[string][]provider.Startup {
	batch := make(map[string][]provider.Startup)
	for _, s := range startups {
		batch[s.Batch] = append(batch[s.Batch], s)
	}
	return batch
}
func sortBatches(batches []string) []string {
	sort.Slice(batches, func(i, j int) bool {
		yearI := batches[i][1:]
		yearJ := batches[j][1:]
		if yearI != yearJ {
			return yearI < yearJ
		}
		return batches[i][0] < batches[j][0]
	})
	return batches
}
func BuildHeatmap(startups []provider.Startup) HeatmapResult {
	if len(startups) == 0 {
		return HeatmapResult{MarketStatus: "stable", GrowthFactor: 1.0}
	}
	batches := make(map[string][]provider.Startup)
	batches = groupByBatch(startups)
	keys := make([]string, 0, len(batches))
	for k := range batches {
		keys = append(keys, k)
	}
	sorts := sortBatches(keys)
	trends := make([]BatchTrend, 0, len(sorts))
	for i, batch := range sorts {
		count := len(batches[batch])
		trend := "stable"
		if i > 0 {
			prev := trends[i-1].Count
			if count > prev {
				trend = "growing"
			} else if count < prev {
				trend = "shrinking"
			}
		}
		trends = append(trends, BatchTrend{Batch: batch, Count: count, Trend: trend})
	}
	firstCount := trends[0].Count
	lastCount := trends[len(trends)-1].Count
	growth := float64(lastCount) / float64(firstCount)
	var status string
	if growth >= 1.5 {
		status = "hot"
	} else if growth > 1.0 {
		status = "growing"
	} else if growth == 1.0 {
		status = "stable"
	} else {
		status = "declining"
	}

	return HeatmapResult{
		Batches:      trends,
		MarketStatus: status,
		GrowthFactor: growth}
}
