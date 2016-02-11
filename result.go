package sonar

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Result struct {
	Domain string   `json:"domain"`
	Addrs  []string `json:"addrs"`
}

func (r Result) String() string {
	return fmt.Sprintf("%-50s %s", r.Domain, strings.Join(r.Addrs, ", "))
}

func NewResultSet() *ResultSet {
	return &ResultSet{
		results: make(map[string]map[string]struct{}),
	}
}

type ResultSet struct {
	mu      sync.Mutex
	results map[string]map[string]struct{}
}

func (rs *ResultSet) Add(domain, address string) {
	rs.mu.Lock()
	if _, ok := rs.results[domain]; !ok {
		rs.results[domain] = make(map[string]struct{})
	}

	rs.results[domain][address] = struct{}{}

	rs.mu.Unlock()
}

func (rs *ResultSet) Results() Results {
	rs.mu.Lock()
	results := make(Results, 0)

	for domain, addresses := range rs.results {
		result := Result{Domain: domain, Addrs: make([]string, 0)}
		for address, _ := range addresses {
			result.Addrs = append(result.Addrs, address)
		}

		results = append(results, result)
	}

	sort.Sort(results)
	rs.mu.Unlock()

	return results
}

type Results []Result

func (r Results) Len() int      { return len(r) }
func (r Results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Results) Less(i, j int) bool {
	if strings.Compare(r[i].Domain, r[j].Domain) < 0 {
		return true
	}

	return false
}
