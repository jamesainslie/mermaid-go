package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderTimeline(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Timeline
	g.TimelineTitle = "History"
	g.TimelineSections = []*ir.TimelineSection{
		{
			Title: "Early",
			Periods: []*ir.TimelinePeriod{
				{Title: "2002", Events: []*ir.TimelineEvent{{Text: "LinkedIn"}}},
				{Title: "2004", Events: []*ir.TimelineEvent{{Text: "Facebook"}}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "History") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "LinkedIn") {
		t.Error("missing event text")
	}
	if !strings.Contains(svg, "2002") {
		t.Error("missing period label")
	}
}
