# Releasing

This project uses [Semantic Versioning](https://semver.org/) and immutable releases.

## Version Format

```
vMAJOR.MINOR.PATCH
```

- **MAJOR** — breaking changes
- **MINOR** — new features, backward compatible
- **PATCH** — bug fixes, backward compatible

## Creating a Release

1. Ensure `main` branch is clean and tested
2. Create and push a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

3. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Create a GitHub Release
   - Upload artifacts with checksums

## Immutability Policy

**Releases are immutable.** Once a version is released:

- ❌ Never delete a tag
- ❌ Never force-push to move a tag
- ❌ Never edit release artifacts
- ✅ Create a new version to fix issues

### If a release has a bug

```bash
# Wrong: don't do this
git tag -d v0.1.0
git push origin :refs/tags/v0.1.0

# Right: release a new version
git tag v0.1.1
git push origin v0.1.1
```

### If a release was made from wrong commit

Release a new patch version from the correct commit. Document the issue in the new release notes if needed.

## Pre-release Versions

For testing releases before making them official:

```bash
git tag v0.2.0-rc.1
git push origin v0.2.0-rc.1
```

## Checking Existing Tags

```bash
git tag -l "v*"
git describe --tags --abbrev=0
```
