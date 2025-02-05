---
title: Deno+TypeScript Port of Skimmer
createDate: 2025-02-04
kewords:
  - Deno
  - TypeScript
  - Feeds
  - Reader
---

# Deno+TypeScript Port of Skimmer

TypeScript seems a good fit for evolving Skimmer and this is an exploration of porting the Go codebase to TypeScript to be compiled by Deno. Just as the Go code relies and external package to do the heavy lifting I can rely on [Dave Winer's JavaScript modules](https://github.com/scripting) to do the heavy listing.

## Help Packages

The follow packages are NPM based but can be used with Deno via add using the "npm:" prefix.

- [reallysimple](https://github.com/scripting/reallysimple), support various common feed formats, RSS, Atom and JSON feed
- [OPML](https://github.com/scripting/opmlPackage), Dave's OPML package for lists of feeds

## Working with deno

The following is how I added these packages for use in Skimmer

~~~shell
deno add npm:reallysimple
deno add npm:opml
~~~

## Configuration

The current Skimmer uses a feed file format derived from Newsboat feed lists.  I will need to port the Go implementation of that to TypeScript if I wish to continue to support that simple text format. OPML import and export makes allot of sense as that would be interoperatable with other software.  YAML would also be a useful configuration format given the readability, easy of typing and support across programming languages.

## Human rendered views

Skimmer's utilies has the ability to turn the harvested feeds into various formats like Markdown. There are several approach.

- [Handlebarsjs](https://handlebarsjs.com), the Handlebars template engine (NPM import)
- [@libs/markdown](https://jsr.io/@libs/markdown), markdown to HTML (jsr import)

## Future ideas to explore after port

Today Skimmer is focused on reading and support a public reading site like [Antenna](https://rsdoiel.github.io/antenna).  Tomorrows Skimmer should be fully read write enabling anyone to have their own Antenna as a turn key solution.  An OPML list can be used to identify collections and the order they are to be listed on the landing page of the Antenna.  Each collection itself could be an OPML file inguested into the channels tables of an SQLite3 database (each database represents a collection).  Active collections would be use to produce an Antenna like website using a simple generic design.

Extending Skimmer to be a feed builder.  A blog is a list of POSTs. POST are a text file which might have metadata beyond the text content (e.g. title, pubDate, createDate, etc). The Metadata of a POST maps to the attributes of an RSS 2 item. As sugject there should be a way of pointing at a Markdown document and directly generate an RSS item in a feed. The channel of the item could be used to build a local copy of the a markdown blog in the file system.  This would be similar to how `pttk blogit` works but tied into the SQLite3 database. Items could be rendered out to a local file system providing the basis for a statically hosted blog. 

Workflow for publishable feed:

- Write a Markdonw document some where on local disk, `edit LOCAL_DOC_NAME`
- Use skimit to ingest the Markdown document, gest it with skimmer, `skimit LOCAL_DOC_NAME TARGET_FEED_NAME`
- Render into a blog structured, `skim2blog TARGET_FEED_NAME [OPTIONS]`
- Render RSS for blog, `skim2rss TARGET_FEED_NAME [OPTIONS]`

Why is this important? A small social network can be formed from two or more end points that both use RSS.  The Antenna project shows that a single person can aggregate their own news and not overburden the remote website.  It's small ammount of procesing that can be done on an SBC like a Raspberry Pi 3 or 4.

A means of also browsing content by hash tags and "at" tag would make it more like a social web engine, similarly leveraging the at tags. OPML can be used to ready build a person's Antenna.

NOTE: Blue Sky and Mastodon already support RSS output.


## Things to document

- Discovering a basic feed source for a website (e.g. RSS or Atom)
  - Wordpress
  - Blue Sky
  - Mastodon
  - Medium
  - Substack
  - Viewing the HTML source
- Using a list of feeds
  - Newsboat text file
  - OPML
- Generating a "local" feed
- Running a blog from a "local" feed
- Quoting another POST in a "local" feed
- Hash tag clouds and facet
- @ tag clouds and facets
