# KeyConfigurator

## Description
Command line tool to change the keymap of "[C4NDY KeyVLM and STK](https://github.com/yamamo2shun1/C4NDY)".

## How to use
Download the latest version from [here](https://github.com/yamamo2shun1/KeyConfigurator/releases) and run it from Command Prompt (Windows) or Terminal.app (macOS).
```
-version
        Show the version of the tool installed.
        ex) keyconfig -version
-check
        Show information on C4NDY KeyVLM/STK connected to PC/Mac.
        ex) keyconfig -check
-list
        Show connected device name list.
        ex) keyconfig -list
-id [int]
        Select connected device ID(ID can be checked in -check/-list).
        This option is available when using the following command options.
        If ID is not specified, 0 is the default.
-load
        Show the current key names of the keyboard.
        ex) keyconfig -load
            keyconfig -load -id 1
-remap
        Apply the keymap infomation from layouts.toml by default.
        ex) keyconfig -remap
            keyconfig -remap -id 0
-file [string]
        Specify .toml file to be read.
        This option is available when using the '-remap' option.
        ex) keyconfig -remap -file layout_STK.toml
            keyconfig -id 0 -remap -file layout_KeyVLM.toml
-save
        Save the keymap written by "remap" to the memory area
        ex) keyconfig -save
            keyconfig -id 0 -save
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
> go build -o keyconfig
```
