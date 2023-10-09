package skimmer

var (
	// SQLCreateTables provides the statements that are use to create our tables
	// It has two percent s, first is feed list name, second is datetime scheme
	// was generated.
	SQLCreateTables = `-- This is the scheme used for %s's SQLite 3 database
-- %s
CREATE TABLE IF NOT EXISTS channels (
	link PRIMARY KEY,
	title TEXT,
	description TEXT,
	feed_link TEXT,
	links JSON,
	updated DATETIME,
	published DATETIME,
	authors JSON,
	language TEXT,
	copyright TEXT,
	generator TEXT,
	categories JSON,
	feed_type TEXT,
	feed_version TEXT
);

CREATE TABLE IF NOT EXISTS items (
	link PRIMARY KEY,
	title TEXT,
	description TEXT,
	authors JSON,
    updated DATETIME,
	published DATETIME,
	feedLabel TEXT,
	channel TEXT,
	retrieved DATETIME DEFAULT CURRENT_TIMESTAMP
);
`
	// SQLResetChannels clear the channels talbe
	SQLResetChannels = `DELETE FROM channels;`

	// Update the channels in the skimmer file
	SQLUpdateChannel = `REPLACE INTO channels (
link, title, description, feed_link, links,
updated, published, 
authors, language, copyright, generator,
categories, feed_type, feed_version
) VALUES (
?, ?, ?, ?, ?, 
?, ?,
?, ?, ?, ?,
?, ?, ?
);`

	// Update a feed item in the items table
	SQLUpdateItem = `REPLACE INTO items (
link, title, description, updated, published, feedLabel)
VALUES (?, ?, ?, ?, ?, ?);`

	// SQLItemCount returns a list of items in the items table
	SQLItemCount = `-- Count the items in the feed_items table.
SELECT COUNT(*) FROM items;`

	// SQLDisplayItems returns a list of items in decending chronological order.
	SQLDisplayItems = `-- Basic SQL to retrieve an ordered list of items from all feeds.
SELECT link, title, description, updated, published, feedLabel AS label
FROM items
WHERE description != ""
ORDER BY updated DESC;`

	// SQLPruneItems will prune our items table for all items that have easier
	// a updated or publication date early than the timestamp provided.
	SQLPruneItems = `DELETE FROM items 
WHERE (updated IS NULL AND publish IS NULL) 
   OR ((updated >= ? OR published >= ?) AND
   	(updated < ? AND published < ?))
`

)
