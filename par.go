// TODO:
//  - flag: proc timeout
//  - flag: output to separate files (named with a pattern or based on line ran)
//  - flag: output to a single file
//  - flag: number of workers
//  - flag: include/exclude stderr
//  - flag: crash/bypass on error
//  - flag: buffered/unbuffered output
//  - tests :(

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func main() {
	template := parseCommand()

	lines := make(chan string)
	workers := &sync.WaitGroup{}
	numWorkers := maxInt(runtime.NumCPU(), 4)

	// Spin up the workers
	fmt.Printf("Running with %d workers\n", numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers.Add(1)
		go work(workers, template, lines)
	}

	// Read standard input (while available) and queue it up on the lines channel
	readFromStdin(lines)
	close(lines)

	// Wait for the workers to finish.
	workers.Wait()
}

func parseCommand() string {
	flag.Parse()

	// Template of command to be run for each line of stdin
	if flag.NArg() == 0 {
		log.Fatal("Sorry, you must supply a command to run as the first argument.")
	}
	return flag.Arg(0)
}

func maxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func readFromStdin(lines chan<- string) {
	reader := bufio.NewReader(os.Stdin)

	line, err := reader.ReadString('\n')
	for err == nil {
		lines <- line
		line, err = reader.ReadString('\n')
	}
}

func work(wg *sync.WaitGroup, template string, lines <-chan string) {
	defer wg.Done()

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for line := range lines {
		line = strings.TrimSpace(line)
		command := strings.Replace(template, "%line", line, -1)
		writer.WriteString(fmt.Sprintf("[%s] starting\n", command))

		split_command := strings.Split(command, " ")
		cmd := exec.Command(split_command[0], split_command[1:]...)
		stdout, err := cmd.StdoutPipe()
		if err != nil { // error getting output, no need to crash the whole app
			log.Println(err)
		}
		if err := cmd.Start(); err != nil {
			log.Println(err)
			continue
		}

		reader := bufio.NewReader(stdout)

		outputLine, err := reader.ReadString('\n')
		for err == nil {
			writer.WriteString(fmt.Sprintf("[%s] %s", command, outputLine))
			outputLine, err = reader.ReadString('\n')
		}
		stdout.Close()

		if err := cmd.Wait(); err != nil {
			writer.WriteString(fmt.Sprintf("[%s] err: %v\n", command, err))
		}
	}
}
