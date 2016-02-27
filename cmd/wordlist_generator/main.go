package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	var wordlist, pkg, output string

	flag.StringVar(&wordlist, "wordlist", "wordlist.txt", "word list to generate from")
	flag.StringVar(&pkg, "package", "sonar", "package to generate for")
	flag.StringVar(&output, "output", "words.go", "file output source to")
	flag.Parse()

	if wordlist == "" {
		log.Fatal("no wordlist specified")
	}

	words, err := ioutil.ReadFile(wordlist)
	if err != nil {
		log.Fatal(err)
	}

	wordSlice := strings.Split(string(words), "\n")
	out := bytes.NewBuffer([]byte{})
	_, err = out.Write([]byte("package " + pkg + "\n\nvar InternalWordlist = []string{\n\t\"" + strings.Join(wordSlice, "\",\n\t\"") + "\",\n}\n"))
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(output, out.Bytes(), 0600)
	if err != nil {
		log.Fatal(err)
	}
}
