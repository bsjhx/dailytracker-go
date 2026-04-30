# Release Process

This document describes how to create a new release of DailyTracker.

## Versioning

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR** version (v1.0.0 → v2.0.0): Breaking changes
- **MINOR** version (v1.0.0 → v1.1.0): New features, backwards compatible
- **PATCH** version (v1.0.0 → v1.0.1): Bug fixes, backwards compatible

## Creating a Release

### 1. Update Version

Update the version constant in `internal/version/version.go`:

```go
const Version = "0.2.0"
```

### 2. Update Changelog

Add changes to `CHANGELOG.md`:

```markdown
## [0.2.0] - 2026-04-30

### Added
- New feature X
- New feature Y

### Fixed
- Bug fix Z

### Changed
- Improvement A
```

### 3. Commit Changes

```bash
git add internal/version/version.go CHANGELOG.md
git commit -m "Bump version to v0.2.0"
```

### 4. Create Git Tag

```bash
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin develop  # or main
git push origin v0.2.0
```

### 5. Create GitHub Release

1. Go to https://github.com/bsjhx/dailytracker-go/releases
2. Click "Draft a new release"
3. Choose the tag you just pushed (v0.2.0)
4. Title: `v0.2.0`
5. Description: Copy from CHANGELOG.md
6. Click "Publish release"

The version badge in README.md will automatically update!

## Automated Releases (Optional)

Consider setting up GitHub Actions to automate this:
- Create `.github/workflows/release.yml`
- Automatically build binaries
- Attach binaries to GitHub releases
- Build and push Docker images

## First Release (Now)

To create your first release:

```bash
# Tag the current state as v0.1.0
git tag -a v0.1.0 -m "Initial release v0.1.0

Features:
- Daily tracking for work and personal scores
- Weekly statistics
- User authentication
- File-based migrations
- Docker support"

git push origin develop
git push origin v0.1.0
```

Then create a GitHub Release from this tag.
