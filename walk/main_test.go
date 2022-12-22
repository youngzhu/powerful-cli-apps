package main

import (
	"bytes"
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
