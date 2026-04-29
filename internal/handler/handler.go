package handler

import (
	"context"

	"connectrpc.com/connect"
	competev1 "github.com/simrantanwani226/compete-finder/gen/compete/v1"
	"github.com/simrantanwani226/compete-finder/internal/heatmap"
	"github.com/simrantanwani226/compete-finder/internal/matcher"
	"github.com/simrantanwani226/compete-finder/internal/provider"
)

type Handler struct {
	provider provider.Provider
}

func NewHandler(p provider.Provider) *Handler {

	return &Handler{provider: p}
}
func (h *Handler) FindCompetitors(ctx context.Context, req *connect.Request[competev1.FindCompetitorsRequest]) (*connect.Response[competev1.FindCompetitorsResponse], error) {

	startups, error := h.provider.Fetch(ctx)
	if error != nil {
		return nil, connect.NewError(connect.CodeInternal, error)
	}
	filtered := heatmap.FilterBySector(startups, req.Msg.Sector)
	result := matcher.Match(req.Msg.Description, filtered, int(req.Msg.Limit))
	competitors := make([]*competev1.Competitor, 0, len(result))
	for _, r := range result {
		comp := &competev1.Competitor{
			Name:        r.Startup.Name,
			Description: r.Startup.Description,
			Industries:  r.Startup.Industries,
			Batch:       r.Startup.Batch,
			TeamSize:    int32(r.Startup.TeamSize),
			Status:      r.Startup.Status,
			Url:         r.Startup.URL,
			MatchScore:  r.Score,
		}
		competitors = append(competitors, comp)

	}
	return connect.NewResponse(&competev1.FindCompetitorsResponse{
		Competitors:   competitors,
		TotalInSector: int32(len(filtered)),
	}), nil
}
func (h *Handler) GetMarketHeatmap(ctx context.Context, req *connect.Request[competev1.GetMarketHeatmapRequest]) (*connect.Response[competev1.GetMarketHeatmapResponse], error) {
	startups, error := h.provider.Fetch(ctx)
	if error != nil {
		return nil, connect.NewError(connect.CodeInternal, error)
	}
	filtered := heatmap.FilterBySector(startups, req.Msg.Sector)
	trend := heatmap.BuildHeatmap(filtered)
	batchtrend := make([]*competev1.BatchTrend, 0, len(trend.Batches))
	for _, bt := range trend.Batches {
		tren := &competev1.BatchTrend{
			Batch:        bt.Batch,
			StartupCount: int32(bt.Count),
			Trend:        bt.Trend,
		}
		batchtrend = append(batchtrend, tren)
	}
	return connect.NewResponse(&competev1.GetMarketHeatmapResponse{
		BatchTrends:  batchtrend,
		MarketStatus: trend.MarketStatus,
		GrowthFactor: trend.GrowthFactor,
	}), nil

}
