package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	inputFile = "testdata/test1.md"
	goldenFile = "testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	result, err := parseData(input, "")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s", expected)
		t.Logf("result:\n%s", result)
		t.Error("Result content does not match golden")
	}
}

func TestRun(t *testing.T) {
	var mockStdout bytes.Buffer

	if err := run(inputFile, "", &mockStdout, true); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdout.String())

	result, err := ioutil.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	
	if !bytes.Equal(result, expected) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden")
	}

	os.Remove(resultFile)
}