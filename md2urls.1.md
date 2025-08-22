%md2urls(1) user manual | version 0.0.25 9005bb6
% R. S. Doiel
% 2025-08-20

# NAME 

md2urls

# SYNOPSIS

md2urls [OPTIONS] MARKDOWN_FILE [MARKDOWN_FILE ...]

# DESCRIPTION

md2urls reads in one or more ComomonMark/Markdown files containing a list of links and
returns a new URL file written to instandard output.

# EXAMPLE

Here's an example Markdown file.

~~~Markdown

This is a conceptual Markdown document holding links

- [one](https://one.example.edu)
- [two](https://two.example.edu)
- [three](https://three.example.edu)

~~~

Running the follow will result a URL file formated output for a Markdown
file name "myfile.md".

~~~shell
md2urls myfile.md >myfile.urls
~~~

This is the resulting "myfile.urls"

~~~text
https://one.example.edu "~one"
https://two.example.edu "~two"
https://three.example.edu "~three"
~~~


