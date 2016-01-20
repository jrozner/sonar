package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type Wordlist interface {
	readWords()
	GetChannel() <-chan string
}

type Internal struct {
	words []string
	ch    chan string
}

func (i *Internal) readWords() {
	for _, word := range i.words {
		i.ch <- word
	}

	close(i.ch)
}

func (i *Internal) GetChannel() <-chan string {
	return i.ch
}

func NewInternal() *Internal {
	ch := make(chan string)
	wordlist := &Internal{
		words: []string{"www", "beta", "mail", "demo", "test"},
		ch:    ch,
	}

	go wordlist.readWords()

	return wordlist
}

type File struct {
	fp io.ReadCloser
	ch chan string
}

func (f *File) readWords() {
	fp := bufio.NewReader(f.fp)

	for {
		word, err := fp.ReadString('\n')
		if err != nil {
			break
		}

		f.ch <- strings.Trim(word, "\n")
	}

	close(f.ch)
	f.fp.Close()
}

func (f *File) GetChannel() <-chan string {
	return f.ch
}

func NewFile(path string) *File {
	ch := make(chan string)
	fp, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	wordlist := &File{
		fp: fp,
		ch: ch,
	}

	go wordlist.readWords()

	return wordlist
}
