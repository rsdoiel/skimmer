
# Action Items

## Bugs


## Next

- [ ] I would like a tool that you can point at a URL and it'll return the feed link if one is found in the markup or usually places
    - [ ] check in WP location for feed
    - [ ] Microblog's location for feed
    - [ ] Check substack location for feed
    - [ ] Check for Mastodon's RSS feed if linked in page
    - [ ] Check for Bluesky feed if Bluesky starts supporting RSS
    - [ ] Check medium location of feed
    - [ ] check for index.xml, rss.xml, atom.xml or page basename plus .xml for feed
    - [ ] check for the JSON feed names 
- [ ] There are feed like pages and API (e.g. Weather Forcast as DWML documents) I would like to include in news reading, need to figure it if this is in skimmer or related tool, it might even be as simple as writing a suitable Web Component
- [ ] I need a way to take the saved content and render a new RSS feed, e.g. a skim2rss
- [ ] I need a way to take a webpage with links and render an RSS 2 feed, there is existing software that does this, html2rss, checkout https://github.com/html2rss/html2rss
- [ ] I need a way to save adhoc lists of items to an RSS feed
- [ ] I need a simple tool that can take a list saved URLs and generate RSS from them by retrieving them and populating an item for each
    - [ ] See if there is an existing tool like html2rss that does this or if I could just use a Pandoc markdown doc as a list and use html2rss
    - [ ] figure description and link
    - [ ] figure out title (optional)
- [ ] I need a way to take a web page, or a list of web pages and transform page links into RSS
    - [ ] Look at [html2rss](html2rss.github.io) works but without Ruby
that has a list of links and transform it into an RSS feed (e.g. for blogs that don't provide RSS but do provide a list of posts in HTML)
- [ ] I should think about a skim2html that doesn't rely on Pandoc if I get a decent HTML structure worked out
- [ ] Look implementing full text search via SQLite on items table via FTS4 or FTS5
	- expose as command line option and in interactive mode
- [ ] Evaludate Go Lua implementations, extended skimmer will lua filters on input (pre-feed parsing) and on gofeed.Feed struct
      and to filter items in or out from a feed (e.g. flag items for "read" or "save" items that have some identifable element)
      - [x] https://github.com/arnodel/golua a Lua 5.4 implementation
      	- No releases, last commit Febraury 2023
      - [ ] https://github.com/RyouZhang/go-lua (a Lua Jit and embedable environment)
      	- Has releases, last commit April 2023, last release was April 2023
      - [ ] https://github.com/yuin/gopher-lua Lua 5.1 implementation
      	- Has releases, last commit Oct 2023, two releases
      - [ ] https://github.com/Shopify/go-lua A Lua 5.2 imeplementation
      	- No releases, last commit October 2022, used by Shopify since 2014
      - [ ] https://github.com/vlorc/lua-vm
      	- Has releases, last commit in 2021, last release Nov. 2020
- [x] Reviews newsboat document on urls file, make sure I cover what is supported, see https://wiki.archlinux.org/title/Newsboat
- [x] Document how skimmer+skim2md+pandoc can be used to create a personal aggregation page (NOTE: See Antenna project, that's done there)
- [x] Per feed I need the option to provide specific headers and user agent (NOTE: user agent done, not sure I need the header really)
- [x] Add an option to list items for a specific feed (NOTE: this can be done via SQL query of the items table)
- [x] User-Agents may need to be set per feed, this may require a change in format of the urls file, CSV is starting to make more sense to manage a list of urls
	- [x] NOAA Weather API suggests including a contact email in header string, for problem responses, should provide for that
- [x] stats should use the current saved/read queues in reflecting current stats or trigger updates for read/saved then run stats
- [x] Add button to add metadata frontmatter to generated Markdown output in skim2md (needed by Antenna to make archived pages PageFind search friendly)
- [x] Add an option to include a "save to pocket" button for each RSS item displayed in skim2md
    - See https://getpocket.com/publisher/button_docs (shows how to explicitly add a perma link to for use by save to pocket)
- [x] Add a stats option to show the items in the database


## Someday, maybe

- [ ] Add feed detection so I can point at a URL and auto-magically return the feed URL if avaialble
	- [ ] RSS
	- [ ] JSON feed
	- [ ] Atom feed
- [ ] Add SQL based filter options for viewing and for the actions of marking read, saved or ignoring
- [ ] Add a full text search option to look for specific items
- [ ] Add a "goto item" mechanism so I can reset what I am viewing in the result list
    - [ ] investigate using SQL LIMIT on main query to achieve restarting the query just before the desired item
- [ ] Add support to read Gopher URLs in skimmer
- [ ] Add support to read Gemini URLs in skimmer
- [ ] Add support to query the item database, sort of like dsquery in dataset
- [ ] Add a way to output content to a local staging directory and search it with pagefind on localhost
- [ ] Add a "open" option in interactive mode
    - [ ] use a similar setup as newsboat, document how to create a bash/bat file to invoke a GUI browser for macOS/Windows
- [ ] Add an "review editor" mode which willl bring the item into a editor session so I can write a blog post about it 
    - [ ] Use a wrapping bash/bat file for GUI editors like newsboat handles opening a GUI web browser
- [ ] Add lua support to add feed automation in a manner like Pandoc filters
- [ ] Add a send via Pandoc that takes the saved items and builds a web page that then can be opened in the web browser
    - [ ] see if I can add buttons for save link to Pocket
