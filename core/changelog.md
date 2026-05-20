# Changelog

All notable changes to radlang are documented here.
Format based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [v0.1.0] - 2026-05-19

### Features
- add tree-walking interpreter (interpreter) (#10)
- add float literal support, split NUMBER into INT and FLOAT (lexer) (#9)
- adds semantic analysis pass (semantic) (#7)
- support underscores in identifiers (lexer) (#6)
- implements parsing logic (parser) (#4)
- adds parser helpers (parser) (#2)

### Bug Fixes
- fixes case when underscore in identifier is treated as illegal (lexer) (#6)

### Refactors
- split monolith into packages (#11)
- moves shared token definitions to token.go (#3)
- use int for TokenType, fix string/number/illegal token handling (lexer) (#1)

### Tests
- adds CI validation tests (#5)

### Chores / Build
- add PR template and wire analyzer into compile test (#8)
