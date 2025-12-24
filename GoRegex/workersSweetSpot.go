package main
import "runtime"
import "fmt"
import "time"
import "rand"


/// this file is used to find the sweet spot for number of workers (number of goroutines) to use in the regex matching worker pool

// Easy benchmark function
func benchmarkRegex(lines []string, workers int) time.Duration {
    start := time.Now()
    // ... run your worker pool ...
    return time.Since(start)
}

func generateTestLines(count int) []string {
    lines := make([]string, count)
    for i := range lines {
        lines[i] = lines[i] = fmt.Sprintf("log entry %d: ERROR user%d@host%d failed auth %d times", i, rand.Intn(1000), rand.Intn(100), rand.Intn(10))
    }
    return lines
}

var testLines = generateTestLines(1000000) // 1 million lines

func main() {
    fmt.Println("CPU cores:", runtime.NumCPU())
    fmt.Println("Goroutine limit:", runtime.GOMAXPROCS(0))
// Test: 1, 2, 4, 8, 16, 32 workers
    
    for _, w := range []int{1, 2, 4, 8, 16, 32} {
        t := benchmarkRegex(testLines, w)
        fmt.Printf("Workers=%2d: %v\n", w, t)
    }
}