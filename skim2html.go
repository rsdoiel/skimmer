package skimmer

import (
    "database/sql"
    "fmt"
    "io"
    "os"
    "strings"
    "time"

    // 3rd Party Packages
    "gopkg.in/yaml.v3"
)

// Skim2Html supports the skim2md cli.
type Skim2Html struct {
    // AppName holds the name of the application
    AppName string `json:"app_name,omitempty" yaml:"app_name,omitempty"`

    // DbName holds the path to the SQLite3 database
    DBName string `json:"db_name,omitempty" yaml:"db_name,omitempty"`

    // Title if this is set the title will be included
    // when generating the markdown of saved items
    Title string `json:"title,omitempty" yaml:"title,omitempty"`
    
    // Description, included as metadata in head element
    Description string `json:"description,omitempty" yaml:"description,omitempty"`

    // CSS is the path to a CSS file
    CSS string `json:"css,omitempty" yaml:"css,omitempty"`

    // Modules is a list for ES6 diles
    Modules []string `json:"modules,omitempty" yaml:"modules,omitempty"`

    // Header hold the HTML markdup of the Header element. If not included
    // then it will be generated using the Title and timestamp
    Header string `json:"header,omitempty" yaml:"header,omitempty"` 

    // Nav holds the HTML markup for navigation
    Nav string `json:"nav,omitempty" yaml:"nav,omitempty"`

    // Footer holds the HTML markup for the footer
    Footer string `json:"footer,omitempty" yaml:"footer,omitempty"`

    out  io.Writer
    eout io.Writer
}

// NewSkim2Html initialized a new Skim2Html struct
func NewSkim2Html(appName string) (*Skim2Html, error) {
    app := new(Skim2Html)
    app.AppName = appName
    return app, nil
}

func (app *Skim2Html) DisplayItem(link string, title string, description string, enclosures string, updated string, published string, label string, tags string) error {
    // Setup expressing update time.
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
    // Setup the label (feed name) of the item 
    if strings.HasPrefix(label, `"~`) {
        label = strings.Trim(label, `"~`)
    }
    // Setup the Title
    if title == "" {
        title = fmt.Sprintf("<h1>@%s</h1><br>\n\n(date: %s, from: %s)", label, pressTime, label)
    } else {
        title = fmt.Sprintf("<h1>%s</h1>\n\n(date: %s, from: %s)", title, pressTime, label)
    }

    fmt.Fprintf(app.out, `
    <article data-published=%q data-link=%q>
      %s
      <p>
      %s
      <address>
        <a href=%q>%s</a>
      </address>
    </article>
`, published, link, title, description, link, link)
    return nil
}

func (app *Skim2Html) writeHeadElement() {
    fmt.Fprintln(app.out, "<head>");
    defer fmt.Fprintln(app.out, "</head>")
    // Write out charset
    fmt.Fprintln(app.out, "<meta charset=\"UTF-8\" />")
    // Write title
    if app.Title != "" {
      fmt.Fprintf(app.out, "  <title>%s</title>\n", app.Title)
    }
    fmt.Fprintln(app.out, "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />")
    if app.CSS != "" {
        fmt.Fprintf(app.out, "  <link rel=\"stylesheet\" href=\"%s\" />\n", app.CSS)
    }
    if (app.Modules != nil) {
        for module := range app.Modules {
            fmt.Fprintf(app.out, "  <script type=\"module\" src=\"%s\"></script>\n", module)
        }
    }
}

// Write, display the contents from database
func (app *Skim2Html) Write(db *sql.DB) error {
    // Create the outer elements of a page.
    fmt.Fprintf(app.out, `<!DOCS html>
<html lang="en-US">`);
    defer fmt.Fprintf(app.out, `</html>`)
    // Setup the metadata in the head element
    app.writeHeadElement()
    // Setup body element
    fmt.Fprintln(app.out, "<body>")
    defer fmt.Fprintln(app.out, "</body>")
    // Setup header element
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    if app.Header != "" {
        fmt.Fprintf(app.out, "  <header>%s</header>\n", strings.TrimSpace(app.Header))
    } else if app.Title != "" {
        fmt.Fprintf(app.out, `  <header>
    <h1>%s</h1>

    (date: %s)

  </header>
`, app.Title, timestamp)
    } else {
        fmt.Fprintf(app.out, `  <header>
(date: %s)
  </header>
`, timestamp)
    }
    // Setup nav element
    if app.Nav != "" {
        fmt.Fprintf(app.out, `  <nav>
  %s
  </nav>
`, app.Nav)
    }
    // Setup section
    fmt.Fprintln(app.out, "  <section>")
    stmt := SQLDisplayItems
    rows, err := db.Query(stmt, "saved")
    if err != nil {
        return err
    }
    defer rows.Close()
    // Setup and write out the body
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
    fmt.Fprintln(app.out, "  </section>")
    if app.Footer != "" {
        fmt.Fprintf(app.out, "  <footer>%s</footer>\n", strings.TrimSpace(app.Footer))
    }
    // close the body
    return nil
}

func getDsnAndCfgName(args []string) (string, string) {
    dsn := args[0]
    if (len(args) == 2) {
        return args[0], args[1]
    }
    // Figure out if we have a YAML config or not
    cfgName := strings.TrimSuffix(dsn, ".skim") + ".yaml"
    if _, err := os.Stat(cfgName); err != nil {
        return dsn, ""
    }
    return dsn, cfgName
}

func (app *Skim2Html) LoadCfg(cfgName string) error {
    src, err := os.ReadFile(cfgName)
    if err != nil {
        return err
    }
    obj := Skim2Html{}
    if err := yaml.Unmarshal(src, &obj); err != nil {
        return err
    }
    // Pull in the configuration elements
    if obj.Title != "" {
        app.Title = obj.Title
    }
    if obj.Description != "" {
        app.Description = obj.Description
    }
    if obj.CSS != "" {
        app.CSS = obj.CSS
    }
    if obj.Modules != nil && len(obj.Modules) > 0 {
        app.Modules = obj.Modules[:]
    }
    if obj.Header != "" {
        app.Header = obj.Header
    }
    if obj.Nav != "" {
        app.Nav = obj.Nav
    }
    if obj.Footer != "" {
        app.Footer = obj.Footer
    }
    return nil
}

func (app *Skim2Html) Run(out io.Writer, eout io.Writer, args []string) error {
    cfgName := ""
    if len(args) < 1 {
        return fmt.Errorf("missing skimmer database file")
    }
    if len(args) > 2 {
        return fmt.Errorf("expected only skimmer database file and optional YAML configuration %+v", args)
    }
    if len(args) > 1 {
        cfgName = args[1]
    }
    app.out = out
    app.eout = eout
    dsn, cfgName := getDsnAndCfgName(args)
    if cfgName != "" {
        if err := app.LoadCfg(cfgName); err != nil {
            return err
        }
    }
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
