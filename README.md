# GoRegex

GoRegex is a high-performance, concurrent CLI tool written in Go. It allows you to search for Regular Expression (Regex) patterns inside large text files using a worker pool architecture.

By leveraging Go's concurrency primitives (`goroutines` and `channels`), this tool processes files significantly faster than traditional single-threaded scripts.

## 🚀 Features

* **Concurrent Processing:** Uses a worker pool (default: 8 workers) to scan lines in parallel.
* **Memory Efficient:** Streams file content line-by-line instead of loading the whole file into RAM.
* **CLI Flags:** Fully customizable file paths and regex patterns via command-line arguments.
* **Performance Tuning:** Includes a utility (`workersSweetSpot.go`) to test different worker counts.

## 📂 Project Structure

```text
GoRegex/
├── data/                 # Folder to store your input text files
├── main.go               # The main CLI application
├── workersSweetSpot.go   # Script to benchmark and find optimal worker count
├── go.mod                # Go module definition
└── README.md             # This documentation
```

## 📊 Performance Report

Benchmarking with `workersSweetSpot.go` on 1M synthetic lines (20 CPU cores) reveals workload-dependent scaling:

* **Light tasks** (simple patterns like `[A-Za-z]`, short lines): More workers slow performance due to goroutine overhead and channel contention. Sequential (1 worker) is fastest (~320ms).
* **Heavy tasks** (complex patterns like email regex `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, longer processing): Parallelism shines. 4 workers halved time (~806ms vs. 1.75s for 1 worker), but beyond core count, overhead dominates.

**Key Insights:**
* Optimal workers ≈ CPU cores (4-8 here) for CPU-bound regex.
* Increase channel buffers (e.g., 1000) for high worker counts.
* Code is efficient; speedup depends on pattern complexity and data size.

## Getting Started

### Prerequisites

* Go (version 1.18 or higher recommended)

### Installation

* Clone the repository :

```
git clone [https://github.com/yourusername/GoRegex.git](https://github.com/yourusername/GoRegex.git)
cd GoRegex
```

* Initialize the module :

```
go mod tidy
```

## Usage

### The main search tool (`main.go`)

You can run the tool directly using `go.run`.
Basic Usage (uses defaults):

* Default File: `data/LongLog.txt`
* Default Pattern: `[A-Za-z]`

```
go run main.go

**Custom search** : Use flags to specify your own file and pattern
- `--file | -f` : path to text file
- `--pattern | -p` : the regex string to search for
- `--workers | -w : specify the number of workers`
```

**Example : finding emails**

```
go run main.go --file=data/my_log.txt --pattern="[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}"
```

**Example : finding errors**

```
go run main.go --pattern="ERROR"
```

### The benchmark script (`workersSweetSpot.go`)

This script benchmarks different numbers of workers to find the optimal count for your regex pattern and data size.

Basic Usage (default suite: 1,2,4,8,16,32 workers):

```
go run workersSweetSpot.go
```

Custom Usage:

* `--file | -f` : path to file to benchmark (optional, defaults to 1M generated lines)
* `--workers | -w` : test a specific number of workers (0 for default suite)
* `--pattern | -p` : regex pattern to benchmark (default: email regex `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

Example: Test 4 workers with a custom pattern on a real file

```
go run workersSweetSpot.go -f data/LongLog.txt -w 4 -p "ERROR"
```

*Note : Ensure the searched file is large enough, the lines long enough and the pattern complex enough to get meaningful results*
