package skimmer

	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	// 3rd Party Packages
	_ "github.com/glebarez/go-sqlite"
	"github.com/mmcdole/gofeed"
)

const (
	SkimmerURLs   = "skimmer.urls"
	SkimmerScheme = "skimmer.scheme"
	SkimmerDB     = "skimmer.db"
)

// ParseURLList takes a filename and byte slice source, parses the contents
// returning a map of urls to labels and an error value.
func ParseURLList(fName string, src []byte) (map[string]string, error) {
	urls := map[string]string{}
	// Parse the url value collecting our keys and values
	s := bufio.NewScanner(bytes.NewBuffer(src))
	key, val := "", ""
	line := 1
	for s.Scan() {
		txt := strings.TrimSpace(s.Text())
		if !strings.HasPrefix(txt, "#") {
			parts := strings.SplitN(txt, " ", 2)
			switch len(parts) {
			case 1:
				key, val = parts[0], parts[0]
			case 2:
				key, val = parts[0], parts[1]
			}
			urls[key] = strings.Trim(val, `"~`)
			key, val = "", ""
		}
		line++
	}
	return urls, nil
}

// Skimmer is the application structure that holds configuration
// and ties the app to the runner for the cli.
type Skimmer struct {
	// AppName holds the name of the application
	AppName string `json:"app_name,omitempty"`
	// AppDir holds the path to the directory where the urls and feed_items.db are held
	AppDir string `json:"app_dir,omitempty"`
	// Fetch indicates that items need to be retrieved from the url list and stored in the database
	Fetch bool `json:"fetch,omitempty"`
	// Display indicates the contents of the database needs to be streamed to output
	Display bool `json:"read,omitempty"`
	// Urls are the map of urls to labels to be fetched or read
	Urls map[string]string `json:"urls,omitempty"`

	// Limit contrains the number of items shown
	Limit int `json:"limit,omitempty"`

	// Prune contains the date to use to prune the database.
	Prune string `json:"prune,omitempty"`

	// Map in some private data to keep things easy to work with.
	out  io.Writer
	eout io.Writer
}

func NewSkimmer(out io.Writer, eout io.Writer, appName string) (*Skimmer, error) {
	app := new(Skimmer)
	app.out = out
	app.eout = eout
	app.AppName = appName
	return app, nil
}

// Setup checks to see if anything needs to be setup (or fixed) for skimmer to run.
func (app *Skimmer) Setup(appDir string) error {
	// Check if we have an appDir
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		//fmt.Fprintf(app.eout, "Creating %s\n", appDir)
		if err := os.MkdirAll(appDir, 0750); err != nil {
			return err
		}
	}
	// Check to see if we have a url list, if not create a default one.
	fName := path.Join(appDir, SkimmerURLs)
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		//fmt.Fprintf(app.eout, "Creating %s\n", fName)
		src := []byte(`# This is an example url list. 
https://laist.com/index.atom "~The LAist"
`)
		// Create a sample urls file
		if err := os.WriteFile(fName, src, 0660); err != nil {
			return err
		}
	}
	// Check if we have a SQL schema file
	fName = path.Join(appDir, SkimmerScheme)
	stmt := fmt.Sprintf(SQLCreateTables, app.AppName, time.Now().Format("2006-01-02"))
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		//fmt.Fprintf(app.eout, "Creating %s\n", fName)
		if err := os.WriteFile(fName, []byte(stmt), 0660); err != nil {
			return err
		}
	}
	fName = path.Join(appDir, SkimmerDB)
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		dsn := fName
		//fmt.Printf("SQLite 3 dsn -> %q\n", dsn)
		db, err := sql.Open("sqlite", dsn)
		if err != nil {
			return err
		}
		defer db.Close()
		if db == nil {
			return fmt.Errorf("%s opened and returned nil", fName)
		}
		_, err = db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	// If we got this far then the appDir is the one we want to remember.
	app.AppDir = appDir
	return nil
}

// ReadUrls reads the $HOME/.skimmer/urls file and populates the app.Urls map.
// Newsboat's url file format is `<URL><SPACE>"~<LABEL>"` one entry per line
// The hash mark, "#" at the start of the line indicates a comment line.
func (app *Skimmer) ReadUrls() error {
	fName := path.Join(app.AppDir, SkimmerURLs)
	src, err := os.ReadFile(fName)
	if err != nil {
		return err
	}
	app.Urls, err = ParseURLList(fName, src)
	if err != nil {
		return err
	}
	if len(app.Urls) == 0 {
		return fmt.Errorf("no urls found")
	}
	return nil
}

// webget retrieves a feed and parses it.
// Uses mmcdole's gofeed, see docs at https://pkg.go.dev/github.com/mmcdole/gofeed
func webget(url string) (*gofeed.Feed, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", res.Status)
	}
	fp := gofeed.NewParser()
	feed, err := fp.Parse(res.Body)
	if err != nil {
		return nil, fmt.Errorf("feed error for %q, %s", url, err)
	}
	return feed, nil
}

func saveItem(db *sql.DB, feedLabel string, item *gofeed.Item) error {
	var (
		published string
		updated   string
	)
	if item.UpdatedParsed != nil {
		updated = item.UpdatedParsed.Format("2006-01-02 15:04:05")
	}
	if item.PublishedParsed != nil {
		published = item.PublishedParsed.Format("2006-01-02 15:04:05")
	}
	stmt := SQLUpdateItem
	_, err := db.Exec(stmt,
		item.Link, item.Title,
		item.Description, updated, published,
		feedLabel)
	if err != nil {
		return err
	}
	return nil
}

// Download the contents from app.Urls
func (app *Skimmer) Download() error {
	eCnt := 0
	dsn := path.Join(app.AppDir, SkimmerDB)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	for k, v := range app.Urls {
		feed, err := webget(k)
		if err != nil {
			eCnt++
			fmt.Fprintf(app.eout, "failed to get %q, %s\n", k, err)
			continue
		}
		// Setup a progress output
		t0 := time.Now()
		rptTime := time.Now()
		reportProgress := false
		tot := feed.Len()
		// FIXME: the logger should really be writing to app.out or app.eout
		log.Printf("processing %d items from %s\n", tot, v)
		i := 0
		for _, item := range feed.Items {
			// Add items from feed to database table
			if err := saveItem(db, v, item); err != nil {
				return err
			}
			if rptTime, reportProgress = CheckWaitInterval(rptTime, (20 * time.Second)); reportProgress {
				// FIXME: the logger should really be writing to app.out or app.eout
				log.Printf("(%d/%d) %s", i, tot, ProgressETA(t0, i, tot))
			}
			i++
		}
		// FIXME: the logger should really be writing to app.out or app.eout
		log.Printf("processed %d/%d from %s\n", i, tot, v)

	}
	if eCnt > 0 {
		return fmt.Errorf("%d errors encounter downloading feeds", eCnt)
	}
	return nil
}

// ItemCount returns the total number items in the database.
func (app *Skimmer) ItemCount() (int, error) {
	dsn := path.Join(app.AppDir, SkimmerDB)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return -1, err
	}
	defer db.Close()
	stmt := SQLItemCount
	rows, err := db.Query(stmt)
	cnt := 0
	for rows.Next() {
		if err := rows.Scan(&cnt); err != nil {
			return -1, err
		}
	}
	return cnt, nil
}

// PruneItems takes a timestamp and performs a row delete on the table
// for items that are older than the timestamp.
func (app *Skimmer) PruneItems(dt time.Time) error {
	dsn := path.Join(app.AppDir, SkimmerDB)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt := SQLPruneItems
	timestamp := dt.Format("2006-01-02 15:04:05")
	_, err = db.Exec(stmt, timestamp, timestamp, timestamp)
	if err != nil {
		return err
	}
	return err	
}

// Display the contents from database
func (app *Skimmer) Write(out io.Writer) error {
	dsn := path.Join(app.AppDir, SkimmerDB)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt := SQLDisplayItems
	if app.Limit > 0 {
		stmt = fmt.Sprintf("%s LIMIT %d", stmt, app.Limit)
	}
	rows, err := db.Query(stmt)
	for rows.Next() {
		var (
			link        string
			title       string
			description string
			updated     string
			published   string
			label       string
		)
		if err := rows.Scan(&link, &title, &description, &updated, &published, &label); err != nil {
			fmt.Fprint(app.eout, err)
			continue
		}
		pressTime := published
		if updated != "" {
			pressTime = updated
		}
		if title == "" {
			fmt.Fprintf(app.out, `
--

## %s @%s

    <%s>

%s 
`, pressTime, label, link, description)
		} else {
			fmt.Fprintf(app.out, `
--

## %s %s

    <%s>

%s 

`, pressTime, title, link, description)
		}
	}
	return nil
}

// Run provides the runner for skimmer. It allows for testing of much of the cli functionality
func (app *Skimmer) Run(out io.Writer, eout io.Writer, args []string) error {
	app.out = out
	app.eout = eout

	// Find our "home" directory.
	homeDir := os.Getenv("HOME")
	// Find out application directory
	appDir := path.Join(homeDir, "."+app.AppName)
	// See if we need to setup things.
	if err := app.Setup(appDir); err != nil {
		return err
	}
	// See what the user wants to do.
	if (app.Fetch == false) && (app.Display == false) {
		cnt, err := app.ItemCount()
		if err != nil {
			fmt.Fprintf(app.eout, "fail to count items, %s", err)
		}
		fmt.Fprintf(app.out, "\n%d items available to read\n", cnt) 
		if cnt == 0 {
			fmt.Fprintf(app.out, "Try %s -fetch to populate the feed database", app.AppName)
		}
		app.Display = true
	}
	if app.Fetch {
		if err := app.ReadUrls(); err != nil {
			return err
		}
		if err := app.Download(); err != nil {
			return err
		}
	}
	if app.Prune != "" {
		fmtStr := "2006-01-02"
		switch {
			case app.Prune == "today":
				fmtStr = "2006-01-02"
				app.Prune = time.Now().Format(fmtStr)
			case app.Prune == "now":
				fmtStr = "2006-01-02 15:04:05"
				app.Prune = time.Now().Format(fmtStr)
			case len(app.Prune) == 10:
				fmtStr = "2006-01-02"
			case len(app.Prune) == 12:
				fmtStr = "2006-01-02 15"
			case len(app.Prune) == 15:
				fmtStr = "2006-01-02 15:04"
			case len(app.Prune) == 18:
				fmtStr = "2006-01-02 15:04:05"
			default:
				return fmt.Errorf("bad prune date %q", app.Prune)
		}
		dt, err := time.Parse(fmtStr, app.Prune)
		if err != nil {
			return err
		}
		if err := app.PruneItems(dt); err != nil {
			return err
		}
	}
	if app.Display {
		if err := app.Write(out); err != nil {
			return err
		}
	}
	return nil
}
