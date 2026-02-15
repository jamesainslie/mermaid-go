package ir

// GitCommitType represents the visual type of a git commit.
type GitCommitType int

const (
	GitCommitNormal GitCommitType = iota
	GitCommitReverse
	GitCommitHighlight
)

// String returns the mermaid keyword for the commit type.
func (t GitCommitType) String() string {
	switch t {
	case GitCommitReverse:
		return "REVERSE"
	case GitCommitHighlight:
		return "HIGHLIGHT"
	default:
		return "NORMAL"
	}
}

// GitAction is a sealed interface for git operations in a gitGraph diagram.
type GitAction interface {
	gitAction()
}

// GitCommit represents a commit operation.
type GitCommit struct {
	ID   string
	Tag  string
	Type GitCommitType
}

func (*GitCommit) gitAction() {}

// GitBranch represents a branch creation.
type GitBranch struct {
	Name  string
	Order int // Display order (-1 = unset)
}

func (*GitBranch) gitAction() {}

// GitCheckout represents switching to a branch.
type GitCheckout struct {
	Branch string
}

func (*GitCheckout) gitAction() {}

// GitMerge represents merging a branch into the current branch.
type GitMerge struct {
	Branch string
	ID     string
	Tag    string
	Type   GitCommitType
}

func (*GitMerge) gitAction() {}

// GitCherryPick represents cherry-picking a commit.
type GitCherryPick struct {
	ID     string // Source commit ID
	Parent string // Parent ID for merge commits
}

func (*GitCherryPick) gitAction() {}
