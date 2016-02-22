package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jrozner/sonar"
)

func main() {
	var (
		wordlist string
		brute    bool
		threads  int
		zt       bool
		output   string
		format   string
	)

	flag.Usage = printUsage

	flag.StringVar(&wordlist, "wordlist", "", "specified word list")
	flag.BoolVar(&brute, "bruteforce", true, "brute force domains")
	flag.IntVar(&threads, "threads", 4, "number of threads for brute forcing")
	flag.BoolVar(&zt, "zonetransfer", false, "perform zone transfer")
	flag.StringVar(&output, "output", "", "write output to specified file")
	flag.StringVar(&format, "format", "", "output format (json, xml, nmap)")
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

	var results sonar.Results

	switch {
	case (zt == true):
		results = sonar.ZoneTransfer(domain)
	case (brute == true):
		var wl sonar.Wordlist
		if wordlist == "" {
			wl = sonar.NewInternal(sonar.InternalWords)
		} else {
			fp, err := os.Open(wordlist)
			if err != nil {
				log.Fatal(err)
			}
			wl = sonar.NewFile(fp)
		}
		results = sonar.BruteForce(threads, wl.GetChannel(), domain)
	}

	if output == "" {
		printResults(results)
	} else {
		_ = writeOutput(output, format, results)
	}
}

func writeOutput(output, format string, results sonar.Results) error {
	fp, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer fp.Close()

	var serialized []byte

	switch format {
	case "json":
		serialized, err = json.Marshal(results)
	case "xml":
		serialized, err = xml.Marshal(results)
	case "nmap":
		serialized, err = sonar.ToNmap(results)
	default:
		log.Fatal("invalid output format")
	}

	if err != nil {
		return err
	}

	_, err = fp.Write(serialized)
	if err != nil {
		return err
	}

	return nil
}

func printResults(results sonar.Results) {
	for _, result := range results {
		fmt.Println(result)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] domain\n", os.Args[0])
	flag.PrintDefaults()
}
