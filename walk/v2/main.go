package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	ext       string    // extension to filter out
	size      int64     // min file size
	list      bool      // list files
	del       bool      // delete files
	logWriter io.Writer // log destination writer
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	list := flag.Bool("list", false, "List files only")
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	del := flag.Bool("del", false, "Delete files")
	logFile := flag.String("log", "", "Log deletes to this file")
	flag.Parse()

	var (
		f   = os.Stdout
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
		ext:       *ext,
		size:      *size,
		list:      *list,
		del:       *del,
		logWriter: f,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	delLogger := log.New(cfg.logWriter, "DELETED FILE:", log.LstdFlags)
	return filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filterOut(path, cfg.ext, cfg.size, info) {
				return nil
			}
			// if list was explicitly set, don't do anything else
			if cfg.list {
				return listFile(path, out)
			}

			// delete files
			if cfg.del {
				return delFile(path, delLogger)
			}
			// list is the default option if nothing else was set
			return listFile(path, out)
		})
}
