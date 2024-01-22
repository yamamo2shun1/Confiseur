# KeyVLM_Configurator

## Description
C4NDY KeyVLMのキーマップを変更するためのコマンドラインツールです。

## How to use
```
> vlmconfig.exe ver                  // インストールされているツールのバージョンを表示します
> vlmconfig.exe check                // PC/Macに接続しているC4NDY KeyVLMの情報を表示します
> vlmconfig.exe load                 // キーボードの現在のキーマップを表示します
> vlmconfig.exe remap                // キーボードにlayouts.tomlで設定したキーマップを書き込みます
> vlmconfig.exe remap -f custom.toml // キードードに指定のtomlで設定したキーマップを書き込みます
> vlmconfig.exe save                 // remapで書き込んだキーマップをメモリ領域に保存します
```

## How to build from Source Code
```
> go build -o vlmconfig.exe
```
