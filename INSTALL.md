
Installation
============

*skimmer* is a command line program run from a shell like Bash. You can find compiled
version in the [releases](https://github.com/rsdoiel/skimmer/releases/latest) 

## Quick install with curl

The following curl command can be used to run the installer on most
POSIX systems. Programs are installed into `$HOME/bin`. `$HOME/bin` will
need to be in your path. From a shell (or terminal session) run the
following.

~~~
curl https://rsdoiel.github.io/skimmer/installer.sh | sh
~~~

## Compiled version

This is generalized instructions for a release. 

Compiled versions are available for Mac OS X (Intel and M1 processors, macOS-x86_64, macOS-arm64), 
Linux Intel and ARM (x86_64, aarch64, armv7l ), Windows 11 Intel/ARM (x86_64, aarch64), 
Raspberry Pi (Linux and arml7).


VERSION_NUMBER is a [semantic version number](http://semver.org/) (e.g. v0.1.2)


For all the released version go to the project page on GitHub and click latest release

>    https://github.com/rsdoiel/skimmer/releases/latest


| Platform     | Zip Filename                              |
|--------------|-------------------------------------------|
| Windows 11   | skimmer-VERSION_NUMBER-Windows-x86_64.zip |
| Windows 11   | skimmer-VERSION_NUMBER-Windows-arm64.zip  |
| Mac OS X     | skimmer-VERSION_NUMBER-macOS-x86_64.zip   |
| Mac OS X     | skimmer-VERSION_NUMBER-macOS-arm64.zip    |
| Linux/Intel  | skimmer-VERSION_NUMBER-Linux-x86_64.zip   |
| Linux/ARM 64 | skimmer-VERSION_NUMBER-Linux-aarch64.zip  |
| Linux/ARM 32 | skimmer-VERSION_NUMBER-Linux-armv7l.zip    | 


## The basic recipe

- Find the Zip file listed matching the architecture you're running and download it
    - Example: if you're on a Windows 11 desktop or laptop using an Intel style
      CPU and you would look for the name with "Windows-x86-64"
    - Example: if you're on a Windows 11 Surface tablet with a arm64 CPU or the
      "Windows ARM Developer Kit" you'd choose the Zip file with "Windows-arm64" in the name
- Download the zip file and unzip the file.  
- Copy the contents of the folder named "bin" to a folder that is in your path 
    - (e.g. "$HOME/bin" is common).
- Adjust your PATH if needed
    - (e.g. export PATH="$HOME/bin:$PATH")
- Test by displaying the version string

### macOS

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to $HOME/bin (or a folder in your path)
4. Make sure the new location in in our path
5. Test

Here's an example of the commands run in the Terminal App after downloading the 
zip file.

```shell
    cd Downloads/
    unzip skimmer-*-macOS-x86_64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    skimmer -version
```

### Windows 11 Intel

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to $HOME/bin (or a folder in your path)
4. Test

Here's an example of the commands run in from the Bash shell on Windows 10 after
downloading the zip file.

```shell
    cd Downloads/
    unzip skimmer-*-Windows-x86_64.zip
    .\bin\skimmer.exe -version
```

### Windows 11 ARM

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to $HOME/bin (or a folder in your path)
4. Test

Here's an example of the commands run in from the Bash shell on Windows 10 after
downloading the zip file.

```shell
    cd Downloads/
    unzip skimmer-*-Windows-arm64.zip
    .\bin\skimmer.exe -version
```


### Linux Intel

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to $HOME/bin (or a folder in your path)
4. Test

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd Downloads/
    unzip skimmer-*-Linux-x86_64.zip
    mkdir -p $HOME/bin
    cp -v bin/skimmer $HOME/bin/
    export PATH=$HOME/bin:$PATH
    skimmer -version
```

### Linux ARM64

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to $HOME/bin (or a folder in your path)
4. Test

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd Downloads/
    unzip skimmer-*-Linux-aarch64.zip
    mkdir -p $HOME/bin
    cp -v bin/skimmer $HOME/bin/
    export PATH=$HOME/bin:$PATH
    skimmer -version
```



### Raspberry Pi OS

Released version is for a Raspberry Pi 2 or later use (i.e. requires ARM 7 support).

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to $HOME/bin (or a folder in your path)
4. Test

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd Downloads/
    unzip skimmer-*-Linux-armv7l.zip
    mkdir -p $HOME/bin
    cp -v bin/skimmer $HOME/bin/
    export PATH=$HOME/bin:$PATH
    skimmer -version
```

### Windows 11

The general steps

1. Download the zip file into your "Downloads" folder with your web browser
2. Setup your path to where you will install the `.exe` file
3. Unzip the zip file
4. Copy the executable to the "bin" directory in to someplace where Windows cmd shell file find it
5. Test

```shell
    mkdir %userprofile%\bin
    set PATH=%PATH%;%userprofile$\bin
    powershell Expand-Archive Downloads\dataset-*-Windows-*.zip Dataset
    copy Dataset\bin\*.exe %userprofile%\bin\
    skimmer -version
```

Compiling from source
---------------------

_skimmer_ is "go get-able".  Use the "go get" command to download the dependent packages
as well as _skimmer_'s source code.

```shell
    go get -u github.com/rsdoiel/skimmer/...
```

### Requirements for compiling

1. Git
2. GNU Make
3. SQLite 3
4. Pandoc > 3
5. Go >= 1.21.1


### Compiling on a POSIX system

If you have all the required software (e.g. Git, GNU Make, Pandoc, SQLite3, Go)
you can clone the repository and then compile in the traditional POSIX manner.


```shell
    cd
    git clone https://github.com/rsdoiel/skimmer src/github.com/rsdoiel/skimmer
    cd src/github.com/rsdoiel/skimmer
    make
    make test
    make install
```

#### Compiling on a Windows machine

On a Windows box in the command shell these are the steps I would take

```shell
	cd %userprofile%
	set PATH=%PATH%;%userprofile%\bin
	go build cmd\skimmer\skimmer.go
	copy skimmer.exe %userprofile%\bin\
	skimmer -version
```

