# Phase 3: Sequence Diagrams — Design

## Goal

Add full-parity mermaid.js sequence diagram support: parser, timeline-based layout, and SVG renderer. Covers all participant types, message arrows, activations, notes, frames, boxes, autonumber, create/destroy lifecycle, links, and properties.

## Architecture

Sequence diagrams are fundamentally different from graph-based diagrams (flowchart, class, ER, state). They use a **timeline-based layout** with fixed participant columns on the horizontal axis and events ordered top-to-bottom on the vertical axis. This means a new custom layout algorithm, not the Sugiyama pipeline.

The implementation follows the same package structure as Phase 2: IR types in `ir/`, parser in `parser/`, layout in `layout/`, renderer in `render/`.

## IR Types

New file `ir/sequence.go`. Sequence-specific fields added to `ir.Graph`.

### Participants

```go
type SeqParticipantKind int

const (
    ParticipantBox SeqParticipantKind = iota
    ActorStickFigure
    ParticipantBoundary
    ParticipantControl
    ParticipantEntity
    ParticipantDatabase
    ParticipantCollections
    ParticipantQueue
)

type SeqParticipant struct {
    ID          string
    Alias       string
    Kind        SeqParticipantKind
    Links       []SeqLink
    Properties  map[string]string
    IsCreated   bool
    IsDestroyed bool
}

type SeqLink struct {
    Label string
    URL   string
}
```

### Messages

```go
type SeqMessageKind int

const (
    MsgSolid SeqMessageKind = iota         // ->
    MsgDotted                               // -->
    MsgSolidArrow                           // ->>
    MsgDottedArrow                          // -->>
    MsgSolidCross                           // -x
    MsgDottedCross                          // --x
    MsgSolidOpen                            // -)
    MsgDottedOpen                           // --)
    MsgBiSolid                              // <<->>
    MsgBiDotted                             // <<-->>
)

type SeqMessage struct {
    From             string
    To               string
    Text             string
    Kind             SeqMessageKind
    ActivateTarget   bool
    DeactivateSource bool
}
```

### Events

The core data structure is an ordered event list. The parser emits events in declaration order; the layout walks them top-to-bottom.

```go
type SeqEventKind int

const (
    EvMessage SeqEventKind = iota
    EvNote
    EvActivate
    EvDeactivate
    EvFrameStart
    EvFrameMiddle   // else, and, option
    EvFrameEnd
    EvCreate
    EvDestroy
)

type SeqEvent struct {
    Kind    SeqEventKind
    Message *SeqMessage
    Note    *SeqNote
    Frame   *SeqFrame
    Target  string          // for activate/deactivate/create/destroy
}
```

### Notes

```go
type SeqNotePosition int

const (
    NoteLeft SeqNotePosition = iota
    NoteRight
    NoteOver
)

type SeqNote struct {
    Position     SeqNotePosition
    Participants []string
    Text         string
}
```

### Frames

```go
type SeqFrameKind int

const (
    FrameLoop SeqFrameKind = iota
    FrameAlt
    FrameOpt
    FramePar
    FrameCritical
    FrameBreak
    FrameRect
)

type SeqFrame struct {
    Kind  SeqFrameKind
    Label string
    Color string  // for rect frames
}
```

### Boxes

```go
type SeqBox struct {
    Label        string
    Color        string
    Participants []string
}
```

### Graph Fields

```go
// On ir.Graph:
Participants []*SeqParticipant
Events       []*SeqEvent
Boxes        []*SeqBox
Autonumber   bool
```

## Parser

New file `parser/sequence.go`. Single-pass line-by-line parser using `preprocessInput()` and regex-based matching.

### Parse Order

1. `sequenceDiagram` header (detected by `detectDiagramKind`)
2. Participant/actor declarations: `participant A as Alice`, `actor B`, JSON syntax `participant API@{ "type": "boundary" }`
3. Box blocks: `box Color Label` ... `end`
4. `autonumber` flag
5. Links/properties: `link Alice: Dashboard @ https://...`, `links Alice: {...}`, `properties Alice: {...}`
6. Messages: all 10 arrow types with `+`/`-` activation shorthand
7. Activation directives: `activate Alice`, `deactivate Alice`
8. Notes: `Note right of Alice: text`, `Note over Alice,Bob: text`
9. Frame keywords: `loop`, `alt`, `else`, `opt`, `par`, `and`, `critical`, `option`, `break`, `rect`, `end`
10. Create/destroy: `create participant Carl`, `destroy Carl`

### Implicit Participants

Any participant referenced in a message but not explicitly declared is auto-created at first mention, in order of appearance.

### Line Break Handling

`<br/>` tags in message text and note text are replaced with newlines, producing multi-line TextBlocks.

### JSON Participant Types

The `@{ ... }` syntax is parsed with `encoding/json` to extract `type` and optional `alias` fields.

## Layout

New file `layout/sequence.go`. Custom timeline-based layout algorithm (not Sugiyama).

### Algorithm

Single pass top-to-bottom:

1. **Measure participants**: compute header/footer text widths, assign column X positions with horizontal spacing. Boxes add grouping padding around their contained participants.
2. **Walk events** with a Y cursor:
   - **Message**: advance Y by message spacing. X endpoints from participant columns. Self-messages get a right-bump loop shape.
   - **Note**: measure text, advance Y by note height + padding.
   - **Activate/Deactivate**: push/pop activation stack per participant. Track Y ranges.
   - **Frame start/middle/end**: push/pop frame stack. Frames record Y-start; on end, compute height. Alt/par/critical get horizontal divider Y positions at each else/and/option.
   - **Create**: position new participant header at current Y (not top).
   - **Destroy**: mark destruction Y, stop lifeline.
3. **Finalize**: set lifeline heights, add footer headers at bottom, compute bounding box.

### Layout Data

```go
type SequenceData struct {
    Participants  []SeqParticipantLayout
    Lifelines     []SeqLifeline
    Messages      []SeqMessageLayout
    Activations   []SeqActivationLayout
    Notes         []SeqNoteLayout
    Frames        []SeqFrameLayout
    Boxes         []SeqBoxLayout
    Autonumber    bool
    DiagramHeight float32
}
func (SequenceData) diagramData() {}
```

### Config

```go
type SequenceConfig struct {
    ParticipantSpacing float32  // default 80
    MessageSpacing     float32  // default 40
    ActivationWidth    float32  // default 16
    NoteMaxWidth       float32  // default 200
    BoxPadding         float32  // default 12
    FramePadding       float32  // default 10
    HeaderHeight       float32  // default 40
    SelfMessageWidth   float32  // default 40
}
```

## Renderer

New file `render/sequence.go`. Rendering follows visual stacking order:

1. **Boxes** — background rectangles behind participant groups
2. **Frames** — loop/alt/par/etc. rectangles with label tabs. Alt/par/critical get horizontal dashed dividers at else/and/option boundaries. Rect frames use specified color as fill.
3. **Lifelines** — vertical dashed lines from header bottom to footer top
4. **Activations** — narrow filled rectangles on lifelines
5. **Messages** — horizontal arrows with labels. Self-messages as right-bump loops. Autonumber badges as filled circles with sequence numbers.
6. **Notes** — positioned rectangles with text
7. **Participant headers + footers** — rendered by kind:
   - `ParticipantBox`: rounded rect
   - `ActorStickFigure`: UML stick figure
   - `Boundary`/`Control`/`Entity`/`Database`/`Collections`/`Queue`: UML stereotyped shapes

### New SVG Markers

- Open arrow (async messages)
- Cross end (termination messages)
- Bidirectional arrows

Reuse existing closed triangle for solid arrows.

### Theme Colors (already present)

ActorBorder, ActorBackground, ActorTextColor, ActorLineColor, SignalColor, SignalTextColor, ActivationBorderColor, ActivationBackground, SequenceNumberColor, LoopTextColor, NoteBackground, NoteBorderColor, NoteTextColor.

## Testing

### Parser Tests (`parser/sequence_test.go`)

Table-driven tests for:
- Participant/actor declarations with aliases
- All 10 message types
- Activation shorthand (+/-)
- Notes (left, right, over spanning)
- Each frame type (loop, alt/else, opt, par/and, critical/option/break, rect)
- Boxes with participant grouping
- Autonumber
- Create/destroy lifecycle
- Links and properties
- JSON participant types
- Implicit participant creation
- `<br/>` line breaks

### Layout Tests (`layout/sequence_test.go`)

- Participant X positions monotonically increasing
- Message Y positions advance top-to-bottom
- Self-message width matches config
- Activation bars on correct lifeline X
- Frame bounds contain child events
- Created participants have header Y > 0

### Renderer Tests (`render/sequence_test.go`)

SVG structural checks:
- Lifeline dashed lines present
- Actor boxes/stick figures rendered
- Arrow markers in defs
- Activation rectangles present
- Frame labels rendered

### Integration Tests (`mermaid_test.go`)

- Simple two-participant exchange
- Complex diagram with frames, activations, notes

### Fixtures (`testdata/fixtures/`)

- `sequence-simple.mmd`
- `sequence-activations.mmd`
- `sequence-frames.mmd`
- `sequence-full.mmd`

### Benchmarks (`mermaid_bench_test.go`)

- `BenchmarkRenderSequenceSimple`
- `BenchmarkRenderSequenceComplex`
