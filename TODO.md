
# Action Items

## Bugs

- [ ] Some website refuse the get from skimmer, but seem to provide the data via a web browser fine, I might need to retry with a different user agent string
- [x] Some feeds use relative links in the item or channel information, need to normalize those to full links if possible

## Next

- [ ] There are feed like things (e.g. Weather Forcast as DWML documents) I would like to include in news reading, need to figure it if this is in skimmer or related tool
- [ ] Per feed I need the option to provide specific headers and user agent
- [ ] User-Agents may need to be set per feed, this may require a change in format of the urls file, CSV is starting to make more sense to manage a list of urls
	- [ ] NOAA Weather API suggests including a contact email in header string, for problem responses, should provide for that
- [ ] stats should use the current saved/read queues in reflecting current stats or trigger updates for read/saved then run stats
- [x] Add a stats option to show the items in the database
- [ ] Look implementing full text search on items table via FTS4 or FTS5
	- expose as command line option and in interactive mode
- [ ] Add an option to list items for a specific feed
- [ ] Document how skimmer+skim2md+pandoc can be used to create a personal aggregation page
- [ ] Reviews newsboat document on urls file, make sure I cover what is supported, see https://wiki.archlinux.org/title/Newsboat
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


## Someday, maybe

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
