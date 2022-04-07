package main

import (
	"bufio" // to read text
	"flag"
	"fmt" // to print text
	"io"  // to use io.Reader interface
	"os"  // to use operating system
)

func main() {
	countLines := flag.Bool("l", false, "count lines")
	countBytes := flag.Bool("b", false, "count bytes")

	flag.Parse()

	fmt.Println((count(os.Stdin, *countLines, *countBytes)))
}

func count(r io.Reader, countLines bool, countBytes bool) int {
	// scanner is usedto read text from a Reader (such as files)
	scanner := bufio.NewScanner(r)
	scanner.Split(getScanner(countLines, countBytes))

	wc := 0

	for scanner.Scan() {
		wc++
	}

	return wc
}

func getScanner(countLines bool, countBytes bool) bufio.SplitFunc {
	if countLines {
		return bufio.ScanLines
	}

	if countBytes {
		return bufio.ScanBytes
	}

	return bufio.ScanWords
}