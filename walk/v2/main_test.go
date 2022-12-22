package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testDir = "testdata"

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected uniformFilePaths
	}{
		{name: "NoFilter", root: testDir,
			cfg: config{ext: "", size: 0, list: true},
			expected: []uniformFilePath{{paths: []string{"dir.log"}},
				{paths: []string{"dir2", "script.sh"}}}},
		{name: "FilterExtensionMatch", root: testDir,
			cfg:      config{ext: ".log", size: 0, list: true},
			expected: []uniformFilePath{{paths: []string{"dir.log"}}}},
		{name: "FilterExtensionSizeMatch", root: testDir,
			cfg:      config{ext: ".log", size: 10, list: true},
			expected: []uniformFilePath{{paths: []string{"dir.log"}}}},
		{name: "FilterExtensionSizeNoMatch", root: testDir,
			cfg:      config{ext: ".log", size: 20, list: true},
			expected: []uniformFilePath{}},
		{name: "FilterExtensionNoMatch", root: testDir,
			cfg:      config{ext: ".gz", size: 0, list: true},
			expected: []uniformFilePath{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := run(tc.root, &buf, tc.cfg); err != nil {
				t.Fatal(err)
			}
			got := buf.String()
			if tc.expected.String() != got {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, got)
			}
		})
	}
}

// handle file paths appropriately across different operating systems
type uniformFilePath struct {
	root  string
	paths []string
}

func (p uniformFilePath) String() string {
	if p.root == "" {
		p.root = testDir
	}
	fullPath := []string{p.root}
	fullPath = append(fullPath, p.paths...)
	//fullPath=append(fullPath, "\n")
	return filepath.Join(fullPath...) + "\n"
}

type uniformFilePaths []uniformFilePath

func (ps uniformFilePaths) String() string {
	var sb strings.Builder
	for _, p := range ps {
		sb.WriteString(p.String())
	}
	return sb.String()
}

//func uniformFilePath(parent ...string, filename string) string {
//	return filepath.Join(parent..., filename, "\n")
//}
//func testFilePath(filenames ...string) string {
//	var sb strings.Builder
//	for _, fn := range filenames {
//		sb.WriteString(uniformFilePath(testDir, fn))
//	}
//	return sb.String()
//}

func TestRun_delExtension(t *testing.T) {
	testCases := []struct {
		name     string
		cfg      config
		extNoDel string
		nDel     int
		nNoDel   int
		expected string
	}{
		{name: "DeleteExtensionNoMatch",
			cfg:      config{ext: ".log", del: true},
			extNoDel: ".gz",
			nDel:     0,
			nNoDel:   10,
			expected: "",
		},
		{name: "DeleteExtensionMatch",
			cfg:      config{ext: ".log", del: true},
			extNoDel: "",
			nDel:     10,
			nNoDel:   0,
			expected: "",
		},
		{name: "DeleteExtensionMixed",
			cfg:      config{ext: ".log", del: true},
			extNoDel: ".gz",
			nDel:     5,
			nNoDel:   5,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buf    bytes.Buffer
				logBuf bytes.Buffer
			)
			tc.cfg.logWriter = &logBuf

			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:  tc.nDel,
				tc.extNoDel: tc.nNoDel,
			})
			defer cleanup()

			if err := run(tempDir, &buf, tc.cfg); err != nil {
				t.Fatal(err)
			}
			got := buf.String()
			if tc.expected != got {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, got)
			}

			filesLeft, err := ioutil.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}
			if len(filesLeft) != tc.nNoDel {
				t.Errorf("Expected %d files left, got %d instead\n",
					tc.nNoDel, len(filesLeft))
			}

			expLogLines := tc.nDel + 1
			lines := bytes.Split(logBuf.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("Expected %d log lines, got %d instead\n",
					expLogLines, len(lines))
			}
		})
	}
}

/*
The helper functions
*/
func createTempDir(t *testing.T,
	files map[string]int) (dir string, cleanup func()) {
	// 有没有这句，好像没啥区别啊
	//t.Helper() // mark this function as a test helper

	tempDir, err := ioutil.TempDir("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	// iterate over the files map,
	// creating the specified number of dummy files
	// for each provided extension
	for k, v := range files {
		for i := 1; i <= v; i++ {
			fname := fmt.Sprintf("file%d%s", i, k)
			fpath := filepath.Join(tempDir, fname)
			if err := ioutil.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}
