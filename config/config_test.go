package config

import "testing"

func TestDefaultLayout(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.NodeSpacing != 50 {
		t.Errorf("NodeSpacing = %f, want 50", cfg.NodeSpacing)
	}
	if cfg.RankSpacing != 70 {
		t.Errorf("RankSpacing = %f, want 70", cfg.RankSpacing)
	}
	if cfg.LabelLineHeight <= 0 {
		t.Error("LabelLineHeight should be > 0")
	}
}

func TestDefaultLayoutHasClassConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Class.CompartmentPadX <= 0 {
		t.Error("Class.CompartmentPadX should be > 0")
	}
	if cfg.State.CompositePadding <= 0 {
		t.Error("State.CompositePadding should be > 0")
	}
	if cfg.ER.AttributeRowHeight <= 0 {
		t.Error("ER.AttributeRowHeight should be > 0")
	}
}

func TestDefaultLayoutHasSequenceConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Sequence.ParticipantSpacing <= 0 {
		t.Error("Sequence.ParticipantSpacing should be > 0")
	}
	if cfg.Sequence.MessageSpacing <= 0 {
		t.Error("Sequence.MessageSpacing should be > 0")
	}
	if cfg.Sequence.ActivationWidth <= 0 {
		t.Error("Sequence.ActivationWidth should be > 0")
	}
	if cfg.Sequence.HeaderHeight <= 0 {
		t.Error("Sequence.HeaderHeight should be > 0")
	}
}

func TestDefaultLayoutPieConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Pie.Radius != 150 {
		t.Errorf("Pie.Radius = %f, want 150", cfg.Pie.Radius)
	}
	if cfg.Pie.TextPosition != 0.75 {
		t.Errorf("Pie.TextPosition = %f, want 0.75", cfg.Pie.TextPosition)
	}
	if cfg.Pie.PaddingX != 20 {
		t.Errorf("Pie.PaddingX = %f, want 20", cfg.Pie.PaddingX)
	}
}

func TestDefaultLayoutQuadrantConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Quadrant.ChartWidth != 400 {
		t.Errorf("Quadrant.ChartWidth = %f, want 400", cfg.Quadrant.ChartWidth)
	}
	if cfg.Quadrant.ChartHeight != 400 {
		t.Errorf("Quadrant.ChartHeight = %f, want 400", cfg.Quadrant.ChartHeight)
	}
	if cfg.Quadrant.PointRadius != 5 {
		t.Errorf("Quadrant.PointRadius = %f, want 5", cfg.Quadrant.PointRadius)
	}
}
