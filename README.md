# PageRank Algorithm Implementation in Go

A clean, efficient implementation of Google's PageRank algorithm written in Go. This implementation extracts the core ranking logic without external dependencies, making it perfect for understanding the algorithm or integrating into other projects.

##  What is PageRank?

PageRank is a link analysis algorithm developed by Google founders Larry Page and Sergey Brin at Stanford University. It measures the importance of web pages based on the structure of links between them.

### The Core Concept
Think of PageRank as modeling a "random surfer" who:
1. **Randomly jumps** to any page with probability `(1 - damping_factor)`
2. **Follows links** from the current page with probability `damping_factor`
3. The PageRank score represents the **probability** of finding the surfer on any given page

##  The Mathematics

The PageRank algorithm uses this iterative formula:

```
PR(page) = (1-d)/N + d * Σ(PR(backlink)/outlinks(backlink))
```

Where:
- `PR(page)` = PageRank score of the page
- `d` = damping factor (typically 0.85)
- `N` = total number of pages
- `(1-d)/N` = random jump probability
- `Σ(PR(backlink)/outlinks(backlink))` = sum of link contributions

### Why This Works
- **Higher backlinks from important pages** = higher PageRank
- **Damping factor prevents rank sinks** (pages with no outlinks)
- **Random jumps ensure connectivity** across the entire web graph

##  Features

- **Clean, readable Go code** with comprehensive comments
- **Configurable parameters** (damping factor, iterations)
- **Efficient implementation** using maps for O(1) lookups
- **No external dependencies** - uses only Go standard library
- **Production-ready** with proper error handling
- **Flexible data structures** supporting any URL/page identifiers

##  Code Structure

```
pagerank.go
├── PageRankResult          # Stores URL and rank score
├── PageRankCalculator      # Main algorithm implementation
│   ├── NewPageRankCalculator()    # Constructor with defaults
│   ├── SetDamping()               # Configure damping factor
│   ├── SetIterations()            # Set iteration count
│   ├── Calculate()                # Core PageRank computation
│   └── sortResults()              # Sort and format results
└── PrintResults()          # Pretty-print ranked results
```

## Usage

### Basic Usage

```go
package main

import "fmt"

func main() {
    // Define link structure
    backlinks := map[string][]string{
        "page-a": {"page-b", "page-c"},     // page-a gets links from b,c
        "page-b": {"page-c"},               // page-b gets link from c
        "page-c": {"page-a"},               // page-c gets link from a
    }
    
    outlinksCount := map[string]int{
        "page-a": 1,  // page-a links to 1 page
        "page-b": 1,  // page-b links to 1 page  
        "page-c": 2,  // page-c links to 2 pages
    }
    
    // Calculate PageRank
    calculator := NewPageRankCalculator()
    results := calculator.Calculate(backlinks, outlinksCount)
    
    // Display results
    PrintResults(results, 10)
}
```

### Advanced Configuration

```go
// Custom damping factor and iterations
calculator := NewPageRankCalculator().
    SetDamping(0.9).        // Higher damping = more link influence
    SetIterations(100)      // More iterations = better convergence

results := calculator.Calculate(backlinks, outlinksCount)
```

## Data Structures

### Input Format

**Backlinks Map**: `map[string][]string`
```go
backlinks := map[string][]string{
    "target-page": {"source1", "source2", "source3"},  // Pages linking TO target-page
    "other-page":  {"source1"},                        // Pages linking TO other-page
}
```

**Outlinks Count**: `map[string]int`
```go
outlinksCount := map[string]int{
    "source1": 5,  // source1 has 5 outgoing links
    "source2": 2,  // source2 has 2 outgoing links
}
```

### Output Format

**PageRankResult**: Sorted slice of results
```go
type PageRankResult struct {
    URL  string   // Page identifier
    Rank float64  // PageRank score (0.0 to 1.0+)
}
```

## Algorithm Details

### Initialization
Each page starts with equal probability:
```go
initial_rank = 1.0 / total_pages
```

### Iteration Process
For each iteration:
1. **Calculate random jump contribution**: `(1 - damping) / total_pages`
2. **Sum link contributions**: For each backlink, add `backlink_rank / backlink_outlinks`
3. **Apply damping**: `random_jump + damping * link_contributions`
4. **Update all pages simultaneously**

### Convergence
The algorithm runs for a fixed number of iterations (default: 50). In practice:
- **10 iterations**: Rough approximation
- **50 iterations**: Good convergence for most graphs
- **100+ iterations**: High precision for large, complex graphs

## Performance Characteristics

### Time Complexity
- **Per iteration**: O(E) where E = total number of edges/links
- **Total**: O(I × E) where I = iterations
- **Sorting**: O(N log N) where N = number of pages

### Space Complexity
- **O(N)** for PageRank scores
- **O(E)** for storing link structure
- **Total**: O(N + E)

### Practical Performance
- **Small graphs** (< 1000 pages): Milliseconds
- **Medium graphs** (< 100K pages): Seconds  
- **Large graphs** (1M+ pages): Minutes

## Real-World Applications

### Web Search Engines
- **Original use**: Google's web page ranking
- **Modern SEO**: Still influences search rankings
- **Link analysis**: Detecting spam and authority sites

### Social Networks
- **Influence ranking**: Most influential users/accounts  
- **Recommendation systems**: Suggest connections
- **Community detection**: Find important nodes

### Citation Analysis
- **Academic papers**: Most cited/influential research
- **Patent analysis**: Key innovations and dependencies
- **Knowledge graphs**: Important entities and relationships

### Other Applications
- **Transportation**: Important hubs and routes
- **Finance**: Systemic risk analysis
- **Biology**: Protein interaction networks
- **Marketing**: Viral spread and influence

## Configuration Guidelines

### Damping Factor (d)
- **0.85** (default): Standard value, good for most use cases
- **0.5-0.8**: More emphasis on random jumps, flattens rankings
- **0.9-0.95**: More emphasis on links, amplifies authority differences
- **Never 1.0**: Would create infinite loops in disconnected components

### Iterations
- **10-20**: Quick approximation, good for prototyping
- **50** (default): Balanced accuracy and performance
- **100+**: High precision for critical applications

### When to Use More Iterations
- **Large, complex graphs**: More connections = slower convergence
- **High precision requirements**: Financial or academic analysis
- **Debugging**: Understanding algorithm behavior

## Limitations and Considerations

### Known Issues
- **Rank sink problem**: Pages with no outlinks (solved by random jumps)
- **Spider traps**: Cycles that trap the random surfer (mitigated by damping)
- **Dangling nodes**: Pages not linked by others (handled gracefully)

### Performance Considerations
- **Memory usage**: Grows with graph size
- **I/O intensive**: Reading large link graphs
- **CPU intensive**: Matrix operations for large graphs

### Alternative Approaches
- **Personalized PageRank**: Biased random jumps to specific topics
- **Topic-Sensitive PageRank**: Multiple vectors for different topics  
- **TrustRank**: Anti-spam variant focusing on trusted seed pages

## Running the Code

### Prerequisites
- Go 1.16 or higher
- No external dependencies required

### Execution
```bash
# Run with default example
go run pagerank.go

# Build executable
go build -o pagerank pagerank.go
./pagerank

# Run tests (if implemented)
go test -v
```

### Expected Output
```
Total URLs processed: 4

Top 4 PageRank Results:
==================================================
page-d                                   | 0.37285156
page-c                                   | 0.28515625  
page-a                                   | 0.20703125
page-b                                   | 0.13496094

PageRank convergence demonstration:
After  1 iterations - Top page: page-d (0.325000)
After  5 iterations - Top page: page-d (0.372314)
After 10 iterations - Top page: page-d (0.372803)
After 50 iterations - Top page: page-d (0.372852)
```

## Contributing

Contributions are welcome! Areas for improvement:
- **Performance optimization** for large graphs
- **Parallel processing** for multi-core systems
- **Convergence detection** instead of fixed iterations
- **Graph visualization** tools
- **Benchmark suite** for performance testing

## References

- **Original Paper**: "The PageRank Citation Ranking: Bringing Order to the Web" by Page, Brin, Motwani, and Winograd
- **Google Patent**: US Patent 6,285,999
- **Stanford CS Course**: CS246 - Mining Massive Data Sets
- **Book**: "Introduction to Information Retrieval" by Manning, Raghavan, and Schütze

## TF-IDF Processor
Its job is to compute TF-IDF (Term Frequency-Inverse Document Frequency) after the Indexer has processed the crawled web pages. The TF-IDF Processor takes the indexed data from MongoDB, calculates the TF-IDF scores for each term in the documents, and stores the results back in MongoDB for fast retrieval by other services. This is essential for ranking search results based on the relevance of terms in the context of the entire collection of documents.

