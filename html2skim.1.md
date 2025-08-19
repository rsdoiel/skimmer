%html2skim(1) html2skim user manual | version 0.0.24 9591fb8
% R. S. Doiel
% 2025-08-19

# NAME 

html2skim

# SYNOPSIS

html2skim [OPTIONS] SKIM_DB_FILE URL [CSS_SELECTOR]

# DESCRIPTION

html2skim provides a way to treat a webpage containing links as a feed
storing the results in a skimmer database. A SKIM_DB_FILE will be
created if it does not already exist. Feed channel information can
be specified via html2skim options. The URL should be for an 
HTML page that html2skim will scrape for links. An optional
CSS_SELECTOR can be included to filter a specific section of the
HTML document for links. If none is provided then the selector
`a[href]` will be used.

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
html2skim  \
     myfeeds.skim \
	 https://example.edu/ \
     '.college-news > a[href]'
~~~


