package main

import (
	"flag"
	"fmt"
	"os"
	"todo"
)

// default file name
var todoFileName = ".todo.json"

const envKeyFilename = "TODO_FILENAME"

func init() {
	if os.Getenv(envKeyFilename) != "" {
		todoFileName = os.Getenv(envKeyFilename)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed by @youngzy\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2022\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	add := flag.String("add", "", "Task to be added to the todo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	ls := &todo.List{}

	if err := ls.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// decide what to do based on the provided flags
	switch {
	case *list:
		// list current to-do items
		fmt.Print(ls)
	case *complete > 0:
		// complete the given item
		if err := ls.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// save the new list
		if err := ls.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add != "":
		// add the task
		ls.Add(*add)

		// save the new list
		if err := ls.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "invalid option")
		os.Exit(1)
	}
}
