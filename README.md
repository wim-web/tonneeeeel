# tonneeeeel

This CLI is an easy way to handle some ecs-related commands. (e.g. ecs exec, ssm start-session)

Normally, it is necessary to specify the cluster name, task ID, and container ID. However, this CLI tool allows you to specify them interactively.

## Demo

![](./_img/demo.gif)

## Main Features

```
Usage:
  tonneeeeel [command]

Available Commands:
  completion        Generate the autocompletion script for the specified shell
  exec              like ecs execute-command
  help              Help about any command
  portforward       like start-session --document-name AWS-StartPortForwardingSession
  remoteportforward like start-session --document-name AWS-StartPortForwardingSessionToRemote

Flags:
  -h, --help   help for tonneeeeel
```

1. `exec` command: Log in to the specified container, similar to AWS ECS Exec.
2. `portforward` command: Perform port forwarding to the specified container.
3. `remoteportforward` command: Perform port forwarding to a remote host via the specified container.

## Installation

Install the CLI tool with the following command:

```bash
go install github.com/wim-web/tonneeeeel
```

Or get binary from https://github.com/wim-web/tonneeeeel/releases

## License

This project is published under the MIT License. For more information, see the [LICENSE](LICENSE) file.
