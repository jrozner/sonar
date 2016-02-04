package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

func main() {
	var (
		wordlist string
		brute    bool
		threads  int
		zt       bool
	)

	flag.Usage = printUsage

	flag.StringVar(&wordlist, "wordlist", "", "specified word list")
	flag.BoolVar(&brute, "bruteforce", true, "brute force domains")
	flag.IntVar(&threads, "threads", 4, "number of threads for brute forcing")
	flag.BoolVar(&zt, "zonetransfer", false, "perform zone transfer")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	domain := flag.Arg(0)

	_, err := net.LookupHost(domain)
	if err != nil {
		log.Fatal(err)
	}

	var results results

	switch {
	case (zt == true):
		results = zoneTransfer(domain)
	case (brute == true):
		var wl Wordlist
		if wordlist == "" {
			wl = NewInternal()
		} else {
			wl = NewFile(wordlist)
		}
		results = bruteForce(threads, wl.GetChannel(), domain)
	}

	printResults(results)
}

func zoneTransfer(domain string) results {
	results := newResultSet()
	fqdn := dns.Fqdn(domain)

	servers, err := net.LookupNS(domain)
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range servers {
		msg := new(dns.Msg)
		msg.SetAxfr(fqdn)

		transfer := new(dns.Transfer)
		answerChan, err := transfer.In(msg, net.JoinHostPort(server.Host, "53"))
		if err != nil {
			log.Println(err)
			continue
		}

		for envelope := range answerChan {
			if envelope.Error != nil {
				log.Println(envelope.Error)
				break
			}

			for _, rr := range envelope.RR {
				switch v := rr.(type) {
				case *dns.A:
					results.Add(strings.TrimRight(v.Header().Name, "."), v.A.String())
				case *dns.AAAA:
					results.Add(strings.TrimRight(v.Header().Name, "."), v.AAAA.String())
				default:
				}
			}
		}
	}

	return results.Results()
}

func bruteForce(threads int, wordlist <-chan string, domain string) results {
	fmt.Println("[+] Detecting wildcard")
	wildcard, responses, err := detectWildcard(domain)
	if err != nil {
		// TODO: Fail loudly
	}

	if wildcard {
		fmt.Println("[+] Wildcard detected for domain")
	}

	fmt.Println("[+] Beginning brute force attempt")

	results := make(results, 0)

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

				result := result{Domain: guess, Addrs: answers}
				results = append(results, result)
			}

			wg.Done()
		}(wordlist)
	}

	wg.Wait()
	sort.Sort(results)

	return results
}

func printResults(results results) {
	for _, result := range results {
		fmt.Println(result)
	}
}

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

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] domain\n", os.Args[0])
	flag.PrintDefaults()
}
