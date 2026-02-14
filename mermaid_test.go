package mermaid

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	svg, err := Render("flowchart LR; A-->B-->C")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("missing </svg>")
	}
}

func TestRenderWithOptions(t *testing.T) {
	opts := Options{}
	svg, err := RenderWithOptions("flowchart TD; X-->Y", opts)
	if err != nil {
		t.Fatalf("RenderWithOptions() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
}

func TestRenderWithTiming(t *testing.T) {
	result, err := RenderWithTiming("flowchart LR; A-->B", Options{})
	if err != nil {
		t.Fatalf("RenderWithTiming() error: %v", err)
	}
	if !strings.Contains(result.SVG, "<svg") {
		t.Error("missing <svg")
	}
	if result.TotalUs() <= 0 {
		t.Error("TotalUs should be > 0")
	}
}

func TestRenderInvalidInput(t *testing.T) {
	_, err := Render("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestRenderContainsNodeLabels(t *testing.T) {
	svg, err := Render("flowchart LR\n  A[Start] --> B[End]")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "Start") {
		t.Error("missing label 'Start'")
	}
	if !strings.Contains(svg, "End") {
		t.Error("missing label 'End'")
	}
}
