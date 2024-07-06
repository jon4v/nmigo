package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// search root directory for files containing sensitive keywords and delete them
func deleteSensitiveFiles(root string, keywords []string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && containsKeywords(path, keywords) {
			fmt.Println("Deleting:", path)
			return os.Remove(path)
		}
		return nil
	})
}

// verify if a file contains any of the sensitive keywords
func containsKeywords(path string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(filepath.Base(path)), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func main() {
	// define flags
	root := flag.String("root", ".", "Root directory to scan")
	keywordList := flag.String("keywords", "", "Comma-separated list of keywords to search for")

	// parse flags
	flag.Parse()

	// validate flags
	if *keywordList == "" {
		fmt.Println("Error: You must provide a comma-separated list of keywords using the -keywords flag.")
		return
	}

	// Convert the comma-separated list of keywords into a slice
	keywords := strings.Split(*keywordList, ",")

	// delete sensitive files
	if err := deleteSensitiveFiles(*root, keywords); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Sensitive files deleted successfully.")
	}
}
