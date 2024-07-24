# KeyConfigurator

## Description
Command line tool to change the keymap of "[C4NDY KeyVLM and STK](https://github.com/yamamo2shun1/C4NDY)".

## How to use
Download the latest version from [here](https://github.com/yamamo2shun1/KeyConfigurator/releases) and run it from Command Prompt (Windows) or Terminal.app (macOS).
```
> keyconfig -version                 // Show the version of the tool installed
> keyconfig -check                   // Show information on C4NDY KeyVLM/STK connected to PC/Mac
> keyconfig -load                    // Show the current key names of the keyboard
> keyconfig -remap                   // Write the keyboard with the keymap set in layouts.toml
> keyconfig -remap -file custom.toml // Write the keymap set in the specified .toml to the keydoad
> keyconfig -save                    // Save the keymap written by "remap" to the memory area
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
