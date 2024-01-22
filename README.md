# KeyVLM_Configurator

## Description
Command line tool to change the keymap of "[C4NDY KeyVLM](https://github.com/yamamo2shun1/C4NDY)".

## How to use
```
> vlmconfig.exe ver                  // Show the version of the tool installed
> vlmconfig.exe check                // Show information on C4NDY KeyVLM connected to PC/Mac
> vlmconfig.exe load                 // Show the current keymap(ScanCode) of the keyboard
> vlmconfig.exe remap                // Write the keyboard with the keymap set in layouts.toml
> vlmconfig.exe remap -f custom.toml // Write the keymap set in the specified .toml to the keydoad
> vlmconfig.exe save                 // Save the keymap written by "remap" to the memory area
```

## How to build from Source Code
```
> go build -o vlmconfig.exe
```
