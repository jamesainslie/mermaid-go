package theme

// Overrides specifies selective theme field overrides.
// Non-nil pointer fields replace the corresponding field in the base theme.
type Overrides struct {
	FontFamily *string
	FontSize   *float32
	Background *string

	PrimaryColor       *string
	PrimaryBorderColor *string
	PrimaryTextColor   *string

	SecondaryColor       *string
	SecondaryBorderColor *string
	SecondaryTextColor   *string

	TertiaryColor       *string
	TertiaryBorderColor *string

	LineColor *string
	TextColor *string

	ClusterBackground *string
	ClusterBorder     *string
	NodeBorderColor   *string

	NoteBackground  *string
	NoteBorderColor *string
	NoteTextColor   *string
}

// WithOverrides creates a new Theme by copying base and applying non-nil overrides.
// If base is nil, Modern() is used.
func WithOverrides(base *Theme, o Overrides) *Theme {
	if base == nil {
		base = Modern()
	}
	// Shallow copy the base theme.
	t := *base
	// Deep-copy slice fields to prevent mutation.
	t.PieColors = copyStrings(base.PieColors)
	t.TimelineSectionColors = copyStrings(base.TimelineSectionColors)
	t.GanttSectionColors = copyStrings(base.GanttSectionColors)
	t.GitBranchColors = copyStrings(base.GitBranchColors)
	t.XYChartColors = copyStrings(base.XYChartColors)
	t.RadarCurveColors = copyStrings(base.RadarCurveColors)
	t.MindmapBranchColors = copyStrings(base.MindmapBranchColors)
	t.SankeyNodeColors = copyStrings(base.SankeyNodeColors)
	t.TreemapColors = copyStrings(base.TreemapColors)
	t.BlockColors = copyStrings(base.BlockColors)
	t.JourneySectionColors = copyStrings(base.JourneySectionColors)

	// Apply non-nil overrides.
	if o.FontFamily != nil {
		t.FontFamily = *o.FontFamily
	}
	if o.FontSize != nil {
		t.FontSize = *o.FontSize
	}
	if o.Background != nil {
		t.Background = *o.Background
	}
	if o.PrimaryColor != nil {
		t.PrimaryColor = *o.PrimaryColor
	}
	if o.PrimaryBorderColor != nil {
		t.PrimaryBorderColor = *o.PrimaryBorderColor
	}
	if o.PrimaryTextColor != nil {
		t.PrimaryTextColor = *o.PrimaryTextColor
	}
	if o.SecondaryColor != nil {
		t.SecondaryColor = *o.SecondaryColor
	}
	if o.SecondaryBorderColor != nil {
		t.SecondaryBorderColor = *o.SecondaryBorderColor
	}
	if o.SecondaryTextColor != nil {
		t.SecondaryTextColor = *o.SecondaryTextColor
	}
	if o.TertiaryColor != nil {
		t.TertiaryColor = *o.TertiaryColor
	}
	if o.TertiaryBorderColor != nil {
		t.TertiaryBorderColor = *o.TertiaryBorderColor
	}
	if o.LineColor != nil {
		t.LineColor = *o.LineColor
	}
	if o.TextColor != nil {
		t.TextColor = *o.TextColor
	}
	if o.ClusterBackground != nil {
		t.ClusterBackground = *o.ClusterBackground
	}
	if o.ClusterBorder != nil {
		t.ClusterBorder = *o.ClusterBorder
	}
	if o.NodeBorderColor != nil {
		t.NodeBorderColor = *o.NodeBorderColor
	}
	if o.NoteBackground != nil {
		t.NoteBackground = *o.NoteBackground
	}
	if o.NoteBorderColor != nil {
		t.NoteBorderColor = *o.NoteBorderColor
	}
	if o.NoteTextColor != nil {
		t.NoteTextColor = *o.NoteTextColor
	}
	return &t
}

func copyStrings(s []string) []string {
	if s == nil {
		return nil
	}
	cp := make([]string, len(s))
	copy(cp, s)
	return cp
}
