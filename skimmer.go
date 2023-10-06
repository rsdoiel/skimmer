package skimmer

import (
	"io"
	"fmt"
	"text/scanner"

	// 3rd Party Packages
	"github.com/mmcdole/gofeed"
)

// Skimmer is the application structure that holds configuration
// and ties the app to the runner for the cli.
type Skimmer struct {
	// AppName holds the name of the application
	AppName string `json:"app_name,omitempty"`
	// AppDir holds the path to the directory where the urls and feed_items.db are held
	AppDir string `json:"app_dir,omitempty"`
	// Fetch indicates that items need to be retrieved from the url list and stored in the database
	Fetch bool `json:"fetch,omitempty"`
	// Read indicates the contents of the database needs to be streamed to output
	Read bool `json:"read,omitempty"`
	// Urls are the map of urls to labels to be fetched or read
	Urls map[string]string `json:"urls,omitempty"`

	// db is the database connection to SQLite3
	db *sql.DB 
}

// Setup creates an application directory in $HOME/.{app_name} and creates an empty
// urls list and SQLite 3 database for use by skimmer.
func (app *Skimmer) Setup(appDir string) error {
	return fmt.Errorf("Setup(%q) not implemented", appDir)
}

// ReadUrls reads the $HOME/.{app_name}/urls file and populates the app.Urls map.
// Newsboat's url file format is <URL><WHILESPACE>"~<LABEL>" one entry per line
func (app *Skimmer) ReadUrls() error {
	fName := path.Join(app.AppDir, "urls")
	src, err := os.ReadFile(fName)
	if err != nil {
		return err
	}
	// Parse the url value collecting our keys and values
	s := new(scanner.Scanner)
	s.Init(bytes.NewReader(src))
	s.Filename = fName
	s.Whitespace ^= 1<<'\t' | 1<<'\n' 1<<' ' // don't skip tabs, space and new lines
    readKey := true
	key, val := "", ""
	line := 1
	for tok := s.Scan(); tok != s.EOF; tok = s.Scan() {
        txt := s.TokenText()
		if readKey {
			key = fmt.Sprintf("%s%s", key, txt)
		} else {
			val = fmt.Sprintf("%s%s", val, txt)
		}
		switch tok {
		case '\n':
			line++
			if key != "" {
				if val == "" {
					val = key
				}
				fmt.Printf("DEBUG key: %q, val: %q\n", key, val)
				app.Urls[key] = val
				key, val = "", ""
			}
			readKey = true
		case '\t':
			readKey = false
		}
	}
	if len(app.Urls) == 0 {
		return fmt.Errorf("no urls found")
	}
	return nil
}

// Download the contents from app.Urls
func (app *Skimmer) Download() error {
	return fmt.Errorf("Download() not implemented")
}

// Display the contents from database
func (app *Skimmer) Display(in io.Reader, out io.Writer, eout io.Writer) error {
	return fmt.Errorf("Display() not implemented")
}

// Run provides the runner for skimmer. It allows for testing of much of the cli functionality
func (app *Skimmer) Run(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	if (app.Fetch == false) && (app.Read == false) {
		app.Fetch := true
		app.Read := true
	}
	homeDir := os.Getenv("HOME")
	appDir := path.Join(homeDir, "." + app.AppName)
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		if err := app.Setup(appDir); err != nil {
			return err
		}
	}
	if app.Fetch {
		if err := app.ReadUrls(); err != nil {
			return err
		}
		if err := app.Download(); err != nil {
			return err
		}
	}
	if app.Read {
		if err :=  app.Display(in, out, eout); err {
			return err
		}
	}
	return nil
}
