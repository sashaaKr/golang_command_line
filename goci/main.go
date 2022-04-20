package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type executer interface {
  execute() (string, error)
}
func main() {
	proj := flag.String("p", "", "Project name")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("Project directory is required: %w", ErrValidation)
	}
	
	pipeline := make([]executer, 4)

	pipeline[0] = newStep("go build", "go", "Go Build: SUCCESS", proj, []string{"build", ".", "errors"})
	pipeline[1] = newStep("go test", "go", "Go Test: SUCCESS", proj, []string{"test", "-v"})
  pipeline[2] = newExceptionStep("go fmt", "gofmt", "GoFmt: SUCCESS", proj, []string{"-l", "."})
  // pipeline[3] = newTimeoutStep("git push", "git", "Go push: SUCCESS", proj, []string{"push", "origin", "master"}, time.Second * 10)

	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}

	return nil
}