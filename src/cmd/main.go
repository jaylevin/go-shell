package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jaylevin/go-shell/src/pkg/manager"
)

var mgr = &manager.Manager{}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		// Read the keyboad input.
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// ErrNoPath is returned when 'cd' was called without a second argument.
var ErrNoPath = errors.New("path required for cd")

func execInput(input string) error {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	// Split the input separate the command and the arguments.
	args := strings.Split(input, " ")

	// Check for built-in commands.
	switch args[0] {
	case "cd":
		// 'cd' to home with empty path not yet supported.
		if len(args) < 2 {
			return os.Chdir(os.Getenv("HOME"))
		}
		// Change the directory and return the error.
		return os.Chdir(args[1])

	case "in":
		mgr.Init()
		return nil

	case "cr":
		if len(args) < 2 {
			log.Print("A priority level must be specified, either 0, 1, or 2.")
		}
		priority, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Encountered error while parsing string to integer: %s", err.Error())
		}
		if err := mgr.Create(priority); err != nil {
			log.Print(err.Error())
		}
		return nil

	case "dr":
		if len(args) < 2 {
			log.Print("A process index number must be specified.")
		}
		index, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Encountered error while parsing string to integer: %s", err.Error())
		}
		numDeleted, err := mgr.Destroy(index)
		if err != nil {
			log.Print(err.Error())
		}
		log.Printf("%v processes destroyed", numDeleted)
		return nil

	case "exit":
		os.Exit(0)
	}

	// Prepare the command to execute.
	cmd := exec.Command(args[0], args[1:]...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
