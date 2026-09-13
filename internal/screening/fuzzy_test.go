package screening

import "testing"

func TestLevenshtein(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"", "", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"kitten", "sitting", 3}, // the textbook example
		{"same", "same", 0},
	}
	for _, tc := range cases {
		if got := levenshtein(tc.a, tc.b); got != tc.want {
			t.Errorf("levenshtein(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestSimilarity_ThresholdBehavior(t *testing.T) {
	if s := similarity("Acme Corp", "acme corp"); s != 1.0 {
		t.Errorf("expected case-insensitive exact match to score 1.0, got %v", s)
	}
	if s := similarity("Northwind Trading Consortium", "Northwind Trading Consortum"); s < matchThreshold {
		t.Errorf("expected a single dropped letter to still clear the match threshold, got %v", s)
	}
	if s := similarity("Acme Exports Ltd", "Example Trading Co"); s >= matchThreshold {
		t.Errorf("expected unrelated names to score below the match threshold, got %v", s)
	}
}
