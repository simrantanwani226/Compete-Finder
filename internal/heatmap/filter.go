package heatmap

import (
	"strings"

	"github.com/simrantanwani226/compete-finder/internal/provider"
)

func FilterBySector(startups []provider.Startup, sector string) []provider.Startup {
	if sector == "" {
		return startups
	}
	industry := make([]provider.Startup, 0, len(startups))
	for _, s := range startups {
		for _, i := range s.Industries {
			if strings.EqualFold(i, sector) {
				industry = append(industry, s)
				break
			}
		}
	}
	return industry
}
