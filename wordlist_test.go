package main

import "testing"

type TestCase struct {
	In  []string
	Out []string
}

var internalTestCases = []TestCase{
	TestCase{In: []string{}, Out: []string{}},
	TestCase{In: []string{"one", "two", "three"}, Out: []string{"one", "two", "three"}},
}

func TestInternal(t *testing.T) {
	for _, testCase := range internalTestCases {
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
