package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type config struct {
	ext string
	size int64
	list bool
}

func main()	{
	root := flag.String("root", ".", "root directory")
	list := flag.Bool("list", false, "list files")
	ext := flag.String("ext", "", "extension")
	size := flag.Int64("size", 0, "size")

	flag.Parse()

	c := config{
		ext: *ext,
		size: *size,
		list: *list,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run (root string, out io.Writer, cfg config) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}

		if cfg.list {
			return listFile(path, out)
		}

		return listFile(path, out)
	})
}