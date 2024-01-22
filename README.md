# KeyVLM_Configurator

## Description
C4NDY KeyVLMのキーマップを変更するためのコマンドラインツールです。

## How to use
```
> vlmconfig ver                  // インストールされているツールのバージョンを表示します
> vlmconfig check                // PC/Macに接続しているC4NDY KeyVLMの情報を表示します
> vlmconfig load                 // キーボードの現在のキーマップを表示します
> vlmconfig remap                // キーボードにlayouts.tomlで設定したキーマップを書き込みます
> vlmconfig remap -f custom.toml // キードードに指定のtomlで設定したキーマップを書き込みます
> vlmconfig save                 // remapで書き込んだキーマップをメモリ領域に保存します
```

## How to build from Source Code
```
> go build -o kvmconfig
```
