package main

import (
	"fmt"
	"sort"
)

type PageRankResult struct {
	URL  string
	Rank float64
}

type PageRankCalculator struct {
	damping    float64
	iterations int
}

func NewPageRankCalculator() *PageRankCalculator {
	return &PageRankCalculator{
		damping:    0.85,
		iterations: 50,
	}
}

func (prc *PageRankCalculator) SetDamping(damping float64) *PageRankCalculator {
	if damping > 0 && damping < 1 {
		prc.damping = damping
	}
	return prc
}

func (prc *PageRankCalculator) SetIterations(iterations int) *PageRankCalculator {
	if iterations > 0 {
		prc.iterations = iterations
	}
	return prc
}

func (prc *PageRankCalculator) Calculate(backlinks map[string][]string, outlinksCount map[string]int) []PageRankResult {
	allUrls := make(map[string]bool)
	for url := range backlinks {
		allUrls[url] = true
	}
	for url := range outlinksCount {
		allUrls[url] = true
	}

	totalUrls := len(allUrls)
	if totalUrls == 0 {
		return []PageRankResult{}
	}

	pageRank := make(map[string]float64)
	for url := range allUrls {
		pageRank[url] = 1.0 / float64(totalUrls)
	}

	for i := 0; i < prc.iterations; i++ {
		newPageRank := make(map[string]float64)
		for url := range allUrls {
			newPageRank[url] = (1.0 - prc.damping) / float64(totalUrls)
			if backlinksForUrl, exists := backlinks[url]; exists {
				var linkContribution float64
				for _, backlink := range backlinksForUrl {
					outlinkCount, hasOutlinks := outlinksCount[backlink]
					if !hasOutlinks || outlinkCount == 0 {
						outlinkCount = 1
					}
					if backlinkRank, exists := pageRank[backlink]; exists {
						linkContribution += backlinkRank / float64(outlinkCount)
					}
				}
				newPageRank[url] += prc.damping * linkContribution
			}
		}
		pageRank = newPageRank
	}

	return prc.sortResults(pageRank)
}

func (prc *PageRankCalculator) sortResults(pageRank map[string]float64) []PageRankResult {
	results := make([]PageRankResult, 0, len(pageRank))
	for url, rank := range pageRank {
		results = append(results, PageRankResult{
			URL:  url,
			Rank: rank,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Rank > results[j].Rank
	})
	return results
}

func PrintResults(results []PageRankResult, limit int) {
	fmt.Printf("Top %d PageRank Results:\n", limit)
	fmt.Println("=" + fmt.Sprintf("%50s", "="))
	count := limit
	if len(results) < limit {
		count = len(results)
	}
	for i := 0; i < count; i++ {
		fmt.Printf("%-40s | %.8f\n", results[i].URL, results[i].Rank)
	}
}

func main() {
	backlinks := map[string][]string{
		"page-a": {"page-b", "page-c"},
		"page-b": {"page-c"},
		"page-c": {"page-a", "page-d"},
		"page-d": {"page-b", "page-c", "page-a"},
	}

	outlinksCount := map[string]int{
		"page-a": 2,
		"page-b": 1,
		"page-c": 2,
		"page-d": 1,
	}

	calculator := NewPageRankCalculator().
		SetDamping(0.85).
		SetIterations(50)

	results := calculator.Calculate(backlinks, outlinksCount)

	fmt.Printf("Total URLs processed: %d\n\n", len(results))
	PrintResults(results, 10)

	fmt.Println("\nPageRank convergence demonstration:")
	for _, iterations := range []int{1, 5, 10, 50} {
		calc := NewPageRankCalculator().SetIterations(iterations)
		results := calc.Calculate(backlinks, outlinksCount)
		fmt.Printf("After %2d iterations - Top page: %s (%.6f)\n",
			iterations, results[0].URL, results[0].Rank)
	}
}
