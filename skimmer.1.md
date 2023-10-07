%skimmer(1) skimmer user manual | version 0.0.1 5126b2e
% R. S. Doiel
% 2023-10-06

# NAME

skimmer

# SYNOPSIS

skimmer [OPTIONS]

# DESCRIPTION

skimmer is a extremely lightweight feed reader. It can do two things. Retrieve
feeds from a list of urls and store the contents in a SQLite 3 database. It can
write the contents of the database to standard output (the default behavior).
It displays the harvested content in reverse chronological order. The content
display is not "paged" so typically you would use skimmer in conjunction with
the POSIX command 'more' or GNU command 'less'.

The url list file uses the newsboat url file format and skimmer looks in the 
`$HOME/.skimmer/skimmer.urls` directory for that file.

By default (i.e. no options specified) skimmer will check the database
and see if there are less than 100 items. If so it'll perform a "fetch"
operation then do a "display" operation. If there are items
in the database a "display" option is attempted.

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

-display
: display the downloaded contents

-limit N
: display the N most recent items.

-prune TIMESTAMP
: Remove items in database older then TIMESTAMP. If TIMESTAMP "today" the
it removes items before today's date, if it is "now" prunes everything with
a timestamp before the current time. Timestamps can be day or day hour in
the YYYY-MM-DD or YYYY-MM-DD HH:MM:SS formats.

# EXAMPLE

skimmer read your feed items and page items with "less -R".

~~~
skimmer | less -R
~~~

The fetch read the url file and download any feed items found.

~~~
skimmer -fetch
~~~

Display the current contents of feed item database in reverse chronological order.

~~~
skimmer -display
~~~

skimmer can prune it's own database and also limit the count of items displayed.
In this example we're pruning all the items older than today and displaying the recent
five items.

~~~
skimmer -prune today -limit 5
~~~

If I limit the number of items I am reading to about 100 or so I've found
that this combination works nice.

~~~
skimmer -limit 100 | pandoc -f markdown -t plain | less -R
~~~


# Acknowledgments

This project is an experiment it would not be possible without the authors of
newsboat, SQLite3, Pandoc and mmcdole's gofeed package have provided excellent
free software to the planet. - RSD, 2023-10-07


