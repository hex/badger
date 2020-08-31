# Badger
A command-line tool that adds labels to your app icon
## Installation
### macOS
#### Homebrew
```sh
brew tap hex/formulae
brew install badger
```
### Windows
#### Scoop
```sh
scoop bucket add badger https://github.com/hex/Badger
scoop install badger
```
## Usage

```sh
Usage:
  badger [options] <text> <icon>

Arguments:
  <text>    Set label text
  <icon>    Set path to icon with format .png | .jpg | .jpeg | .appiconset

Options:
  -c, --color <color>              Set badge color with a hexadecimal color code [default: goldenrod1]
  -t, --text-color <text-color>    Set badge text color with a hexadecimal color code [default: white]
  -a, --angle <angle>              Set badge rotation [default: 0]
  -p, --position <position>        Set badge position [default: bottom]
  -r, --replace                    Replace input icon [default: False]
  --version                        Show version information
  -?, -h, --help                   Show help and usage information
```
