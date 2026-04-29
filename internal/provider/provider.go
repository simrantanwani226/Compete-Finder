// Package provider
package provider

import "context"

func NewTags(params ...string) []string {
	return params
}

type Startup struct {
	Name        string
	Description string
	Industries  []string
	Batch       string
	TeamSize    int
	Status      string
	URL         string
}
type Provider interface {
	Name() string
	Fetch(context.Context) ([]Startup, error)
}
