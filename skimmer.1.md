%skimmer(1) skimmer user manual | version 0.0.2 a6c6bc9
% R. S. Doiel
% 2023-10-09

# NAME

skimmer

# SYNOPSIS

skimmer [OPTIONS] FILENAME [TIME_RANGE]

# DESCRIPTION

skimmer is a extremely lightweight feed reader. It trys to do two things well.

1. Download feed content
2. Display feed content

To download feed content you need a list of URLs. A list of URLs can be provided
in [newsboat's](https://newsboat.org) urls file format or as an [OPML](http://opml.org/) 
file. 

If either type of these files is provided on the command line then the file will be read
and a similarly named SQLite3 database will be created with a `.skim` extension.
skimmer will then display the downloaded content.

After populating your skimmer database you can update it using the `-fetach`
option or read it by providing the skimmer file instead of a urls file or OPML file.


skimmer followed by a skimmer file displays the harvested content in reverse 
chronological order. The content displayed is not "paged" unless you use the 
`-interactive` option. Typically you would use skimmer in 
conjunction with the POSIX command 'more' or GNU command 'less'. If you provide
a time range then only items published or updated in that time range will be display.
If you only include one timestamp then the items starting with that published or updates
times will be display.

The output format uses Pandoc's style of markdown markup. 
- "--" starts the item record
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

-fetch
: download the latest feed content

-limit N
: display the N most recent items.


-prune 
: The deletes items from the items table for the skimmer file provided. If a time range is provided
then the items in the time range will be deleted. If a single time is provided everything older than
that time is deleted.  A time can be specified in several ways. An alias of "today" would remove all
items older than today (this is the oposite behavior of reading items). If "now" is specified then
all items older then the current time would be removed. Otherwise time can be specified as a date
in YYYY-MM-DD format or timestamp YYYY-MM-DD HH:MM:SS format.

-i, -interactive
: display an item and prompt for next action. e.g. (n)ext, (p)rev, (s)ave, (/)search, (d)elete

-urls
: Output the contents of the SQLite 3 database channels table as a newsboat URLs list

-opml
: Output the contents of the SQLite 3 database channels table as an OPML file.

# EXAMPLE

Create a "my-news.skim" database from "my-news.opml".

~~~
{app_name] my-news.opml
~~~

Now that my-news.skim exists we can read it with

~~~
skimmer my-news.skim
~~~

Update and read the my-news.skim file.

~~~
skimmer -fetch my-news.skim
skimmer my-news.skim
~~~


skimmer can prune it's own database and also limit the count of items displayed.
In this example we're pruning all the items older than today and displaying the recent
five items.

~~~
skimmer -prune today -limit 5 my-news.skim
~~~

If I limit the number of items I am reading to about 100 or so I've found
that this combination works nice.

~~~
skimmer -limit 100 my-news.skim | pandoc -f markdown -t plain | less -R
~~~


# Acknowledgments

This project is an experiment it would not be possible without the authors of
newsboat, SQLite3, Pandoc and mmcdole's gofeed package have provided excellent
free software to the planet. - RSD, 2023-10-07


