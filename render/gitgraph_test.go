package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderGitGraph(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"
	g.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1", Tag: "v1.0"},
		&ir.GitBranch{Name: "develop"},
		&ir.GitCheckout{Branch: "develop"},
		&ir.GitCommit{ID: "c2"},
		&ir.GitCheckout{Branch: "main"},
		&ir.GitMerge{Branch: "develop"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "v1.0") {
		t.Error("missing tag label")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing commit circles")
	}
	if !strings.Contains(svg, "main") {
		t.Error("missing branch label")
	}
}
