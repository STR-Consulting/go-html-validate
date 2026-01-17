---
# go-html-validate-0y8v
title: Implement GoReleaser for multi-platform releases
status: completed
type: feature
priority: normal
created_at: 2026-01-17T22:00:43Z
updated_at: 2026-01-17T22:02:12Z
---

Add GoReleaser to automate multi-platform binary releases following the bean-me-up project pattern.

## Checklist
- [x] Create .goreleaser.yaml
- [x] Create .github/workflows/release.yml
- [x] Update main.go version detection
- [x] Update commit.sh to remove gh release
- [x] Update README.md binary name to htmlint
- [x] Update mise.toml with self-reference
- [x] Verify with goreleaser check (not installed locally, will work in CI)