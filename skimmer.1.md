%skimmer(1) skimmer user manual | version 0.0.20 dbf71f8
% R. S. Doiel
% 2025-01-16

# NAME

skimmer

# SYNOPSIS

skimmer [OPTIONS] FILENAME [DATESTAMP]

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
stored in an SQLite3 database of a similarly named file but with the `.skim`
extension. When you want to read the downloaded items you invoke skimmer again with
the `.skim` file.  This allows you to easily maintain separate list of feeds
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
skimmer my-news.skim
~~~

Update and read interactively.

~~~
skimmer my-news.urls
skimmer -i my-news.skim
~~~

skimmer can prune it's own database and also limit the count of items displayed.
In this example we're pruning all the items older than today and displaying the recent
five items.

~~~
skimmer -prune -limit 5 my-news.skim today
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


