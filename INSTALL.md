Installation for development of **skimmer**
===========================================

**skimmer** Skimmer is a lightweight feed reader inspired by [newsboat](https://newsboat.org). skimmer is very minimal and deliberately lacks features.  Skimmer's best feature is what it doesn't do. Skimmer tries to do two things well.

1. Read a list of URLs, fetch the feeds and write the items to an SQLite 3 database
2. Display the items in the SQLite 3 database in reverse chronological order

That's it. That is skimmer secret power. It does only two things. There is no elaborate user interface beyond standard input, standard output and standard error found on POSIX type operating systems. Even if you invoke it in "interactive" mode your choices are limited, press enter and go to next item, press "n" and mark the item read, press "s" and save the item, press "q" and quit interactive mode.

Quick install with curl or irm
------------------------------

There is an experimental installer.sh script that can be run with the following command to install latest table release. This may work for macOS, Linux and if you’re using Windows with the Unix subsystem. This would be run from your shell (e.g. Terminal on macOS).

~~~shell
curl https://rsdoiel.github.io/skimmer/installer.sh | sh
~~~

This will install the programs included in skimmer in your `$HOME/bin` directory.

If you are running Windows 10 or 11 use the Powershell command below.

~~~ps1
irm https://rsdoiel.github.io/skimmer/installer.ps1 | iex
~~~

### If your are running macOS or Windows

You may get security warnings if you are using macOS or Windows. See the notes for the specific operating system you're using to fix issues.

- [INSTALL_NOTES_macOS.md](INSTALL_NOTES_macOS.md)
- [INSTALL_NOTES_Windows.md](INSTALL_NOTES_Windows.md)

Installing from source
----------------------

### Required software

- Go &gt;&#x3D; 1.25.0
- CMTools &gt;&#x3D; 0.0.40

### Steps

1. git clone https://github.com/rsdoiel/skimmer
2. Change directory into the `skimmer` directory
3. Make to build, test and install

~~~shell
git clone https://github.com/rsdoiel/skimmer
cd skimmer
make
make test
make install
~~~

