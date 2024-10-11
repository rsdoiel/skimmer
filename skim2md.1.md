%skim2md(1) skim2md user manual | version 0.0.17 d91cbc7
% R. S. Doiel
% 2024-10-11

# NAME 

skim2md

# SYNOPSIS

skim2md [OPTIONS] SKIM_DB_FILENAME

# DESCRIPTION

skim2md reads a skimmer DB and writes the saved or tagged items 
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
then use skim2md to generate a webpage of saved or tagged items.

~~~
skimmer myfeeds.urls
skimmer -i myfeeds.skim
skim2md myfeeds.skim >save_items.md
~~~

You could then process the "saved_items.md" file further with Pandoc.


