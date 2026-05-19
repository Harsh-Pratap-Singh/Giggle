package pagerank

import (
	"math"
	"testing"
)

func sampleGraph() (map[string][]string, map[string]int) {
	backlinks := map[string][]string{
		"page-a": {"page-b", "page-c"},
		"page-b": {"page-c"},
		"page-c": {"page-a", "page-d"},
		"page-d": {"page-b", "page-c", "page-a"},
	}
	outlinks := map[string]int{
		"page-a": 2,
		"page-b": 2,
		"page-c": 3,
		"page-d": 1,
	}
	return backlinks, outlinks
}

func TestCalculateProducesAllPages(t *testing.T) {
	backlinks, outlinks := sampleGraph()
	results := New().Calculate(backlinks, outlinks)

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	seen := map[string]bool{}
	for _, r := range results {
		seen[r.URL] = true
	}
	for _, want := range []string{"page-a", "page-b", "page-c", "page-d"} {
		if !seen[want] {
			t.Errorf("missing url %q in results", want)
		}
	}
}

func TestRanksSumNearOne(t *testing.T) {
	backlinks, outlinks := sampleGraph()
	results := New().SetIterations(100).Calculate(backlinks, outlinks)

	var sum float64
	for _, r := range results {
		sum += r.Rank
	}
	if math.Abs(sum-1.0) > 0.01 {
		t.Errorf("ranks should sum to ~1.0, got %.6f", sum)
	}
}

func TestResultsAreSortedDescending(t *testing.T) {
	backlinks, outlinks := sampleGraph()
	results := New().Calculate(backlinks, outlinks)

	for i := 1; i < len(results); i++ {
		if results[i-1].Rank < results[i].Rank {
			t.Fatalf("results not sorted at %d: %v < %v", i, results[i-1], results[i])
		}
	}
}

func TestEmptyGraph(t *testing.T) {
	results := New().Calculate(nil, nil)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestDanglingNodeHandled(t *testing.T) {
	backlinks := map[string][]string{
		"page-a": {"page-b"},
	}
	outlinks := map[string]int{}
	results := New().Calculate(backlinks, outlinks)
	if len(results) == 0 {
		t.Fatal("dangling node graph should still produce results")
	}
	for _, r := range results {
		if math.IsNaN(r.Rank) || math.IsInf(r.Rank, 0) {
			t.Errorf("invalid rank for %s: %v", r.URL, r.Rank)
		}
	}
}

func TestSetDampingRejectsInvalid(t *testing.T) {
	c := New().SetDamping(1.5).SetDamping(-0.1).SetDamping(0.9)
	if c.Damping() != 0.9 {
		t.Errorf("expected damping 0.9, got %v", c.Damping())
	}
}

func TestSetIterationsRejectsInvalid(t *testing.T) {
	c := New().SetIterations(-5).SetIterations(0).SetIterations(25)
	if c.Iterations() != 25 {
		t.Errorf("expected iterations 25, got %v", c.Iterations())
	}
}

func TestConvergenceStabilizes(t *testing.T) {
	backlinks, outlinks := sampleGraph()
	r50 := New().SetIterations(50).Calculate(backlinks, outlinks)
	r100 := New().SetIterations(100).Calculate(backlinks, outlinks)

	rankByURL := func(rs []Result) map[string]float64 {
		m := map[string]float64{}
		for _, r := range rs {
			m[r.URL] = r.Rank
		}
		return m
	}
	a, b := rankByURL(r50), rankByURL(r100)
	for url, v := range a {
		if math.Abs(v-b[url]) > 1e-4 {
			t.Errorf("rank for %s did not converge: 50=%v 100=%v", url, v, b[url])
		}
	}
}
