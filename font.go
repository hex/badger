// ABOUTME: Embeds and loads bundled fonts for badge text rendering.
// ABOUTME: Supports loading by name ("inter", "roboto") or by file path.
package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/Inter-Bold.ttf
var interBoldTTF []byte

//go:embed fonts/Roboto-Bold.ttf
var robotoBoldTTF []byte

var bundledFonts = map[string][]byte{
	"inter":  interBoldTTF,
	"roboto": robotoBoldTTF,
}

func loadFont(name string, size float64) (font.Face, error) {
	key := strings.ToLower(name)

	data, ok := bundledFonts[key]
	if !ok {
		var err error
		data, err = os.ReadFile(name)
		if err != nil {
			available := make([]string, 0, len(bundledFonts))
			for k := range bundledFonts {
				available = append(available, k)
			}
			return nil, fmt.Errorf("unknown font %q (bundled: %s; or provide a path to a .ttf/.otf file)", name, strings.Join(available, ", "))
		}
	}

	parsed, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing font %q: %w", name, err)
	}

	face, err := opentype.NewFace(parsed, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("creating font face %q at size %.1f: %w", name, size, err)
	}

	return face, nil
}
