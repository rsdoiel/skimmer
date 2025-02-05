export function fmtHelp(
    txt: string,
    appName: string,
    version: string,
    releaseDate: string,
    releaseHash: string,
): string {
return txt.replaceAll("{app_name}", appName).replaceAll("{version}", version)
    .replaceAll("{release_date}", releaseDate).replaceAll(
    "{release_hash}",
    releaseHash,
    );
}
  
export const skimmerHelpText: string = `%{app_name}(1) skimmer user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] FILENAME [DATESTAMP]

# DESCRIPTION

skimmer is a lightweight feed reader inspired by [newsboat](https://newsboat.org) and
[yarnc](https://git.mills.io/yarnsocial/yarn). skimmer is very minimal and lacks features.
That is skimmer's best feature. skimmer tries to do two things well.

- Read a list of URLs and fetch the items and saving them to an SQLite 3 database
- Display the contents of the SQLite3 database in reverse chronological order

That's it. That is skimmer secret power. It does only two things. There is no elaborate
user interface beyond standard input, standard output and standard error found on POSIX
type operating systems.

skimmer needs to know what feed items to download and display. This done by providing a
newsboat style URLs file. The feeds are read and the channel and item information is
stored in an SQLite3 database of a similarly named file but with the `+"`"+`.skim`+"`"+`
extension. When you want to read the downloaded items you invoke skimmer again with
the `+"`"+`.skim`+"`"+` file.  This allows you to easily maintain separate list of feeds
to skim and potentially re-use the feed output.

Presently skimmer is focused on reading RSS 2, Atom and jsonfeeds.

The output format uses Pandoc's style of markdown markup. 
- "---" starts the item record
- This is followed by "##" (aka H2), followed by a date 
(updated or published), followed by title or if none an
"@<LABEL>" where LABEL is the feed's title. 
- Next is the link line (link is in angle brackets)
- Finally this is wrapped up by the description content.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version and release hash

-limit N
: display the N most recent items.

-prune 
: The deletes items from the items table in the skimmer database that are older than the date
or timestamp provided after the skimmer filename.  Timestamp can be specified in several ways.
An alias of "today" would remove all items older than today (this is the oposite behavior of
reading items). If "now" is specified then all items older then the current time would be
removed. Otherwise time can be specified as a date in YYYY-MM-DD format or
timestamp YYYY-MM-DD HH:MM:SS format.

-i, -interactive
: display an item and prompt for next action. e.g. (n)ext, (s)ave, (q)uit. If you press
enter the next item will be displayed without marking changing the items state (e.g. marking it
read). If you press "n" the item will be marked as read before displaying the next item. If you
press "s" the saved and next item will be displayed.  Pressing "q" will quit interactive mode
without changing the last item's state.

-urls
: Output the contents of the SQLite 3 database channels table as a newsboat style URLs list

# EXAMPLE

Create a "my-news.skim" database from "my-news.urls".

~~~
{app_name] my-news.urls
~~~

Now that my-news.skim exists we can read it with

~~~
{app_name} my-news.skim
~~~

Update and read interactively.

~~~
skimmer my-news.urls
skimmer -i my-news.skim
~~~

{app_name} can prune it's own database and also limit the count of items displayed.
In this example we're pruning all the items older than today and displaying the recent
five items.

~~~
{app_name} -prune -limit 5 my-news.skim today
~~~

If I limit the number of items I am reading to about 100 or so I've found
that this combination works nice.

~~~
{app_name} -limit 100 my-news.skim | pandoc -f markdown -t plain | less -R
~~~


# Acknowledgments

This project is an experiment it would not be possible without the authors of
newsboat, SQLite3, Pandoc and mmcdole's gofeed package have provided excellent
free software to the planet. - RSD, 2023-10-07

`;

export const skim2mdHelpText = `%{app_name}(1) {app_name} user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME 

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] SKIM_DB_FILENAME

# DESCRIPTION

{app_name} reads a skimmer DB and writes the saved or tagged items 
to the display in a Markdown friendly way.  This includes embedded
audio and video elements from Podcasts and Audio casts.

# OPTIONS

-help
: display help

-version
: display version info

-license
: display license

-title
: Set a page title to be included in the output of saved items

-frontmatter
: add frontmatter to Markdown output

-pocket
: add "save to pocket" link for each RSS Item displayed

# EXAMPLE

In the example we fetch URL content, read some feeds, save some interactively
then use {app_name} to generate a webpage of saved or tagged items.

~~~
skimmer myfeeds.urls
skimmer -i myfeeds.skim
{app_name} myfeeds.skim >save_items.md
~~~

You could then process the "saved_items.md" file further with Pandoc.

`;

export const html2skimHelpText: string = `%{app_name}(1) {app_name} user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME 

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] SKIM_DB_FILE URL [CSS_SELECTOR]

# DESCRIPTION

{app_name} provides a way to treat a webpage containing links as a feed
storing the results in a skimmer database. A SKIM_DB_FILE will be
created if it does not already exist. Feed channel information can
be specified via {app_name} options. The URL should be for an 
HTML page that {app_name} will scrape for links. An optional
CSS_SELECTOR can be included to filter a specific section of the
HTML document for links. If none is provided then the selector
`+"`"+`a[href]`+"`"+` will be used.

# OPTIONS

-help
: display help

-version
: display version info

-license
: display license

-title
: this is the channel title to use for the links scraped in the page

-description
: this is the channel description to use for links scraped in the page

-link
: this is the link to associated with the channel. You'd set this if
you were going to take the scaped links and turn them into RSS 2.0
documents.

# EXAMPLE

This is an example of scaping a web page identified
as "https://example.edu/" filter the links found in the 
".college-news" element.

~~~
{app_name}  \
     myfeeds.skim \
	 https://example.edu/ \
     '.college-news > a[href]'
~~~

`

