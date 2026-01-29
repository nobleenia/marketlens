# Git workflow

## Branch model
- `main`: always releasable
- feature branches: `feat/<area>-<short>`
- fixes: `fix/<issue>`
- docs: `docs/<topic>`

## PR rules
- CI must pass
- Include tests for core logic changes
- Update docs if behaviour changes
- Squash merge

## Releases
- Tags follow SemVer: `v0.x.y` during MVP
