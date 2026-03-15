// Package provider defines the core domain model and data source interface.
//
// WHY THIS EXISTS:
// Every data source (YC, GitHub, HN) returns data in a different format.
// We don't want the rest of our app to care about that.
// So we define ONE shape (Startup) and ONE interface (Provider).
// Each data source maps its data into Startup.
// The matcher, store, heatmap — they only see Startup. They never know
// where the data came from.
//
// This is called the "ports and adapters" pattern. The Provider interface
// is a "port" — it defines WHAT we need. The YC provider is an "adapter" —
// it implements HOW to get it from a specific source.
package provider

import "context"

// Startup is the core domain model. Every field here matters for matching
// or display. We keep it flat and simple — no nested structs, no pointers.
type Startup struct {
	Name        string   // Company name
	Description string   // What they do (used for TF-IDF matching)
	Industries  []string // e.g. ["Fintech", "B2B"] — used for sector filtering
	Batch       string   // YC batch, e.g. "W24" — used for heatmap trends
	TeamSize    int      // 0 if unknown
	Status      string   // "Active", "Dead", "Acquired"
	URL         string   // Company website
}

// Provider is the interface for any data source that supplies startups.
//
// WHY an interface and not just a function?
// 1. Testing — we can create a fake provider with hardcoded data
// 2. Extensibility — add GitHub/HN providers later without changing existing code
// 3. Name() tells us which source the data came from (useful for logging)
type Provider interface {
	Name() string
	Fetch(ctx context.Context) ([]Startup, error)
}
