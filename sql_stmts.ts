export const SQLCreateTables = `-- This is the scheme used for %s's SQLite 3 database
-- {{database_name}}
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
	enclosures JSON DEFAULT '',
	authors JSON,
	updated DATETIME,
	published DATETIME,
	label TEXT,
	tags JSON DEFAULT '',
	channel TEXT,
	retrieved DATETIME DEFAULT CURRENT_TIMESTAMP,
	status TEXT DEFAULT '',
	dc_ext JSON
);
`,
	// SQLResetChannels clear the channels talbe
	SQLResetChannels = `DELETE FROM channels;`,

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
);`,

	// Update a feed item in the items table
	SQLUpdateItem = `INSERT INTO items (
	link, title, description, enclosures, updated, published, label, authors, dc_ext, tags)
VALUES (
	?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10
) ON CONFLICT (link) DO
  UPDATE SET title = ?2, description = ?3, enclosures = ?4, updated = ?5,
      published = ?6, label = ?7, tags = ?10;`,

	// Return link and title for Urls formatted output
	SQLChannelsAsUrls = `SELECT link, title FROM channels ORDER BY link;`,

	// SQLItemCount returns a list of items in the items table
	SQLItemCount = `SELECT COUNT(*) FROM items;`,

	// SQLItemStats returns a list of rows with totals per status
	SQLItemStats = `SELECT IIF(status = '', 'unread', status) AS status, COUNT(*) FROM items GROUP BY status ORDER BY status`,

	// SQLDisplayItems returns a list of items in decending chronological order.
	SQLDisplayItems = `SELECT link, title, description, enclosures,
	updated, published, label, tags
FROM items
WHERE (description != "" OR title != "") AND status = ?
ORDER BY published DESC, updated DESC;`,

	SQLMarkItem = `UPDATE items SET status = ? WHERE link = ?;`,

	SQLTagItem = `UPDATE items SET tags = ? WHERE link = ?;`,

	// SQLPruneItems will prune our items table for all items that have easier
	// a updated or publication date early than the timestamp provided.
	SQLPruneItems = `DELETE FROM items 
WHERE ((updated IS NULL AND published IS NULL) OR
   (updated == '' AND published == '')
   OR (updated < ? AND published < ?))
   AND status = '';
`;
