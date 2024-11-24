package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	// In this case,
	// we’re appending the suffix .exe to the binary name so Go finds the
	// executable during tests.
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// The -o flag is used to specify the name of the output binary file.
	// In this case, the output file is called todo.
	// If you don’t specify the -o flag, the output file will be called main.
	build := exec.Command("go", "build", "-o", binName)

	// we are running the binary build process.
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests....")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	// task to be added
	task := "test task number 1"

	// get the current parent directory
	dir, err := os.Getwd()

	if err != nil {
		t.Fatal(err)
	}

	// joins the directory with the binary name
	cmdPath := filepath.Join(dir, binName)

	// this test runs the default block: adds a new task.
	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")

		cmdStdIn, err := cmd.StdinPipe()

		if err != nil {
			t.Fatal(err)
		}

		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// this test runs the reads from the stdout & stderr then compares it the test inserted
	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")

		out, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf( "  1: %s\n  2: %s\n", task, task2)

		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
