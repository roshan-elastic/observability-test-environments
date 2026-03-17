# Slackbot

It's possible to start the tool as a Slack bot, which has been previously added to the Elastic workspace (see helpdesk case **# 00910462**, `ref:_00Db0H5KI._5008X24pQfZ:ref _`).

## Slack app details

The settings for the slack bot used in production can be seen in https://api.slack.com/apps/A035GN7PGAV/general

For local development, you can create your own Slack app and workplace so you can test your changes in isolation and with reassurance you won't break anything.

## Developing new features

The Slackbot is using the [`github.com/slack-go/slack`](http://github.com/slack-go/slack) Go package, which provides a Go implementation for Slack APIs.

It contains an `Examples` section, where you can find multiple implementations for the majority of the code in this bot: https://github.com/slack-go/slack/tree/master/examples


### Prerequisites

* `Golang` version `1.17` (defined in the `go.mod` file). Mandatory
* `Vault` to be able to access the existing secrets. (Optional)

### Local deployment
It's possible to start the bot in your local machine, and connect it to Slack. Simply follow these steps:

1. Run `make create-env` that will create the  `.env` file in the root directory of the tool with the values stored in the `secret/k8s/elastic-apps/apm/oblt-robot` Vault secret.

```
SLACK_APP_TOKEN="YourSlackSocketToken"
SLACK_AUTH_TOKEN="YourBotUserOAuthToken"
SLACK_CHANNEL_ID="YourSlackChannelID"
```

2. Run the application from your your local machine.

    a. Run `go run main.go --dry-run` to start the bot in dry-run mode, or
    b. If you are using **VSCode** you can benefit from a debug profile. Simply paste the below snippet into the `.vscode/launch.json` file in the root directory of the repository. Then go to the `Debug` icon in the left panel, and click on the `Debug oblt-cli Slack bot` option from the **RUN AND DEBUG** selector. For further information about debugging in VSCode, please read https://code.visualstudio.com/docs/languages/go#_debugging

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug oblt-cli Slack bot",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/tools/oblt-robot",
            "args": ["--verbose", "--dry-run"]
        }
    ]
}
```

In any case, you should see the application to start in Slack mode, and pings to the socket mode should appear in your terminal in a continuous manner:

> For the debug mode, the terminal would be the `DEBUG CONSOLE` from VSCode.

```
socketmode: 2022/04/01 10:33:48 socket_mode_managed_conn.go:493: Finished to receive message
socketmode: 2022/04/01 10:33:48 socket_mode_managed_conn.go:439: Starting to receive message
socketmode: 2022/04/01 10:33:48 socket_mode_managed_conn.go:336: Received WebSocket message: {"type":"hello","num_connections":1,"debug_info":{"host":"applink-78ccf745b6-m6hmk","build_number":3,"approximate_connection_time":18060},"connection_info":{"app_id":"A035GN7PGAV"}}
socketmode: 2022/04/01 10:33:52 socket_mode_managed_conn.go:561: WebSocket ping message received: Ping from applink-78ccf745b6-m6hmk
```

Now you are connected to Slack!

1. Go to the [`#observablt-bots`](https://elastic.slack.com/archives/C01APBTSM0V) channel and interact with the `@oblt-robot` bot: mentions and slash commands are available.


### Good reads
Please read the following resources to understand Slack APIs and the implementation of the bot.

- Slack APIs documentation: https://api.slack.com/docs
- Go examples: https://github.com/slack-go/slack/tree/master/examples

#### Websockets and Socket Mode

> The bot uses websockets and Socket Mode to communicate with Slack.

- https://api.slack.com/apis/connections/socket
- https://api.slack.com/apis/connections/socket-implement

#### Modal forms
> The bot creates a modal form for interacting with the user: ex. the creation of a new Cross-Cluster-Search cluster.

- https://api.slack.com/surfaces
- https://api.slack.com/surfaces/modals/using

#### User interactions

- https://api.slack.com/interactivity
- https://api.slack.com/interactivity/handling


### Use cases

#### Add a new mention event

If you'd like to add a new mention event, something like `@oblt-robots my-new-action`, please follow
the below steps:

1) Create `myNewAction.go` file under `slack/mentions`
  a) Implement `func MyNewAction(client *slack.Client, channelID string, user slack.User) error {`
  b) Implement `func IsMyNewActionMention(text string) bool {`
  c) Implement `func MyNewActionTitle() string {`
  d) Implement `func MyNewActionDescription() string {`
2) After implementing the above logic needed, then add the support in the `slack/events_messages.go`:
  a) Within the `handleAppMentionEvent` function add
```
  else if mentions.IsMyNewActionMention(event.Text) {
		return mentions.MyNewAction(client, event.Channel, *user)
	}
```
3) Add support for the `help` mention, so the new event is also shown when people ask for help.
  a) Add `messageText += fmt.Sprintf("\n*%s*: %s", MyNewActionTitle(), MyNewActionDescription())` in `slack/mentions/help.go`


### Release

The `oblt-robot` project is built and released to the internal docker registry for every commit merged in `main`, additionally, you can release any version locally or from the CI running the Workflow with the flag `true to release oblt-robot`.

### Deployment

This bot runs in Elastic Apps, so you can re-deploy a new version by running:

```bash
$ make deploy
```

If the new version contains new secrets, then it's better to `undeploy` and `deploy`, but be sure the environment and secrets are in place in Vault.
