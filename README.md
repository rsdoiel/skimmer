
# skimmer

skimmer is a lightweight feed reader inspired by [newsboat](https://newsboat.org) and
[yarnc](https://git.mills.io/yarnsocial/yarn). skimmer is minimal and lacks features. 
That is skimmer's best feature. skimmer can do two things. 

- fetch a list of URLs and download their items to an SQLite3 database
- Display the contents of the SQLite3 database in reverse chronological order

That's it.  No paging, now UI other than the command line options and what is sent
to standard output.

The URL file format is based on newsboat's URL list. That is because I use 
Newsboat for my interactive feed reading. The skimmer application
stores the list of URLs and the SQLite 3 database in your home
directory under `.skimmer`.


## OPTIONS

-help
: display a help page

-license
: display license

-version
: display version number and build hash

-fetch
: Download items from the list of URLs

-display
: Display the contents of the SQLite 3 database

-limit N
: Limit the display the N most recent items

-prune TIMESTAMP
: The deletes items from the database that are older than TIMESTAMP.
TIMESTAMP can be "now","today", a day in YYYY-MM-DD format or a full
timestamp in YYYY-MM-DD HH:MM:SS format.


# Examples

Fetch and read some news

~~~
skimmer
~~~

Download some news to read later

~~~
skimmer -fetch
~~~

Display the downloaded news

~~~
skimmer -display
~~~

Limit the number of items sent to the screen.

~~~
skimmer -display -limit 25
~~~

Or my favorite is to run the output through Pandoc
and page with less.

~~~
skimmer -display -limit 25 | \
    pandoc -f markdown -t plain | \
    less -R
~~~


Prune the items in the database older than today.

~~~
skimmer -prune today
~~~

## Installation instructions

- [INSTALL](INSTALL.md) contains the general steps to install binary releases
- You can download a release from <https://github.com/rsdoiel/skimmer/releases>

## Installation From Source

### Requirements

skimmer is an experimental. The precompiled binaries are not tested.
To compile from source you need to have git, make, Pandoc SQLite3 and Go.

- Git >= 2
- Make >= 3.8 (GNU Make)
- Pandoc > 3
- Go >= 1.21.1
- SQLite3 > 3.4

### Steps to compile and install

Installation process I used to setup skimmer on a new machine.

~~~
git clone https://github.com/rsdoiel/skimmer
cd skimmer
make
make install
~~~

A default URLs list is provided as an example URLs list. The first time
you run skimmer. You should edit `$HOME/.skimmer/skimmer.urls` to fit your needs.


## Acknowledgments

This experiment would not be possible with the authors of newsboat, SQLite3,
Pandoc and the [gofeed](https://github.com/mmcdole/gofeed) package.

