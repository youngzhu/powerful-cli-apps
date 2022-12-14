package main

import (
	"fmt"
	"os"
	"strings"
	"todo"
)

const todoFileName = ".todo.json"

func main() {
	ls := &todo.List{}

	if err := ls.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// decide what to do based on the number of arguments provided
	switch {
	// for no extra arguments, print the list
	case len(os.Args) == 1:
		for _, item := range *ls {
			fmt.Println(item.Task)
		}
	// concatenate all provided arguments with a space and
	// add to the list as an item
	default:
		item := strings.Join(os.Args[1:], " ")
		ls.Add(item)

		if err := ls.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
