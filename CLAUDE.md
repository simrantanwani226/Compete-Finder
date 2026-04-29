# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## What This Is

A 6-module learn-by-building Go course. Students build a competitor finder (scrape YC startups, score with TF-IDF, serve via Connect RPC) while learning Go from first principles. Each module is launched with `/module-N` and follows a teach → quiz → exercises flow.

## Commands

```bash
go test ./...                                  # Run all tests
go test ./internal/provider/...                # Run tests for a package
go run tools/progress/main.go                  # Show student progress
go run tools/progress/main.go start N          # Start module N
go run tools/progress/main.go quiz N S T       # Record quiz (score S out of T)
go run tools/progress/main.go exercise N E     # Record exercise E in module N
go run tools/progress/main.go complete N       # Mark module N complete
```

## Architecture

Students build code incrementally across modules:

```
internal/provider/         → Startup struct, Provider interface (Module 1)
internal/provider/yc/      → YC API adapter, HTTP client (Module 2)
internal/matcher/           → Tokenizer, TF-IDF scoring (Module 3)
internal/heatmap/           → Sector filtering, market heatmap (Module 4)
proto/compete/v1/           → Protobuf definitions (Module 5)
gen/compete/v1/             → Generated Connect RPC code (Module 5)
internal/handler/           → Connect service handlers (Module 6)
cmd/server/                 → Main entry point, server wiring (Module 6)
```

## Module Teaching System

- `course/module_0N_*.md` — course content (concepts, quiz answers, exercises)
- `.claude/commands/module-N.md` — Claude Code slash commands defining the teaching flow
- `tools/progress/main.go` — progress tracking, state in `progress.json` (gitignored)
- All learning infrastructure is gitignored — only the student's Go code gets committed

## Teaching Flow (for module commands)

Each `/module-N` command follows this sequence:

1. Read course content and existing code
2. Teach section by section using the **student-builds-everything** approach
3. Quiz: 5 questions, one at a time, 80% to pass
4. Exercises: guided implementation, progress tracked per exercise
5. Completion: `go run tools/progress/main.go complete N`

### Student-Builds-Everything Approach

**The student writes ALL implementation code. The assistant writes tests and explains concepts.**

For each section:

1. **Explain** the concept clearly (what it is, why it matters, how it works)
2. **Write a test** (`_test.go` file) that defines expected behavior
3. **Tell the student what to build** — file path, exports, types, behavior in plain English
4. **Wait** for the student to write the code and run `go test`
5. **If tests pass** — briefly discuss, move on
6. **If tests fail** — give a hint (not the answer), let them fix it

Rules:

- **NEVER** write implementation files for the student — only `_test.go` files
- **NEVER** show complete function bodies — show signatures, types, describe logic in words
- **ONE section at a time** — do not proceed until tests pass
- **Wait for student input** between every section
- Tests use standard library `testing` package only — no external frameworks

## Code Conventions

- Standard Go project layout (`internal/`, `cmd/`, `gen/`)
- `go test ./...` for all tests
- Table-driven tests where appropriate
- Error wrapping with `fmt.Errorf("context: %w", err)`
- Context propagation for HTTP calls
- No external dependencies except connectrpc and protobuf
