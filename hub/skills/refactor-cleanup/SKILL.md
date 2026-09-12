---
name: ox-refactor-cleanup
description: Use when refactoring messy files, simplifying bloated functions, and deleting legacy debt.
author: TheOwlOps
version: 1.0.0
---

# OX Refactoring & Dead-Code Elimination

Eliminate boilerplate, streamline architectures, and remove technical debt while preserving exact behavior.

## Execution Rules

1. **Delete Before Add**
   - Identify unreferenced functions, dead variables, deprecated packages.
   - Delete obsolete structs/types and their test leftovers.

2. **Flatten Control Flow**
   - Replace deep nested `if-else` trees with guard clauses / early returns.
   - Extract convoluted helper logic into small, pure functions (< 30 lines).

3. **Interface Simplification**
   - Eliminate interfaces with only 1 implementation unless required for mocking external I/O.
   - Accept interfaces, return concrete structs (Go idiom).

4. **Safety Net**
   - Ensure existing tests pass before touching any code.
   - Re-run test suite after every small edit step (`go test ./...` or `pnpm test`).
