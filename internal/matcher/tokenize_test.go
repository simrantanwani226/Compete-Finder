package matcher

import (
	"testing"
)

func TestTokenizeBasic(t *testing.T) {
	got := Tokenize("We build fintech tools for small businesses")
	want := []string{"build", "fintech", "tools", "small", "businesses"}

	if len(got) != len(want) {
		t.Fatalf("expected %d tokens, got %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("token %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestTokenizePunctuation(t *testing.T) {
	got := Tokenize("Hello, world! This is a test.")
	want := []string{"hello", "world", "test"}

	if len(got) != len(want) {
		t.Fatalf("expected %d tokens, got %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("token %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestTokenizeEmpty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 tokens, got %d: %v", len(got), got)
	}
}
