package main

import (
	"flag"
	"fmt"

	"github.com/Harsh-Pratap-Singh/Search_Engine/pagerank"
)

func main() {
	damping := flag.Float64("damping", 0.85, "damping factor (0 < d < 1)")
	iterations := flag.Int("iterations", 50, "number of iterations")
	limit := flag.Int("top", 10, "number of top results to display")
	flag.Parse()

	backlinks := map[string][]string{
		"page-a": {"page-b", "page-c"},
		"page-b": {"page-c"},
		"page-c": {"page-a", "page-d"},
		"page-d": {"page-b", "page-c", "page-a"},
	}
	outlinksCount := map[string]int{
		"page-a": 2,
		"page-b": 2,
		"page-c": 3,
		"page-d": 1,
	}

	calc := pagerank.New().
		SetDamping(*damping).
		SetIterations(*iterations)

	results := calc.Calculate(backlinks, outlinksCount)

	fmt.Printf("Total URLs processed: %d\n\n", len(results))
	pagerank.Print(results, *limit)

	fmt.Println("\nPageRank convergence demonstration:")
	for _, n := range []int{1, 5, 10, 50} {
		r := pagerank.New().SetIterations(n).Calculate(backlinks, outlinksCount)
		fmt.Printf("After %2d iterations - Top page: %s (%.6f)\n", n, r[0].URL, r[0].Rank)
	}
}
