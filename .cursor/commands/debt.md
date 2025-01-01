# Debt

Mark a file as reviewed in `debt.md` (Go under package dirs, frontend under `web/`).

## File

Use the path the user typed after `/debt`. If they gave no path, use the currently focused / open editor file.

Normalize to a repo-relative path (strip leading `./` or absolute workspace prefix). Prefer the form already used in `debt.md` (e.g. `locale/locale.go`, `api/health.go`, `web/src/App.vue`).

## Steps

1. Resolve `FILE` from the user's input or the active editor path.
2. In `debt.md`, find the checklist line for that file:
   - `- [ ] FILE` → change to `- [x] FILE`
   - already `- [x] FILE` → leave as-is and say so
3. If the file is missing from `debt.md`, add `- [x] FILE` under the matching `##` section (create the section if needed), using the same style as neighboring entries.
4. Do not edit anything else in `debt.md`.
5. If this review revealed a durable convention (prefer X, avoid helper Y, lib gotcha), append a short note to root `SKILL.md` under the right section. Skip trivia.
6. Confirm with the path marked done.

## Examples

- `/debt` with `locale/locale.go` focused → marks `- [x] locale/locale.go`
- `/debt api/me.go` → marks `- [x] api/me.go`
- `/debt ./db/audit.go` → marks `- [x] db/audit.go`
- `/debt web/src/lib/auth.js` → marks `- [x] web/src/lib/auth.js`

## Notes

- Do not ask for confirmation.
- Only flip the checkbox for the resolved file; never mass-check a directory.
- Match the exact path string already in `debt.md` when present.
