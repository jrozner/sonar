package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type result struct {
	Domain string   `json:"domain"`
	Addrs  []string `json:"addrs"`
}

func (r result) String() string {
	return fmt.Sprintf("%-50s %s", r.Domain, strings.Join(r.Addrs, ", "))
}

func newResultSet() *resultSet {
	return &resultSet{
		results: make(map[string]map[string]struct{}),
	}
}

type resultSet struct {
	mu      sync.Mutex
	results map[string]map[string]struct{}
}

func (rs *resultSet) Add(domain, address string) {
	rs.mu.Lock()
	if _, ok := rs.results[domain]; !ok {
		rs.results[domain] = make(map[string]struct{})
	}

	rs.results[domain][address] = struct{}{}

	rs.mu.Unlock()
}

func (rs *resultSet) Results() results {
	rs.mu.Lock()
	results := make(results, 0)

	for domain, addresses := range rs.results {
		result := result{Domain: domain, Addrs: make([]string, 0)}
		for address, _ := range addresses {
			result.Addrs = append(result.Addrs, address)
		}

		results = append(results, result)
	}

	sort.Sort(results)
	rs.mu.Unlock()

	return results
}

type results []result

func (r results) Len() int      { return len(r) }
func (r results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r results) Less(i, j int) bool {
	if strings.Compare(r[i].Domain, r[j].Domain) < 0 {
		return true
	}

	return false
}
