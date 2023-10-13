
# Action Items

## Bugs

- [ ] Some website refuse the get from skimmer, but seem to provide the data via a web browser fine, I might need to retry with a different user agent string
- [x] Some feeds use relative links in the item or channel information, need to normalize those to full links if possible

## Next

- [x] Add a stats option to show the items in the database
- [ ] Add an option to list items for a specific feed
- [ ] Document how skimmer+skim2md+pandoc can be used to create a personal aggregation page
- [ ] Reviews newsboat document on urls file, make sure I cover what is supported, see https://wiki.archlinux.org/title/Newsboat

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
