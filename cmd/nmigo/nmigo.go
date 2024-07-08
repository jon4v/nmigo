package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// containsKeywords checks if a file name contains any of the sensitive keywords.
func containsKeywords(path string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(filepath.Base(path)), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// deleteFile attempts to delete the specified file and logs the result.
func deleteFile(path string, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	fmt.Println("Deleting:", path)
	if err := os.Remove(path); err != nil {
		errChan <- err
	}
}

// deleteSensitiveFiles searches the root directory for files containing sensitive keywords and deletes them using concurrency.
func deleteSensitiveFiles(root string, keywords []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	doneChan := make(chan bool)

	go func() {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				errChan <- err
				return err
			}
			if !info.IsDir() && containsKeywords(path, keywords) {
				wg.Add(1)
				go deleteFile(path, &wg, errChan)
			}
			return nil
		})
		if err != nil {
			return
		}
		wg.Wait()
		close(doneChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

func main() {
	// Define flags for root directory and keywords.
	root := flag.String("root", ".", "Root directory to scan")
	keywordList := flag.String("keywords", "", "Comma-separated list of keywords to search for")

	// Parse command-line flags.
	flag.Parse()

	// Validate provided keywords.
	if *keywordList == "" {
		fmt.Println("Error: You must provide a comma-separated list of keywords using the -keywords flag.")
		return
	}

	// Convert comma-separated keywords into a slice.
	keywords := strings.Split(*keywordList, ",")

	// Delete sensitive files using the provided root directory and keywords.
	if err := deleteSensitiveFiles(*root, keywords); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Sensitive files deleted successfully.")
	}
}
