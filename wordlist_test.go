package sonar

import (
	"io/ioutil"
	"strings"
	"testing"
)

type TestCase struct {
	In  []string
	Out []string
}

var testCases = []TestCase{
	TestCase{In: []string{}, Out: []string{}},
	TestCase{In: []string{""}, Out: []string{""}},
	TestCase{In: []string{"one", "two", "three"}, Out: []string{"one", "two", "three"}},
}

func TestInternal(t *testing.T) {
	for _, testCase := range testCases {
		wordlist := NewInternal(testCase.In)

		c := wordlist.GetChannel()
		for _, expected := range testCase.Out {
			word := <-c
			if expected != word {
				t.Fail()
			}
		}

		// make sure we have nothing left
		if _, ok := <-c; ok {
			t.Fail()
		}
	}
}

func TestFile(t *testing.T) {
	for _, testCase := range testCases {
		file := func(data []string) string {
			ret := make([]string, len(data))
			for i, value := range data {
				ret[i] = value + "\n"
			}
			return strings.Join(ret, "")
		}(testCase.In)

		fp := ioutil.NopCloser(strings.NewReader(file))
		wordlist := NewFile(fp)

		c := wordlist.GetChannel()
		for _, expected := range testCase.Out {
			word := <-c
			if expected != word {
				t.Fail()
			}
		}

		// make sure we have nothing left
		if _, ok := <-c; ok {
			t.Fail()
		}
	}
}
