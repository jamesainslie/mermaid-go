package theme

import "testing"

func TestWithOverridesNilBase(t *testing.T) {
	th := WithOverrides(nil, Overrides{Background: strPtr("#000")})
	if th == nil {
		t.Fatal("WithOverrides with nil base returned nil")
	}
	// Should use Modern as fallback base with override applied.
	if th.Background != "#000" {
		t.Errorf("Background = %q, want #000", th.Background)
	}
}

func TestWithOverridesBackground(t *testing.T) {
	base := Modern()
	th := WithOverrides(base, Overrides{Background: strPtr("#111111")})
	if th.Background != "#111111" {
		t.Errorf("Background = %q, want #111111", th.Background)
	}
	// Other fields should remain unchanged.
	if th.FontFamily != base.FontFamily {
		t.Errorf("FontFamily changed unexpectedly: %q", th.FontFamily)
	}
}

func TestWithOverridesPrimaryColor(t *testing.T) {
	base := Modern()
	th := WithOverrides(base, Overrides{PrimaryColor: strPtr("#FF0000")})
	if th.PrimaryColor != "#FF0000" {
		t.Errorf("PrimaryColor = %q, want #FF0000", th.PrimaryColor)
	}
}

func TestWithOverridesMultipleFields(t *testing.T) {
	base := MermaidDefault()
	th := WithOverrides(base, Overrides{
		FontFamily: strPtr("Fira Code"),
		FontSize:   float32Ptr(12),
		LineColor:  strPtr("#AABBCC"),
	})
	if th.FontFamily != "Fira Code" {
		t.Errorf("FontFamily = %q, want Fira Code", th.FontFamily)
	}
	if th.FontSize != 12 {
		t.Errorf("FontSize = %f, want 12", th.FontSize)
	}
	if th.LineColor != "#AABBCC" {
		t.Errorf("LineColor = %q, want #AABBCC", th.LineColor)
	}
	// Unchanged field.
	if th.TextColor != base.TextColor {
		t.Errorf("TextColor changed unexpectedly")
	}
}

func TestWithOverridesEmptyIsNoop(t *testing.T) {
	base := Modern()
	th := WithOverrides(base, Overrides{})
	if th.Background != base.Background {
		t.Errorf("Background = %q, want %q", th.Background, base.Background)
	}
	if th.PrimaryColor != base.PrimaryColor {
		t.Errorf("PrimaryColor changed unexpectedly")
	}
}

func TestWithOverridesDoesNotMutateBase(t *testing.T) {
	base := Modern()
	origBg := base.Background
	_ = WithOverrides(base, Overrides{Background: strPtr("#999999")})
	if base.Background != origBg {
		t.Error("WithOverrides mutated the original base theme")
	}
}

func strPtr(s string) *string       { return &s }
func float32Ptr(f float32) *float32 { return &f }
