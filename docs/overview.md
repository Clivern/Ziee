# Ziee overview

Ziee is the merge layer for agent-scale delivery. Pull requests are queued, checked, and merged without a human on the merge button.

## Merge queue

Speculative CI runs on temporary stack branches. Nothing lands on real `main` until the front of the line is green.

- `mq/1` = `main` + PR #1
- `mq/2` = `main` + PR #1 + PR #2
- and so on, up to `max_parallel_checks`

If `mq/1` is green, #1 merges. If it fails, #1 is dequeued and later stacks are rebuilt without it.

Modes:

| Mode | Behavior |
| --- | --- |
| `serial` | Stacked branches, advance one PR at a time |
| `parallel` | Same stacks, CI for several slots at once |
| `isolated` | Each PR tested against `main` alone |

## Queue state labels

| Label | Meaning |
| --- | --- |
| `state/queued` | Waiting in line |
| `state/checking` | Speculative CI is running |
| `state/dequeued` | Kicked out or removed from the queue |

## Commands

Comment on a pull request:

- `@ziee queue` — enter the first matching queue rule
- `@ziee queue hotfix` — enter the hotfix line (SRE only)
- `@ziee dequeue` — leave the queue
- `@ziee requeue` — dequeue then queue again
- `@ziee rebase` — rebase onto latest `main`
- `@ziee update` — merge latest `main` into the PR branch
- `@ziee refresh` — re-run triage and re-evaluate `queue_when`

## Priority

- `hotfix` jumps the line and rebuilds in-flight speculative checks
- Human PRs sit at medium
- Bot PRs (`ziee-bot`) wait behind humans
