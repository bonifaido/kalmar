package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	colors = []string{
		"\033[0;31m", // Red
		"\033[0;32m", // Green
		"\033[0;33m", // Yellow
		"\033[0;34m", // Blue
		"\033[0;35m", // Magenta
		"\033[0;36m", // Cyan
	}
)

func executeCommand(command string, longestNameLen int, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("bash", "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %s\n", err)
		return
	}

	cmd.Start()

	color := colors[rand.Intn(len(colors))]

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		timestamp := time.Now().Format("15:04:05")
		paddedCommand := fmt.Sprintf("%-*s", longestNameLen, command)
		fmt.Printf("%s %s%s\033[0m | %s\n", timestamp, color, paddedCommand, scanner.Text()) // Timestamp + Random color
	}

	if err := cmd.Wait(); err != nil {
		timestamp := time.Now().Format("15:04:05")
		paddedCommand := fmt.Sprintf("%-*s", longestNameLen, command)
		fmt.Printf("%s %s%s\033[0m | Error: %s\n", timestamp, color, paddedCommand, err) // Timestamp + Random color
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kalmar [command1] [command2] ...")
		return
	}

	commands := os.Args[1:]

	// Find the length of the longest program name
	longestNameLen := 0
	for _, cmd := range commands {
		if len(cmd) > longestNameLen {
			longestNameLen = len(cmd)
		}
	}

	var wg sync.WaitGroup

	for _, cmd := range commands {
		wg.Add(1)
		go executeCommand(cmd, longestNameLen, &wg)
	}

	wg.Wait()
}
