package main

import (
	"fmt"
	"regexp"
	"sync"
	"bufio"
	"os"
	"flag"
)

// What does a piece of work look like?
type Job struct {
	LineNumber int
	Text       string
}

// What does a "Found Match" look like?
type Result struct {
	LineNumber int
	Match      string
}

func main() {
	// Define your regex pattern (e.g., finding emails, specific words)
	// pattern := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`) 

	// --- Define flags ---
	filePth := flag.String("file", "data/LongLog.txt", "File to process")
	patternStr := flag.String("pattern", `[a-Z]`, "Regex pattern to search for")

	flag.Parse()
	pattern := regexp.MustCompile(*patternStr)

	// --- Use flags ---
	fmt.Println("Processing file:", *filePth)
	fmt.Println("Using pattern:", *patternStr)
	
	// Create the "Conveyor Belts" (Channels)
	// jobs: Where we send lines of text to be processed
	jobs := make(chan Job, 100)
	// results: Where workers drop their findings
	results := make(chan Result, 100)
	
	// A "WaitGroup" is like a clipboard to track attendance.
	// It ensures we don't quit until all 8 interns are done.
	var wg sync.WaitGroup

	// --- 3. RECRUITING THE INTERNS (Start Workers) ---
	numWorkers := 8
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // Add one intern to the attendance sheet
		go worker(i, pattern, jobs, results, &wg)
	}

	// --- 4. THE MANAGER (Send Work) ---
	// In a separate goroutine, we start reading the file and filling the 'jobs' channel.
	// We do this in the background so the main function can focus on gathering results.
	go func() {
		//open file
		file, err := os.Open(*filePth)
		if err != nil {
			fmt.Println("Error opening file:", err)
			close(jobs)
			return
		}
		// Make sure to close the file when done
		defer file.Close()

		// create a scanner to read the file line by line
		scanner := bufio.NewScanner(file)
		lineNumber := 1
		for scanner.Scan() {
			line := scanner.Text()
			// Send each line as a job to the workers
			jobs <- Job{
				LineNumber: lineNumber,
				Text:       line,
			}
			lineNumber++
		}
		// Check for errors during scanning
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}
		// Close the jobs channel to signal no more work
		close(jobs) // "No more work! Go home after you finish your current page."
	}()

	// --- 5. THE COLLECTOR (Receive Results) ---
	// While workers are working, we wait for them to finish in a separate routine
	// so we can close the results channel when they are all done.
	go func() {
		wg.Wait()      // Wait for all 8 interns to sign out
		close(results) // Close the results box so we stop waiting for more
	}()

	// --- 6. REPORTING ---
	// Read from the results channel until it's closed
	for res := range results {
		fmt.Printf("Found match on line %d: %s\n", res.LineNumber, res.Match)
	}
}

// --- 7. THE INTERN (Worker Function) ---
// This function runs 8 times simultaneously.
func worker(id int, pattern *regexp.Regexp, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done() // When this function finishes, cross name off attendance sheet

	// Keep grabbing jobs until the channel is closed
	for job := range jobs {
		if pattern.MatchString(job.Text) {
			// Found one! Send it to results.
			results <- Result{
				LineNumber: job.LineNumber,
				Match:      job.Text, // Or the specific regex match
			}
		}
	}
}