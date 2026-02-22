# Releasing `ordo`

## Automated release flow

Releases are tag-driven via GitHub Actions.

1. Ensure `main` is green (tests pass in CI).
2. Create and push a semver tag (`vX.Y.Z`).
3. GitHub Actions builds binaries and creates a GitHub Release.

## Quick command

Use the helper script:

```bash
scripts/release-tag.sh v0.2.0
```

## What gets published

- `ordo` tarballs for:
  - `linux/amd64`
  - `linux/arm64`
  - `darwin/amd64`
  - `darwin/arm64`
- `SHA256SUMS` checksum file
- `install.sh` installer script
- Auto-generated GitHub release notes

## How users install from releases

See `INSTALL.md` for quick install and pinned versions.

Current one-liner:

```bash
curl -fsSL https://github.com/edsonjaramillo/ordo/releases/latest/download/install.sh | sh
```
