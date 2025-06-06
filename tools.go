package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Run standard terminal commands
func execute(variation, task string, args ...string) []byte {
	lpath, err := exec.LookPath(task)
	inspect(err)
	osCmd := exec.Command(lpath, args...)
	switch variation {
	case "-e":
		// Execute straight away
		exec.Command(lpath, args...).CombinedOutput()
	case "-c":
		// Capture and return the output as a byte
		both, err := osCmd.CombinedOutput()
		inspect(err)
		err = nil
		return both
	case "-v":
		// Execute verbosely
		osCmd.Stdout = os.Stdout
		osCmd.Stderr = os.Stderr
		err := osCmd.Run()
		inspect(err)
	}
	return nil
}

// Navigate to specific directory
func changeDIR(goal string) {
	os.Chdir(goal)
}

// Read any file and return the contents as a byte variable
func read(file string) []byte {
	mission, err := os.Open(file)
	inspect(err)
	outcome, err := io.ReadAll(mission)
	inspect(err)
	defer mission.Close()
	return outcome
}

// Write a passed variable to a named file
func document(name string, d []byte) {
	inspect(os.WriteFile(name, d, 0666))
}

// Check if a file is present in the supplied path
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Check for errors, print the result if found
func inspect(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Provide and highlight an informational message
func banner(message string) {
	fmt.Println(yellowFG)
	fmt.Println("**", automatic, message, yellowFG, "**", automatic)
}
