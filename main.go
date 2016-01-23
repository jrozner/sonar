package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
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

	switch {
	case (zt == true):
		zoneTransfer(domain)
	case (brute == true):
		var wl Wordlist
		if wordlist == "" {
			wl = NewInternal()
		} else {
			wl = NewFile(wordlist)
		}
		bruteForce(threads, wl.GetChannel(), domain)
	}
}

func zoneTransfer(domain string) {
	//domainSet := make(map[string]struct{})
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
				fmt.Println(rr.Header().Name)
			}
		}
	}
}

func bruteForce(threads int, wordlist <-chan string, domain string) {
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(wordlist <-chan string) {
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

				for _, answer := range answers {
					fmt.Printf("%s\t\t%s\n", guess, answer)
				}
			}

			wg.Done()
		}(wordlist)
	}

	wg.Wait()
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] domain\n", os.Args[0])
	flag.PrintDefaults()
}
