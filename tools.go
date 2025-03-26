package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Tell the program what to do based on the results of a --dry-run
func direct(answer, nav string) {
	if strings.ToLower(answer) == "y" {
		proceed(nav)
	} else {
		os.Exit(0)
	}
}

// Get user confirmation after completion of a --dry-run
func confirm() string {
	answer := solicit("Does this output seem acceptable? Shall we continue without the --dry-run flag? (y/n) ")
	return answer
}

// Execute the functions without a --dry-run condition
func proceed(action string) {
	switch action {
	case "ac":
		copyAssets()
	case "hf":
		fixProtocol()
	}
}

// Get user input via screen prompt
func solicit(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	response, _ := reader.ReadString('\n')
	return strings.TrimSpace(response)
}

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
