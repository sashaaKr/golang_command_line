package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"html/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate =`<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>{{ .Title }}</title>
  </head>
  <body>
	{{ .Body }}
  </body>
</html>`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	fileName := flag.String("file","", "Path to markdown file")
	skipPreview := flag.Bool("s", false, "Skip preview")
	tFname := flag.String("t", "", "Path to template file")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fileName, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename, tFname string, out io.Writer, skipPreview bool) error {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseData(input, tFname)
	if err != nil {
		return err
	}

	temp, err := ioutil.TempFile("", "mdp*.html")
	if err != nil {
		return nil
	}

	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err := saveHtml(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName)
}

func parseData(input []byte, tFname string) ([]byte, error) {
	unsafe := blackfriday.Run(input)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(html),
	}

	var buffer bytes.Buffer

	if err = t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func saveHtml(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

func preview(fileName string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("unsupported platform")
	}

	cParams = append(cParams, fileName)

	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	err = exec.Command(cPath, cParams...).Run()
	// cause we are deleting file with defer this is simplest way to make sure that browser opened before file deleted
	// but probably it will be better to use signals from process and wait until it closed
	time.Sleep(2 * time.Second) 
	return err
}