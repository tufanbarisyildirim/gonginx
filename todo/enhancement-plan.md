# Gonginx Enhancement Tasks

## Scope
- [x] Improve reliability, correctness, API safety, and maintainability based on scan findings and repros.

## Primary Goals
- [x] Remove crash paths and resource leaks.
- [x] Make parser and dumper behavior deterministic and side-effect free.
- [x] Harden public APIs and extension points against invalid inputs and custom wrappers.
- [x] Add regression coverage for all high-risk behaviors.

## Baseline Validation
- [x] Confirm `go test ./...` is green before changes.
- [x] Confirm `go vet ./...` is green before changes.
- [x] Keep baseline issue-fixture behavior validated while refactoring.

## Milestone P0: Safety and Correctness

### P0.1 Include Cycle Handling and Parser File Lifecycle Safety (Critical)
- [x] Replace include recursion tracking keyed by `*config.Include` with canonical include path keys.
- [x] Introduce include state tracking (`visiting`, `done`) keyed by cleaned absolute include path.
- [x] Canonicalize include paths via `filepath.Clean` and `filepath.Abs` before map lookup.
- [x] Prevent recursion when include path is already in `visiting` state.
- [x] Decide cycle behavior:
- [x] Option A: skip cyclic branch and continue parse.
- [x] Option B: return explicit include-cycle error.
- [x] Cache parsed include configs by canonical path to avoid duplicate reparses.
- [x] Refactor parser close behavior so file descriptors close on all parse paths.
- [x] Add `defer`-based close handling in `Parse()`.
- [x] Ensure close errors are surfaced without masking parse errors.
- [x] Verify nested include parser instances do not leak open files on parse errors.
- [x] Add regression test: `TestParser_IncludeCycle_DoesNotLoop`.
- [x] Add regression test: `TestParser_IncludeCycle_NoFDExhaustion`.
- [x] Add regression test: `TestParser_IncludeDuplicate_UsesCache`.
- [x] Add regression test: `TestParser_FileClosedOnParseError`.
- [x] Validate no `too many open files` in cycle fixture.

### P0.2 Replace Panic-Based Failures with Error Returns (Critical)
- [x] Remove panic path for unterminated quoted strings in lexer.
- [x] Remove panic path for unterminated Lua code blocks in lexer.
- [x] Thread lexical scan errors through parser return path with line/column metadata.
- [x] Refactor `config.NewInclude` to validate args before indexing.
- [x] Return descriptive errors instead of panic for include arity violations.
- [x] Return descriptive errors instead of panic for include block misuse.
- [x] Refactor `config.NewUpstream` to validate required upstream name before indexing.
- [x] Replace panic in dumper include write path with explicit error return.
- [x] Audit parser/config/dumper for any remaining panic paths reachable from malformed input.
- [x] Add regression test: unterminated quote returns error and does not panic.
- [x] Add regression test: unterminated Lua block returns error and does not panic.
- [x] Add regression test: `include;` returns validation error and does not panic.
- [x] Add regression test: `include a.conf b.conf;` returns validation error and does not panic.
- [x] Add regression test: missing upstream name returns validation error.
- [x] Confirm malformed user input cannot panic host application.

### P0.3 Directive Validation Consistency and Message Quality (High)
- [x] Standardize min/max parameter validation before indexing in directive constructors.
- [x] Ensure all directive validation errors include directive name.
- [x] Fix low-quality/incomplete errors (for example in `config/location.go`).
- [x] Normalize wording and style of constructor errors for consistency.
- [x] Keep error phrasing stable enough for regression tests.
- [x] Add constructor-level tests for invalid arity across include/upstream/location.
- [x] Add parser integration tests for stable, descriptive error messages.
- [x] Confirm no index-out-of-range remains in constructor argument handling.

## Milestone P1: Determinism and Data Integrity

### P1.1 Non-Mutating Sort Behavior in Dumper (High)
- [x] Stop sorting directive slices in-place in `DumpBlock`.
- [x] Copy directives before sorting (`append([]config.IDirective(nil), ...)`).
- [x] Preserve current sorted output format while avoiding AST mutation.
- [x] Add regression test flow:
- [x] Dump unsorted output.
- [x] Dump sorted output.
- [x] Dump unsorted output again.
- [x] Assert first and third unsorted outputs are identical.
- [x] Confirm `SortDirectives` affects only rendered output, not in-memory order.

### P1.2 Lua Formatting Safety Without Semantic Corruption (High)
- [x] Remove global `#` -> `--` replacement across entire Lua text.
- [x] Implement comment-marker conversion that only applies outside string literals.
- [x] Use reversible strategy so only converted comments are restored.
- [x] Ensure literal strings containing `--` or `#` remain unchanged after dump.
- [x] Preserve existing indentation behavior for formatted Lua output.
- [x] Ensure formatter failure falls back to original Lua text safely.
- [x] Add regression test: Lua strings with `"--"` and `"#"` round-trip unchanged.
- [x] Add regression test: Lua comments still format correctly.
- [x] Add regression test: invalid Lua does not corrupt emitted content.
- [x] Validate existing issue fixtures involving Lua remain green.

### P1.3 Type Assertion Hardening for Extension Points (High)
- [x] Replace hard `.(*Type)` assertions in public paths with checked assertions.
- [x] Return explicit errors for type mismatch rather than panic.
- [x] Update include wrapper parse flow to validate wrapper output type before cast.
- [x] Update `FindUpstreams` behavior for mixed/custom directive trees:
- [x] Option A: skip non-upstream entries safely.
- [x] Option B: add strict variant that returns typed error.
- [x] Update dumper include write path to safely handle non-include entries.
- [x] Audit parser/config/dumper for remaining unsafe assertions in runtime paths.
- [x] Add regression test: wrapper returns non-include type and parser errors cleanly.
- [x] Add regression test: mixed upstream trees do not panic.
- [x] Add regression test: dumper include write path handles unexpected types safely.

## Milestone P2: API Consistency and Developer Experience

### P2.1 Parent Pointer Model Consistency (Medium)
- [x] Define explicit parent semantics for all directive/block relationships.
- [x] Set root-level directive parent to `nil` (not self).
- [x] Ensure nested directive parent points to enclosing wrapper/directive.
- [x] Remove self-parent assignment behavior in parser.
- [x] Validate parent assignments for HTTP, Server, Location, Upstream, Include.
- [x] Add regression test: root-level parent behavior.
- [x] Add regression test: nested parent chains for common directive types.
- [x] Add regression test: no self-parent cycles.

### P2.2 Documentation and Migration Notes (Medium)
- [x] Update `README.md` with parser error model (error return, no panic contract).
- [x] Update `GUIDE.md` with include cycle handling and include dedupe behavior.
- [x] Update `GUIDE.md` with sort behavior (output-only, non-mutating).
- [x] Update docs with parent semantics expectations.
- [x] Add migration notes for behavior changes that may impact consumers.
- [x] Update `CONTRIBUTING.md` where workflow/tests changed.

## Test and Verification
- [x] Add targeted unit tests for each issue class.
- [x] Add targeted integration tests for parser include and Lua paths.
- [x] Keep issue fixture tests (`17`, `20`, `22`, `50`) passing.
- [x] Run `go test ./...` after each milestone.
- [x] Run `go vet ./...` after each milestone.
- [x] Add CI gate for `go test ./... -race` (at least on default branch).
- [x] Add parser fuzz test for malformed snippet/token-edge robustness.

## Delivery Plan (PR Task Checklist)

### PR-1 (P0.1 + P0.2 + P0.3)
- [x] Implement include cycle and file lifecycle safety changes.
- [x] Implement panic-to-error refactors in lexer/parser/config/dumper.
- [x] Implement directive validation consistency improvements.
- [x] Add and pass all PR-1 regression tests.
- [x] Run `go test ./...` and `go vet ./...`.

### PR-2 (P1.1 + P1.3)
- [x] Implement non-mutating sort behavior.
- [x] Implement type assertion hardening.
- [x] Add and pass all PR-2 regression tests.
- [x] Run `go test ./...` and `go vet ./...`.

### PR-3 (P1.2)
- [x] Implement Lua formatting safety changes.
- [x] Add and pass Lua semantics/regression tests.
- [x] Run `go test ./...` and `go vet ./...`.

### PR-4 (P2.1 + P2.2)
- [x] Implement parent model consistency changes.
- [x] Update docs and migration notes.
- [x] Add and pass all PR-4 regression tests.
- [x] Run `go test ./...` and `go vet ./...`.

## Timeline Tracking
- [x] Target total effort: 6.25 to 7.5 engineering days.
- [x] Re-estimate after PR-1 based on implementation complexity discovered.
- [x] Re-estimate after PR-3 based on Lua formatter edge-case burden.

## Risk and Decision Tracking
- [x] Finalize include cycle policy (skip vs explicit error): skip by default, opt-in error with `WithIncludeCycleErr()`.
- [x] Finalize Lua comment style policy (preserve source vs normalize): preserve source style where safe; fallback to original code on formatter failure.
- [x] Finalize parent semantics compatibility strategy (hard switch vs compatibility option): hard switch, no compatibility flag.
- [x] Finalize `FindUpstreams` behavior policy (strict vs permissive): permissive `FindUpstreams()` and strict `FindUpstreamsStrict()`.

## Done Criteria
- [x] No known malformed-input repro causes panic in parser/config/dumper.
- [x] Include cycles are bounded and safe.
- [x] Parser closes files on all parse paths.
- [x] Sort mode is side-effect free.
- [x] Lua dump preserves literal semantics.
- [x] All new regression tests pass.
- [x] Existing tests remain green.
- [x] Documentation reflects shipped behavior.
