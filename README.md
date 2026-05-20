# Radlang

> A compiled systems programming language inspired by Go and Python.

![Status](https://img.shields.io/badge/status-experimental-orange)
![Language](https://img.shields.io/badge/language-Go-blue)
![License](https://img.shields.io/badge/license-MIT-green)

RadLang combines what Go gets right — simplicity, strong static typing, concurrency — with the ergonomics that make Python pleasant to write.

Not a Go clone. Not a Python clone. A practical language that feels lightweight and expressive without becoming a bloated research project.

---

## Table of Contents

- [Why RadLang](#why-radlang-exists)
- [Syntax](#syntax)
- [Language Goals](#language-goals)
- [Current Pipeline](#current-state)
- [Roadmap](#roadmap)
- [Philosophy](#philosophy)
- [Inspiration](#inspiration)

---

## Why RadLang Exists

**Go** is clean, fast, and productive — but some parts feel unnecessarily rigid:

- map syntax is verbose
- array/slice ergonomics are limited
- some trivial Python operations require boilerplate

**Python** is expressive and ergonomic — but:

- dynamic typing becomes painful at scale
- performance is limited
- concurrency is messy
- deployment is heavier than it needs to be

RadLang sits in the middle:

| Feature       | Go                  | RadLang                         | Python                     |
| ------------- | ------------------- | ------------------------------- | -------------------------- |
| Typing        | Static              | Static                          | Dynamic                    |
| Execution     | Compiled            | Compiled                        | Interpreted                |
| Concurrency   | Goroutines          | Lightweight concurrency planned | Async/threading            |
| Collections   | Verbose             | Python-inspired, typed          | Very ergonomic             |
| Syntax        | Minimal             | Minimal                         | Highly expressive          |

---

## Syntax

RadLang uses Go-style syntax with type inference support.

### Variable Declaration

```rad
var age int
var name string = "rad"
```

### Type Inference

```rad
x := 10
name := "radlang"
```

### Typed Collections

```rad
var nums [int]

nums[-1]        // negative indexing
nums[1:5]       // slicing
nums.append(10)
```

Collections are statically typed — no runtime surprises.

---

## Language Goals

### Strong Static Typing

Strict compile-time type checking. No hidden magic. No runtime guessing.

### Ergonomic Collections

Python-inspired array operations with static typing. Negative indexing, slicing, and built-in methods — without sacrificing type safety.

### Simplicity First

One of Go's best qualities: most codebases look similar. RadLang follows the same philosophy — minimal syntax, readable code, small feature surface, no "clever" language tricks.

### Concurrency

Go's goroutines are a core inspiration. Lightweight concurrency primitives are planned as a first-class language feature.

---

## Current State

RadLang is in early development. Current pipeline:

```
Source
  ↓
Lexer
  ↓
Parser
  ↓
Semantic Analyzer
  ↓
Tree-Walking Interpreter
```

The interpreter validates language design decisions quickly before moving to proper code generation. LLVM-based compilation is planned.

---

## Roadmap

### MK 1 — Foundation ✅
- [x] Lexer
- [x] Parser
- [x] AST generation
- [x] Semantic analysis
- [x] Tree-walking interpreter

### MK 2 — Language Completeness
- [ ] Full type system (structs, interfaces)
- [ ] Functions with multiple return values
- [ ] Error handling primitives
- [ ] Standard library (I/O, strings, math)
- [ ] Module/import system

### MK 3 — Code Generation
- [ ] IR / bytecode emission
- [ ] LLVM backend
- [ ] Cross-platform compilation targets

### MK 4 — Concurrency
- [ ] Lightweight goroutine-style primitives
- [ ] Channels
- [ ] Scheduler

### MK 5 — Tooling
- [ ] Language server (LSP)
- [ ] Formatter
- [ ] Package manager
- [ ] Playground / REPL

---

## Philosophy

RadLang intentionally avoids becoming:

- an overengineered academic language
- "Rust but harder"
- a syntax experiment
- a framework disguised as a language

Target:

> "A practical systems language that is pleasant to write daily."

---

## Inspiration

- [Go](https://golang.org)
- [Python](https://python.org)
- [Kotlin](https://kotlinlang.org)
- [Rust](https://rust-lang.org)

---

## Status

Experimental. Under active development. Expect breaking changes.
