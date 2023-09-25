package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Press 'Enter' to quit.")

	// Disable terminal input buffering and echo
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	defer func() {
		// Restore terminal settings when the program exits
		exec.Command("stty", "-F", "/dev/tty", "cooked").Run()
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	}()

	for {
		char := make([]byte, 3) // Read up to 3 bytes
		_, err := os.Stdin.Read(char)
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			return
		}

		switch string(char) {
		case "\n": // Enter key
			fmt.Println("Quitting...")
			return
		case "\x1b[A": // Up arrow key
			fmt.Println("Up arrow key pressed")
		case "\x1b[B": // Down arrow key
			fmt.Println("Down arrow key pressed")
		case "\x1b[5~": // Page Up key
			fmt.Println("Page Up key pressed")
		case "\x1b[6~": // Page Down key
			fmt.Println("Page Down key pressed")
		case "\x1b[2~": // Insert key
			fmt.Println("Insert key pressed")
		default:
			fmt.Printf("You pressed: %s\n", char)
		}
	}
}
