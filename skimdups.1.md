%skimdups(1) user manual | version 0.0.25 f3d94d0
% R. S. Doiel
% 2026-03-30

# NAME 

skimdups

# SYNOPSIS

skimdups [OPTIONS] URL_FILE [URL_FILE ...]

# DESCRIPTION

skimdups reads in one or more URL files and scans them for duplicate URLs. It writes a report
to standard output indication the URL duplicated along with the files and line number where the
duplication is found.

# EXAMPLE

Look for duplicate URLs in page1.urls and page2.urls.

~~~shell
skimdups page1.urls page2.urls
~~~


