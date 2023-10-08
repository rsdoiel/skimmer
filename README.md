
# skimmer

skimmer is a lightweight feed reader inspired by [newsboat](https://newsboat.org) and
[yarnc](https://git.mills.io/yarnsocial/yarn). skimmer is very minimal and lacks features.
That is skimmer's best feature. skimmer tries to do two things well.

- fetch a list of URLs and download their items to an SQLite3 database
- Display the contents of the SQLite3 database in reverse chronological order

That's it.  No elaborate UI beyond what is easily accomplished using standard input,
standard output and standard err.

skimmer needs to know what feeds to download and display. That is done by 
reading either a newsboat style url file or and [OPML](http://opml.org/) file.
These are read and stored in an SQLite 3 database with the same base name of the file read
but with a file extension of `.skim`.  This allows you to easily maintain separate
list of feeds to skim and potentially re-use the feed output.

Presently skimmer is focused on reading RSS 2, Atom and JSONfeeds. If this
experiment evolves further than I hope to add support for txtxt as well as
support for reading feeds from Gopher, Gemini and SFTP sites.

# SYNOSIS

~~~
skimmer [OPTIONS] FILENAE [TIME_RANGE]
~~~

## OPTIONS

-help
: display a help page

-license
: display license

-version
: display version number and build hash

-fetch
: Download items from the list of URLs

-limit N
: Limit the display the N most recent items

-prune 
: The deletes items from the items table for the skimmer file provided. If a time range is provided
then the items in the time range will be deleted. If a single time is provided everything older than
that time is deleted.  A time can be specified in several ways. An alias of "today" would remove all
items older than today. If "now" is specified then all items older then the current time would be 
removed. Otherwise time can be specified as a date in YYYY-MM-DD format or timestamp 
YYYY-MM-DD HH:MM:SS format.

-i, -interactive
: display an item and prompt for next action. e.g. (n)ext, (p)rev, (s)ave, (/)search, (d)elete

# Examples

Fetch and read my newsboat feeds from `.newsboat/urls`. This will create a `.newsboat/urls.skim`.

~~~shell
skimmer .newsboat/urls
~~~

Fetch and read the feeds from `my-news.opml`. This will create a `my-news.skim` file.

~~~shell
skimmer my-news.opml
~~~

Get the latest items for the skimmer file "my-news.skim"
Download some news to read later

~~~shell
skimmer -fetch my-news.skim
~~~

Read the last downloaded content from `my-news.skim`

Display the downloaded news

~~~shell
skimmer my-news.skim
~~~

Limit the number of items sent to the screen.

~~~shell
skimmer -display -limit 25 my-news.skim
~~~

Or my favorite is to run the output through Pandoc
and page with less.

~~~shell
skimmer -display -limit 25 my-news.skim | \
    pandoc -f markdown -t plain | \
    less -R
~~~

Prune the items in the database older than today.

~~~shell
skimmer -prune my-news.skim today
~~~

Prune the items from the month of September 2023.

~~~shell
skimmer -prune my-news.skim \
    "2023-09-01 00:00:00" "2023-09-30 23:59:59"
~~~

## Installation instructions

- [INSTALL.md](INSTALL.md) contains the general steps to install binary releases
- You can download a release from <https://github.com/rsdoiel/skimmer/releases>

## Installation From Source

### Requirements

skimmer is an experimental. The precompiled binaries are not necessarily tested.
To compile from source you need to have git, make, Pandoc SQLite3 and Go.

- Git >= 2
- Make >= 3.8 (GNU Make)
- Pandoc > 3
- SQLite3 > 3.4
- Go >= 1.21.1

### Steps to compile and install

Installation process I used to setup skimmer on a new machine.

~~~
git clone https://github.com/rsdoiel/skimmer
cd skimmer
make
make install
~~~

## Acknowledgments

This experiment would not be possible with the authors of newsboat, SQLite3,
Pandoc and the [gofeed](https://github.com/mmcdole/gofeed) package.
