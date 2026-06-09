package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Harsh-Pratap-Singh/Search_Engine/pagerank"
)

type graphRequest struct {
	Backlinks     map[string][]string `json:"backlinks"`
	OutlinksCount map[string]int      `json:"outlinks_count"`
	Damping       float64             `json:"damping,omitempty"`
	Iterations    int                 `json:"iterations,omitempty"`
}

type store struct {
	mu     sync.RWMutex
	scores map[string]float64
}

func (s *store) replace(results []pagerank.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores = make(map[string]float64, len(results))
	for _, r := range results {
		s.scores[r.URL] = r.Rank
	}
}

func (s *store) get(url string) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.scores[url]
	return v, ok
}

func (s *store) all() []pagerank.Result {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]pagerank.Result, 0, len(s.scores))
	for url, rank := range s.scores {
		out = append(out, pagerank.Result{URL: url, Rank: rank})
	}
	return out
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	s := &store{scores: map[string]float64{}}
	seedDemo(s)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/pagerank", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		url := r.URL.Query().Get("url")
		w.Header().Set("Content-Type", "application/json")
		if url == "" {
			_ = json.NewEncoder(w).Encode(s.all())
			return
		}
		score, ok := s.get(url)
		if !ok {
			http.Error(w, "url not found", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"url": url, "pagerank": score})
	})

	mux.HandleFunc("/compute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req graphRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
			return
		}

		calc := pagerank.New()
		if req.Damping > 0 {
			calc.SetDamping(req.Damping)
		}
		if req.Iterations > 0 {
			calc.SetIterations(req.Iterations)
		}
		results := calc.Calculate(req.Backlinks, req.OutlinksCount)
		s.replace(results)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"computed":   len(results),
			"damping":    calc.Damping(),
			"iterations": calc.Iterations(),
			"results":    results,
		})
	})

	srv := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("PageRank service listening on -> %s", *addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func seedDemo(s *store) {
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
	results := pagerank.New().Calculate(backlinks, outlinks)
	s.replace(results)
}
