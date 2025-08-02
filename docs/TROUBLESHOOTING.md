# Troubleshooting

## Logging

The language server produces detailed logs (including every LSP request
and response) which are send to stderr by default.

It may be helpful to share these logs when reporting bugs.

### stderr (default)

Most clients provide a way of inspecting these logs when server is launched
in the default "stdio" mode where stdout & stdin are used as communication
channels for LSP. For example:

<!-- TODO: Update this link when we get a better display and itemName. See https://github.com/opentofu/vscode-opentofu/issues/30 -->

[**OpenTofu VS Code Extension**](https://marketplace.visualstudio.com/items?itemName=opentofu.vscode-opentofu)

1. `View` -> `Output`\
   ![vscode-view-output-menu](./images/vscode-view-output-menu.png)
2. `Output` -> `OpenTofu`\
   ![vscode-output-pane](./images/vscode-output-pane.png)

### Logging to Files

Server logs can also be directed to files using **`-log-file=<filepath>`**
of the `serve` command.

```sh
$ tofu-ls serve -log-file='/tmp/tofu-ls-{{pid}}.log'
```

Clients which manage LS installation typically allow passing extra arguments.

### OpenTofu CLI Execution Logs

Given that the server may also execute OpenTofu itself, it may be useful
to collect logs from all these executions too. This is equivalent
to setting [`TF_LOG_PATH` variable](https://opentofu.org/docs/internals/debugging/).

This can be enabled via [`tofu.logFilePath` LSP settings](./SETTINGS.md#logfilepath-string).

Clients which manage LS installation typically expose this as a dedicated setting option.
For example:

[**OpenTofu VS Code Extension**](https://marketplace.visualstudio.com/items?itemName=opentofu.vscode-opentofu)

1. Open the command palette via `⌘/Ctrl + Shift + P`
2. ![vscode-open-settings-json](./images/vscode-open-settings-json.png)
3. Set `"tofu-ls.tofu.logFilePath"` to a file path, such as `/tmp/tf-exec-{{lsPid}}-{{method}}-{{timestamp}}.log`

### How To Share Logs

It is recommended to avoid pasting logs into the body of an issue,
unless you are trying to draw attention to a selected line or two.

It's always better to upload the log as [GitHub Gist](https://gist.github.com/)
and attach the link to your issue/comment, or [attach the file to your issue/comment](https://docs.github.com/en/github/managing-your-work-on-github/file-attachments-on-issues-and-pull-requests).

### Sensitive Data

Be careful and review your logs before submitting the issue. Logs may contain sensitive data (such as content of the files being edited in the editor).

### Log Rotation

Keep in mind that the language server itself does not have any log rotation facility,
but the destination path will be truncated before being logged into.

Static paths may produce large files over the lifetime of the server and
templated paths (as described below) may produce many log files over time.

### Log Path Templating

Log paths support template syntax. This allows for separation of logs while accounting for:

- multiple server instances
- multiple clients
- multiple OpenTofu executions which may happen in parallel

**`-log-file`** flag supports the following functions:

- `timestamp` - current timestamp (formatted as [`Time.Unix()`](https://golang.org/pkg/time/#Time.Unix), i.e. the number of seconds elapsed since January 1, 1970 UTC)
- `pid` - process ID of the language server
- `ppid` - parent process ID (typically editor's or editor plugin's PID)

**`tofu.logFilePath`** option supports the following functions:

- `timestamp` - current timestamp (formatted as [`Time.Unix()`](https://golang.org/pkg/time/#Time.Unix), i.e. the number of seconds elapsed since January 1, 1970 UTC)
- `lsPid` - process ID of the language server
- `lsPpid` - parent process ID of the language server (typically editor's or editor plugin's PID)
- `method` - [`tofu-exec`](https://pkg.go.dev/github.com/opentofu/tofu-exec) method (e.g. `Format` or `Version`)

The path is interpreted as [Go template](https://golang.org/pkg/text/template/), e.g. `/tmp/tofu-ls-{{timestamp}}.log`.

## CPU Profiling

If the bug you are reporting is related to high CPU usage it may be helpful
to collect and share CPU profile which can be done via `cpuprofile` flag.

For example you can modify the launch arguments in your editor to:

```sh
$ tofu-ls serve \
	-cpuprofile=/tmp/tofu-ls-cpu.prof
```

Clients which manage LS installation typically allow passing extra arguments.
See [Logging to Files](#logging-to-files) section above for examples.

The target file will be truncated before being written into.

### Path Templating

Path supports template syntax. This allows for separation of logs while accounting for multiple server or client instances.

**`-cpuprofile`** supports the following functions:

- `timestamp` - current timestamp (formatted as [`Time.Unix()`](https://golang.org/pkg/time/#Time.Unix), i.e. the number of seconds elapsed since January 1, 1970 UTC)
- `pid` - process ID of the language server
- `ppid` - parent process ID (typically editor's or editor plugin's PID)

The path is interpreted as [Go template](https://golang.org/pkg/text/template/), e.g. `/tmp/tofu-ls-cpuprofile-{{timestamp}}.log`.

## Memory Profiling

If the bug you are reporting is related to high memory usage it may be helpful
to collect and share memory profile which can be done via `memprofile` flag.

For example you can modify the launch arguments in your editor to:

```sh
$ tofu-ls serve \
	-memprofile=/tmp/tofu-ls-mem.prof
```

Clients which manage LS installation typically allow passing extra arguments.
See [Logging to Files](#logging-to-files) section above for examples.

The target file will be truncated before being written into.

### Path Templating

Path supports template syntax. This allows for separation of logs while accounting for multiple server or client instances.

**`-memprofile`** supports the following functions:

- `timestamp` - current timestamp (formatted as [`Time.Unix()`](https://golang.org/pkg/time/#Time.Unix), i.e. the number of seconds elapsed since January 1, 1970 UTC)
- `pid` - process ID of the language server
- `ppid` - parent process ID (typically editor's or editor plugin's PID)

The path is interpreted as [Go template](https://golang.org/pkg/text/template/), e.g. `/tmp/tofu-ls-memprofile-{{timestamp}}.log`.
