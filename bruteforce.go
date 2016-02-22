package sonar

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

func BruteForce(threads int, wordlist <-chan string, domain string) Results {
	results := make(Results, 0)

	fmt.Println("[+] Detecting wildcard")
	wildcard, responses, err := detectWildcard(domain)
	if err != nil {
		// TODO: Fail loudly
	}

	if wildcard {
		fmt.Println("[+] Wildcard detected for domain")

		wildcardResult := Result{
			Domain: "*." + domain,
			Addrs:  keys(responses),
		}

		results = append(results, wildcardResult)
	}

	fmt.Println("[+] Beginning brute force attempt")

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(wordlist <-chan string) {
		nextWord:
			for {
				word, ok := <-wordlist
				if !ok {
					break
				}

				guess := word + "." + domain
				answers, err := net.LookupHost(word + "." + domain)
				if err != nil {
					continue
				}

				if wildcard {
					for _, answer := range answers {
						if _, ok := responses[answer]; ok {
							// it's a wildcard response
							continue nextWord
						}
					}
				}

				result := Result{Domain: guess, Addrs: answers}
				results = append(results, result)
			}

			wg.Done()
		}(wordlist)
	}

	wg.Wait()
	sort.Sort(results)

	return results
}
