# Phase 2: Core Graph Variants Design

**Date:** 2026-02-14
**Status:** Approved
**Scope:** Class diagrams, state diagrams (v2), ER diagrams — full mermaid-js parity

## Constraints

- Full syntax parity with mermaid-js for all three diagram types
- Class and ER diagrams reuse Sugiyama graph layout pipeline
- State diagrams get recursive layout for composite/nested states
- Class rendering uses UML compartment boxes (name/attributes/methods)
- TDD throughout

## 1. Approach

Extend the existing architecture (Approach A from brainstorming):

- Single `ir.Graph` struct gains diagram-specific optional fields
- New IR type files: `ir/class.go`, `ir/state.go`, `ir/er.go`
- One parser file per diagram type: `parser/class.go`, `parser/state.go`, `parser/er.go`
- Layout reuse: class/ER use `computeGraphLayout` with sizing tweaks; state gets `computeStateLayout` with recursive layout for composites
- One renderer file per diagram type: `render/class.go`, `render/state.go`, `render/er.go`

## 2. IR Extensions

### 2.1 Class Diagram Types (`ir/class.go`)

```go
type Visibility int
const (
    VisPublic    Visibility = iota // +
    VisPrivate                     // -
    VisProtected                   // #
    VisPackage                     // ~
    VisNone                        // no prefix
)

type MemberClassifier int
const (
    ClassifierNone     MemberClassifier = iota
    ClassifierAbstract                          // *
    ClassifierStatic                            // $
)

type ClassMember struct {
    Name       string
    Type       string           // attribute type or return type
    Params     string           // method parameters (empty for attributes)
    IsMethod   bool             // true if has parentheses
    Visibility Visibility
    Classifier MemberClassifier
    Generic    string           // generic type parameter e.g. "T"
}

type ClassMembers struct {
    Attributes []ClassMember
    Methods    []ClassMember
}

type Namespace struct {
    Name    string
    Classes []string // class node IDs
}
```

**Fields added to `ir.Graph`:**

```go
Members      map[string]*ClassMembers // node ID -> members
Annotations  map[string]string        // node ID -> stereotype text
Namespaces   []*Namespace
ClassNotes   []*DiagramNote           // shared note type
```

### 2.2 State Diagram Types (`ir/state.go`)

```go
type StateAnnotation int
const (
    StateChoice StateAnnotation = iota
    StateFork
    StateJoin
)

type CompositeState struct {
    ID        string
    Label     string
    Inner     *Graph            // nested state machine
    Regions   []*Graph          // concurrent regions (if any)
    Direction *Direction        // optional inner direction
}

type DiagramNote struct {
    Text     string
    Position string // "right of", "left of"
    Target   string // node ID
}
```

**Fields added to `ir.Graph`:**

```go
CompositeStates  map[string]*CompositeState
StateDescriptions map[string]string
StateAnnotations map[string]StateAnnotation
StateNotes       []*DiagramNote
```

### 2.3 ER Diagram Types (`ir/er.go`)

```go
type AttributeKey int
const (
    KeyNone    AttributeKey = iota
    KeyPrimary                     // PK
    KeyForeign                     // FK
    KeyUnique                      // UK
)

type EntityAttribute struct {
    Type    string
    Name    string
    Keys    []AttributeKey
    Comment string
}

type Entity struct {
    ID         string
    Label      string     // display name (alias)
    Attributes []EntityAttribute
}
```

**Fields added to `ir.Graph`:**

```go
Entities map[string]*Entity
```

### 2.4 New Edge Arrowheads

Added to existing `EdgeArrowhead` enum in `ir/shapes.go`:

```go
const (
    OpenTriangle     EdgeArrowhead = iota // existing
    ClassDependency                       // existing
    ClosedTriangle                        // inheritance, realization
    FilledDiamond                         // composition
    OpenDiamond                           // aggregation
    Lollipop                              // provided interface
)
```

## 3. Parsers

### 3.1 Class Diagram Parser (`parser/class.go`)

**Header detection:** `classDiagram` or `classDiagram-v2`

**Syntax features (full parity):**
- Class definitions: implicit (via relationship), explicit (`class Name`), with body (`class Name { ... }`)
- Class labels: `class Name["Display Label"]`
- Members: visibility prefix, type, name, method parameters, return type, `*`/`$` classifiers
- Generic types: `~T~` notation, nested generics
- Relationships: 8 types + bidirectional + reversed (`<|--`, `*--`, `o--`, `-->`, `..>`, `..|>`, `--`, `..`)
- Cardinality: `"1" --> "*"` syntax
- Lollipop interfaces: `()--` and `--()`
- Annotations: `<<interface>>`, `<<abstract>>`, etc. (inline and standalone)
- Namespaces: `namespace Name { ... }` blocks
- Notes: `note "text"` and `note for ClassName "text"`
- Directives: direction, classDef, cssClass/`:::`, style, click/callback/link
- Comments: `%%`

**Approach:** Line-by-line processing. Brace-delimited blocks (class bodies, namespaces) tracked with depth counter. Regex-based relationship parsing with arrow pattern matching.

### 3.2 State Diagram Parser (`parser/state.go`)

**Header detection:** `stateDiagram` or `stateDiagram-v2`

**Syntax features (full parity):**
- State definitions: bare ID, colon syntax (`s1 : desc`), `state "desc" as s1`
- Transitions: `s1 --> s2`, `s1 --> s2 : label`
- Start/end: `[*] --> s1`, `s1 --> [*]`
- Composite states: `state Name { ... }` with recursive parsing of inner content
- Annotations: `<<choice>>`, `<<fork>>`, `<<join>>` (and `[[choice]]` etc.)
- Concurrent regions: `--` separator within composite states
- Notes: block (`note right of S1 ... end note`) and inline (`note left of S1 : text`)
- Direction within composites
- Directives: classDef, class/`:::`, hide empty description
- Comments: `%%`

**Approach:** Recursive descent for composite states. Top-level line-by-line, drops into recursive parser when opening brace of `state Name {` is encountered. Concurrent regions split inner content on `--` lines.

### 3.3 ER Diagram Parser (`parser/er.go`)

**Header detection:** `erDiagram`

**Syntax features (full parity):**
- Entity definitions: bare name, with attributes `{ ... }`
- Entity aliases: `e["Display Name"]`, `e[Alias]`
- Attributes: `type name`, `type name PK`, `type name "comment"`, `type name PK,FK "comment"`
- Relationships: cardinality markers + line style + label
- Cardinality: `||`, `|o`, `o|`, `}|`, `|{`, `}o`, `o{` (all combinations)
- Line styles: `--` (identifying), `..` (non-identifying)
- Labels: required after `:` on relationship lines
- Directives: direction, classDef, class/`:::`, style
- Comments: `%%`

**Approach:** Line-by-line. Entity blocks tracked with brace depth. Relationship regex matches cardinality + line + cardinality pattern.

## 4. Layout

### 4.1 Class Diagram Layout

Reuses `computeGraphLayout` with class-specific sizing:

- `sizeClassNode()`: measures UML compartment box
  - Header: class name + annotation height
  - Attributes section: sum of attribute line heights
  - Methods section: sum of method line heights
  - Width: max(header width, max attribute width, max method width) + padding
  - Height: sum of sections + divider lines + padding
- Namespace groups rendered like subgraphs (bounding box around contained classes)
- Note positioning: offset from target node

**DiagramData:** `layout.ClassData` stores per-node compartment dimensions (header height, attribute section height, method section height).

### 4.2 State Diagram Layout (Specialized)

`computeStateLayout()` — recursive layout:

1. Identify top-level states (not inside composites)
2. For each composite state, recursively layout inner graph
3. For composites with concurrent regions, layout each region independently, stack vertically
4. Size composite nodes to contain their inner layout + padding + label
5. Run Sugiyama pipeline on top-level graph with sized composites
6. Translate inner layout coordinates relative to composite position

**Special node types:**
- `[*]` start: small filled circle (fixed size)
- `[*]` end: bullseye (fixed size)
- `<<fork>>`/`<<join>>`: horizontal bar (ForkJoin shape)
- `<<choice>>`: diamond (Diamond shape)

**DiagramData:** `layout.StateData` stores nested layout maps (composite ID -> inner Layout), region boundaries.

### 4.3 ER Diagram Layout

Reuses `computeGraphLayout` with ER-specific sizing:

- `sizeEREntity()`: measures entity box
  - Header: entity name height
  - Attributes: rows with type/name/key columns
  - Width: sum of column widths + padding
  - Height: header + (num attributes * row height) + padding
- Edge endpoints adjusted for crow's foot decoration clearance

**DiagramData:** `layout.ERData` stores entity column widths, row heights.

## 5. Rendering

### 5.1 Class Diagram Renderer (`render/class.go`)

- UML compartment boxes:
  - Rounded rectangle with three sections
  - Header: class name centered, annotation above in guillemets
  - Horizontal divider lines between sections
  - Attributes: left-aligned, visibility symbol prefix
  - Methods: left-aligned, visibility symbol prefix, `()` suffix
  - Abstract members: italic text
  - Static members: underlined text
- Relationship arrows:
  - Inheritance: solid line + closed triangle (filled white)
  - Composition: solid line + filled diamond
  - Aggregation: solid line + open diamond
  - Association: solid line + open arrow
  - Dependency: dashed line + open arrow
  - Realization: dashed line + closed triangle
  - Lollipop: line + circle
- Cardinality labels: near arrow endpoints
- Namespace boxes: dashed rectangle background with label
- Notes: rectangle with folded corner + text

### 5.2 State Diagram Renderer (`render/state.go`)

- States: rounded rectangles with optional description divider
- Start state: filled black circle
- End state: filled black circle with outer ring (bullseye)
- Fork/join: wide horizontal bar (filled)
- Choice: diamond shape
- Composite states: larger rounded rectangle with label + inner content
- Concurrent regions: horizontal dashed divider lines within composites
- Transitions: arrows with optional label
- Notes: rectangles positioned relative to target state

### 5.3 ER Diagram Renderer (`render/er.go`)

- Entity boxes: rectangle with header + attribute rows
- Header: bold entity name, filled background
- Attribute rows: type | name | key columns, alternating row backgrounds
- Crow's foot decorations at relationship endpoints:
  - `||` exactly one: two perpendicular lines
  - `o|` / `|o` zero or one: circle + line
  - `|{` / `}|` one or more: line + crow's foot
  - `o{` / `}o` zero or more: circle + crow's foot
- Solid lines (identifying) vs dashed lines (non-identifying)
- Relationship labels at edge midpoints

### 5.4 New SVG Markers (`render/svg.go`)

Added to `<defs>` section:
- `marker-closed-triangle`: filled white triangle (inheritance/realization)
- `marker-filled-diamond`: filled diamond (composition)
- `marker-open-diamond`: open diamond (aggregation)
- `marker-crowsfoot-many`: three-line crow's foot
- `marker-crowsfoot-one`: perpendicular line
- `marker-crowsfoot-zero`: circle

## 6. Theme & Config

### 6.1 Theme Additions

```go
// Added to theme.Theme
ClassHeaderBg    string // class name section background
ClassBodyBg      string // members section background
ClassBorder      string // class box border

StateFill        string // state box fill
StateBorder      string // state box border
StateStartEnd    string // start/end state fill (black)
CompositeHeaderBg string // composite state header

EntityHeaderBg   string // entity header background
EntityBodyBg     string // entity attribute rows
EntityBorder     string // entity box border
```

### 6.2 Config Additions

```go
type ClassConfig struct {
    MemberFontSize    float32
    CompartmentPadX   float32
    CompartmentPadY   float32
}

type StateConfig struct {
    CompositePadding  float32
    RegionSeparatorPad float32
    StartEndRadius    float32
    ForkBarWidth      float32
    ForkBarHeight     float32
}

type ERConfig struct {
    AttributeRowHeight float32
    ColumnPadding      float32
    HeaderPadding      float32
}
```

## 7. Testing

### Per diagram type:
- **Parser tests:** Table-driven, one test case per syntax feature
- **Layout tests:** Verify node dimensions, rank assignment, positioning
- **Renderer tests:** Golden SVG tests, structural assertions
- **Integration tests:** `mermaid.Render()` end-to-end

### Test fixtures (`testdata/fixtures/`):
- `class-simple.mmd` — basic classes with members
- `class-relationships.mmd` — all 8 relationship types
- `class-annotations.mmd` — stereotypes, namespaces
- `class-generics.mmd` — generic type parameters
- `state-simple.mmd` — states and transitions
- `state-composite.mmd` — nested states
- `state-concurrent.mmd` — parallel regions
- `state-choice-fork.mmd` — choice, fork, join
- `er-simple.mmd` — entities and relationships
- `er-attributes.mmd` — typed attributes with keys
- `er-cardinality.mmd` — all cardinality combinations

## 8. Implementation Order

1. IR extensions (types only, no behavior)
2. Class diagram: parser -> layout sizing -> renderer -> integration
3. State diagram: parser -> recursive layout -> renderer -> integration
4. ER diagram: parser -> layout sizing -> renderer -> integration
5. Theme & config additions (woven in as needed)
6. Cross-cutting: new SVG markers, shared helpers
