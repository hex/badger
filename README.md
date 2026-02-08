![GitHub](https://img.shields.io/github/license/hex/badger?style=flat-square)
![Language](https://img.shields.io/badge/language-C%23-blue?style=flat-square)
![Platform](https://img.shields.io/badge/platform-windows%20%7C%20macos%20%7C%20linux-lightgrey?style=flat-square)

<p style="text-align: center">
<img src="Assets/badger.png" title="Badger" alt="Badger">
</p>

# Badger

A cross-platform command-line tool that adds labels to your app icon.

Badger is powered by [SixLabors.ImageSharp](https://github.com/SixLabors/ImageSharp)

## Installation

### macOS/Linux

#### Homebrew

```sh
brew tap hex/tap
brew install badger
```

### Windows

#### Scoop

```sh
scoop bucket add badger https://github.com/hex/Badger
scoop install badger
```

## Options

| Option                 | Description                                                         | Default   |                                     Format                                     |
|:-----------------------|:--------------------------------------------------------------------|:----------|:------------------------------------------------------------------------------:|
| `--text`               | Text to be displayed on the badge                                   |           |                                                                                |
| `--icon`               | Path to the icon file / directory                                   |           |                                                                                |
| `--font-name`          | Name of the font to be used                                         | `Arial`   |                                                                                |
| `--height`             | Height as percentage                                                | `20`      |                                   `0 - 100`                                    |
| `--width`              | Width as percentage                                                 | `100`     |                                   `0 - 100`                                    |
| `--color`              | Background color                                                    | `#4096EE` |                                                                                |
| `--opacity`            | Opacity                                                             | `1`       |                                    `0 - 1`                                     |
| `--text-color`         | Text color                                                          | `#F9F7ED` |                                                                                |
| `--text-alignment`     | Text alignment                                                      | `center`  |                             `left, center, right`                              |
| `-r, --angle`          | Rotation angle                                                      | `0`       |                                   `0 - 360`                                    |
| `-x, --offsetx`        | X-axis offset                                                       | `0`       |                                                                                |
| `-y, --offsety`        | Y-axis offset                                                       | `0`       |                                                                                |
| `--badge-pivot`        | Badge pivot point                                                   | `bottom`  | `top, left, bottom, right, topLeft, topRight, bottomLeft, bottomRight, center` |
| `--horizontal-padding` | Text horizontal padding                                             | `5`       |                                                                                |
| `--vertical-padding`   | Text vertical padding                                               | `0`       |                                                                                |
| `--horizontal-pivot`   | Text horizontal pivot                                               | `center`  |                             `left, center, right`                              |
| `--vertical-pivot`     | Text vertical pivot                                                 | `center`  |                             `top, center, bottom`                              |
| `-o, --overwrite`      | Replace input icon. **WARNING**: This will overwrite the input icon | `false`   |                                                                                |

## Usage

### Examples

<p style="text-align: left">
<img src="Assets/ex1.png" alt="Badger" width="256">
</p>

```sh
badger --text ALPHA --icon icon.png --badge-height 25 --angle -45 --horizontal-padding 60 --offsetx 65 --offsety 65
```

<p style="text-align: left">
<img src="Assets/ex2.png" alt="Badger"  width="256">
</p>

```sh
badger --text BETA --icon icon.png --color "#FFFD88" --text-color "#C79811" --offsety -25
```

<p style="text-align: left">
<img src="Assets/ex3.png" alt="Badger"  width="256">
</p>

```sh
badger --text DEV --icon icon.png --width 50 --color "#363A3D" --text-color "#CDEB8B" --offsety -100 --badge-pivot bottomRight
```

```sh
Usage: badger [options...]

Options:
  --text <String>                 Set badge text (Required)
  --icon <String>                 Icon path.[.png | .jpg | .jpeg | .appiconset] (Required)
  --font-name <String>            Font name (Default: Arial)
  --width <Int32>                 Badge width in percentage. 0 - 100  (Default: 100)
  --height <Int32>                Badge height in percentage. 0 - 100  (Default: 20)
  --color <String>                Set badge background color with a hexadecimal color code (Default: #4096EE)
  --opacity <Single>              Badge opacity (Default: 1)
  --text-color <String>           Set badge text color with a hexadecimal color code (Default: #F9F7ED)
  --text-alignment <String>       Set badge text alignment. left | center | right (Default: center)
  -r, --angle <Int32>             Set badge rotation (Default: 0)
  -x, --offsetx <Int32>           Set badge x-axis offset (Default: 0)
  -y, --offsety <Int32>           Set badge y-axis offset (Default: 0)
  --badge-pivot <String>          Set badge pivot point. top | left | bottom | right | topLeft | topRight | bottomLeft | bottomRight (Default: bottomLeft)
  --horizontal-padding <Int32>    Set badge text horizontal padding (Default: 5)
  --vertical-padding <Int32>      Set badge text vertical padding (Default: 0)
  --horizontal-pivot <String>     Set badge text horizontal pivot. left | center | right (Default: center)
  --vertical-pivot <String>       Set badge text vertical pivot. top | center | bottom (Default: center)
  -o, --overwrite                 Replace input icon. WARNING: This will overwrite the input icon. (Optional)

Commands:
  help       Display help.
  version    Display version
```

## License

Badger is released under the MIT license. See [LICENSE](https://github.com/hex/badger/blob/master/LICENSE) for more
information.
