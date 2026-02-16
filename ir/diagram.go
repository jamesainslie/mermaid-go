package ir

type DiagramKind int

const (
	Flowchart DiagramKind = iota
	Class
	State
	Sequence
	Er
	Pie
	Mindmap
	Journey
	Timeline
	Gantt
	Requirement
	GitGraph
	C4
	Sankey
	Quadrant
	ZenUML
	Block
	Packet
	Kanban
	Architecture
	Radar
	Treemap
	XYChart
)

// String returns a human-readable name for the diagram kind.
func (k DiagramKind) String() string {
	switch k {
	case Flowchart:
		return "Flowchart"
	case Class:
		return "Class"
	case State:
		return "State"
	case Sequence:
		return "Sequence"
	case Er:
		return "ER"
	case Pie:
		return "Pie"
	case Mindmap:
		return "Mindmap"
	case Journey:
		return "Journey"
	case Timeline:
		return "Timeline"
	case Gantt:
		return "Gantt"
	case Requirement:
		return "Requirement"
	case GitGraph:
		return "GitGraph"
	case C4:
		return "C4"
	case Sankey:
		return "Sankey"
	case Quadrant:
		return "Quadrant"
	case ZenUML:
		return "ZenUML"
	case Block:
		return "Block"
	case Packet:
		return "Packet"
	case Kanban:
		return "Kanban"
	case Architecture:
		return "Architecture"
	case Radar:
		return "Radar"
	case Treemap:
		return "Treemap"
	case XYChart:
		return "XYChart"
	default:
		return "Unknown"
	}
}

type Direction int

const (
	TopDown Direction = iota
	LeftRight
	BottomTop
	RightLeft
)

func DirectionFromToken(token string) (Direction, bool) {
	switch token {
	case "TD", "TB":
		return TopDown, true
	case "LR":
		return LeftRight, true
	case "RL":
		return RightLeft, true
	case "BT":
		return BottomTop, true
	default:
		return TopDown, false
	}
}
