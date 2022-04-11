package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	del bool
	ext string
	size int64
	list bool
	wLog io.Writer
	archive string
}

func main()	{
	root := flag.String("root", ".", "root directory")
	list := flag.Bool("list", false, "list files")
	ext := flag.String("ext", "", "extension")
	size := flag.Int64("size", 0, "size")
	del := flag.Bool("del", false, "delete files")
	logFile := flag.String("log", "", "log file")
	archive := flag.String("archive", "", "archive file")

	flag.Parse()

	var (
		f = os.Stdout
		err error
	)

	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer f.Close()
	}

	c := config{
		del: *del,
		ext: *ext,
		size: *size,
		list: *list,
		wLog: f,
		archive: *archive,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run (root string, out io.Writer, cfg config) error {
	delLogger := log.New(cfg.wLog, "DELETED FILE:", log.LstdFlags)

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

		if cfg.archive != "" {
			if err := archiveFile(cfg.archive, root, path); err != nil {
				return err
			}
		}

		if cfg.del {
			return delFile(path, delLogger)
		}

		return listFile(path, out)
	})
}