---
title: skimmer
abstract: "Skimmer is a lightweight feed reader inspired by [newsboat](https://newsboat.org). skimmer is very minimal and deliberately lacks features.  Skimmer&#x27;s best feature is what it doesn&#x27;t do. Skimmer tries to do two things well.

1. Read a list of URLs, fetch the feeds and write the items to an SQLite 3 database
2. Display the items in the SQLite 3 database in reverse chronological order

That&#x27;s it. That is skimmer secret power. It does only two things. There is no elaborate user interface beyond standard input, standard output and standard error found on POSIX type operating systems. Even if you invoke it in &quot;interactive&quot; mode your choices are limited, press enter and go to next item, press &quot;n&quot; and mark the item read, press &quot;s&quot; and save the item, press &quot;q&quot; and quit interactive mode."
authors:
  - family_name: Doiel
    given_name: R. S.
    id: https://orcid.org/0000-0003-0900-6903



repository_code: git+https://github.com/rsdoiel/skimmer
version: 0.0.23
license_url: https://spdx.org/licenses/AGPL-3.0-or-later
operating_system:
  - Windows
  - macOS
  - Linux
  - Raspberry Pi OS

programming_language:
  - Go

keywords:
  - RSS
  - Website genertor

date_released: 2025-08-17
---

About this software
===================

## skimmer 0.0.23

- Removed support for getpocket.com, Mozilla ended it.
- Add new **skim2html** tool, this supercedes skim2md for website generation
  - includes improved HTML structure around section olding article elements
  - configurable CSS and ES6 module support
  - configurable header, nav and footer support
- Updated **skimmer** to read a "skimmer.yaml" in the current directory 
  and setup the configuration.

### Authors

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>






Skimmer is a lightweight feed reader inspired by [newsboat](https://newsboat.org). skimmer is very minimal and deliberately lacks features.  Skimmer&#x27;s best feature is what it doesn&#x27;t do. Skimmer tries to do two things well.

1. Read a list of URLs, fetch the feeds and write the items to an SQLite 3 database
2. Display the items in the SQLite 3 database in reverse chronological order

That&#x27;s it. That is skimmer secret power. It does only two things. There is no elaborate user interface beyond standard input, standard output and standard error found on POSIX type operating systems. Even if you invoke it in &quot;interactive&quot; mode your choices are limited, press enter and go to next item, press &quot;n&quot; and mark the item read, press &quot;s&quot; and save the item, press &quot;q&quot; and quit interactive mode.

- License: <https://spdx.org/licenses/AGPL-3.0-or-later>
- GitHub: <git+https://github.com/rsdoiel/skimmer>
- Issues: <https://github.com/rsdoiel/skimmer/issues>

### Programming languages

- Go


### Operating Systems

- Windows
- macOS
- Linux
- Raspberry Pi OS


### Software Requirements

- Go >= 1.25.0


### Software Suggestions

- Git &gt;&#x3D; 2.3
- GNU Make &gt;&#x3D; 3.8
- Pandoc &gt;&#x3D; 3.1
- SQLite3 &gt;&#x3D; 3.43


