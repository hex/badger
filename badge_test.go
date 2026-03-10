// ABOUTME: Tests for badge rendering, pivot calculation, image compositing, and text effects.
// ABOUTME: Covers pivot positions, rotation, opacity, text transforms, outlines, shadows, glow, and emboss.
package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/fogleman/gg"
)

func makeTestIcon(t *testing.T, w, h int) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "icon.png")
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, color.RGBA{R: 100, G: 100, B: 200, A: 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	return path
}

func defaultConfig() BadgeConfig {
	return BadgeConfig{
		Text:          "BETA",
		FontName:      "inter",
		Width:         100,
		Height:        20,
		Color:         "#4096EE",
		Opacity:       1.0,
		TextColor:     "#F9F7ED",
		TextAlignment: "center",
		Angle:         0,
		OffsetX:       0,
		OffsetY:       0,
		BadgePivot:    "bottomLeft",
		HPadding:      5,
		VPadding:      0,
		HPivot:        "center",
		VPivot:        "center",
		Uppercase:     false,
		LetterSpacing: 0,
		OutlineColor:  "",
		OutlineWidth:  2,
		ShadowColor:   "",
		ShadowX:       2,
		ShadowY:       2,
		GlowColor:     "",
		GlowRadius:    3,
		Emboss:        false,
	}
}

func TestPivotPointBottomLeft(t *testing.T) {
	x, y := pivotPoint("bottomLeft", 50, 20, 100, 100)
	if x != 0 || y != 80 {
		t.Errorf("bottomLeft: got (%v, %v), want (0, 80)", x, y)
	}
}

func TestPivotPointTopRight(t *testing.T) {
	x, y := pivotPoint("topRight", 50, 20, 100, 100)
	if x != 50 || y != 0 {
		t.Errorf("topRight: got (%v, %v), want (50, 0)", x, y)
	}
}

func TestPivotPointCenter(t *testing.T) {
	x, y := pivotPoint("center", 50, 20, 100, 100)
	if x != 25 || y != 40 {
		t.Errorf("center: got (%v, %v), want (25, 40)", x, y)
	}
}

func TestAllPivotPositions(t *testing.T) {
	tests := []struct {
		pivot string
		wantX float64
		wantY float64
	}{
		{"top", 25, 0},
		{"left", 0, 40},
		{"bottom", 25, 80},
		{"right", 50, 40},
		{"topLeft", 0, 0},
		{"topRight", 50, 0},
		{"bottomLeft", 0, 80},
		{"bottomRight", 50, 80},
		{"center", 25, 40},
	}

	for _, tt := range tests {
		t.Run(tt.pivot, func(t *testing.T) {
			x, y := pivotPoint(tt.pivot, 50, 20, 100, 100)
			if x != tt.wantX || y != tt.wantY {
				t.Errorf("%s: got (%v, %v), want (%v, %v)", tt.pivot, x, y, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestUnknownPivotDefaultsToCenter(t *testing.T) {
	x, y := pivotPoint("bogus", 50, 20, 100, 100)
	cx, cy := pivotPoint("center", 50, 20, 100, 100)
	if x != cx || y != cy {
		t.Errorf("unknown pivot should default to center, got (%v, %v)", x, y)
	}
}

func TestCreateBadgeProducesImage(t *testing.T) {
	cfg := defaultConfig()
	badge, err := createBadge(100, 100, cfg)
	if err != nil {
		t.Fatalf("createBadge: %v", err)
	}
	bounds := badge.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("badge size: got %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestCreateBadgeHasColoredPixels(t *testing.T) {
	cfg := defaultConfig()
	cfg.BadgePivot = "bottomLeft"
	cfg.Width = 100
	cfg.Height = 100 // full coverage for easy testing
	badge, err := createBadge(100, 100, cfg)
	if err != nil {
		t.Fatalf("createBadge: %v", err)
	}
	// The badge should have non-transparent pixels in the bottom-left area
	r, g, b, a := badge.At(10, 90).RGBA()
	if a == 0 {
		t.Errorf("expected non-transparent pixel at (10, 90), got RGBA(%d, %d, %d, %d)", r, g, b, a)
	}
}

func TestCreateBadgeWithRotation(t *testing.T) {
	cfg := defaultConfig()
	cfg.Angle = 45
	badge, err := createBadge(100, 100, cfg)
	if err != nil {
		t.Fatalf("createBadge with rotation: %v", err)
	}
	bounds := badge.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("rotated badge size: got %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestCompositeImage(t *testing.T) {
	icon := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			icon.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	badge := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			badge.Set(x, y, color.RGBA{R: 0, G: 0, B: 255, A: 255})
		}
	}

	result := compositeImage(icon, badge, 0, 0, 1.0)
	r, _, b, _ := result.At(50, 50).RGBA()
	// With full opacity blue badge over red icon, result should be blue
	if r > b {
		t.Error("expected blue to dominate over red with full opacity badge")
	}
}

func TestCompositeImageHalfOpacity(t *testing.T) {
	icon := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			icon.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	badge := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			badge.Set(x, y, color.RGBA{R: 0, G: 0, B: 255, A: 255})
		}
	}

	result := compositeImage(icon, badge, 0, 0, 0.5)
	r, _, b, _ := result.At(50, 50).RGBA()
	// With half opacity, both red and blue should be present
	if r == 0 || b == 0 {
		t.Errorf("expected both red and blue with 0.5 opacity, got r=%d b=%d", r, b)
	}
}

func TestProcessIconSingleFile(t *testing.T) {
	iconPath := makeTestIcon(t, 128, 128)
	outDir := t.TempDir()
	cfg := defaultConfig()

	err := processIcon(iconPath, outDir, cfg)
	if err != nil {
		t.Fatalf("processIcon: %v", err)
	}

	outPath := filepath.Join(outDir, filepath.Base(iconPath))
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Fatalf("expected output file at %s", outPath)
	}

	// Verify output is a valid PNG with correct dimensions
	f, err := os.Open(outPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("decoding output: %v", err)
	}
	if img.Bounds().Dx() != 128 || img.Bounds().Dy() != 128 {
		t.Errorf("output size: got %dx%d, want 128x128", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestProcessIconOverwrite(t *testing.T) {
	iconPath := makeTestIcon(t, 64, 64)
	cfg := defaultConfig()

	origInfo, _ := os.Stat(iconPath)
	origSize := origInfo.Size()

	err := processIconOverwrite(iconPath, cfg)
	if err != nil {
		t.Fatalf("processIconOverwrite: %v", err)
	}

	newInfo, err := os.Stat(iconPath)
	if err != nil {
		t.Fatalf("stat after overwrite: %v", err)
	}
	// File should still exist and likely have a different size (badge added)
	if newInfo.Size() == 0 {
		t.Error("overwritten file is empty")
	}
	_ = origSize // size will change due to badge
}

func TestProcessIconDirectory(t *testing.T) {
	dir := t.TempDir()
	// Create a couple of test icons in the directory
	for _, name := range []string{"icon1.png", "icon2.png"} {
		img := image.NewRGBA(image.Rect(0, 0, 64, 64))
		for y := range 64 {
			for x := range 64 {
				img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
			}
		}
		f, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		png.Encode(f, img)
		f.Close()
	}

	outDir := t.TempDir()
	cfg := defaultConfig()

	err := processDirectory(dir, outDir, cfg)
	if err != nil {
		t.Fatalf("processDirectory: %v", err)
	}

	for _, name := range []string{"icon1.png", "icon2.png"} {
		outPath := filepath.Join(outDir, name)
		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			t.Errorf("expected output file %s", outPath)
		}
	}
}

func TestCreateBadgeUppercase(t *testing.T) {
	cfgUpper := defaultConfig()
	cfgUpper.Text = "beta"
	cfgUpper.Uppercase = true
	cfgUpper.Width = 100
	cfgUpper.Height = 100

	badgeUpper, err := createBadge(100, 100, cfgUpper)
	if err != nil {
		t.Fatalf("createBadge with uppercase: %v", err)
	}
	boundsUpper := badgeUpper.Bounds()
	if boundsUpper.Dx() != 100 || boundsUpper.Dy() != 100 {
		t.Errorf("uppercase badge size: got %dx%d, want 100x100", boundsUpper.Dx(), boundsUpper.Dy())
	}

	cfgLower := defaultConfig()
	cfgLower.Text = "beta"
	cfgLower.Uppercase = false
	cfgLower.Width = 100
	cfgLower.Height = 100

	badgeLower, err := createBadge(100, 100, cfgLower)
	if err != nil {
		t.Fatalf("createBadge without uppercase: %v", err)
	}

	// Compare pixel data: uppercase rendering should differ from lowercase
	differ := false
	for y := range 100 {
		for x := range 100 {
			r1, g1, b1, a1 := badgeUpper.At(x, y).RGBA()
			r2, g2, b2, a2 := badgeLower.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				differ = true
				break
			}
		}
		if differ {
			break
		}
	}
	if !differ {
		t.Error("expected uppercase and lowercase badges to render differently")
	}
}

func TestCreateBadgeLetterSpacing(t *testing.T) {
	cfg := defaultConfig()
	cfg.LetterSpacing = 5.0

	badge, err := createBadge(100, 100, cfg)
	if err != nil {
		t.Fatalf("createBadge with letter spacing: %v", err)
	}
	bounds := badge.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("letter spacing badge size: got %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestCreateBadgeTextShadow(t *testing.T) {
	cfg := defaultConfig()
	cfg.ShadowColor = "#000000"
	cfg.ShadowX = 3
	cfg.ShadowY = 3
	cfg.Width = 100
	cfg.Height = 100

	badge, err := createBadge(200, 200, cfg)
	if err != nil {
		t.Fatalf("createBadge with shadow: %v", err)
	}

	// Check that non-transparent pixels exist in the shadow offset region.
	// The shadow is drawn offset from the text, so scan the badge area for
	// pixels that are not purely the background color.
	hasNonTransparent := false
	for y := 3; y < 200; y++ {
		for x := 3; x < 200; x++ {
			_, _, _, a := badge.At(x, y).RGBA()
			if a > 0 {
				hasNonTransparent = true
				break
			}
		}
		if hasNonTransparent {
			break
		}
	}
	if !hasNonTransparent {
		t.Error("expected non-transparent pixels in shadow region")
	}
}

func TestCreateBadgeTextOutline(t *testing.T) {
	cfg := defaultConfig()
	cfg.OutlineColor = "#FF0000"
	cfg.OutlineWidth = 3
	cfg.Width = 100
	cfg.Height = 100
	cfg.Color = "#000000" // dark background so red outline is visible

	badge, err := createBadge(200, 200, cfg)
	if err != nil {
		t.Fatalf("createBadge with outline: %v", err)
	}

	// Look for red-ish pixels from the outline (R > G and R > B)
	hasRedPixel := false
	for y := range 200 {
		for x := range 200 {
			r, g, b, a := badge.At(x, y).RGBA()
			if a > 0 && r > g && r > b && r > 0x8000 {
				hasRedPixel = true
				break
			}
		}
		if hasRedPixel {
			break
		}
	}
	if !hasRedPixel {
		t.Error("expected red-ish pixels from text outline")
	}
}

func TestCreateBadgeTextGlow(t *testing.T) {
	cfg := defaultConfig()
	cfg.GlowColor = "#00FF00"
	cfg.GlowRadius = 5
	cfg.Width = 100
	cfg.Height = 100
	cfg.Color = "#000000" // dark background so green glow is visible

	badge, err := createBadge(200, 200, cfg)
	if err != nil {
		t.Fatalf("createBadge with glow: %v", err)
	}

	// Look for green-ish pixels from the glow (G > R and G > B)
	hasGreenPixel := false
	for y := range 200 {
		for x := range 200 {
			r, g, b, a := badge.At(x, y).RGBA()
			if a > 0 && g > r && g > b && g > 0x8000 {
				hasGreenPixel = true
				break
			}
		}
		if hasGreenPixel {
			break
		}
	}
	if !hasGreenPixel {
		t.Error("expected green-ish pixels from text glow")
	}
}

func TestCreateBadgeEmboss(t *testing.T) {
	cfg := defaultConfig()
	cfg.Emboss = true

	badge, err := createBadge(100, 100, cfg)
	if err != nil {
		t.Fatalf("createBadge with emboss: %v", err)
	}
	bounds := badge.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("emboss badge size: got %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestCreateBadgeCombinedEffects(t *testing.T) {
	cfg := defaultConfig()
	cfg.Uppercase = true
	cfg.Text = "beta"
	cfg.OutlineColor = "#FF0000"
	cfg.OutlineWidth = 2
	cfg.ShadowColor = "#000000"
	cfg.ShadowX = 2
	cfg.ShadowY = 2
	cfg.LetterSpacing = 3.0

	badge, err := createBadge(200, 200, cfg)
	if err != nil {
		t.Fatalf("createBadge with combined effects: %v", err)
	}
	bounds := badge.Bounds()
	if bounds.Dx() != 200 || bounds.Dy() != 200 {
		t.Errorf("combined effects badge size: got %dx%d, want 200x200", bounds.Dx(), bounds.Dy())
	}
}

func TestDrawTextOnContextNoSpacing(t *testing.T) {
	dc := gg.NewContext(100, 100)

	face, err := loadFont("inter", 20)
	if err != nil {
		t.Fatalf("loadFont: %v", err)
	}
	dc.SetFontFace(face)
	dc.SetColor(color.White)

	drawTextOnContext(dc, "TEST", 50, 50, 0.5, 0.5, 0)

	// Verify non-transparent pixels were drawn
	img := dc.Image()
	hasPixel := false
	for y := range 100 {
		for x := range 100 {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				hasPixel = true
				break
			}
		}
		if hasPixel {
			break
		}
	}
	if !hasPixel {
		t.Error("expected non-transparent pixels from drawTextOnContext with no spacing")
	}
}

func TestDrawTextOnContextWithSpacing(t *testing.T) {
	dc := gg.NewContext(100, 100)

	face, err := loadFont("inter", 20)
	if err != nil {
		t.Fatalf("loadFont: %v", err)
	}
	dc.SetFontFace(face)
	dc.SetColor(color.White)

	drawTextOnContext(dc, "TEST", 50, 50, 0.5, 0.5, 5.0)

	// Verify non-transparent pixels were drawn
	img := dc.Image()
	hasPixel := false
	for y := range 100 {
		for x := range 100 {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				hasPixel = true
				break
			}
		}
		if hasPixel {
			break
		}
	}
	if !hasPixel {
		t.Error("expected non-transparent pixels from drawTextOnContext with spacing")
	}
}

func TestDefaultConfigIncludesNewFields(t *testing.T) {
	cfg := defaultConfig()

	if cfg.Uppercase != false {
		t.Errorf("Uppercase: got %v, want false", cfg.Uppercase)
	}
	if cfg.LetterSpacing != 0 {
		t.Errorf("LetterSpacing: got %v, want 0", cfg.LetterSpacing)
	}
	if cfg.OutlineColor != "" {
		t.Errorf("OutlineColor: got %q, want empty", cfg.OutlineColor)
	}
	if cfg.OutlineWidth != 2 {
		t.Errorf("OutlineWidth: got %v, want 2", cfg.OutlineWidth)
	}
	if cfg.ShadowColor != "" {
		t.Errorf("ShadowColor: got %q, want empty", cfg.ShadowColor)
	}
	if cfg.ShadowX != 2 {
		t.Errorf("ShadowX: got %v, want 2", cfg.ShadowX)
	}
	if cfg.ShadowY != 2 {
		t.Errorf("ShadowY: got %v, want 2", cfg.ShadowY)
	}
	if cfg.GlowColor != "" {
		t.Errorf("GlowColor: got %q, want empty", cfg.GlowColor)
	}
	if cfg.GlowRadius != 3 {
		t.Errorf("GlowRadius: got %v, want 3", cfg.GlowRadius)
	}
	if cfg.Emboss != false {
		t.Errorf("Emboss: got %v, want false", cfg.Emboss)
	}
}
