// ABOUTME: Renders text badges and composites them onto app icon images.
// ABOUTME: Supports configurable positioning, colors, fonts, and rotation.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

type BadgeConfig struct {
	Text          string
	FontName      string
	Width         int // percentage of icon width (0-100)
	Height        int // percentage of icon height (0-100)
	Color         string
	Opacity       float64
	TextColor     string
	TextAlignment string
	Angle         int
	OffsetX       int
	OffsetY       int
	BadgePivot    string
	HPadding      int
	VPadding      int
	HPivot        string
	VPivot        string
	Uppercase     bool
	LetterSpacing float64
	OutlineColor  string
	OutlineWidth  int
	ShadowColor   string
	ShadowX       int
	ShadowY       int
	GlowColor     string
	GlowRadius    int
	Emboss        bool
}

func pivotPoint(pivot string, badgeW, badgeH float64, iconW, iconH int) (float64, float64) {
	iw := float64(iconW)
	ih := float64(iconH)

	switch pivot {
	case "top":
		return (iw - badgeW) / 2, 0
	case "left":
		return 0, (ih - badgeH) / 2
	case "bottom":
		return (iw - badgeW) / 2, ih - badgeH
	case "right":
		return iw - badgeW, (ih - badgeH) / 2
	case "topLeft":
		return 0, 0
	case "topRight":
		return iw - badgeW, 0
	case "bottomLeft":
		return 0, ih - badgeH
	case "bottomRight":
		return iw - badgeW, ih - badgeH
	case "center":
		return (iw - badgeW) / 2, (ih - badgeH) / 2
	default:
		return (iw - badgeW) / 2, (ih - badgeH) / 2
	}
}

func parseHexColor(hex string) color.RGBA {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b uint8
	switch len(hex) {
	case 6:
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	case 3:
		fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r *= 17
		g *= 17
		b *= 17
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func textAnchor(pivot string) float64 {
	switch strings.ToLower(pivot) {
	case "left", "top":
		return 0
	case "center":
		return 0.5
	case "right", "bottom":
		return 1
	default:
		return 0.5
	}
}

func drawTextOnContext(dc *gg.Context, text string, x, y, ax, ay, letterSpacing float64) {
	if letterSpacing <= 0 {
		dc.DrawStringAnchored(text, x, y, ax, ay)
		return
	}
	runes := []rune(text)
	totalWidth := 0.0
	charWidths := make([]float64, len(runes))
	for i, ch := range runes {
		w, _ := dc.MeasureString(string(ch))
		charWidths[i] = w
		totalWidth += w
	}
	totalWidth += letterSpacing * float64(len(runes)-1)

	startX := x - ax*totalWidth
	currentX := startX
	for i, ch := range runes {
		dc.DrawStringAnchored(string(ch), currentX, y, 0, ay)
		currentX += charWidths[i] + letterSpacing
	}
}

func createBadge(iconW, iconH int, cfg BadgeConfig) (image.Image, error) {
	text := cfg.Text
	if cfg.Uppercase {
		text = strings.ToUpper(text)
	}

	badgeW := float64(iconW) * float64(cfg.Width) / 100
	badgeH := float64(iconH) * float64(cfg.Height) / 100
	px, py := pivotPoint(cfg.BadgePivot, badgeW, badgeH, iconW, iconH)

	dc := gg.NewContext(iconW, iconH)

	// Draw badge background
	bgColor := parseHexColor(cfg.Color)
	dc.SetColor(bgColor)
	dc.DrawRectangle(px, py, badgeW, badgeH)
	dc.Fill()

	// Load font at a base size, measure, then scale to fit
	baseFace, err := loadFont(cfg.FontName, 10)
	if err != nil {
		return nil, fmt.Errorf("loading font: %w", err)
	}
	dc.SetFontFace(baseFace)
	tw, th := dc.MeasureString(text)

	hPad := float64(cfg.HPadding)
	vPad := float64(cfg.VPadding)
	scaleX := badgeW / (tw + hPad)
	scaleY := badgeH / (th + vPad)
	scale := math.Min(scaleX, scaleY)
	targetSize := scale * 10 // base size was 10

	scaledFace, err := loadFont(cfg.FontName, targetSize)
	if err != nil {
		return nil, fmt.Errorf("loading scaled font: %w", err)
	}
	dc.SetFontFace(scaledFace)

	// Draw text with corrected vertical centering.
	// gg's DrawStringAnchored with ay=0.5 centers on full line height
	// (including descender space), which pushes uppercase text too low.
	// We compute the desired baseline from font metrics, then convert back
	// to gg's expected input coordinate (which uses ay=0.5 internally).
	textColor := parseHexColor(cfg.TextColor)
	dc.SetColor(textColor)

	ax := textAnchor(cfg.HPivot)
	centerX := px + badgeW/2
	centerY := py + badgeH/2

	metrics := scaledFace.Metrics()
	ascent := float64(metrics.Ascent) / 64
	descent := float64(metrics.Descent) / 64
	fHeight := float64(metrics.Height) / 64

	// Compute desired baseline position, then convert to gg's input Y.
	// gg with ay=0.5: baseline = inputY + 0.5*fontHeight
	// So: inputY = desiredBaseline - 0.5*fontHeight
	var desiredBaseline float64
	switch strings.ToLower(cfg.VPivot) {
	case "top":
		desiredBaseline = py + ascent
	case "bottom":
		desiredBaseline = py + badgeH - descent
	default: // center — visual center of cap-height text at badge center
		desiredBaseline = centerY + (ascent-descent)/2
	}
	drawY := desiredBaseline - 0.5*fHeight

	// Glow effect: concentric rings with decreasing alpha
	if cfg.GlowColor != "" {
		glowBase := parseHexColor(cfg.GlowColor)
		for r := cfg.GlowRadius; r >= 1; r-- {
			alpha := uint8(80 / r)
			glowColor := color.RGBA{R: glowBase.R, G: glowBase.G, B: glowBase.B, A: alpha}
			dc.SetColor(glowColor)
			for dx := -r; dx <= r; dx++ {
				for dy := -r; dy <= r; dy++ {
					if dx*dx+dy*dy <= r*r {
						drawTextOnContext(dc, text, centerX+float64(dx), drawY+float64(dy), ax, 0.5, cfg.LetterSpacing)
					}
				}
			}
		}
	}

	// Shadow effect
	if cfg.ShadowColor != "" {
		shadowColor := parseHexColor(cfg.ShadowColor)
		dc.SetColor(shadowColor)
		drawTextOnContext(dc, text, centerX+float64(cfg.ShadowX), drawY+float64(cfg.ShadowY), ax, 0.5, cfg.LetterSpacing)
	}

	// Outline effect
	if cfg.OutlineColor != "" {
		outlineColor := parseHexColor(cfg.OutlineColor)
		dc.SetColor(outlineColor)
		w := cfg.OutlineWidth
		for dx := -w; dx <= w; dx++ {
			for dy := -w; dy <= w; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				if dx*dx+dy*dy <= w*w {
					drawTextOnContext(dc, text, centerX+float64(dx), drawY+float64(dy), ax, 0.5, cfg.LetterSpacing)
				}
			}
		}
	}

	// Emboss effect
	if cfg.Emboss {
		dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 80})
		drawTextOnContext(dc, text, centerX-1, drawY-1, ax, 0.5, cfg.LetterSpacing)
		dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 80})
		drawTextOnContext(dc, text, centerX+1, drawY+1, ax, 0.5, cfg.LetterSpacing)
	}

	// Main text
	dc.SetColor(textColor)
	drawTextOnContext(dc, text, centerX, drawY, ax, 0.5, cfg.LetterSpacing)

	badgeImg := dc.Image()

	// Apply rotation if needed
	if cfg.Angle != 0 {
		dc2 := gg.NewContext(iconW, iconH)
		dc2.RotateAbout(gg.Radians(float64(cfg.Angle)), float64(iconW)/2, float64(iconH)/2)
		dc2.DrawImage(badgeImg, 0, 0)
		badgeImg = dc2.Image()
	}

	return badgeImg, nil
}

func compositeImage(icon *image.RGBA, badge image.Image, offsetX, offsetY int, opacity float64) *image.RGBA {
	result := image.NewRGBA(icon.Bounds())
	draw.Draw(result, icon.Bounds(), icon, image.Point{}, draw.Src)

	if opacity >= 1.0 {
		draw.DrawMask(result, badge.Bounds().Add(image.Pt(offsetX, offsetY)),
			badge, image.Point{}, nil, image.Point{}, draw.Over)
	} else {
		alphaMask := image.NewUniform(color.Alpha{A: uint8(opacity * 255)})
		draw.DrawMask(result, badge.Bounds().Add(image.Pt(offsetX, offsetY)),
			badge, image.Point{}, alphaMask, image.Point{}, draw.Over)
	}

	// Clip result to original icon's alpha channel so the badge
	// doesn't bleed outside rounded corners or transparent areas.
	iconPix := icon.Pix
	resultPix := result.Pix
	for i := 3; i < len(resultPix); i += 4 {
		a := iconPix[i]
		if a == 0 {
			resultPix[i-3] = 0
			resultPix[i-2] = 0
			resultPix[i-1] = 0
			resultPix[i] = 0
		} else if a < 255 {
			resultPix[i-3] = uint8(uint16(resultPix[i-3]) * uint16(a) / 255)
			resultPix[i-2] = uint8(uint16(resultPix[i-2]) * uint16(a) / 255)
			resultPix[i-1] = uint8(uint16(resultPix[i-1]) * uint16(a) / 255)
			resultPix[i] = uint8(uint16(resultPix[i]) * uint16(a) / 255)
		}
	}

	return result
}

func loadImage(path string) (*image.RGBA, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return nil, "", fmt.Errorf("decoding %s: %w", path, err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	return rgba, format, nil
}

func saveImage(img image.Image, path, format string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating %s: %w", path, err)
	}
	defer f.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
	default:
		return png.Encode(f, img)
	}
}

func processIcon(iconPath, outDir string, cfg BadgeConfig) error {
	icon, format, err := loadImage(iconPath)
	if err != nil {
		return err
	}

	badge, err := createBadge(icon.Bounds().Dx(), icon.Bounds().Dy(), cfg)
	if err != nil {
		return err
	}

	result := compositeImage(icon, badge, cfg.OffsetX, cfg.OffsetY, cfg.Opacity)

	outPath := filepath.Join(outDir, filepath.Base(iconPath))
	fmt.Printf("Writing to: %s\n", outPath)
	return saveImage(result, outPath, format)
}

func processIconOverwrite(iconPath string, cfg BadgeConfig) error {
	icon, format, err := loadImage(iconPath)
	if err != nil {
		return err
	}

	badge, err := createBadge(icon.Bounds().Dx(), icon.Bounds().Dy(), cfg)
	if err != nil {
		return err
	}

	result := compositeImage(icon, badge, cfg.OffsetX, cfg.OffsetY, cfg.Opacity)

	fmt.Printf("Writing to: %s\n", iconPath)
	return saveImage(result, iconPath, format)
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}

func processDirectory(dir, outDir string, cfg BadgeConfig) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !isImageFile(path) {
			return nil
		}
		return processIcon(path, outDir, cfg)
	})
}

func processDirectoryWithOverwrite(dir, outDir string, cfg BadgeConfig, overwrite bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !isImageFile(path) {
			return nil
		}
		if overwrite {
			return processIconOverwrite(path, cfg)
		}
		return processIcon(path, outDir, cfg)
	})
}
