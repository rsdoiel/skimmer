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
	published DATETIME,
	updated DATETIME,
	feed_type TEXT,
	feed_version TEXT,
	copyright TEXT,
	language TEXT,
	authors JSON,
	categories JSON,
	dublin_core JSON,
	feed_link TEXT,
	links JSON,
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
ORDER BY updated DESC
`
	// SQLPruneItems will prune our items table for all items that have easier
	// a updated or publication date early than the timestamp provided.
	SQLPruneItems = `DELETE FROM items 
WHERE (updated < ?) OR ((published < ?) AND (updated < ?)) 
   OR (published = "" AND updated = "")
`

)
