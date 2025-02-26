import { parseArgs } from "@std/cli";
import { licenseText, releaseDate, releaseHash, version } from "./version.ts";
import { fmtHelp, skimmerHelpText } from './helptext.ts';
import * as stmts from './sql_stmts.ts';

import * as reallysimple from 'npm:reallysimple';
import * as opml from "npm:opml";
import { Database } from '@db/sqlite';

export class Skimmer {
	appName: string = '';
	userAgent: string = '';
	options: {[key: string]: any} = {};
	errors: string[] = [];

	constructor(appName: string) {
		this.appName = appName;
		// Setup our options
		let userAgent : string | undefined = Deno.env.get("SKIM_USER_AGENT");
		if (userAgent !== undefined && userAgent !== "") {
			this.userAgent = userAgent;
		}
	};

	display_errors(): string {
		return this.errors.join("\n");
	};

	async harvest(filename: string, options: {[key: string]: any}): Promise<boolean> {
		// FIXME: Check if SQLite3 database and tables exist. If not create them.
		// Get Channel(s)
		const stmt = stmts.SQLChannelsAsUrls;
		const db = new Database(filename);
		if (db === undefined) {
			this.errors.push(`failed to open ${filename} database`);
			return false;
		}
		//const query = db.prepare<[string, string]>(stmt);
		const query = db.prepare(stmt);
		let ok: boolean = true;
		const self = this;
		for (const row of query.iter()) {
			console.log(`DEBUG fetch row -> ${row.link}, ${row.title}`);
			const urlFeed: string = row.link;
			reallysimple.readFeed(urlFeed, function(err: {[key:string]: string}, theFeed: string): any {
				if (err) {
					self.errors.push(err.message);
					ok = false;
				} else {
					//FIXME: push into the items table.
					console.log(JSON.stringify(theFeed, undefined, 2));
				}
				return;
			});
		}
		query.finalize();
		db.close();
		return ok;
	};

	async update_channel(filename: string, options: {[key: string]: any}): Promise<boolean> {
		// FIXME: filename without extension will be the name of of the SQLite3 database
		// holding the collection.
		if (filename.indexOf('.opml') > -1) {
			this.errors.push(`don't know how to OPML yet with ${filename}.`);
			return false;
		} else if (filename.indexOf('.txt') > -1) {
			this.errors.push(`don't know how to process .txt file yet with ${filename}.`);
			return false;
		}
		this.errors.push(`do not know how to process ${filename}`);
		return false;
	};

	async run(options: {[key: string]: any}, args: string[]): Promise<boolean> {
		// Open <database>.channel and loop through the feeds using reallysimple
		// updating <database>.items
		if (args.length === 0) {
			this.errors.push(`missing filename`);
			return false;
		}
		console.log(`DEBUG args in run -> ${args}`);
		for (const val of args) {
			const arg = `${val}`;
			console.log(`DEBUG arg -> ${arg}`);
			if (arg.indexOf('.skim') > -1) {
				if (! await this.harvest(arg, options)) {
					return false;
				}
			} else if ((arg.indexOf('.txt') > -1) || (arg.indexOf('.opml') > -1)) {
				if (! await this.update_channel(arg, options)) {
					return false;
				}
			} else {
				this.errors.push(`do not know how to process ${arg}`);
				return false;
			}
		}
		return true;
	}
}

async function main() {
	const appName = 'skimmer';
	const options = parseArgs(Deno.args, {
		alias: {
			help: "h",
			version: "v",
			interactive: "i",
			urls: "u",
			prune: "p",
			limit: "l"
		},
		default: {
			help: false,
			version: false,
			interactive: false,
			urls: false,
			prune: false,
			limit: 0
		}
	});
	const args: string[] = [];
	for (const arg of options._) {
		if (typeof(arg) === 'number') {
			args.push(`${arg}`);
		} else {
			args.push(arg);
		}
	}

	if (options.help) {
		console.log(fmtHelp(skimmerHelpText, appName, version, releaseDate, releaseHash));
		Deno.exit(0);
	}

	if (options.license) {
		console.log(`${licenseText}`);
		Deno.exit(0);
	}

	if (options.version) {
		console.log(`${appName} ${version} ${releaseHash}\n`);
		Deno.exit(0);
	}

	const skimmer = new Skimmer(appName);
	if (! skimmer.run(options, args)) {
		console.error(skimmer.display_errors());
		Deno.exit(1);
	}
}

if (import.meta.main) await main();
