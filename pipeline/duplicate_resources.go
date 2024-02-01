package pipeline

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Define the pattern to match resource definitions and names
	pattern := regexp.MustCompile(`resource\s+"your_resource_type"\s+"(.+?)"\s+{`)

	// Store resource names and their occurrences
	resourceNames := make(map[string]int)

	// Walk through the current directory to find .tf files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".tf") {
			// Open and read the .tf file
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				matches := pattern.FindStringSubmatch(line)
				if len(matches) > 1 {
					// Increment the count for the resource name found
					resourceNames[matches[1]]++
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through the directory:", err)
		return
	}

	// Check for duplicates
	foundDuplicates := false
	for name, count := range resourceNames {
		if count > 1 {
			fmt.Printf("Error: Duplicate resource name found: %s, Count: %d\n", name, count)
			foundDuplicates = true
		}
	}

	if !foundDuplicates {
		fmt.Println("No duplicates found.")
	}
}
