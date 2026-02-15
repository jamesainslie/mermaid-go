package ir

import "testing"

func TestGanttTaskDefaults(t *testing.T) {
	task := &GanttTask{
		ID:       "t1",
		Label:    "Design",
		StartStr: "2024-01-01",
		EndStr:   "10d",
		Tags:     []string{"crit"},
	}
	if task.ID != "t1" {
		t.Errorf("ID = %q, want %q", task.ID, "t1")
	}
	if len(task.Tags) != 1 || task.Tags[0] != "crit" {
		t.Errorf("Tags = %v, want [crit]", task.Tags)
	}
}

func TestGanttSectionDefaults(t *testing.T) {
	s := &GanttSection{
		Title: "Development",
		Tasks: []*GanttTask{
			{ID: "d1", Label: "Code"},
		},
	}
	if s.Title != "Development" {
		t.Errorf("Title = %q", s.Title)
	}
	if len(s.Tasks) != 1 {
		t.Errorf("Tasks = %d, want 1", len(s.Tasks))
	}
}

func TestGraphGanttFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Gantt
	g.GanttTitle = "Project"
	g.GanttDateFormat = "YYYY-MM-DD"
	g.GanttAxisFormat = "%Y-%m-%d"
	g.GanttExcludes = []string{"weekends"}
	g.GanttSections = append(g.GanttSections, &GanttSection{
		Title: "Dev",
		Tasks: []*GanttTask{{ID: "t1", Label: "Code"}},
	})

	if g.GanttTitle != "Project" {
		t.Errorf("GanttTitle = %q", g.GanttTitle)
	}
	if g.GanttDateFormat != "YYYY-MM-DD" {
		t.Errorf("GanttDateFormat = %q", g.GanttDateFormat)
	}
	if len(g.GanttExcludes) != 1 {
		t.Errorf("GanttExcludes = %d, want 1", len(g.GanttExcludes))
	}
	if len(g.GanttSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(g.GanttSections))
	}
}
