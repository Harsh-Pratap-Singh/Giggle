package pagerank

import (
	"fmt"
	"sort"
)

type Result struct {
	URL  string
	Rank float64
}

type Calculator struct {
	damping    float64
	iterations int
}

func New() *Calculator {
	return &Calculator{
		damping:    0.85,
		iterations: 50,
	}
}

func (c *Calculator) SetDamping(d float64) *Calculator {
	if d > 0 && d < 1 {
		c.damping = d
	}
	return c
}

func (c *Calculator) SetIterations(n int) *Calculator {
	if n > 0 {
		c.iterations = n
	}
	return c
}

func (c *Calculator) Damping() float64 { return c.damping }
func (c *Calculator) Iterations() int  { return c.iterations }

func (c *Calculator) Calculate(backlinks map[string][]string, outlinksCount map[string]int) []Result {
	allURLs := make(map[string]bool)
	for url := range backlinks {
		allURLs[url] = true
	}
	for url := range outlinksCount {
		allURLs[url] = true
	}
	for _, sources := range backlinks {
		for _, src := range sources {
			allURLs[src] = true
		}
	}

	total := len(allURLs)
	if total == 0 {
		return []Result{}
	}

	rank := make(map[string]float64, total)
	for url := range allURLs {
		rank[url] = 1.0 / float64(total)
	}

	teleport := (1.0 - c.damping) / float64(total)

	for i := 0; i < c.iterations; i++ {
		next := make(map[string]float64, total)
		for url := range allURLs {
			next[url] = teleport
			sources, ok := backlinks[url]
			if !ok {
				continue
			}
			var contrib float64
			for _, src := range sources {
				out, hasOut := outlinksCount[src]
				if !hasOut || out == 0 {
					out = 1
				}
				if r, exists := rank[src]; exists {
					contrib += r / float64(out)
				}
			}
			next[url] += c.damping * contrib
		}
		rank = next
	}

	return sortResults(rank)
}

func sortResults(rank map[string]float64) []Result {
	out := make([]Result, 0, len(rank))
	for url, r := range rank {
		out = append(out, Result{URL: url, Rank: r})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Rank == out[j].Rank {
			return out[i].URL < out[j].URL
		}
		return out[i].Rank > out[j].Rank
	})
	return out
}

func Print(results []Result, limit int) {
	if limit > len(results) {
		limit = len(results)
	}
	fmt.Printf("Top %d PageRank Results:\n", limit)
	fmt.Println("==================================================")
	for i := 0; i < limit; i++ {
		fmt.Printf("%-40s | %.8f\n", results[i].URL, results[i].Rank)
	}
}
