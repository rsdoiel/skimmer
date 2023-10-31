
# Url Lists

The skimmer URL file format is based on Newsboat's urls file format. Newsboat urls files
are are structed like the following lines.

~~~
# <COMMENT IS HERE>
<URL>
<URL> "~<LABEL>"
~~~

Comment line start with a pound sign. Lines are expected to have a URL for a feed (one per line).
They may be followed by a label that is separated from the URL by one or more spaces, followed by
a double quote and tilde. The text of the label. The is closed by another double quote.

skimmer extends this syntax in one way. If you have a label you can follow that by a user agent
string. The reason for this is that some websites that host feeds controll access by white listing
specific user agent strings. Skimmer's default user agent string will not be on the white list since
skimmer is new and very experimental. If this is the case you can include a specific user agent
string after the label. Here's an example listing where the feed source will accept a user agent
string of "curl/8.4.0"

~~~
https://example.io/feed/ "~Some blog content" curl/8.4.0
~~~

skimmer's URL list parser looks ` "` os the column delimiter. If the column value starts
with a leading tilde it will be striped. If a column is missing (e.g. the label or user agent)
then the value will be assumed to be an empty string. E.g. no label or no user agent string.

