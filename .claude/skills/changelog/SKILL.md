---
name: changelog
description: >
  Populate src/changelog.md by diffing the current main branch against the previous git version tag.
  Fetches commit history via git log and GitHub PR data via `gh` CLI, groups changes by
  Conventional Commits type, and writes a new version section. Use when user says "update
  changelog", "generate changelog", "populate changelog", "/changelog", or prepares a release.
---

# Changelog Skill

Populate `src/changelog.md` with a new version entry based on commits since the last git tag.

## Process

1. **Determine version range**
   - Run `git tag --sort=-version:refname | head -1` to find the previous tag (e.g. `v0.2.0`)
   - If no tags exist, diff from the first commit (`git rev-list --max-parents=0 HEAD`)
   - Ask the user what the new version number should be (e.g. `v0.3.0`) unless they already said

2. **Collect commits**
   ```bash
   git log <prev-tag>..HEAD --oneline --no-merges
   ```
   Also fetch merged PR titles for context:
   ```bash
   gh pr list --state merged --base main --limit 50 --json number,title,mergedAt,labels \
     --search "merged:>=$(git log <prev-tag> -1 --format=%aI 2>/dev/null || echo 1970-01-01)"
   ```

3. **Categorize by Conventional Commits type**

   | Section | Commit prefixes |
   |---|---|
   | Breaking Changes | `feat!`, `fix!`, any `BREAKING CHANGE` footer |
   | Features | `feat` |
   | Bug Fixes | `fix` |
   | Performance | `perf` |
   | Refactors | `refactor` |
   | Tests | `test` |
   | Chores / Build | `chore`, `build`, `ci` |
   | Docs | `docs` |

   Commits that don't match any prefix: place under **Other**.
   Omit sections that have no entries.

4. **Format the new entry**

   ```markdown
   ## [v0.3.0] - YYYY-MM-DD

   ### Breaking Changes
   - description (#PR)

   ### Features
   - description (#PR)

   ### Bug Fixes
   - description (#PR)
   ```

   Rules:
   - Date = today (`date +%Y-%m-%d`)
   - Link PR numbers when available: `(#42)`
   - Strip commit hash prefix from `git log --oneline` output
   - Strip conventional commit prefix from bullet text (`feat(lexer): add X` → `add X` under Features, scope in parens if useful: `add X (lexer)`)
   - Keep bullets short — one line each

5. **Write to src/changelog.md**
   - If `src/changelog.md` doesn't exist, create it with a header first:
     ```markdown
     # Changelog

     All notable changes to radlang are documented here.
     Format based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
     ```
   - Prepend the new version section directly after the header (before any existing entries)
   - Do NOT overwrite existing entries

6. **Show diff and confirm** before writing — print the new section and ask user to confirm.

## Edge Cases

- If `gh` is not authenticated or unavailable: skip PR linking, proceed with git log only
- If HEAD == prev tag (no new commits): tell the user, do nothing
- Merge commits: skip them (`--no-merges` in git log)
- Squash merges that lose conventional prefix: use the PR title as fallback for categorization

## Boundaries

Only writes `src/changelog.md`. Does not create git tags, does not push, does not open PRs.
After writing, remind the user to tag the release: `git tag v<version> && git push --tags`.
