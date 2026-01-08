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
go run workersSweetSpot.go
```

*Note : Ensure the searched file is large to get meaningful results*
