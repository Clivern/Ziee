# Basecamp Chat

Post the user's message to the Basecamp chat integration.

## Message

Use the text the user typed after `/basecamp` as the message to post. Do not include the command name or meta-instructions — only the content they want in Basecamp.

## Steps

1. Take that message verbatim.
2. Run this curl command via the Shell tool with `full_network` permission, replacing `MESSAGE` with the user's text:

```bash
curl --data-urlencode "content=MESSAGE" \
  'https://app.basecamp.com/6175181/integrations/CiPwv9zug5RUxRiZcHBCk72H/buckets/47706023/chats/9995844475/lines'
```

3. Confirm success or report the curl error output.

## Examples

- `/basecamp deploy finished` → content is `deploy finished`
- `/basecamp PR #42 is ready for review` → content is `PR #42 is ready for review`

## Notes

- Do not ask for confirmation before posting unless the message is ambiguous.
- If the message contains HTML entities or special characters, `--data-urlencode` handles encoding.
- Do not modify or summarize the message unless the user asks you to.
