package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bube054/todo"
)

var todoFileName = ".todo.json"

func main() {
	// flag.Usage = func() {
	// 	fmt.Fprintf(flag.CommandLine.Output(),
	// 		"%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
	// 	fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
	// 	fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
	// 	flag.PrintDefaults()
	// }

	// parsing cmd line flags
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// using flags
	switch {
	case *list:
		fmt.Print(l)
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		t, err := GetTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}

	// using os.Args
	// switch {
	// case len(os.Args) == 1:
	// 	for _, item := range *l {
	// 		fmt.Println(item.Task)
	// 	}
	// default:
	// 	item := strings.Join(os.Args[1:], " ")

	// 	l.Add(item)

	// 	if err := l.Save(todoFileName); err != nil {
	// 		fmt.Fprintln(os.Stderr, err)
	// 		os.Exit(1)
	// 	}
	// }
}

func GetTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)

	s.Scan()

	if err := s.Err(); err != nil {
		return "", nil
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return s.Text(), nil
}
