---
name: ox-code-reviewer
description: Use when conducting thorough, multi-axis code reviews before merging or pushing changes.
author: TheOwlOps
version: 1.0.0
---

# OX High-Velocity Code Reviewer

Ultra-terse, high-signal code review standard focused on correctness, memory allocation, edge cases, and maintainability.

## Review Axes

1. **YAGNI & Complexity**
   - Is there unused abstraction, unneeded factory/interface, or premature scaffolding?
   - Delete dead code before merging new additions. Shortest working diff wins.

2. **Concurrency & Thread Safety**
   - Go: Race detector `go test -race ./...`. Goroutines must have context cancellation or exit channels; no goroutine leaks. Check mutex lock/unlock defer semantics.
   - Node/TS: Promise rejection handling, race conditions in async loops, concurrency batching (`p-limit`).

3. **Error Handling & Resilience**
   - No swallowed errors (`_ = err` or empty `catch {}`).
   - Wrap context into errors (`fmt.Errorf("reading config: %w", err)`).
   - Fast-fail on invalid parameters at public API boundaries.

4. **Performance & Allocations**
   - Preallocate slice/map capacities when length is known (`make([]T, 0, len(items))`).
   - String concatenation in tight loops: use `strings.Builder`.
   - Minimize memory allocations in hot paths.

## Output Structure

- **Verdict**: APPROVE | REQUEST_CHANGES | BLOCK
- **Summary**: 1-2 bullet points max.
- **Issues (if any)**:
  - `[severity]` `file:line` - Reason. Required fix.
