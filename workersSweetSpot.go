package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"sync"
	"time"
)

// this file is used to find the sweet spot for number of workers (number of goroutines) to use in the regex matching worker pool

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

// benchWorker function for benchmark processing jobs
func benchWorker(id int, pattern *regexp.Regexp, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		if pattern.MatchString(job.Text) {
			results <- Result{
				LineNumber: job.LineNumber,
				Match:      job.Text,
			}
		}
	}
}

// Benchmark function that runs the worker pool
func benchmarkRegex(pattern *regexp.Regexp, lines []string, workers int) time.Duration {
	start := time.Now()

	jobs := make(chan Job, 1000)
	results := make(chan Result, 1000)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go benchWorker(i, pattern, jobs, results, &wg)
	}

	// Send jobs
	go func() {
		for i, line := range lines {
			jobs <- Job{LineNumber: i + 1, Text: line}
		}
		close(jobs)
	}()

	// Wait for workers and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Drain results
	for range results {
		// discard results, just measure time
	}

	return time.Since(start)
}

func generateTestLines(count int) []string {
	lines := make([]string, count)
	for i := range lines {
		lines[i] = fmt.Sprintf("log entry %d: ERROR user%d@host%d.com failed auth %d times", i, rand.Intn(1000), rand.Intn(100), rand.Intn(10))
	}
	return lines
}

var testLines = generateTestLines(1000000) // 1 million lines

func loadLines(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	fmt.Println("CPU cores:", runtime.NumCPU())
	fmt.Println("Goroutine limit:", runtime.GOMAXPROCS(0))

	// Parse flags
	var file string
	flag.StringVar(&file, "file", "", "File to read lines from (if not set, uses generated test lines)")
	flag.StringVar(&file, "f", "", "File to read lines from (if not set, uses generated test lines)")

	var patternStr string
	flag.StringVar(&patternStr, "pattern", `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, "Regex pattern to search for")
	flag.StringVar(&patternStr, "p", `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, "Regex pattern to search for")

	var workers int
	flag.IntVar(&workers, "workers", 0, "Number of workers to benchmark (0 for default test suite: 1,2,4,8,16,32)")
	flag.IntVar(&workers, "w", 0, "Number of workers to benchmark (0 for default test suite: 1,2,4,8,16,32)")
	flag.Parse()

	pattern := regexp.MustCompile(patternStr)

	// Load lines
	var lines []string
	if file != "" {
		var err error
		lines, err = loadLines(file)
		if err != nil {
			fmt.Println("Error loading file:", err)
			os.Exit(1)
		}
		fmt.Printf("Loaded %d lines from %s\n", len(lines), file)
	} else {
		lines = testLines
		fmt.Printf("Using %d generated test lines\n", len(lines))
	}

	if workers == 0 {
		// Default test suite
		for _, w := range []int{1, 2, 4, 8, 16, 32} {
			t := benchmarkRegex(pattern, testLines, w)
			fmt.Printf("Workers=%2d: %v\n", w, t)
		}
	} else {
		// Benchmark the specified number of workers
		t := benchmarkRegex(pattern, testLines, workers)
		fmt.Printf("Workers=%2d: %v\n", workers, t)
	}
}
