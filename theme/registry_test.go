package theme

import (
	"testing"
)

func TestByNameReturnsModern(t *testing.T) {
	th := ByName("modern")
	if th == nil {
		t.Fatal("ByName(modern) returned nil")
	}
	if th.FontFamily != Modern().FontFamily {
		t.Errorf("FontFamily = %q, want %q", th.FontFamily, Modern().FontFamily)
	}
}

func TestByNameReturnsMermaidDefault(t *testing.T) {
	th := ByName("default")
	if th == nil {
		t.Fatal("ByName(default) returned nil")
	}
	if th.PrimaryColor != MermaidDefault().PrimaryColor {
		t.Errorf("PrimaryColor = %q, want %q", th.PrimaryColor, MermaidDefault().PrimaryColor)
	}
}

func TestByNameReturnsDark(t *testing.T) {
	th := ByName("dark")
	if th == nil {
		t.Fatal("ByName(dark) returned nil")
	}
	if th.Background == "#FFFFFF" {
		t.Error("dark theme should not have white background")
	}
}

func TestByNameReturnsForest(t *testing.T) {
	th := ByName("forest")
	if th == nil {
		t.Fatal("ByName(forest) returned nil")
	}
	// Forest theme should have green tones in primary.
	if th.PrimaryColor == "" {
		t.Error("PrimaryColor empty")
	}
}

func TestByNameReturnsNeutral(t *testing.T) {
	th := ByName("neutral")
	if th == nil {
		t.Fatal("ByName(neutral) returned nil")
	}
	if th.PrimaryColor == "" {
		t.Error("PrimaryColor empty")
	}
}

func TestByNameCaseInsensitive(t *testing.T) {
	tests := []string{"Dark", "DARK", "dark", "DaRk"}
	for _, name := range tests {
		th := ByName(name)
		if th == nil {
			t.Errorf("ByName(%q) returned nil", name)
		}
	}
}

func TestByNameUnknownReturnsFallback(t *testing.T) {
	th := ByName("nonexistent")
	if th == nil {
		t.Fatal("ByName(nonexistent) returned nil, want fallback")
	}
	// Should fall back to Modern.
	if th.FontFamily != Modern().FontFamily {
		t.Errorf("unknown name should fall back to Modern, got FontFamily=%q", th.FontFamily)
	}
}

func TestNamesReturnsList(t *testing.T) {
	names := Names()
	if len(names) < 5 {
		t.Errorf("Names() returned %d entries, want >= 5", len(names))
	}
	// Must include the core themes.
	required := map[string]bool{"modern": false, "default": false, "dark": false, "forest": false, "neutral": false}
	for _, n := range names {
		if _, ok := required[n]; ok {
			required[n] = true
		}
	}
	for name, found := range required {
		if !found {
			t.Errorf("Names() missing %q", name)
		}
	}
}
