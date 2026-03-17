# Slack Message Action

This action sends a message to a Slack channel.

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

      - uses: ./.github/actions/slack-msg
        id: slack-01
        with:
          slack-channel: "@UCKPL50JY"
          slack-bot-token: ${{ secrets.SLACK_BOT_TOKEN }}
          message: "test 01"

      - uses: ./.github/actions/slack-msg
        id: slack-02
        with:
          slack-channel: "@UCKPL50JY"
          slack-bot-token: ${{ secrets.SLACK_BOT_TOKEN }}
          message: "test 02"
          thread-timestamp: ${{ steps.slack-01.outputs.thread-timestamp }}

      - uses: ./.github/actions/slack-msg
        id: slack-03
        with:
          slack-channel: "@UCKPL50JY"
          slack-bot-token: ${{ secrets.SLACK_BOT_TOKEN }}
          message: |
            test 03
            test 03
            test 03
          thread-timestamp: ${{ steps.slack-01.outputs.thread-timestamp }}
```
