package ir

import "testing"

func TestTimelineEventDefaults(t *testing.T) {
	e := &TimelineEvent{Text: "Launch"}
	if e.Text != "Launch" {
		t.Errorf("Text = %q, want %q", e.Text, "Launch")
	}
}

func TestTimelinePeriodDefaults(t *testing.T) {
	p := &TimelinePeriod{
		Title:  "2024 Q1",
		Events: []*TimelineEvent{{Text: "Start"}, {Text: "Hire"}},
	}
	if p.Title != "2024 Q1" {
		t.Errorf("Title = %q, want %q", p.Title, "2024 Q1")
	}
	if len(p.Events) != 2 {
		t.Errorf("Events = %d, want 2", len(p.Events))
	}
}

func TestGraphTimelineFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Timeline
	g.TimelineTitle = "Project"
	g.TimelineSections = append(g.TimelineSections, &TimelineSection{
		Title: "Phase 1",
		Periods: []*TimelinePeriod{
			{Title: "Jan", Events: []*TimelineEvent{{Text: "Kickoff"}}},
		},
	})
	if g.TimelineTitle != "Project" {
		t.Errorf("TimelineTitle = %q, want %q", g.TimelineTitle, "Project")
	}
	if len(g.TimelineSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(g.TimelineSections))
	}
	if len(g.TimelineSections[0].Periods) != 1 {
		t.Fatalf("Periods = %d, want 1", len(g.TimelineSections[0].Periods))
	}
}
