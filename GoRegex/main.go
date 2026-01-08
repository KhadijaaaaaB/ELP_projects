package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sync"
)

// Job request channel format
type Job struct {
	LineNumber int
	Text       string
}

// Job result channel format
type Result struct {
	LineNumber int
	Match      string
}

func main() {
	// --- Process flags ---
	file, patternStr, numWorkers := parseFlags()
	pattern := regexp.MustCompile(patternStr)

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- Use flags ---
	fmt.Println("Processing file:", file)
	fmt.Println("Using pattern:", patternStr)

	// Create the channels
	// jobs: Where we send lines of text to be processed
	jobs := make(chan Job, 100)
	// results: Where workers drop their findings
	results := make(chan Result, 100)

	// A WaitGroup ensures we don't quit until all interns are done.
	var wg sync.WaitGroup

	// --- Start Workers ---
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, pattern, jobs, results, &wg, ctx)
	}

	// --- Send Work ---
	// In a separate goroutine, we start reading the file and filling the 'jobs' channel.
	// We do this in the background so the main function can focus on gathering results.
	go func() {
		// open file
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening file:", err)
			cancel()
			close(jobs)
			return
		}
		// Make sure to close the file when done
		defer f.Close()

		// create a scanner to read the file line by line
		scanner := bufio.NewScanner(f)
		lineNumber := 1
		for scanner.Scan() {
			line := scanner.Text()
			// Send each line as a job to the workers
			select {
			case jobs <- Job{
				LineNumber: lineNumber,
				Text:       line,
			}:
			case <-ctx.Done():
				close(jobs)
				return
			}
			lineNumber++
		}

		// Check for errors during scanning
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
			cancel()
		}

		// Close the jobs channel to signal no more work
		close(jobs)
	}()

	// --- Receive Results ---
	// We can close the results channel when the workers are all done.
	go func() {
		wg.Wait()
		close(results)
	}()

	// --- Gathering results ---
	// Read from the results channel until it's closed by the previous goroutine
	for res := range results {
		fmt.Printf("Found match on line %d: %s\n", res.LineNumber, res.Match)
	}

	// Check if operation was canceled
	if ctx.Err() != nil {
		fmt.Println("Operation was canceled due to an error.")
	}
}

// --- Worker Function ---
// This function runs n times simultaneously.
func worker(id int, pattern *regexp.Regexp, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done() // When this function finishes, signal it to the WaitGroup

	// Keep grabbing jobs until the channel is closed or context canceled
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			if pattern.MatchString(job.Text) {
				select {
				case results <- Result{
					LineNumber: job.LineNumber,
					Match:      job.Text,
				}:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// parseFlags handles all command-line flag parsing and returns the parsed values
func parseFlags() (string, string, int) {
	var file string
	flag.StringVar(&file, "file", "data/LongLog.txt", "File to process")
	flag.StringVar(&file, "f", "data/LongLog.txt", "File to process")

	var patternStr string
	flag.StringVar(&patternStr, "pattern", `[A-Za-z]`, "Regex pattern to search for")
	flag.StringVar(&patternStr, "p", `[A-Za-z]`, "Regex pattern to search for")

	var numWorkers int
	flag.IntVar(&numWorkers, "workers", 8, "Number of worker goroutines")
	flag.IntVar(&numWorkers, "w", 8, "Number of worker goroutines")

	// Prints the custom help message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Perform parallel regex matching on a file using goroutines.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -f, --file string\n\tFile to process (default \"data/LongLog.txt\")\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -p, --pattern string\n\tRegex pattern to search for (default \"[A-Za-z]\")\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -w, --workers int\n\tNumber of worker goroutines (default 8)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nExample:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -f data/LongLog.txt -p ERROR -w 4\n", os.Args[0])
	}

	flag.Parse()

	return file, patternStr, numWorkers
}
