package skimmer

import (
	"io"
	"fmt"
	"os"
	"bufio"
	"regexp"
)

// URLLocation holds the file path and line number of a URL occurrence
type URLLocation struct {
	filePath   string
	lineNumber int
}

func parseFile(filePath string) (map[string][]URLLocation, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	urlMap := make(map[string][]URLLocation)
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmedLine := line

		// Skip empty lines and comments
		if len(trimmedLine) == 0 || trimmedLine[0] == '#' {
			continue
		}

		// Regex to match URL and optional label
		re := regexp.MustCompile(`^(\S+)\s*(?:"([^"]*)")?$`)
		matches := re.FindStringSubmatch(trimmedLine)

		if len(matches) >= 2 {
			url := matches[1]
			location := URLLocation{filePath: filePath, lineNumber: lineNumber}
			urlMap[url] = append(urlMap[url], location)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urlMap, nil
}

type SkimDups struct {
}

func NewSkimDups() *SkimDups{
	return new (SkimDups)
}

func (app *SkimDups) Run(out io.Writer, eout io.Writer, args []string)  error {

	filePaths := args[:]
	globalUrlMap := make(map[string][]URLLocation)

	for _, filePath := range filePaths {
		fileUrlMap, err := parseFile(filePath)
		if err != nil {
			fmt.Fprintf(eout, "Error reading file %s: %v\n", filePath, err)
			continue
		}

		for url, locations := range fileUrlMap {
			globalUrlMap[url] = append(globalUrlMap[url], locations...)
		}
	}

	hasDuplicates := false
	for url, locations := range globalUrlMap {
		if len(locations) > 1 {
			hasDuplicates = true
			fmt.Fprintf(out, "Duplicate URL found: %s\n", url)
			for _, loc := range locations {
				fmt.Fprintf(out, "  Located at file: %s, line: %d\n", loc.filePath, loc.lineNumber)
			}
		}
	}

	if !hasDuplicates {
		fmt.Fprintln(out, "No duplicate URLs found.")
	}
	return nil
}