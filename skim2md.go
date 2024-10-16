package skimmer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	// 3rd Party Packages
	"github.com/mmcdole/gofeed"
)

// Skim2Md supports the skim2md cli.
type Skim2Md struct {
	// AppName holds the name of the application
	AppName string `json:"app_name,omitempty"`

	// DbName holds the path to the SQLite3 database
	DBName string `json:"db_name,omitempty"`

	// Title if this is set the title will be included
	// when generating the markdown of saved items
	Title string `json:"title,omitempty"`

	// FrontMatter, if true insert Frontmatter block in Markdown output
	FrontMatter bool `json:"frontmatter,omitempty"`

	// PocketButton, if true insert a "save to pocket" button for each RSS item output
	PocketButton bool

	out  io.Writer
	eout io.Writer
}

// NewSkim2Md initialized a new Skim2Md struct
func NewSkim2Md(appName string) (*Skim2Md, error) {
	app := new(Skim2Md)
	app.AppName = appName
	return app, nil
}

func (app *Skim2Md) DisplayItem(link string, title string, description string, enclosures string, updated string, published string, label string, tags string) error {
	// Then see about formatting things.
	pressTime := published
	if len(pressTime) > 10 {
		pressTime = pressTime[0:10]
	}
	if updated != "" {
		if len(updated) > 10 {
			updated = updated[0:10]
		}
		pressTime += ", updated: " + updated
	}
	if strings.HasPrefix(label, `"~`) {
		label = strings.Trim(label, `"~`)
	}
	if title == "" {
		title = fmt.Sprintf("**@%s** (date: %s, from: %s)", label, pressTime, label)
	} else {
		title = fmt.Sprintf("## %s\n\ndate: %s, from: %s", title, pressTime, label)
	}
	var (
		audioElement string
		videoElement string
		err          error
	)
	if enclosures != "" {
		audioElement, err = enclosuresToAudioElement(enclosures)
		if err != nil {
			fmt.Fprintf(app.eout, "could not make audio element for %s, %s", link, err)
			audioElement = ""
		}
		videoElement, err = enclosuresToVideoElement(enclosures)
		if err != nil {
			fmt.Fprintf(app.eout, "could not make video element for %s, %s", link, err)
			videoElement = ""
		}
	} else {
		audioElement = ""
		videoElement = ""
	}
	mediaElement := strings.Join([]string{
		audioElement,
		videoElement,
	}, "<br>")
	if app.PocketButton {
		fmt.Fprintf(app.out, `---

%s

%s

%s

<span class="feed-item-link">
<a href="%s">%s</a> <a href="https://getpocket.com/save" class="pocket-btn" data-lang="en" data-save-url="%s">Save to Pocket</a>
</span>

`, title, description, mediaElement, link, link, link)
	} else {
		fmt.Fprintf(app.out, `---

%s

%s 

%s 

<%s>

`, title, description, mediaElement, link)
	}
	return nil
}

// Write, display the contents from database
func (app *Skim2Md) Write(db *sql.DB) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if app.FrontMatter {
		fmt.Fprintf(app.out, "---\n")
		if app.Title != "" {
			fmt.Fprintf(app.out, "title: %s\n", app.Title)
		}
		fmt.Fprintf(app.out, "updated: %s\n", timestamp)
		fmt.Fprintf(app.out, "---\n\n")
	}
	if app.Title != "" {
		fmt.Fprintf(app.out, `# %s

(date: %s)

`, app.Title, timestamp)
	} else {
		fmt.Fprintf(app.out, `
(date: %s)

`, timestamp)
	}
	stmt := SQLDisplayItems
	rows, err := db.Query(stmt, "saved")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			link        string
			title       string
			description string
			enclosures  string
			updated     string
			published   string
			label       string
			tags        string
		)
		if err := rows.Scan(&link, &title, &description, &enclosures, &updated, &published, &label, &tags); err != nil {
			fmt.Fprintf(app.eout, "%s\n", err)
			continue
		}
		if err := app.DisplayItem(link, title, description, enclosures, updated, published, label, tags); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if app.PocketButton {
		fmt.Fprintf(app.out, `

<script type="text/javascript">!function(d,i){if(!d.getElementById(i)){var j=d.createElement("script");j.id=i;j.src="https://widgets.getpocket.com/v1/j/btn.js?v=1";var w=d.getElementById(i);d.body.appendChild(j);}}(document,"pocket-btn-js");</script>

`)
	}
	return nil
}

func (app *Skim2Md) Run(out io.Writer, eout io.Writer, args []string, frontMatter bool, pocketButton bool) error {
	if len(args) < 1 {
		return fmt.Errorf("missing skimmer database file")
	}
	if len(args) > 1 {
		return fmt.Errorf("expected only one skimmer database file %+v", args)
	}
	app.FrontMatter = frontMatter
	app.PocketButton = pocketButton
	app.out = out
	app.eout = eout
	dsn := args[0]
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := app.Write(db); err != nil {
		return err
	}
	return nil
}

func enclosuresToAudioElement(enclosures string) (string, error) {
	elements := []*gofeed.Enclosure{}
	if err := json.Unmarshal([]byte(enclosures), &elements); err != nil {
		return "", err
	}
	parts := []string{}
	for _, elem := range elements {
		if strings.Contains(elem.Type, "audio") && elem.Length != "" {
			parts = append(parts, fmt.Sprintf(`<source type="%s" src="%s"></source>`, elem.Type, elem.URL))
			parts = append(parts, fmt.Sprintf(`<a href=%q target="_blank">%s</a>`, elem.URL, elem.Type))
		}
	}
	if len(parts) > 0 {
		return fmt.Sprintf(`<audio crossorigin="anonymous" controls="controls">
%s
</audio>`, strings.Join(parts, "\n\t")), nil
	}
	return "", nil
}

func enclosuresToVideoElement(enclosures string) (string, error) {
	elements := []*gofeed.Enclosure{}
	if err := json.Unmarshal([]byte(enclosures), &elements); err != nil {
		return "", err
	}
	parts := []string{}
	for _, elem := range elements {
		if strings.Contains(elem.Type, "video") && elem.Length != "" {
			parts = append(parts, fmt.Sprintf(`<source type="%s" src="%s"></source>`, elem.Type, elem.URL))
			parts = append(parts, fmt.Sprintf(`<a href=%q target="_blank">%s</a>`, elem.URL, elem.Type))
		}
	}
	if len(parts) > 0 {
		return fmt.Sprintf(`<video crossorigin="anonymous" controls="controls" width="250">
%s
</video>`, strings.Join(parts, "\n\t")), nil
	}
	return "", nil
}
