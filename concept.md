
# skimmer

I have a problem. I like to read my feeds in newsboat but I can't seem to get it working on macOS Sonoma
yet or get the send to pocket working under Windows 11 (native and lsw). There are also times I am using
a device where I don't want to install newsboat I just want to look at some current content while I'm
waiting on something. This has lead me to think about skimmer. Something that works with RSS, Atom
and jsonfeeds in the same way I use `yarnc timeline | less -R`.  My inspiration is Dave Winer's 
river of news but minus the outline modality. It's a stream because it is smaller and more ephemeral
than a river. This has left me with some questions.

- How simple is would it be to write skimmer?
- How much effort would be required to maintain it?
- Can this tool incorporate support for twtxt feeds I follow?

There is a Go package called [gofeed](https://github.com/mmcdole/gofeed). The README describes it
as a "universal" feed reading parser.

## Design issues

The reader tools needs to output to standard out in the same manner as `yarnc timeline` does. The goal isn't
to be newsboat or Lynx but to present a stream of items usefully formatted.

Some design questions

1. Are feeds fetch by the same tool as the reader output I pipe to less?
2. How do I handle articles that are reference in more than one feed? (I really don't need to see them more than once)
3. For a given list of feed URLs do I display in descending timestamp order or by feed? Do I offer a choice?
4. Can this tool consolidate my twtxt and RSS feed reading?

# A thin wrapper around gofeed

I see one immediate design decision. Do I fetch the feeds and process them on each invocation or do
I cache there results?  If two feeds point to the same article
I think this could be used to fetch the feed content, parse it and store a common JSON model in a dataset
collection using a SQL Store. This could be done as a tutorial for using dataset Go package. "Reading" records
would then boil down to a query to retrieve the latest or use some other filter (e.g. like dsquery). Pandoc
should work to render the content to the script in a pleasant way. 

I could take advantage of Go's concurrency model to allow updating the feeds in the db table while also
reading what was available or I could provide them as separate capabilities.

Another idea would be to render out to other text centered formats like Gemini and Gopher. This could be
useful to for publishing Gemini and Gopher content from a static blog site while not simply replicating
the repository. E.g. on sfg.org I could pull the content from rsdoiel.github.io and re-process it into
a proper Gopher site and similarly to a Gemini implementation.



