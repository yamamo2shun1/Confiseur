# C4NDY Confiseur

## Description
Command line tool to change the keymap of "[C4NDY KeyVLM and STK](https://github.com/yamamo2shun1/C4NDY)".

## Installation
The `go install` is available.
```
$ go install https://github.com/yamamo2shun1/Confiseur/cmd/confiseur@latest
```
You can also download binaries from the [Release page](https://github.com/yamamo2shun1/Confiseur/releases).

## How to use
```
-version
        Show the version of the tool installed.
        ex) confiseur -version
-check
        Show information on C4NDY KeyVLM/STK connected to PC/Mac.
        ex) confiseur -check
-list
        Show connected device name list.
        ex) confiseur -list
-id [int]
        Select connected device ID(ID can be checked in -check/-list).
        This option is available when using the following command options.
        If ID is not specified, 0 is the default.
-load
        Show the current key names of the keyboard.
        ex) confiseur -load
            confiseur -load -id 1
-remap
        Apply the keymap infomation from layouts.toml by default.
        ex) confiseur -remap
            confiseur -remap -id 0
-file [string]
        Specify .toml file to be read.
        This option is available when using the '-remap' option.
        ex) confiseur -remap -file layout_STK.toml
            confiseur -id 0 -remap -file layout_KeyVLM.toml
-save
        Save the keymap written by "remap" to the memory area
        ex) confiseur -save
            confiseur -id 0 -save
```

## Preparation to build
First, install the [Go language](https://go.dev/) development environment.  
Some preparation is required to build the code because we are using [go-hid](https://github.com/sstallion/go-hid).

### for Windows
Add "CGO_ENABLED=1" to the environment variable and install a C compiler such as [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) beforehand.

### for macOS/Linux
Add "CGO_ENABLED=1" to your shell configuration file, such as .zshrc.

## How to build from Source Code
```
> go build -o confiseur
```
