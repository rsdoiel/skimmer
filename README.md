
# skimmer

skimmer is a lightweight feed reader inspired by [newsboat](https://newsboat.org) and [yarnc](https://git.mills.io/yarnsocial/yarn). skimmer is very minimal and lacks features.  That is skimmer's best feature. skimmer tries to do two things well.

- Read a list of URLs and fetch the items and saving them to an SQLite 3 database
- Display the contents of the SQLite3 database in reverse chronological order

That's it. That is skimmer secret power. It does only two things. There is no elaborate user interface beyond standard input, standard output and standard error found on POSIX type operating systems.

skimmer needs to know what feed items to download and display. This done by providing a newsboat style URLs file. The feeds are read and the channel and item information is stored in an SQLite3 database of a similarly named file but with the `.skim` extension. When you want to read the downloaded items you invoke skimmer again with the `.skim` file.  This allows you to easily maintain separate list of feeds to skim and potentially re-use the feed output.

Presently skimmer is focused on reading RSS 2, Atom and jsonfeeds.


# SYNOPSIS

~~~
skimmer [OPTIONS] URL_LIST_FILENAME
skimmer [OPTIONS] SKIMMER_DB_FILENAME [TIME_RANGE]
~~~

skimmer have two ways to invoke it. You can fetch the contents from list of URLs in newsboat urls file format. You can read the items from the related skimmer database.

## OPTIONS

-help
: display a help page

-license
: display license

-version
: display version number and build hash

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
: display an item and prompt for next action. e.g. (n)ext, (s)ave, (t)ag, (q)uit. If you press enter the next item will be displayed without marking changing the items state (e.g. marking it read). If you press "n" the item will be marked as read before displaying the next item. If you press "s" the item will be tagged as saved and next item will be displayed. If you press "t" you can tag the items. Tagged items are treated as save but the next item is not fetched. Pressing "q" will quit interactive mode without changing the last items state.


# Examples

Fetch and read my newsboat feeds from `.newsboat/urls`. This will create a `.newsboat/urls.skim` 
if it doesn't exist. Remember invoking skimmer with a URLs file will retrieve feeds and their contents and invoking skimmer with the skimmer database file will let you read them.

~~~shell
skimmer .newsboat/urls
skimmer .newsboat/urls.skim
~~~

This will fetch and read the feeds from`my-news.urls`. This will create a `my-news.skim` file.
When the skimmer database is read a simplistic interactive mode is presented.

~~~shell
skimmer my-news.urls
skimmer -i my-news.skim
~~~

The same method is used to update your `my-news.skim` file and read it.

Export the current state of the skimmer database channels to a urls file. Feeds that failed
to be retrieved will not be in the database channels table channels table. This is an 
easy way to get rid of the cruft and dead feeds.

~~~shell
skimmer -urls my-news.skim >my-news.urls
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

skimmer is an experimental. The compiled binaries are not necessarily tested.
To compile from source you need to have git, make, Pandoc, SQLite3 and Go.

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
