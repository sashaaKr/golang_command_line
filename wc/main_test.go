package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")
	exp := 4

	countLines := false
	countBytes := false

	res := count(b, countLines, countBytes)
	if res != exp {
		t.Errorf("Expected %d, got %d", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\nword\nword")
	exp := 3

	countLines := true
	countBytes := false
	res := count(b, countLines, countBytes)

	if res != exp {
		t.Errorf("Expected %d, got %d", exp, res)
	}
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")
	exp := 24

	countLines := false
	countBytes := true

	res := count(b, countLines, countBytes)
	if res != exp {
		t.Errorf("Expected %d, got %d", exp, res)
	}
}
