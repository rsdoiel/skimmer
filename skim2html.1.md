%skim2html(1) user manual | version 0.0.25 9005bb6
% R. S. Doiel
% 2025-08-20

# NAME 

skim2html

# SYNOPSIS

skim2html [OPTIONS] SKIM_DB_FILENAME [YAML_CONFIG]

# DESCRIPTION

skim2html reads a skimmer DB and related YAML file (configuration)
writes the saved or tagged items to the display in HTML.  This includes embedded
audio and video elements from Podcasts and Audio casts.

# CONFIGURATION

The configuration file is assumed to have a similarly named YAML file where the
file extension is ".yaml" instead of ".skim". If YAML_CONFIG is provided it'll be
used.

The configuration file supports the following attributes.

title
: (optional) This title of the feed and resulting page

description
: (optional) Any additional description (included in the head, meta element)

CSS
: (optional) Path to CSS file (use `@import` to include other CSS)

Modules
: (optional) a list of paths to ES6 modules to include in the page

Header
: (optional) HTML markup for header (if not included one will be generated from title and timestamp)

Nav
: (optional) HTML markup to be included for navigation

Footer
: (optional) HTML markup to be included in the footer


# OPTIONS

-help
: display help

-version
: display version info

-license
: display license

# EXAMPLE

In the example we fetch URL content, read some feeds, save some interactively
then use skim2html to generate a webpage of saved or tagged items.

~~~
skimmer myfeeds.urls
skimmer -i myfeeds.skim
skim2html myfeeds.skim >save_items.html
~~~

You could then process the "saved_items.md" file further with Pandoc.


