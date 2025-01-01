# Registry

Trigger the **Publish to DigitalOcean Registry** workflow (`.github/workflows/registry.yml`) via `gh`.

## Ref

Use the text the user typed after `/registry` as the git ref (branch or tag). If they gave no ref, use the current branch from `git branch --show-current`.

## Steps

1. Set `REF` from the user's input, or from `git branch --show-current` when empty.
2. Run via the Shell tool with `full_network` permission:

```bash
gh workflow run registry.yml --ref "$REF"
```

3. Wait a few seconds, then fetch the run URL:

```bash
gh run list --workflow=registry.yml --limit 1
```

4. Report the ref used, whether the dispatch succeeded, and the run URL. If `gh` fails, show the error output.

## Examples

- `/registry` → runs against the current branch
- `/registry main` → runs against `main`
- `/registry v0.1.0` → runs against tag `v0.1.0`

## Notes

- Do not ask for confirmation before dispatching unless the ref is ambiguous.
- The workflow builds and pushes `registry.digitalocean.com/<DO_REGISTRY_NAME>/actx0` and prunes old tags (keeps 3).
