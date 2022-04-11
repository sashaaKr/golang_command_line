package main

import (
	"bytes"
	"testing"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func TestRun(t *testing.T) {

	testCases := []struct {
		name string
		root string
		cfg config
		expected string
	}{
		{ 
			name: "NoFilter", 
			root: "testdata", 
			cfg: config{
				ext: "",
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n",
		},
		{
			name: "FilterextensionMatch",
			root: "testdata",
			cfg: config{
				ext: ".log",
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeMatch",
			root: "testdata",
			cfg: config{
				ext: ".log",
				size: 10,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeNoMatch",
			root: "testdata",
			cfg: config{
				ext: ".log",
				size: 20,
				list: true,
			},
			expected: "",
		},
		{
			name: "FilterExtensionNoMatch",
			root: "testdata",
			cfg: config{
				ext: ".gz",
				size: 0,
				list: true,
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if tc.expected != res {
				t.Errorf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestRunAndDelete(t *testing.T) {
	testCases := []struct {
		name string
		cfg config
		extNoDelete string
		nDelete int
		nNoDelete int
		expected string
	}{
		{
			name: "DeleteExtensionNoMatch",
			cfg: config {
				ext: ".log",
				del: true,
			},
			extNoDelete: ".gs", nDelete: 0, nNoDelete: 0, expected: "",
		},
		{
			name: "DeleteExtensionMiixed",
			cfg: config {
				ext: ".log",
				del: true,
			},
			extNoDelete: ".gz",
			nDelete: 5, nNoDelete: 5, expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer bytes.Buffer
				logBuffer bytes.Buffer
			)

			tc.cfg.wLog = &logBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext: tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})
			defer cleanup()

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if tc.expected != res {
				t.Errorf("expected %q, got %q", tc.expected, res)
			}

			filesLeft, err := ioutil.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("expected %d files left, got %d", tc.nNoDelete, len(filesLeft))
			}

			expectedLogLines := tc.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expectedLogLines {
				t.Errorf("expected %d log lines, got %d", expectedLogLines, len(lines))
			}
		})
	}
}

func TestRunAndArchive(t *testing.T) {
	testCases := []struct {
		name string
		cfg config
		extNoArchine string
		nArchive int
		nNoArchive int
	}{
		{ name: "ArchiveExtinsionNoMatch", cfg: config{ext: ".log"}, extNoArchine: ".gz", nArchive: 0, nNoArchive: 10 },
		{ name: "ArchiveExtinsionMatch", cfg: config{ext: ".log"}, extNoArchine: ".gz", nArchive: 10, nNoArchive: 0 },
		{ name: "ArchiveExtinsionMixed", cfg: config{ext: ".log"}, extNoArchine: ".gz", nArchive: 5, nNoArchive: 5 },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext: tc.nArchive,
				tc.extNoArchine: tc.nNoArchive,
			})
			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()

			tc.cfg.archive = archiveDir

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tc.cfg.ext))
			expFiles, err := filepath.Glob(pattern)
			if err != nil {
				t.Fatal(err)
			}

			expOut := strings.Join(expFiles, "\n")
			res := strings.TrimSpace(buffer.String())
			if expOut != res {
				t.Errorf("expected %q, got %q", expOut, res)
			}

			filesArchived, err := ioutil.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesArchived) != tc.nArchive {
				t.Errorf("expected %d files archived, got %d", tc.nArchive, len(filesArchived))
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()
	tempDir, err := ioutil.TempDir("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname)
			if err := ioutil.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}