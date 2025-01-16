package skimmer

import (
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"strings"
	"os"
	"time"

	// 3rd Party packages
	"github.com/gocolly/colly/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
)

// Htm2Skim uses the Coly Golang package to scrape a website and turn it into an
// RSS feed.
//
// Html2Skim struct holds the configuration for scraping a webpage and
// and updating a skimmer database populating both the channel table and
// items table based on how the struct is set.
type Html2Skim struct {
	// AppName holds the name of the application
	AppName string `json:"app_name,omitempty"`

	// DbName holds the path to the SQLite3 database
	DBName string `json:"db_name,omitempty"`

	// URL holds the URL to visit to collect items from
	URL string `json:"url,omitempty"`

	// Selector holds the HTML selector to used to retrieve links
	// an empty page will result looking for all href in the page document
	Selector string `json:"selector,omitempty"`

	// Title holds channel title for the psuedo feed created by scraping
	Title string `json:"title,omitempty"`

	// Description holds the channel description for the pseudo feed created by scraping
	Description string `json:"description,omitempty"`

	// Link set the feed link for channel, this is useful if you render a pseudo feed to RSS
	Link string `json:"link,omitempty"`

	// Generator lets you set the generator value for the channel
	Generator string `json:"generator,omitempty"`

	// LastBuildDate sets the date for the channel being built
	LastBuildDate string `json:"last_build_date,omitempty"`

	out  io.Writer
	eout io.Writer
}

// NewHtml2Skim initialized a new Html2Skim struct
func NewHtml2Skim(appName string) (*Html2Skim, error) {
	app := new(Html2Skim)
	app.AppName = appName
	return app, nil
}

// Scrape takes a Skimmer database, a URI (url) and CSS selector pointing
// at anchor elements you want to create a feed with. It then collects those
// links and renders a feed struct and error value.
func (app *Html2Skim) Scrape(db *sql.DB, uri string, selector string) (*gofeed.Feed, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	// NOTE: We need to define a channel to hold our items
	feed := new(gofeed.Feed)
	feed.Title = app.Title
	feed.Description = app.Description
	feed.Link = app.Link
	feed.Generator = app.Generator
	now := time.Now()
	feed.Updated = now.Format("2006-01-02 15:04:05")
	feed.Published = now.Format("2006-01-02 15:04:05")
	feed.UpdatedParsed = &now
	feed.PublishedParsed = &now
	linkErrors := 0

	// Setup our scraper base on what we've defined
	c := colly.NewCollector(
		colly.AllowedDomains(u.Host),
	)
	policy := bluemonday.UGCPolicy()
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" {
			// make sure the link parses
			u, err := url.Parse(e.Request.AbsoluteURL(link))
			if err != nil {
				fmt.Fprintf(app.eout, "link error %q, %s\n", link, err)
				linkErrors += 1
			} else {
				// NOTE: If we got  link then we can proceed.
				item := new(gofeed.Item)
				item.Title = strings.ReplaceAll(policy.Sanitize(e.Text), "\n", " ")
				item.Link = u.String()
				// We need to set the date retrieved
				now = time.Now()
				item.Updated = now.Format("2006-01-02 15:04:05")
				item.Published = now.Format("2006-01-02 15:04:05")
				item.PublishedParsed = &now
				item.UpdatedParsed = &now
				fmt.Fprintf(app.out, "Adding %q -> %s\n", item.Title, item.Link)
				// FIXME: Need to figure out how to generate a description of the linked item
				feed.Items = append(feed.Items, item)
			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Fprintln(app.out, "Visiting", r.URL.String())
	})

	// Handle a page request error
	c.OnError(func(_ *colly.Response, err error) {
		err = fmt.Errorf("Something went wrong: %s", err)
	})

	// Visit and Scrape
	c.Visit(uri)
	if linkErrors > 0 {
		if err == nil {
			err = fmt.Errorf("%d links had parse errors", linkErrors)
		} else {
			err = fmt.Errorf("%s, %d links had parse errors", err, linkErrors)
		}
	}
	return feed, err
}

func (app *Html2Skim) Run(out io.Writer, eout io.Writer, args []string, title string, description string, link string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing skimmer database file")
	}
	if len(args) < 2 {
		return fmt.Errorf("expected a url to scrape %+v", args)
	}
	if len(args) > 3 {
		return fmt.Errorf("expected skimmer database file, url and optional selector, got %+v", args)
	}
	app.Title = title
	app.Description = description
	app.Link = link
	app.out = out
	app.eout = eout
	dsn := args[0]
	// FIXME: Before we open the database it needs to exist and be defined
	createDB := false
	if _, err := os.Stat(dsn); os.IsNotExist(err) {
		createDB = true
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if createDB {
		stmt := fmt.Sprintf(SQLCreateTables, app.AppName, time.Now().Format("2006-01-02"))
		_, err = db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("%s\nstmt: %s", err, stmt)
		}
	}

	app.URL = args[1]

	if len(args) == 3 {
		app.Selector = args[2]
	} else {
		app.Selector = `a[href]`
	}

	feed, err := app.Scrape(db, app.URL, app.Selector)
	if err != nil {
		return err
	}
	fmt.Fprintf(app.out, "%d links collected\n", len(feed.Items))
	fmt.Fprintf(app.eout, "DEBUG: %+v\n", feed)

	// Now we save the feed generated in a skimmer database
	if err := SaveChannel(db, app.Link, app.Title, feed); err != nil {
		return err
	}
	for _, item := range feed.Items {
		fmt.Fprintf(app.out, "Writing item %q -> %+v\n", item.Title, item.Link)
		if err := SaveItem(db, app.Title, item); err != nil {
			return err
		}
	}
	return nil
}
