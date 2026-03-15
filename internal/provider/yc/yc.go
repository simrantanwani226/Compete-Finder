// Package yc fetches startup data from the YC Companies API.
//
// HOW IT WORKS:
// 1. HTTP GET to the YC API endpoint (a static JSON file)
// 2. Decode JSON array into intermediate ycCompany structs
// 3. Map each ycCompany → provider.Startup (our domain model)
// 4. Return the list
//
// WHY the intermediate ycCompany struct?
// The YC API uses field names like "one_liner" and "long_description".
// Our domain model uses "Description". The ycCompany struct handles this
// translation. Nothing outside this package ever sees the YC-specific fields.
//
// This is the "adapter" part of ports-and-adapters. The YC API is the
// external system, ycCompany is the translation layer, and provider.Startup
// is what the rest of our app works with.
package yc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/simrantanwani226/compete-finder/internal/provider"
)

// DefaultURL is the YC Companies API endpoint.
// It's a static JSON file hosted on GitHub Pages — no auth, no rate limits.
const DefaultURL = "https://yc-oss.github.io/api/meta.json"

// ycCompany maps to the raw JSON shape from the YC API.
// These field names match what the API returns.
// This struct is PRIVATE (lowercase 'y') — nothing outside this package sees it.
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

// toStartup converts a YC API response into our domain model.
// KEY LOGIC: If one_liner is empty, fall back to long_description.
// This ensures we always have something for TF-IDF matching.
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

// YCProvider implements provider.Provider for the YC Companies API.
type YCProvider struct {
	url string
}

// New creates a YCProvider. Pass DefaultURL for production,
// or a test server URL for testing.
func New(url string) *YCProvider {
	return &YCProvider{url: url}
}

func (p *YCProvider) Name() string {
	return "yc"
}

// Fetch downloads all YC companies and converts them to Startups.
func (p *YCProvider) Fetch(ctx context.Context) ([]provider.Startup, error) {
	// NewRequestWithContext ties the request to the context.
	// If the context is cancelled (e.g. server shutting down), the HTTP
	// request is also cancelled. This prevents hanging requests.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching YC data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YC API returned status %d", resp.StatusCode)
	}

	// Decode directly from the response body (streaming, memory efficient)
	// rather than reading the entire body into memory first.
	var companies []ycCompany
	if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
		return nil, fmt.Errorf("decoding YC data: %w", err)
	}

	// Convert each YC company to our domain model
	startups := make([]provider.Startup, 0, len(companies))
	for _, c := range companies {
		startups = append(startups, c.toStartup())
	}

	return startups, nil
}
