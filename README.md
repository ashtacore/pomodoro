<p align="center">
	<img src="https://raw.githubusercontent.com/ashtacore/pomodoro/main/pomo.gif" alt="Made with VHS">
	<a href="https://vhs.charm.sh">
		<img src="https://stuff.charm.sh/vhs/badge.svg">
	</a>
	<br>
	<h1 align="center">Pomo</h1>
	<p align="center">A fancier way to do Pomodoro.</p>
</p>

---

Pomo is a small CLI tool that helps you manage your Pomodoro routine. It is a small edit to the `timer` tool by [caarlos0](https://github.com/caarlos0/timer)

## Usage

```sh
pomo <Focus Duration> <Break Duration>
pomo -f <Focus Title> -b <Break Title> <Focus Duration> <Break Duration>
man pomo
pomo --help
```

Several other flags are available:
- `-a` makes the timer take up the whole terminal
- `-n=false` disables the default system notifications
- `-s` will cause a short beep to play when the timer finishes
  - On Linux sounds require permission to access `/dev/tty0` or `/dev/input/by-path/platform-pcspkr-event-spkr` files for writing, and `pcspkr` module must be loaded. User must be in correct groups, usually `input` and/or `tty`.
  - On macOS, you must enable `Audible bell` in Terminal --> Preferences --> Settings --> Advanced.
  - Note that on some platforms a sound will be included with a desktop notification, making this unnecessary.

You can pause the timer with the space bar.

It is possible to pass a time unit for `<Duration>`.

Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
If no unit is passed, it defaults to seconds ("s").

## Install

**manually**:

Download the pre-compiled binaries from the [releases page][releases] or clone the repo build from source.

[releases]:  https://github.com/ashtacore/pomodoro/releases

# Badges

[![Release](https://img.shields.io/github/release/ashtacore/pomodoro.svg?style=for-the-badge)](https://github.com/ashtacore/pomodoro/releases/latest)

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE.md)

[![Build](https://img.shields.io/github/actions/workflow/status/ashtacore/pomodoro/build.yml?style=for-the-badge)](https://github.com/ashtacore/pomodoro/actions?query=workflow%3Abuild)

[![Go Report Card](https://goreportcard.com/badge/github.com/ashtacore/pomodoro?style=for-the-badge)](https://goreportcard.com/report/github.com/ashtacore/pomodoro)

[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)

