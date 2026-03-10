// ABOUTME: Entry point for the badger CLI tool.
// ABOUTME: Parses command-line flags and delegates to badge processing.
package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var version = "dev"

func main() {
	var cfg BadgeConfig
	var icon string
	var overwrite bool
	var showVersion bool

	flag.StringVar(&cfg.Text, "text", "", "Badge text")
	flag.StringVar(&icon, "icon", "", "Icon path (.png | .jpg | .jpeg | .appiconset)")
	flag.StringVar(&cfg.FontName, "font-name", "inter", "Font name (inter, roboto) or path to .ttf/.otf")
	flag.IntVar(&cfg.Width, "width", 100, "Badge width as percentage of icon width (0-100)")
	flag.IntVar(&cfg.Height, "height", 20, "Badge height as percentage of icon height (0-100)")
	flag.StringVar(&cfg.Color, "color", "#4096EE", "Badge background color (hex)")
	flag.Float64Var(&cfg.Opacity, "opacity", 1, "Badge opacity (0-1)")
	flag.StringVar(&cfg.TextColor, "text-color", "#F9F7ED", "Badge text color (hex)")
	flag.StringVar(&cfg.TextAlignment, "text-alignment", "center", "Text alignment: left | center | right")
	flag.IntVarP(&cfg.Angle, "angle", "r", 0, "Badge rotation angle (0-360)")
	flag.IntVarP(&cfg.OffsetX, "offsetx", "x", 0, "Badge x-axis offset")
	flag.IntVarP(&cfg.OffsetY, "offsety", "y", 0, "Badge y-axis offset")
	flag.StringVar(&cfg.BadgePivot, "badge-pivot", "bottomLeft", "Badge pivot: top | left | bottom | right | topLeft | topRight | bottomLeft | bottomRight | center")
	flag.IntVar(&cfg.HPadding, "horizontal-padding", 5, "Text horizontal padding")
	flag.IntVar(&cfg.VPadding, "vertical-padding", 0, "Text vertical padding")
	flag.StringVar(&cfg.HPivot, "horizontal-pivot", "center", "Text horizontal pivot: left | center | right")
	flag.StringVar(&cfg.VPivot, "vertical-pivot", "center", "Text vertical pivot: top | center | bottom")
	flag.BoolVar(&cfg.Uppercase, "uppercase", false, "Transform text to uppercase")
	flag.Float64Var(&cfg.LetterSpacing, "letter-spacing", 0, "Extra spacing between characters (pixels)")
	flag.StringVar(&cfg.OutlineColor, "text-outline-color", "", "Text outline color (hex, enables outline)")
	flag.IntVar(&cfg.OutlineWidth, "text-outline-width", 2, "Text outline width in pixels")
	flag.StringVar(&cfg.ShadowColor, "text-shadow-color", "", "Text shadow color (hex, enables shadow)")
	flag.IntVar(&cfg.ShadowX, "text-shadow-x", 2, "Text shadow X offset")
	flag.IntVar(&cfg.ShadowY, "text-shadow-y", 2, "Text shadow Y offset")
	flag.StringVar(&cfg.GlowColor, "text-glow-color", "", "Text glow color (hex, enables glow)")
	flag.IntVar(&cfg.GlowRadius, "text-glow-radius", 3, "Text glow radius in pixels")
	flag.BoolVar(&cfg.Emboss, "emboss", false, "Apply emboss effect to text")
	flag.BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite input icon (WARNING: destructive)")
	flag.BoolVar(&showVersion, "version", false, "Show version")

	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if cfg.Text == "" || icon == "" {
		fmt.Fprintln(os.Stderr, "Error: --text and --icon are required")
		flag.Usage()
		os.Exit(1)
	}

	info, err := os.Stat(icon)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if info.IsDir() {
		fmt.Printf("Adding badge to all icons in %s\n", icon)
		outDir := "badgerOutput"
		if overwrite {
			outDir = "" // not used in overwrite mode, but processDirectory handles it
		}

		if !overwrite {
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
				os.Exit(1)
			}
		}

		err = processDirectoryWithOverwrite(icon, outDir, cfg, overwrite)
	} else {
		if overwrite {
			err = processIconOverwrite(icon, cfg)
		} else {
			outDir := "badgerOutput"
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
				os.Exit(1)
			}
			err = processIcon(icon, outDir, cfg)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
