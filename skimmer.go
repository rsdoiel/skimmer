package skimmer

import (
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
				key, val = parts[0], ""
			case 2:
				key, val = parts[0], parts[1]
			}
			urls[key] = val
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

	// DbName holds the path to the SQLite3 database 
	DBName string `json:"db_name,omitempty"`

	// Fetch indicates that items need to be retrieved from the url list and stored in the database
	Fetch bool `json:"fetch,omitempty"`

	// Urls are the map of urls to labels to be fetched or read
	Urls map[string]string `json:"urls,omitempty"`

	// Limit contrains the number of items shown
	Limit int `json:"limit,omitempty"`

	// Prune contains the date to use to prune the database.
	Prune bool `json:"prune,omitempty"`

	// Interactive if true causes Run to display one item at a time with a minimal of input
	Interactive bool `json:"interactive,omitempty"`

	// AsURLs, output the skimmer feeds as a newsboat style url file
	AsURLs bool `json:"urls,omitempty"`

	// AsOPML, output the skimmer feeds as OPML
	AsOPML bool `json:"opml,omittempty"`

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
func (app *Skimmer) Setup(fPath string) error {
	// Check if we have an appDir
	bName := path.Base(fPath)
	xName := path.Ext(bName)
	fName := fPath
	if xName != ".skim" {
		fName = strings.TrimSuffix(fPath, xName) + ".skim"
		// Check to see if we have an existing skimmer file
		if _, err := os.Stat(fName); os.IsNotExist(err) {
			stmt := fmt.Sprintf(SQLCreateTables, app.AppName, time.Now().Format("2006-01-02"))
			dsn := fName
			fmt.Printf("SQLite 3 dsn -> %q\n", dsn)
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
	}
	app.DBName = fName
	return nil
}

// ReadUrls reads urls or OPML file provided and updates the feeds in the skimmer
// skimmer file.
//
// Newsboat's url file format is `<URL><SPACE>"~<LABEL>"` one entry per line
// The hash mark, "#" at the start of the line indicates a comment line.
//
// OPML is documented at http://opml.org
//
func (app *Skimmer) ReadUrls(fName string) error {
	xName := path.Ext(fName)
	src, err := os.ReadFile(fName)
	if err != nil {
		return err
	}
	if xName == ".opml" {
		return fmt.Errorf("opml implort not implemeted, opml packages need to support ToMap")
		/*
		var o *opml.OPML
		if o, err = opml.Parse(src); err != nil {
			return err
		}
		app.Urls, err = o.ToMap(urls)
		if err != nil [
			return err
		}
		*/
	} else {
		app.Urls, err = ParseURLList(fName, src)
		if err != nil {
			return err
		}
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

func saveChannel(db *sql.DB, feedLabel string, channel *gofeed.Feed) error {
/*
link, title, description, feed_link, links,
updated, published, 
authors, language, copyright, generator,
categories, feed_type, feed_version
*/
	var (
		err error
		src []byte
		title string
		linksStr string
		authorsStr string
		categoriesStr string
	)
	if feedLabel == "" {
		title = channel.Title
	} else {
		title = feedLabel
	}
	if channel.Links != nil {
		src, err = JSONMarshal(channel.Links)
		if err != nil {
			return err
		}
		linksStr = fmt.Sprintf("%s", src)
	}
	if channel.Authors != nil {
		src, err = JSONMarshal(channel.Authors)
		if err != nil {
			return err
		}
		authorsStr = fmt.Sprintf("%s", src)
	}
	if channel.Categories != nil {
		src, err = JSONMarshal(channel.Categories)
		if err != nil {
			return err
		}
		categoriesStr = fmt.Sprintf("%s", src)
	}

	stmt := SQLUpdateChannel
	_, err = db.Exec(stmt,
		channel.Link, &title, channel.Description, channel.FeedLink, linksStr, 
		channel.Updated, channel.Published,
		authorsStr, channel.Language, channel.Copyright, channel.Generator,
		categoriesStr, channel.FeedType, channel.FeedVersion)
	if err != nil {
		return err
	}
	return nil
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

func (app *Skimmer) ResetChannels(db *sql.DB) error {
	stmt := SQLResetChannels
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (app *Skimmer) ChannelsToUrls(db *sql.DB) ([]byte, error) {
	return nil, fmt.Errorf("ChannelsToUrls() not implement")
}

// Download the contents from app.Urls
func (app *Skimmer) Download(db *sql.DB) error {
	eCnt := 0
	dsn := app.DBName
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
		if err := saveChannel(db, v, feed); err != nil {
			fmt.Fprintf(app.eout, "failed to save chanel %q, %s\n", k, err)
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
func (app *Skimmer) ItemCount(db *sql.DB) (int, error) {
	stmt := SQLItemCount
	rows, err := db.Query(stmt)
	if err != nil {
		return -1, err
	}
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
func (app *Skimmer) PruneItems(db *sql.DB, startDT time.Time, endDT time.Time) error {
	stmt := SQLPruneItems
	start := startDT.Format("2006-01-02 15:04:05")
	end := endDT.Format("2006-01-02 15:04:05")
	_, err := db.Exec(stmt, start, start, end, end)
	if err != nil {
		return err
	}
	return err
}

// Display the contents from database
func (app *Skimmer) Write(db *sql.DB) error {
	stmt := SQLDisplayItems
	if app.Limit > 0 {
		stmt = fmt.Sprintf("%s LIMIT %d", stmt, app.Limit)
	}
	rows, err := db.Query(stmt)
	if err != nil {
		return err
	}
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

// normalizeTFormat sorts out the format string to used to parse the datestamp/timestamp
func normalizeTFormat(t string) (string, error) {
	fmtStr := "2006-01-02"
	switch {
	case t == "today":
		fmtStr = "2006-01-02"
		t = time.Now().Format(fmtStr)
	case t  == "now":
		fmtStr = "2006-01-02 15:04:05"
		t = time.Now().Format(fmtStr)
	case len(t) == 10:
		fmtStr = "2006-01-02"
	case len(t) == 12:
		fmtStr = "2006-01-02 15"
	case len(t) == 15:
		fmtStr = "2006-01-02 15:04"
	case len(t) == 18:
		fmtStr = "2006-01-02 15:04:05"
	default:
		return "", fmt.Errorf("bad prune date %q", t)
	}
	return fmtStr, nil
}


// Run provides the runner for skimmer. It allows for testing of much of the cli functionality
func (app *Skimmer) Run(out io.Writer, eout io.Writer, args []string) error {
	app.out = out
	app.eout = eout
	if len(args) == 0 {
		return fmt.Errorf("expected a .skim, an OPML or urls file to process")
	}
	fPath := args[0]
	xName := path.Ext(fPath)
	// See if we need to setup things.
	if err := app.Setup(fPath); err != nil {
		return err
	}
	var (
		db *sql.DB
		err error
	)

	dsn := app.DBName
	db, err = sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	// See if we need to update the channels
	if xName != ".skim" {
		if err := app.ReadUrls(fPath); err != nil {
			return err
		}
		if err := app.ResetChannels(db); err != nil {
			return err
		}
		// By Downloading the urls all the channel data will get created/updated.
		if err := app.Download(db); err != nil {
			return err
		}
		cnt, err := app.ItemCount(db)
		if err != nil {
			fmt.Fprintf(app.eout, "fail to count items, %s", err)
		}
		fmt.Fprintf(app.out, "\n%d items available to read\n", cnt)
	} else if app.Fetch {
		if err := app.Download(db); err != nil {
			return err
		}
		cnt, err := app.ItemCount(db)
		if err != nil {
			fmt.Fprintf(app.eout, "fail to count items, %s", err)
		}
		fmt.Fprintf(app.out, "\n%d items available to read\n", cnt)
		return nil
	}
	if app.Prune {
		var err error
		if len(args) < 3 {
			return fmt.Errorf("expected a date or date range for prune")
		}
		start, end := "0001-01-01 00:00:00", time.Now().Format("2006-01-02 15:04:05")
		if len(args) == 2 {
			end = args[1]
		}
		if len(args) == 3 {
			start = args[1]
			end = args[2]
		}
		startFmt, err := normalizeTFormat(start)
		if err != nil {
			return err
		}
		endFmt, err := normalizeTFormat(end)
		if err != nil {
			return err
		}
		startDt, err := time.Parse(startFmt, start)
		if err != nil {
			return err
		}
		endDt, err := time.Parse(endFmt, end)
		if err != nil {
			return err
		}
		if err := app.PruneItems(db, startDt, endDt); err != nil {
			return err
		}
		cnt, err := app.ItemCount(db)
		if err != nil {
			fmt.Fprintf(app.eout, "fail to count items, %s", err)
		}
		fmt.Fprintf(app.out, "\n%d items available to read\n", cnt)
	}
	if app.AsOPML {
		return fmt.Errorf("OPML output not implemented")
	}
	if app.AsURLs {
		src, err := app.ChannelsToUrls(db)
		if err != nil {
			return err
		}
		fmt.Fprintf(app.out, "%s\n", src)
		return nil
	}

	if app.Interactive {
		return fmt.Errorf("interactive not implemented")
	}
	if err := app.Write(db); err != nil {
		return err
	}
	return nil
}
