# C4NDY Confiseur

## Description
Command line tool to change the keyboard settings of "[C4NDY KeyVLM and STK](https://github.com/yamamo2shun1/C4NDY)".

## Installation
The `go install` is available.

```shellscript
$ go install github.com/yamamo2shun1/Confiseur/cmd/confiseur@v0.12.0
```

## How to use
```Less
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
-remap [string]
        Apply the keyboard settings from specified toml file.
        ex) confiseur -remap examples/layout_STK.toml
            confiseur -remap examples/layout_KeyVLM.toml -id 0
-save
        Save the keymap written by "remap" to the memory area.
        ex) confiseur -save
            confiseur -id 0 -save
-led [int(0x000000-0xFFFFFF)]
        Set LED RGB value for checking color.
        ex) confiseur -led 0xFF0000 # red
            confiseur -id 0 -led 0x00FFFF # cyan
-intensity [float(0.0-1.0)]
        Set LED intensity.
        ex) confiseur -intensity 1.0
            confiseur -id 1 -intensity 0.5
-restart
        Restart the keyboard immediately.
        ex) confiseur -restart
            confiseur -restart -id 1
-factoryreset
        Reset all settings to factory defaults.
        ex) confiseur -factoryreset
            confiseur -id 0 -factoryreset
```

## Preparation to build
First, install the [Go language](https://go.dev/) development environment.  
Some preparation is required to build the code because we are using [go-hid](https://github.com/sstallion/go-hid).

### for Windows
Add "CGO_ENABLED=1" to the environment variable and install a C compiler such as [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) beforehand.

### for macOS/Linux
Add "CGO_ENABLED=1" to your shell configuration file, such as .zshrc.

## How to build from Source Code

```shellscript
$ go build -o confiseur
```
