# KeyVLM_Configurator

## Description
Command line tool to change the keymap of "[C4NDY KeyVLM](https://github.com/yamamo2shun1/C4NDY)".

## How to use
Download the latest version from [here](https://github.com/yamamo2shun1/KeyVLM_Configurator/releases) and run it from Command Prompt (Windows) or Terminal.app (macOS).
```
> vlmconfig ver                  // Show the version of the tool installed
> vlmconfig check                // Show information on C4NDY KeyVLM connected to PC/Mac
> vlmconfig load                 // Show the current key names of the keyboard
> vlmconfig remap                // Write the keyboard with the keymap set in layouts.toml
> vlmconfig remap -f custom.toml // Write the keymap set in the specified .toml to the keydoad
> vlmconfig save                 // Save the keymap written by "remap" to the memory area
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
> go build -o vlmconfig
```
