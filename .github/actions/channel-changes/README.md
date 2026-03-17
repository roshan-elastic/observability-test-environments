# Channel Changes

Action to process the changes in the `.ci/update-versions/update_channels.json` file.

```yaml
  channel-changes:
    runs-on: ubuntu-latest
    outputs:
      changedChannels: ${{ steps.changed-channels.outputs.changedChannels }}
      currentChannels: ${{ steps.current-channels.outputs.channels }}
      previousChannels: ${{ steps.previous-channels.outputs.channels }}
      dayOfWeek: ${{ steps.changed-channels.outputs.dayOfWeek }}
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
          fetch-tags: false
          ref: ${{ github.sha || 'main' }}
      - uses: ./github/actions/changed-channels
        id: changed-channels
``````
