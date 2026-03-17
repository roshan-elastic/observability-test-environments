# Slack Message File Action

This action sends a message to a Slack channel from a payload file.

```yaml
---
name: Test Workflow
on:
  workflow_dispatch:

permissions:
  pull-requests: read
  contents: write
  id-token: write
  packages: write

jobs:
  test-1:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - uses: ./.github/actions/slack-msg-file
        id: slack-file-01
        with:
          slack-channel: "@UCKPL50JY"
          slack-bot-token: ${{ secrets.SLACK_BOT_TOKEN }}
          payload-file-path: ${{ github.workspace }}/slack-msg.json

      - uses: ./.github/actions/slack-msg-file
        id: slack-file-02
        with:
          slack-channel: "@UCKPL50JY"
          slack-bot-token: ${{ secrets.SLACK_BOT_TOKEN }}
          payload-file-path: ${{ github.workspace }}/slack-msg.json
          thread-timestamp: ${{ steps.slack-file-01.outputs.thread-timestamp }}
```
