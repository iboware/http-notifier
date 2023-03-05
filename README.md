# http notifier

It reads STDIN and sends new messages every interval. Each line is interpreted as a new message that needs to be notified about. It will keep running until it receives SIGINT, or OS Interrupt.

## Usage

```bash
Usage:
  http-notifier [flags]

Flags:
  -h, --help           help for http-notifier
  -i, --interval int   Notification Interval (default 5)
  -u, --url string     Notification URL
```

## Development Requirements

* [Go v1.20](https://go.dev/dl/)
* [Cobra CLI](https://github.com/spf13/cobra)
* [GoMock](https://github.com/golang/mock)

## Testing

All unit tests can be run by using Makefile

```bash
make test.unit
```
