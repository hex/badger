![GitHub](https://img.shields.io/github/license/hex/badger?style=flat-square)
![Language](https://img.shields.io/badge/language-C%23-blue?style=flat-square)
![Platform](https://img.shields.io/badge/platform-windows%20%7C%20macos%20%7C%20linux-lightgrey?style=flat-square)

<p style="text-align: center">
<img src="Assets/badger.png" title="Badger" alt="Badger">
</p>


# Badger
A cross-platform command-line tool that adds labels to your app icon. 

Badger is powered by [ImageMagick](http://www.imagemagick.org/) through [Magick.NET](https://github.com/dlemstra/Magick.NET)
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
## Options
`-c, --color` Set badge background color (supports hexadecimal color codes and [color names](https://imagemagick.org/script/color.php#color_names))

`-t, --text-color` Set badge label text color (supports hexadecimal color codes and [color names](https://imagemagick.org/script/color.php#color_names))

`-l, --text-alignment` Set badge label text alignment, supported values: `top, left, bottom, right, topLeft, topRight, bottomLeft, bottomRight, center`

`-a, --angle` Set badge rotation

`-x, --offset-x` Set x-axis offset

`-y, --offset-y` Set y-axis offset

`-p, --angle` Set badge position, supported values: `top, left, bottom, right, topLeft, topRight, bottomLeft, bottomRight, center`

`-r, --replace` Set this to `true` if you want to overwrite the input file(s)
## Usage

```sh
$ badger Alpha icon.png
$ badger Beta /assets/icons --angle 15 --color '#6BBA70' --text-color '#F9F7ED' --replace true

Usage:
  badger [options] <text> <icon>

Arguments:
  <text>    Set label text
  <icon>    Set path to icon with format .png | .jpg | .jpeg | .appiconset

Options:
  -c, --color <color>                      Set badge color with a hexadecimal color code [default: #4096EE]
  -t, --text-color <text-color>            Set badge text color with a hexadecimal color code [default: #F9F7ED]
  -l, --text-alignment <text-alignment>    Set badge text alignment [default: center]
  -a, --angle <angle>                      Set badge rotation [default: 0]
  -x, --offset-x <offset-x>                Set badge x-axis offset [default: 0]
  -y, --offset-y <offset-y>                Set badge y-axis offset [default: 0]
  -p, --position <position>                Set badge position [default: bottom]
  -r, --replace                            Replace input icon [default: False]
  --version                                Show version information
  -?, -h, --help                           Show help and usage information
```
## License

Badger is released under the MIT license. See [LICENSE](https://github.com/hex/badger/blob/master/LICENSE) for more information.
