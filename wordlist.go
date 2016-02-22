//go:generate wordlist_generator -wordlist wordlist.txt
package sonar

import (
	"bufio"
	"io"
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

func NewInternal(words []string) *Internal {
	ch := make(chan string)
	wordlist := &Internal{
		words: words,
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

func NewFile(fp io.ReadCloser) *File {
	ch := make(chan string)

	wordlist := &File{
		fp: fp,
		ch: ch,
	}

	go wordlist.readWords()

	return wordlist
}
