package skimmer

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	// 3rd Party Packages
	_ "github.com/glebarez/go-sqlite"
	"github.com/kayako/bluemonday"
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

	// Map in some private data to keep things easy to work with.
	in   io.Reader
	out  io.Writer
	eout io.Writer
}

func NewSkimmer(appName string) (*Skimmer, error) {
	app := new(Skimmer)
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
				return fmt.Errorf("%s\nstmt: %s", err, stmt)
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
func (app *Skimmer) ReadUrls(fName string) error {
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
func webget(href string) (*gofeed.Feed, error) {
	res, err := http.Get(href)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", res.Status)
	}
	fp := gofeed.NewParser()
	feed, err := fp.Parse(res.Body)
	if err != nil {
		return nil, fmt.Errorf("feed error for %q, %s", href, err)
	}
	if feed.Link == "" || feed.Link == "/" {
		u, err := url.Parse(href)
		if err != nil {
			return nil, err
		}
		u.Path = "/"
		feed.Link = u.String()
	}
	return feed, nil
}

func saveChannel(db *sql.DB, link string, feedLabel string, channel *gofeed.Feed) error {
	/*
	   link, title, description, feed_link, links,
	   updated, published,
	   authors, language, copyright, generator,
	   categories, feed_type, feed_version
	*/
	var (
		err           error
		src           []byte
		title         string
		linksStr      string
		authorsStr    string
		categoriesStr string
	)
	linksStr = link
	title = feedLabel
	if feedLabel == "" {
		title = channel.Title
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
		&link, &title, channel.Description, channel.FeedLink, linksStr,
		channel.Updated, channel.Published,
		authorsStr, channel.Language, channel.Copyright, channel.Generator,
		categoriesStr, channel.FeedType, channel.FeedVersion)
	if err != nil {
		return fmt.Errorf("%s\nstmt: %s", err, stmt)
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
		return fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	return nil
}

func (app *Skimmer) ResetChannels(db *sql.DB) error {
	stmt := SQLResetChannels
	_, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	return nil
}

func (app *Skimmer) MarkItem(db *sql.DB, link string, val string) error {
	stmt := SQLMarkItem
	_, err := db.Exec(stmt, val, link)
	if err != nil {
		return fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	return nil
}

func (app *Skimmer) TagItem(db *sql.DB, link string, tag string) error {
	stmt := SQLTagItem
	_, err := db.Exec(stmt, tag, link)
	if err != nil {
		return fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	return nil
}

// ChannelsToUrls converts the current channels table to Urls formated output
// and refreshes app.Urls data structure.
func (app *Skimmer) ChannelsToUrls(db *sql.DB) ([]byte, error) {
	stmt := SQLChannelsAsUrls
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	defer rows.Close()
	lines := []string{}
	if app.Urls == nil {
		app.Urls = map[string]string{}
	}
	for rows.Next() {
		var (
			link  string
			title string
		)
		if err := rows.Scan(&link, &title); err != nil {
			return nil, err
		}
		if strings.HasPrefix(title, `"`) {
			title = strings.Trim(title, `"~`)
		}
		if link != "" {
			if title != "" {
				lines = append(lines, fmt.Sprintf(`%s "~%s"%s`, link, title, "\n"))
				app.Urls[link] = fmt.Sprintf(`"~%s"`, title)
			} else {
				lines = append(lines, fmt.Sprintf(`%s%s`, link, "\n"))
				app.Urls[link] = ""
			}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return []byte(strings.Join(lines, "")), nil
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
		if err := saveChannel(db, k, v, feed); err != nil {
			fmt.Fprintf(app.eout, "failed to save chanel %q, %s\n", k, err)
			continue
		}

		// Setup a progress output
		t0 := time.Now()
		rptTime := time.Now()
		reportProgress := false
		tot := feed.Len()
		fmt.Fprintf(app.out, "processing %d items from %s\n", tot, v)
		i := 0
		for _, item := range feed.Items {
			if strings.HasPrefix(item.Link, "/") {
				item.Link = fmt.Sprintf("%s%s", strings.TrimSuffix(feed.Link, "/"), item.Link)
			}
			// Add items from feed to database table
			if err := saveItem(db, v, item); err != nil {
				return err
			}
			if rptTime, reportProgress = CheckWaitInterval(rptTime, (20 * time.Second)); reportProgress {
				fmt.Fprintf(app.out, "(%d/%d) %s", i, tot, ProgressETA(t0, i, tot))
			}
			i++
		}
		fmt.Fprintf(app.out, "processed %d/%d from %s\n", i, tot, v)
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
		return -1, fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		if err := rows.Scan(&cnt); err != nil {
			return -1, err
		}
	}
	if err := rows.Err(); err != nil {
		return -1, err
	}
	return cnt, nil
}

// PruneItems takes a timestamp and performs a row delete on the table
// for items that are older than the timestamp.
func (app *Skimmer) PruneItems(db *sql.DB, pruneDT time.Time) error {
	stmt := SQLPruneItems
	dt := pruneDT.Format("2006-01-02 15:04:05")
	_, err := db.Exec(stmt, dt, dt)
	if err != nil {
		return fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	return err
}

func (app *Skimmer) DisplayItem(link string, title string, description string, updated string, published string, label string, tags string) error {
	// Then see about formatting things.
	pressTime := published
	if updated != "" {
		pressTime = updated
	}
	if len(pressTime) > 10 {
		pressTime = pressTime[0:10]
	}
	if title == "" {
		title = fmt.Sprintf("@%s (date: %s)", label, pressTime)
	} else {
		title = fmt.Sprintf("## %s\n\ndate: %s", title, pressTime)
	}
	fmt.Fprintf(app.out, `---

%s

%s 

<%s>

%s
`, title, description, link, tags)
	return nil
}

// Display the contents from database
func (app *Skimmer) Write(db *sql.DB) error {
	stmt := SQLDisplayItems
	if app.Limit > 0 {
		stmt = fmt.Sprintf("%s LIMIT %d", stmt, app.Limit)
	}
	rows, err := db.Query(stmt, "")
	if err != nil {
		return fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			link        string
			title       string
			description string
			updated     string
			published   string
			label       string
			tags        string
		)
		if err := rows.Scan(&link, &title, &description, &updated, &published, &label, &tags); err != nil {
			fmt.Fprint(app.eout, "%s\n", err)
			continue
		}
		if err := app.DisplayItem(link, title, description, updated, published, label, tags); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
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
	case t == "now":
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

func displayStats(out io.Writer, db *sql.DB, dbName string) error {
	stmt := SQLItemStats
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	defer rows.Close()
	fmt.Fprintf(out, "stats: %s\n\tstatus\tcount\n", dbName)
	for rows.Next() {
		var (
			status string
			cnt int
		)
		if err := rows.Scan(&status, &cnt); err != nil {
			return fmt.Errorf("%s\nstmt: %s", err, stmt)
		}
		fmt.Fprintf(out, "\t%s\t%d\n", status, cnt)
	}
	return nil
}

// RunInteractive provides a sliver of interactive UI, basically displaying an item then
// prompting for an action.
func (app *Skimmer) RunInteractive(db *sql.DB) error {
	SetupScreen(app.out)
	ClearScreen()
	// Get the item count
	tot, err := app.ItemCount(db)
	if err != nil {
		return err
	}
	padding := int(len(fmt.Sprintf("%d", tot)))
	promptStr := "(n)ext, (s)ave, (q)uit %" + fmt.Sprintf("%d", padding) + fmt.Sprintf("d/%d > ", tot)

	stmt := SQLDisplayItems
	if app.Limit > 0 {
		stmt = fmt.Sprintf("%s LIMIT %d", stmt, app.Limit)
	}
	rows, err := db.Query(stmt, "")
	if err != nil {
		return fmt.Errorf("%s\nstmt: %s", err, stmt)
	}
	i := 0
	readItems := []string{}
	savedItems := []string{}
	tagItems := []map[string]string{}
	// Step 1 don't trust the data, sanitize it with BlueMonday
	//p :=  bluemonday.UGCPolicy()
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowAttrs("href").Matching(regexp.MustCompile(`(?i)mailto|http|https|gopher|ftp?`)).OnElements("a")
	p.AllowElements("p")
	p.AllowElements("b")
	p.AllowElements("i")
	p.AllowElements("ul")
	p.AllowElements("ol")
	p.AllowElements("li")
	p.AllowElements("code")
	p.AllowElements("pre")
	p.AllowElements("blockquote")
	//p.AllowElements("div")
	// Step 2 get data and then sanitize it.
	for rows.Next() {
		var (
			link        string
			title       string
			description string
			updated     string
			published   string
			label       string
			tags        string
		)
		i++
		if err := rows.Scan(&link, &title, &description, &updated, &published, &label, &tags); err != nil {
			fmt.Fprintf(app.eout, "%s\n", err)
			continue
		}
		// Now sanitize each element
		title = html.UnescapeString(p.Sanitize(title))
		description = html.UnescapeString(p.Sanitize(description))
		if err := app.DisplayItem(link, title, description, updated, published, label, tags); err != nil {
			return err
		}
		// Wait for some input
		quit := false
		prompt := true
		for prompt {
			buf := bufio.NewReader(app.in)
			fmt.Fprintf(app.out, promptStr, i)
			src, err := buf.ReadBytes('\n')
			if err != nil {
				fmt.Fprintf(app.eout, "do not understand %q?\n", src)
			} else {
				prompt = false
			}
			answer := strings.ToLower(strings.Trim(string(src), " \t\r\n"))
			switch answer {
			case "":
				prompt = false
				ClearScreen()
			case "o":
				OpenInBrowser(app.in, app.out, app.eout, link)
				prompt = true
			case "n":
				readItems = append(readItems, link)
				prompt = false
				ClearScreen()
			case "s":
				savedItems = append(savedItems, link)
				prompt = false
				ClearScreen()
			case "t":
				fmt.Fprintf(app.out, "Enter tags (separated by commas) > ")
				tagBuf := bufio.NewReader(app.in)
				src, err = tagBuf.ReadBytes('\n')
				if err != nil {
					fmt.Fprintf(app.eout, "failed to read tags, %s\n", err)
				} else {
					tag := string(src)
					if tag != "" {
						tagItems = append(tagItems, map[string]string{
							"link": link,
							"tag":  tag,
						})
					}
					// Tagged items should also be auto-saved.
					savedItems = append(savedItems, link)
				}
				prompt = true
			case "q":
				prompt = false
				quit = true
			case "stats":
				prompt = true
				if err := displayStats(app.out, db, app.DBName); err != nil {
					fmt.Fprintf(app.eout, "%s\n", err)
				}
			default:
				fmt.Fprintf(app.eout, "do not understand %q?\n", answer)
				prompt = true
			}
		}
		if quit {
			break
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	rows.Close()
	if len(savedItems) > 0 {
		fmt.Fprintf(app.out, "saving %d items ...\n", len(savedItems))
		for _, link := range savedItems {
			if err := app.MarkItem(db, link, "saved"); err != nil {
				return err
			}
		}
	}
	if len(readItems) > 0 {
		fmt.Fprintf(app.out, "marking %d items read ...\n", len(readItems))
		for _, link := range readItems {
			if err := app.MarkItem(db, link, "read"); err != nil {
				return err
			}
		}
	}
	if len(tagItems) > 0 {
		fmt.Fprintf(app.out, "tagging %d items ...\n", len(tagItems))
		for _, obj := range tagItems {
			if link, hasLink := obj["link"]; hasLink {
				if tag, hasTag := obj["tag"]; hasTag && tag != "" {
					if strings.Contains(tag, ",") {
						tags := strings.Split(tag, ",")
						for i, s := range tags {
							tags[i] = strings.TrimSpace(s)
						}
						src, err := JSONMarshal(tags)
						if err != nil {
							return err
						}
						tag = string(src)
					} else {
						tag = fmt.Sprintf("[%q]", tag)
					}
					if err := app.TagItem(db, link, tag); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// Run provides the runner for skimmer. It allows for testing of much of the cli functionality
func (app *Skimmer) Run(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	app.in = in
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
		db  *sql.DB
		err error
	)

	datestamp := ""
	if len(args) > 1 {
		datestamp = args[1]
	}

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
			fmt.Fprintf(app.eout, "fail to count items, %s\n", err)
		}
		fmt.Fprintf(app.out, "\n%d items available to read\n", cnt)
		return nil
	}
	if app.Prune {
		var err error
		if len(args) < 2 {
			return fmt.Errorf("expected a date or timestamp for pruning")
		}
		datestampFmt, err := normalizeTFormat(datestamp)
		if err != nil {
			return err
		}
		dt, err := time.Parse(datestampFmt, datestamp)
		if err != nil {
			return err
		}
		if err := app.PruneItems(db, dt); err != nil {
			return err
		}
		cnt, err := app.ItemCount(db)
		if err != nil {
			fmt.Fprintf(app.eout, "fail to count items, %s\n", err)
		}
		fmt.Fprintf(app.out, "\n%d items available to read\n", cnt)
		return nil
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
		return app.RunInteractive(db)
	}
	if err := app.Write(db); err != nil {
		return err
	}
	return nil
}
