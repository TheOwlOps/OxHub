---
name: ox-systematic-debug
description: Use when debugging crashes, unexpected outputs, panics, or flaky test failures.
author: TheOwlOps
version: 1.0.0
---

# OX Systematic Root-Cause Debugger

Follow the 4-phase discipline: Understand, Isolate, Fix, Verify. Never guess-and-check.

## Workflow

1. **Phase 1: Reproduce & Capture Exact Failure**
   - Reproduce with minimal deterministic command.
   - Capture exact stack trace, exit code, and stdout/stderr output.
   - Do NOT edit code before reliable reproduction is established.

2. **Phase 2: Isolate Root Cause (Binary Search / Trace)**
   - Trace back from the panic/crash line to the first corrupted state.
   - Inspect variable bounds, nil pointers, slice bounds, type assertions.
   - Add minimal debug assertion or print if runtime debugger is unavailable.

3. **Phase 3: Minimal Surgical Fix**
   - Address the root flaw, not the surface symptom (don't just add `if val != nil` if `val` was never supposed to be nil).
   - Keep the diff minimal and self-contained.

4. **Phase 4: Regression Test & Verification**
   - Write ONE deterministic test case asserting the fixed condition.
   - Run the entire test suite to ensure zero regressions.
