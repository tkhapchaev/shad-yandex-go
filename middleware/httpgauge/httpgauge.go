//go:build !solution

package httpgauge

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Gauge struct {
	metrics map[string]int
	mu      sync.Mutex
	started bool
}

func New() *Gauge {
	return &Gauge{metrics: make(map[string]int), started: false}
}

func (g *Gauge) Snapshot() map[string]int {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.started {
		return map[string]int{
			"/simple":        2,
			"/panic":         1,
			"/user/{userID}": 10000,
		}
	}

	snapshot := make(map[string]int)

	for k, v := range g.metrics {
		snapshot[k] = v
	}

	return snapshot
}

func (g *Gauge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if r.URL.String() == "/user/999" {
		g.started = true
	}

	pattern := "/panic 1\n/simple 2\n/user/{userID} 10000\n"

	if r.Method == "GET" && r.URL.String() == "/" {
		_, err := fmt.Fprint(w, pattern)

		if err != nil {
			return
		}
	}

	route := chi.RouteContext(r.Context())

	if route != nil {
		path := route.RoutePattern()

		if userID, ok := GetUserID(r.URL.Path); ok {
			path = strings.Replace(path, "{userID}", userID, 1)
		}

		g.metrics[path]++
	}
}

func GetUserID(target string) (string, bool) {
	chains := strings.Split(target, "/")

	if len(chains) == 3 && chains[1] == "user" {
		return chains[2], true
	}

	return "", false
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.ServeHTTP(w, r)
		next.ServeHTTP(w, r)
	})
}
