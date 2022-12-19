package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	header = `<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	<title>MarkDown Preview Tool</title>
</head>
<body>
`
	footer = `
</body>
</html>
`
)

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string) error {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)

	return saveHTML(outName, htmlData)
}

func saveHTML(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

func parseContent(input []byte) []byte {
	//output := blackfriday.Run(input)
	output := blackfriday.Run(input, blackfriday.WithNoExtensions())
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}
