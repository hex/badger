// ABOUTME: Tests for font embedding and loading.
// ABOUTME: Verifies bundled fonts load correctly and unknown fonts produce clear errors.
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadBundledFontInter(t *testing.T) {
	face, err := loadFont("inter", 12)
	if err != nil {
		t.Fatalf("loading bundled inter font: %v", err)
	}
	if face == nil {
		t.Fatal("expected non-nil font face")
	}
}

func TestLoadBundledFontRoboto(t *testing.T) {
	face, err := loadFont("roboto", 12)
	if err != nil {
		t.Fatalf("loading bundled roboto font: %v", err)
	}
	if face == nil {
		t.Fatal("expected non-nil font face")
	}
}

func TestLoadBundledFontCaseInsensitive(t *testing.T) {
	face, err := loadFont("INTER", 12)
	if err != nil {
		t.Fatalf("loading font with uppercase name: %v", err)
	}
	if face == nil {
		t.Fatal("expected non-nil font face")
	}
}

func TestLoadFontFromPath(t *testing.T) {
	// Use one of the bundled fonts as a file path test
	path := filepath.Join("fonts", "Inter-Bold.ttf")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("font file not found at", path)
	}
	face, err := loadFont(path, 16)
	if err != nil {
		t.Fatalf("loading font from path %s: %v", path, err)
	}
	if face == nil {
		t.Fatal("expected non-nil font face")
	}
}

func TestLoadUnknownFontReturnsError(t *testing.T) {
	_, err := loadFont("nonexistent-font", 12)
	if err == nil {
		t.Fatal("expected error for unknown font, got nil")
	}
}

func TestLoadFontDifferentSizes(t *testing.T) {
	sizes := []float64{8, 12, 24, 48, 72}
	for _, size := range sizes {
		face, err := loadFont("inter", size)
		if err != nil {
			t.Fatalf("loading inter at size %.0f: %v", size, err)
		}
		if face == nil {
			t.Fatalf("expected non-nil face at size %.0f", size)
		}
	}
}
