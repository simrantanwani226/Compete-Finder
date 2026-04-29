package yc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/simrantanwani226/compete-finder/internal/provider"
)

type ycCompany struct {
	Name            string   `json:"name"`
	OneLiner        string   `json:"one_liner"`
	LongDescription string   `json:"long_description"`
	Industries      []string `json:"industries"`
	Batch           string   `json:"batch"`
	TeamSize        int      `json:"team_size"`
	Status          string   `json:"status"`
	Website         string   `json:"website"`
}

func (c ycCompany) toStartup() provider.Startup {

	desc := c.OneLiner
	if desc == "" {
		desc = c.LongDescription
	}

	return provider.Startup{
		Name:        c.Name,
		Description: desc,
		Industries:  c.Industries,
		Batch:       c.Batch,
		TeamSize:    c.TeamSize,
		Status:      c.Status,
		URL:         c.Website,
	}

}

type YCProvider struct {
	url string
}

func New(url string) *YCProvider {
	return &YCProvider{url: url}
}
func (p *YCProvider) Name() string {
	return "yc"
}
func (p *YCProvider) Fetch(ctx context.Context) ([]provider.Startup, error) {
	// Created Connection
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return nil, fmt.Errorf("Context %w", err)
	}
	// Connection Executed
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Context %w", err)
	}
	// Connection Closed
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YC api return status %d", resp.StatusCode)
	}
	var companies []ycCompany
	err = json.NewDecoder(resp.Body).Decode(&companies)
	if err != nil {
		return nil, fmt.Errorf("Json not serialized %w", err)
	}
	startups := make([]provider.Startup, 0, len(companies))
	for _, c := range companies {
		startups = append(startups, c.toStartup())
	}

	return startups, nil
}
