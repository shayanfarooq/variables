package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	WorkerCount = 10 // Number of concurrent threads
)

// Regex pattern for detecting variables and their values
var (
	variableRegex = regexp.MustCompile(`(?s)\b(?:var|let|const)\s+([a-zA-Z_$][\w$]*)\s*=\s*([^;]*);`)
)

// Extracts variables and their values from JavaScript content
func extractVariables(content []byte, filterVar string) []string {
	var variables []string
	matches := variableRegex.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) > 2 {
			varName := match[1]
			varValue := strings.TrimSpace(match[2])
			if filterVar == "" || filterVar == varName {
				variables = append(variables, fmt.Sprintf("[%s] [%s]", varName, varValue))
			}
		}
	}
	return variables
}

// Analyzes JavaScript content for variables and their values
func analyzeJSContent(url string, content []byte, filterVar string) {
	variables := extractVariables(content, filterVar)

	if len(variables) > 0 {
		fmt.Println(Green + "[URL]" + Reset, Blue+url+Reset)
		fmt.Println(Yellow + "  Variables and Values:" + Reset)
		for _, variable := range variables {
			fmt.Println(Green + "    " + Reset + variable)
		}
	} else {
		fmt.Println(Yellow + "[URL]" + Reset, Blue+url+Reset)
		fmt.Println(Yellow + "  No matching variables found." + Reset)
	}
}

// Worker function to process each URL
func worker(urls <-chan string, wg *sync.WaitGroup, filterVar string) {
	defer wg.Done()
	for url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(Red + "[Error]" + Reset, "Error fetching", url, ":", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println(Yellow + "[Warning]" + Reset, "Non-OK HTTP status for", url, ":", resp.Status)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(Red + "[Error]" + Reset, "Error reading response body from", url, ":", err)
			continue
		}

		analyzeJSContent(url, body, filterVar)
	}
}

// Main function to manage workers and read URLs from stdin
func main() {
	var filterVar string
	flag.StringVar(&filterVar, "var", "", "Filter by variable name")
	flag.Parse()

	if len(flag.Args()) == 0 {
		// Handle URLs from standard input
		if !handleStdinURLs(filterVar) {
			fmt.Println(Red + "Error:" + Reset + " No URLs provided. Please pipe URLs into the program.")
			os.Exit(1)
		}
	} else {
		// Handle URLs from command-line arguments
		urls := flag.Args()
		if len(urls) == 0 {
			fmt.Println(Red + "Error:" + Reset + " No URLs provided.")
			os.Exit(1)
		}
		processURLs(urls, filterVar)
	}
}

// Handle URLs from standard input
func handleStdinURLs(filterVar string) bool {
	var wg sync.WaitGroup
	urls := make(chan string, WorkerCount)
	
	// Start worker pool
	for i := 0; i < WorkerCount; i++ {
		wg.Add(1)
		go worker(urls, &wg, filterVar)
	}
	
	// Read URLs from stdin and send to worker pool
	scanner := bufio.NewScanner(os.Stdin)
	hasURLs := false
	for scanner.Scan() {
		url := scanner.Text()
		if url != "" {
			hasURLs = true
			urls <- url
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Println(Red + "[Error]" + Reset, "Error reading input:", err)
	}
	
	close(urls)
	wg.Wait()
	return hasURLs
}

// Process URLs from command-line arguments
func processURLs(urls []string, filterVar string) {
	var wg sync.WaitGroup
	urlChan := make(chan string, WorkerCount)

	// Start worker pool
	for i := 0; i < WorkerCount; i++ {
		wg.Add(1)
		go worker(urlChan, &wg, filterVar)
	}

	// Send URLs to worker pool
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)
	wg.Wait()
}
