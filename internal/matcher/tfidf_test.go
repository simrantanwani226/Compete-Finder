package matcher

import (
	"math"
	"testing"

	"github.com/simrantanwani226/compete-finder/internal/provider"
)

func TestTF(t *testing.T) {
	tokens := []string{"fintech", "tools", "fintech", "payments"}
	got := tf(tokens)

	tests := []struct {
		word string
		want float64
	}{
		{"fintech", 0.5},
		{"tools", 0.25},
		{"payments", 0.25},
	}

	for _, tt := range tests {
		if math.Abs(got[tt.word]-tt.want) > 0.001 {
			t.Errorf("tf(%q) = %f, want %f", tt.word, got[tt.word], tt.want)
		}
	}
}

func TestIDF(t *testing.T) {
	docs := [][]string{
		{"fintech", "payments", "tools"},
		{"fintech", "banking", "api"},
		{"healthcare", "tools", "data"},
	}

	got := idf(docs)

	tests := []struct {
		word string
		want float64
	}{
		{"fintech", math.Log(3.0 / 2.0)},     // appears in 2 docs
		{"payments", math.Log(3.0 / 1.0)},    // appears in 1 doc
		{"tools", math.Log(3.0 / 2.0)},       // appears in 2 docs
		{"healthcare", math.Log(3.0 / 1.0)},  // appears in 1 doc
	}

	for _, tt := range tests {
		if math.Abs(got[tt.word]-tt.want) > 0.001 {
			t.Errorf("idf(%q) = %f, want %f", tt.word, got[tt.word], tt.want)
		}
	}
}

func TestTFIDFVec(t *testing.T) {
	tokens := []string{"fintech", "payments", "tools"}
	idfScores := map[string]float64{
		"fintech":  0.405,
		"payments": 1.099,
		"tools":    0.405,
	}

	got := tfidfVec(tokens, idfScores)

	tests := []struct {
		word string
		want float64
	}{
		{"fintech", (1.0 / 3.0) * 0.405},
		{"payments", (1.0 / 3.0) * 1.099},
		{"tools", (1.0 / 3.0) * 0.405},
	}

	for _, tt := range tests {
		if math.Abs(got[tt.word]-tt.want) > 0.001 {
			t.Errorf("tfidfVec(%q) = %f, want %f", tt.word, got[tt.word], tt.want)
		}
	}
}

func TestCosineSim(t *testing.T) {
	a := map[string]float64{"fintech": 0.5, "payments": 0.3}
	b := map[string]float64{"fintech": 0.4, "healthcare": 0.6}

	got := cosineSim(a, b)

	// dot = 0.5*0.4 = 0.2
	// magA = sqrt(0.25 + 0.09) = sqrt(0.34)
	// magB = sqrt(0.16 + 0.36) = sqrt(0.52)
	want := 0.2 / (math.Sqrt(0.34) * math.Sqrt(0.52))

	if math.Abs(got-want) > 0.001 {
		t.Errorf("cosineSim = %f, want %f", got, want)
	}
}

func TestCosineSimNoOverlap(t *testing.T) {
	a := map[string]float64{"fintech": 0.5}
	b := map[string]float64{"healthcare": 0.4}

	got := cosineSim(a, b)
	if got != 0 {
		t.Errorf("expected 0 for no overlap, got %f", got)
	}
}

func TestMatch(t *testing.T) {
	startups := []provider.Startup{
		{Name: "PayCo", Description: "fintech payments processing platform"},
		{Name: "HealthBot", Description: "healthcare data analytics tools"},
		{Name: "BankAPI", Description: "fintech banking api infrastructure"},
	}

	results := Match("fintech payments", startups, 2)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// PayCo should rank first (shares "fintech" and "payments")
	if results[0].Startup.Name != "PayCo" {
		t.Errorf("expected PayCo first, got %q", results[0].Startup.Name)
	}

	// BankAPI should rank second (shares "fintech")
	if results[1].Startup.Name != "BankAPI" {
		t.Errorf("expected BankAPI second, got %q", results[1].Startup.Name)
	}

	// Scores should be between 0 and 1
	for _, r := range results {
		if r.Score < 0 || r.Score > 1 {
			t.Errorf("score out of range: %f", r.Score)
		}
	}
}

func TestMatchLimit(t *testing.T) {
	startups := []provider.Startup{
		{Name: "A", Description: "fintech payments"},
		{Name: "B", Description: "fintech banking"},
		{Name: "C", Description: "fintech tools"},
	}

	results := Match("fintech", startups, 1)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}
