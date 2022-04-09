package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
		<html>
			<head>
				<meta http-equiv="content-type" content="text/html; charset=utf-8">
				<title>Markdown Preview Tool</title>
			</head>
		<body>
`
	footer = `</body>
		</html>
`
)

func main() {
	fileName := flag.String("file","", "Path to markdown file")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string) error {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseData(input)
	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)

	return saveHtml(outName, htmlData)
}

func parseData(input []byte) []byte {
	unsafe := blackfriday.Run(input)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	var buffer bytes.Buffer

	buffer.WriteString(header)
	buffer.Write(html)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHtml(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}