# Giggle (Search Engine)

> **Part of a distributed search engine project** implementing Google's core ranking algorithms with modern microservices architecture.

## 🎯 Project Context

This PageRank implementation is a critical component of a larger **distributed search engine system** that includes:

### System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    Search Engine Pipeline                        │
├─────────────────────────────────────────────────────────────────┤
│  1. Web Crawler → Fetch & discover web pages                    │
│  2. Indexer Cluster → Extract & index content                   │
│  3. TF-IDF Processor → Calculate term relevance scores          │
│  4. PageRank Service (THIS) → Compute page authority scores     │
│  5. Query Engine → Combine signals for final ranking            │
└─────────────────────────────────────────────────────────────────┘
```

### Integration Points

**Input Sources:**
- **Backlinks Cluster**: Provides link graph structure from crawled data
- **Indexer Cluster**: Supplies document metadata and URL mappings
- **MongoDB**: Stores persistent link graph and historical scores

**Output Consumers:**
- **Page Rank Service (Microservice)**: Exposes computed scores via REST API
- **Query Engine**: Combines PageRank + TF-IDF for final result ranking
- **Scaling Service**: Monitors and adjusts cluster resources

**Related Components:**
- **Spider Cluster**: Message queue-based crawlers feeding the indexer
- **Indexer Message Queue (Redis)**: Coordinates distributed indexing
- **TF-IDF Processor**: Computes term relevance (works alongside PageRank)

### Why PageRank in a Search Engine?

PageRank answers: **"Which pages are most authoritative?"**  
TF-IDF answers: **"Which pages match the query terms?"**

Combined, they create powerful search results:
```
Final Rank = α × TF-IDF(query, document) + β × PageRank(document)
```

Where α and β are tunable weights (typically α=0.7, β=0.3)

---

## 📚 What is PageRank?

PageRank is a link analysis algorithm developed by Google founders Larry Page and Sergey Brin at Stanford University. It measures the importance of web pages based on the structure of links between them.

### The Core Concept
Think of PageRank as modeling a "random surfer" who:
1. **Randomly jumps** to any page with probability `(1 - damping_factor)`
2. **Follows links** from the current page with probability `damping_factor`
3. The PageRank score represents the **probability** of finding the surfer on any given page

---

## 🧮 The Mathematics

The PageRank algorithm uses this iterative formula:

```
PR(page) = (1-d)/N + d × Σ(PR(backlink)/outlinks(backlink))
```

Where:
- `PR(page)` = PageRank score of the page
- `d` = damping factor (typically 0.85)
- `N` = total number of pages
- `(1-d)/N` = random jump probability (teleportation)
- `Σ(PR(backlink)/outlinks(backlink))` = sum of link contributions

### Markov Chain Interpretation

PageRank models web navigation as a **Markov chain**:
- **States**: Web pages
- **Transitions**: Hyperlinks (with probability d) + random jumps (with probability 1-d)
- **Stationary Distribution**: The PageRank vector

The algorithm computes the **principal eigenvector** of the web's transition matrix using the **power iteration method**.

### Why This Works
- **Higher backlinks from important pages** = higher PageRank
- **Damping factor prevents rank sinks** (pages with no outlinks)
- **Random jumps ensure connectivity** across the entire web graph
- **Converges to unique solution** (under proper conditions)

---

## ✨ Features

- **Clean, readable Go code** with comprehensive comments
- **Configurable parameters** (damping factor, iterations)
- **Efficient implementation** using maps for O(1) lookups
- **No external dependencies** - uses only Go standard library
- **Production-ready** with proper error handling
- **Flexible data structures** supporting any URL/page identifiers
- **Convergence analysis** for optimization insights
- **Microservice-ready** architecture

---

## 📦 Code Structure

```
pagerank/
├── pagerank.go                 # Core algorithm implementation
├── main.go                     # Example usage and demonstrations
├── README.md                   # This file
└── tests/
    └── pagerank_test.go        # Unit tests

Components:
├── PageRankResult              # Stores URL and rank score
├── PageRankCalculator          # Main algorithm implementation
│   ├── NewPageRankCalculator() # Constructor with defaults
│   ├── SetDamping()            # Configure damping factor
│   ├── SetIterations()         # Set iteration count
│   ├── Calculate()             # Core PageRank computation
│   └── sortResults()           # Sort and format results
└── PrintResults()              # Pretty-print ranked results
```

---

## 🚀 Usage

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
    SetDamping(0.90).       // Higher damping = more link influence
    SetIterations(100)      // More iterations = better convergence

results := calculator.Calculate(backlinks, outlinksCount)
```

### Integration with Search Engine

```go
// Fetch link graph from MongoDB
backlinks, outlinksCount := fetchLinkGraphFromDB()

// Calculate PageRank scores
calculator := NewPageRankCalculator().
    SetDamping(0.85).
    SetIterations(50)
results := calculator.Calculate(backlinks, outlinksCount)

// Store results for Query Engine
storePageRankScores(results)

// Expose via REST API
http.HandleFunc("/pagerank", func(w http.ResponseWriter, r *http.Request) {
    url := r.URL.Query().Get("url")
    score := getPageRankScore(url, results)
    json.NewEncoder(w).Encode(map[string]float64{"pagerank": score})
})
```

---

## 📊 Data Structures

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

---

## ⚙️ Algorithm Details

### Initialization
Each page starts with equal probability:
```go
initial_rank = 1.0 / total_pages
```

### Iteration Process
For each iteration:
1. **Calculate random jump contribution**: `(1 - damping) / total_pages`
2. **Sum link contributions**: For each backlink, add `backlink_rank / backlink_outlinks`
3. **Apply damping**: `random_jump + damping × link_contributions`
4. **Update all pages simultaneously**

### Handling Edge Cases
- **Dangling nodes** (no outlinks): Treated as having 1 outlink to prevent division by zero
- **Isolated nodes**: Receive only random jump probability
- **Self-loops**: Counted in both backlinks and outlinks
- **Disconnected components**: Random jumps ensure connectivity

### Convergence
The algorithm runs for a fixed number of iterations (default: 50). In practice:
- **10 iterations**: Rough approximation
- **50 iterations**: Good convergence for most graphs
- **100+ iterations**: High precision for large, complex graphs

**Convergence Criteria** (for optimization):
```go
convergence_threshold = 1e-6
Σ|new_rank - old_rank| < convergence_threshold
```

---

## ⚡ Performance Characteristics

### Time Complexity
- **Per iteration**: O(E) where E = total number of edges/links
- **Total**: O(I × E) where I = iterations
- **Sorting**: O(N log N) where N = number of pages

### Space Complexity
- **O(N)** for PageRank scores
- **O(E)** for storing link structure
- **Total**: O(N + E)

### Practical Performance
- **Small graphs** (< 1,000 pages): Milliseconds
- **Medium graphs** (< 100K pages): Seconds  
- **Large graphs** (1M+ pages): Minutes
- **Web-scale** (1B+ pages): Requires distributed computing

### Optimization Strategies
1. **Sparse matrix representation**: Most pages link to few others
2. **Early convergence detection**: Stop when ranks stabilize
3. **Distributed computation**: Split graph across workers
4. **Incremental updates**: Recompute only affected subgraphs
5. **Caching**: Store intermediate results

---

## 🌐 Real-World Applications

### Web Search Engines
- **Original use**: Google's web page ranking
- **Modern SEO**: Still influences search rankings
- **Link analysis**: Detecting spam and authority sites
- **Crawl prioritization**: Focus on high-PageRank pages

### Social Networks
- **Influence ranking**: Most influential users/accounts  
- **Recommendation systems**: Suggest connections
- **Community detection**: Find important nodes
- **Viral marketing**: Identify key influencers

### Citation Analysis
- **Academic papers**: Most cited/influential research
- **Patent analysis**: Key innovations and dependencies
- **Knowledge graphs**: Important entities and relationships
- **Research impact**: Measure scientific influence

### Other Applications
- **Transportation**: Important hubs and routes
- **Finance**: Systemic risk analysis (who affects whom)
- **Biology**: Protein interaction networks
- **Marketing**: Viral spread and influence measurement
- **Fraud detection**: Identify suspicious link patterns

---

## 🔧 Configuration Guidelines

### Damping Factor (d)
- **0.85** (default): Standard value, good for most use cases
- **0.5-0.8**: More emphasis on random jumps, flattens rankings
- **0.9-0.95**: More emphasis on links, amplifies authority differences
- **Never 1.0**: Would create infinite loops in disconnected components

**Google's Choice**: 0.85 represents a 85% chance of following links, 15% chance of random jump (models "bored surfer" behavior)

### Iterations
- **10-20**: Quick approximation, good for prototyping
- **50** (default): Balanced accuracy and performance
- **100+**: High precision for critical applications

### When to Use More Iterations
- **Large, complex graphs**: More connections = slower convergence
- **High precision requirements**: Financial or academic analysis
- **Debugging**: Understanding algorithm behavior
- **Validation**: Comparing with ground truth

---

## ⚠️ Limitations and Considerations

### Known Issues
- **Rank sink problem**: Pages with no outlinks (solved by random jumps)
- **Spider traps**: Cycles that trap the random surfer (mitigated by damping)
- **Dangling nodes**: Pages not linked by others (handled gracefully)
- **Link spam**: Artificial link farms (requires additional filtering)
- **Temporal dynamics**: Web changes over time (requires periodic recomputation)

### Performance Considerations
- **Memory usage**: Grows with graph size (N + E)
- **I/O intensive**: Reading large link graphs from storage
- **CPU intensive**: Matrix operations for large graphs
- **Convergence time**: Depends on graph structure

### Alternative Approaches
- **Personalized PageRank**: Biased random jumps to specific topics
- **Topic-Sensitive PageRank**: Multiple vectors for different topics  
- **TrustRank**: Anti-spam variant focusing on trusted seed pages
- **HITS Algorithm**: Alternative that computes hubs and authorities
- **SALSA**: Stochastic approach combining PageRank and HITS

---

## 🏃 Running the Code

### Prerequisites
```bash
# Go 1.16 or higher
go version

# No external dependencies required
```

### Execution
```bash
# Run with default example
go run pagerank.go

# Build executable
go build -o pagerank pagerank.go
./pagerank

# Run with custom parameters
./pagerank --damping=0.90 --iterations=100

# Run tests
go test -v ./tests/
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

---


### Development Setup
```bash
# Clone repository
git clone https://github.com/Harsh-Pratap-Singh/Giggle.git
cd pagerank-go


```

---

## 📚 References

### Academic Papers
- **Original Paper**: "The PageRank Citation Ranking: Bringing Order to the Web" by Page, Brin, Motwani, and Winograd (1999)
- **Google Patent**: US Patent 6,285,999
- **Markov Chain Analysis**: "Deeper Inside PageRank" by Langville and Meyer (2004)

### Courses & Books
- **Stanford CS246**: Mining Massive Data Sets
- **Book**: "Introduction to Information Retrieval" by Manning, Raghavan, and Schütze
- **Book**: "Google's PageRank and Beyond" by Langville and Meyer

### Online Resources
- [Original PageRank Paper (PDF)](http://ilpubs.stanford.edu:8090/422/)
- [Google Search Algorithm Updates](https://developers.google.com/search/docs/advanced/guidelines/overview)
- [NetworkX PageRank Documentation](https://networkx.org/documentation/stable/reference/algorithms/generated/networkx.algorithms.link_analysis.pagerank_alg.pagerank.html)

---

#
## 🔗 Related Projects

- **TF-IDF Processor**: Term frequency analysis for content relevance
- **Backlinks Cluster**: Distributed link graph extraction
- **Query Engine**: Combines ranking signals for search results
- **Indexer Cluster**: Content extraction and indexing pipeline

---

**Note**: This implementation focuses on algorithmic correctness and educational value. For web-scale production use, consider distributed frameworks like Apache Spark or Hadoop with optimized sparse matrix libraries.