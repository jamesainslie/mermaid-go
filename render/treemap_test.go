package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderTreemap(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Treemap
	g.TreemapTitle = "Budget"
	g.TreemapRoot = &ir.TreemapNode{
		Label: "Root",
		Children: []*ir.TreemapNode{
			{Label: "Salaries", Value: 70},
			{Label: "Equipment", Value: 30},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Budget") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Salaries") {
		t.Error("missing leaf label")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing rectangles")
	}
}

func TestRenderTreemapEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Treemap

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
