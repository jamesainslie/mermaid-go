package ir

import "testing"

func TestGitCommitTypeString(t *testing.T) {
	tests := []struct {
		ct   GitCommitType
		want string
	}{
		{GitCommitNormal, "NORMAL"},
		{GitCommitReverse, "REVERSE"},
		{GitCommitHighlight, "HIGHLIGHT"},
	}
	for _, tt := range tests {
		if got := tt.ct.String(); got != tt.want {
			t.Errorf("%d.String() = %q, want %q", tt.ct, got, tt.want)
		}
	}
}

func TestGitActionInterface(t *testing.T) {
	// Verify all action types implement GitAction.
	actions := []GitAction{
		&GitCommit{ID: "c1"},
		&GitBranch{Name: "dev"},
		&GitCheckout{Branch: "dev"},
		&GitMerge{Branch: "dev"},
		&GitCherryPick{ID: "c1"},
	}
	if len(actions) != 5 {
		t.Errorf("actions = %d, want 5", len(actions))
	}
}

func TestGraphGitGraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = GitGraph
	g.GitMainBranch = "main"
	g.GitActions = append(g.GitActions,
		&GitCommit{ID: "init", Tag: "v0.1"},
		&GitBranch{Name: "develop", Order: 1},
		&GitCheckout{Branch: "develop"},
		&GitCommit{ID: "feat1", Type: GitCommitHighlight},
		&GitCheckout{Branch: "main"},
		&GitMerge{Branch: "develop", ID: "merge1", Tag: "v1.0"},
	)

	if g.GitMainBranch != "main" {
		t.Errorf("GitMainBranch = %q", g.GitMainBranch)
	}
	if len(g.GitActions) != 6 {
		t.Errorf("GitActions = %d, want 6", len(g.GitActions))
	}
}
