package ir

// GanttTask represents a single task in a Gantt chart.
type GanttTask struct {
	ID       string   // Optional task ID for dependencies
	Label    string   // Display name
	StartStr string   // Start: date string, "after t1", or empty (follows previous)
	EndStr   string   // End: duration string ("5d", "2w") or date string
	Tags     []string // Status tags: done, active, crit, milestone
	AfterIDs []string // Task IDs this depends on (parsed from "after t1 t2")
	UntilID  string   // Task ID this runs until
}

// GanttSection groups tasks under a named section.
type GanttSection struct {
	Title string
	Tasks []*GanttTask
}
