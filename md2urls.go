package skimmer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

type MdToUrls struct {

}

func NewMdToUrls() *MdToUrls {
	return new (MdToUrls)
}

func (app *MdToUrls) Run(out io.Writer, eout io.Writer, args []string) error {
	for _, filePath := range args {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("Error opening file: %v\n", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		re := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)

		for scanner.Scan() {
			line := scanner.Text()
			matches := re.FindAllStringSubmatch(line, -1)

			for _, match := range matches {
				linkText := match[1]
				url := match[2]
				fmt.Fprintf(out, "%s \"~%s\"\n", url, linkText)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("Error reading file: %v\n", err)
		}
	}
	return nil
}