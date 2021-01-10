package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/giancarlosio/gorainbow"
)

func execInput(input string) error {
	// Remove the new line character
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	// Split the input to separate the command and the arguments.
	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return errors.New("path required")
		}
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	var cmd *exec.Cmd

	// Prepare the command to execute
	// we can't use the this like `exec.Command(input)` because it will throw an error like this:
	// "the `input` executalbe doesn't exists".
	// And thats because the first param on exec.Command should be an executable
	var inputCommand []string
	var name string
	switch runtime.GOOS {
	case "windows":
		inputCommand = append([]string{"/C"}, args...)
		name = "pwsh"
	default:
		inputCommand = append([]string{"-c"}, args...)
		name = "bash"
	}

	cmd = exec.Command(name, inputCommand...)
	// Set the correct output device
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	user, err := user.Current()
	if err != nil {
		printError(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		printError(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		printError(err)
	}

	shellText := fmt.Sprintf("GOSHELL 🐢: %s %s %s> ", cwd, hostname, user.Name)
	shellText = gorainbow.Rainbow(shellText)

	for {
		fmt.Print(shellText)
		input, err := reader.ReadString('\n')
		if err != nil {
			printError(err)
		}
		if err = execInput(input); err != nil {
			printError(err)
		}

	}
}
