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

Comment on a pull request to work the queue:

- `@zieeai queue` — enter the first queue rule whose `queue_when` matches
- `@zieeai queue hotfix` — enter the hotfix line (SRE only)
- `@zieeai dequeue` — leave the queue
- `@zieeai requeue` — dequeue, then queue again with fresh speculative checks
- `@zieeai rebase` — rebase the pull request onto latest `main`
- `@zieeai update` — merge latest `main` into the pull request branch
- `@zieeai refresh` — re-evaluate `queue_when` after approvals or CI change. Triage runs again too.

Comment on a pull request to triage it:

- `@zieeai label bug security` — add the named labels
- `@zieeai unlabel bot` — remove the named labels
- `@zieeai assign clivern` — assign the named users
- `@zieeai unassign clivern` — remove the named assignees
- `@zieeai reviewers maya` — request the named reviewers
- `@zieeai close` — close the pull request (dequeues it if it was in line)
- `@zieeai reopen` — reopen the pull request
- `@zieeai spam` — label spam, close, and block the author

Comment on an issue:

- `@zieeai label bug security` — add the named labels
- `@zieeai unlabel bot` — remove the named labels
- `@zieeai assign clivern` — assign the named users
- `@zieeai unassign clivern` — remove the named assignees
- `@zieeai close` — close the issue
- `@zieeai reopen` — reopen the issue
- `@zieeai spam` — label spam, close, and block the author from further issues
- `@zieeai summarize` — summarize the issue

## Priority

- `hotfix` jumps the line and rebuilds in-flight speculative checks
- Human PRs sit at medium
- Bot PRs (`ziee-bot`) wait behind humans
