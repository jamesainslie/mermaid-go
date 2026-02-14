package render

import (
	"strings"
	"testing"

	"github.com/yaklabco/mermaid-go/config"
	"github.com/yaklabco/mermaid-go/ir"
	"github.com/yaklabco/mermaid-go/layout"
	"github.com/yaklabco/mermaid-go/theme"
)

func simpleLayout() *layout.Layout {
	g := ir.NewGraph()
	g.Kind = ir.Flowchart
	g.Direction = ir.LeftRight
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	g.Edges = []*ir.Edge{{
		From: "A", To: "B", Directed: true, ArrowEnd: true, Style: ir.Solid,
	}}
	th := theme.Modern()
	cfg := config.DefaultLayout()
	return layout.ComputeLayout(g, th, cfg)
}

func TestRenderSVGContainsSVGTags(t *testing.T) {
	l := simpleLayout()
	th := theme.Modern()
	cfg := config.DefaultLayout()
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("missing </svg> tag")
	}
}

func TestRenderSVGContainsNodes(t *testing.T) {
	l := simpleLayout()
	th := theme.Modern()
	cfg := config.DefaultLayout()
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, "<rect") {
		t.Error("missing <rect for node shapes")
	}
	if !strings.Contains(svg, "A") {
		t.Error("missing node label A")
	}
	if !strings.Contains(svg, "B") {
		t.Error("missing node label B")
	}
}

func TestRenderSVGContainsEdge(t *testing.T) {
	l := simpleLayout()
	th := theme.Modern()
	cfg := config.DefaultLayout()
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, "<path") || !strings.Contains(svg, "edgePath") {
		t.Error("missing edge path")
	}
}

func TestRenderSVGHasViewBox(t *testing.T) {
	l := simpleLayout()
	th := theme.Modern()
	cfg := config.DefaultLayout()
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, "viewBox") {
		t.Error("missing viewBox attribute")
	}
}
