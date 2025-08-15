// This defines my expectations of the parse function provide by JSON, yaml and xml.
type ObjectParseType = (arg1: string, arg2?: any) => {[key: string]: any} | unknown;

// FeedSource describes the source of a feed. It includes the URL,
// an optional label, user agent string.
export interface FeedSourceInterface {
	url: string;
    label: string;
    userAgent: string;

    fromObject: (args1: {[key: string]: any}) => boolean;
    parseWith: (args1: string, fn: ObjectParseType) => boolean;
}

export class FeedSource implements FeedSourceInterface {
    url: string = '';
    label: string = '';
    userAgent: string = '';

    fromObject(o: {[key: string]: any}): boolean {
        if (o.url !== undefined && o.url !== '') {
            this.url = o.url;
        }
        if (o.label !== undefined && o.label !== '') {
            this.label = o.label;
        }
        if (o.userAgent !== undefined && o.userAgent !== '') {
            this.userAgent = o.userAgent;
        }
        return true;
    };

    parseWith(s: string, fn: ObjectParseType): boolean {
        return this.fromObject(fn(s) as unknown as {[key: string]: any});
    };
}

// ParseURLList takes a filename and byte slice source, parses the contents
// returning a map of urls to labels and an error value.
export async function ParseURLList(fName: string, src:string ): Promise<{[key: string]: FeedSource}> {
	let urls : {[key: string]: FeedSource} = {};
    /*
	// Parse the url value collecting our keys and values
	s := bufio.NewScanner(bytes.NewBuffer(src))
	key, val, userAgent := "", "", ""
	line := 1
	for s.Scan() {
		txt := strings.TrimSpace(s.Text())
		if strings.HasPrefix(txt, "#") {
			txt = ""
		}
		if txt != "" {
			parts := strings.SplitN(txt, ` "`, 2)
			switch len(parts) {
			case 1:
				key, val, userAgent = parts[0], "", ""
			case 2:
				key, val, userAgent = parts[0], parts[1], ""
				pos := strings.LastIndex(val, `"`)
				if pos > -1 {
					if len(val) > pos {
						userAgent = strings.TrimSpace(val[pos+1:])
					}
					val = strings.TrimSpace(val[0:pos])
				}
				val = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(val, `~`), `"`))
			}
			urls[key] = &FeedSource{
				Url: key,
				Label: val,
				UserAgent: userAgent,
			}
		}
		line++
	}
	return urls, nil
    */
   return urls;
}
