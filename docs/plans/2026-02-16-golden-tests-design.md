# Golden Tests Design

**Goal:** Detect SVG output regressions across all 50 fixtures and 5 themes by comparing rendered output against committed golden files.

**Architecture:** Single test function iterates the fixture/theme matrix (250 combinations), compares byte-for-byte against committed SVG golden files. Parallel agents generate and verify the initial golden baseline using anomaly detection.

## Scope

- 50 `.mmd` fixtures in `testdata/fixtures/`
- 5 themes: modern, default, dark, forest, neutral
- 250 golden `.svg` files in `testdata/golden/`
- Naming: `testdata/golden/<fixture>-<theme>.svg`

## Test Harness

Single file `mermaid_golden_test.go` at the root package level:

- Iterates all `.mmd` files in `testdata/fixtures/`
- For each fixture, renders with all 5 themes via `RenderWithOptions()`
- Compares output byte-for-byte against `testdata/golden/<fixture>-<theme>.svg`
- `-update` flag overwrites golden files with current output
- On mismatch: prints fixture name, theme, and a truncated unified diff

## Agentic Generation Workflow

Adapted from the [yaklab golden test methodology](https://yaklab.org/posts/golden-test-methodology/). Agents don't invent expected output; the renderer's actual output becomes the golden baseline. Agents verify that baseline is sane before committing.

### Phase 0: Scaffold (main agent)

Write the test harness and shared verification rules before dispatching agents.

### Phase 1: Generate + Verify (5 parallel agents)

Each agent handles a group of diagram types:

1. Agent A: Flowchart (3), Class (2), State (2), ER (2) = 9 fixtures, 45 golden files
2. Agent B: Sequence (4), Kanban (2), Packet (2) = 8 fixtures, 40 golden files
3. Agent C: Pie (2), Quadrant (2), Timeline (2), Gantt (2), GitGraph (2) = 10 fixtures, 50 golden files
4. Agent D: XYChart (2), Radar (2), Mindmap (2), Sankey (2), Treemap (2) = 10 fixtures, 50 golden files
5. Agent E: Requirement (2), Block (2), C4 (2), Journey (2), Architecture (2), ZenUML (3) = 13 fixtures, 65 golden files

Each agent:

1. Renders its fixtures across all 5 themes
2. Writes golden `.svg` files to `testdata/golden/`
3. Runs anomaly detection on each file:
   - Valid SVG: starts with `<svg`, ends with `</svg>`, has `viewBox`
   - Has diagram-specific elements (rects for flowcharts, lines for sequences, paths for pies, etc.)
   - Theme differentiation: all 5 themes produce distinct SVG for each fixture
   - Non-trivial size: file > 200 bytes
   - Accessibility: `role="img"` and `aria-label` present
   - Font-family: set on root `<svg>` element
4. If anomaly found: reports the issue, does NOT write the golden file

### Phase 2: Full Verification (main agent)

- `go test -run TestGolden` passes (all 250 files match)
- `go test ./...` passes (no regressions)
- Commit all golden files and the test harness

## Update Flow

When output intentionally changes:

```
go test -run TestGolden -update
git diff testdata/golden/   # review changes
git add testdata/golden/ && git commit
```

## What This Catches

- Layout coordinate changes from Sugiyama refactoring
- Theme color changes
- SVG structure changes (new attributes, reordered elements)
- Font-family, accessibility, marker definition regressions
- Any pipeline stage regression (parser, layout, render)

## What This Doesn't Do

- No pixel/visual comparison
- No float normalization
- No HTML diff report
- No per-element structural diffing
