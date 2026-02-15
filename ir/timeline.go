package ir

// TimelineEvent represents a single event within a time period.
type TimelineEvent struct {
	Text string
}

// TimelinePeriod represents a time period with one or more events.
type TimelinePeriod struct {
	Title  string
	Events []*TimelineEvent
}

// TimelineSection groups time periods under a named section.
type TimelineSection struct {
	Title   string
	Periods []*TimelinePeriod
}
