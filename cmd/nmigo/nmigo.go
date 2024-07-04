package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// delete all files that contain any of the keywords in their name
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

// verify if the file name contains any of the keywords
func containsKeywords(path string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(filepath.Base(path)), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func main() {
	root := "./path/to/scan"                                    // directory to scan
	keywords := []string{"sensitive", "confidential", "secret"} // words to search for in the file names

	if err := deleteSensitiveFiles(root, keywords); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Sensitive files deleted successfully.")
	}
}
