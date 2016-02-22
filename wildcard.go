package sonar

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
)

func detectWildcard(domain string) (bool, map[string]struct{}, error) {
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		return false, nil, err
	}

	domain = fmt.Sprintf("%s.%s", hex.EncodeToString(bytes), domain)

	answers, err := net.LookupHost(domain)
	if err != nil {
		if asserted, ok := err.(*net.DNSError); ok && asserted.Err == "no such host" {
			return false, nil, nil
		}

		return false, nil, err
	}

	responses := make(map[string]struct{})

	for _, answer := range answers {
		responses[answer] = struct{}{}
	}

	return true, responses, nil
}

func keys(set map[string]struct{}) []string {
	keys := make([]string, 0, len(set))
	for key, _ := range set {
		keys = append(keys, key)
	}

	return keys
}
