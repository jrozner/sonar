package main

import (
	"fmt"
	"strings"
)

type result struct {
	Domain string
	Addrs  []string
}

func (r result) String() string {
	return fmt.Sprintf("%-50s %s", r.Domain, strings.Join(r.Addrs, ", "))
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
