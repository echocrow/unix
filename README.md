# unix

> A simple UNIX timestamp and date converter.

Unix is a CLI that allows easy conversion between formatted dates and UNIX timestamps across different timezones and various date formats

## Contents

- [Features](#features)
- [Installation](#installation)
  - [macos](#macos)
- [Examples](#examples)


## Features

- **UNIX timestamps:** convert formatted dates to timestamps, and timestamps to formatted dates
- **Format detection:** automatically detect various date formats
- **Format customization:** specify the desired date output format via Go's [time layout notation](https://golang.org/pkg/time/#pkg-constants) or common [date format directives](https://strftime.org/)
- **Timezones:** detect or manually set a specific input timezone, and optionally convert the output into a different timezone, automatically adjusting for time offsets


## Installation

Below you'll find the recommended ways to install unix.

Alternatively, you can download unix from the [Releases](https://github.com/echocrow/unix/releases) page.

### macOS
Via [Homebrew](https://brew.sh/):
```sh
# Install:
brew install echocrow/tap/unix
# Update:
brew upgrade echocrow/tap/unix
```


## Examples

- Get the current UNIX timestamp:
  ```sh
  unix
  # e.g. 1612345678
  ```
- Convert a date to a UNIX timestamp:
  ```sh
  unix '1983-01-01 13:37:11'
  # 410276231
  ```
- Convert a timestamp into a formatted date:
  ```sh
  unix 410276231
  # Sat Jan  1 13:37:11 UTC 1983
  ```
- Reformat a date or timestamp:
  ```sh
  unix '1983-01-01 13:37:11' -f long
  # Sat Jan  1 13:37:11 UTC 1983

  unix 410276231 -f '%Y-%m-%d %H:%M:%S'
  # 1983-01-01 13:37:11
  ```
- Convert a date to a different timezone:
  ```sh
  unix '2000-01-01 00:00:00' -z vienna -Z toronto -f long
  # Fri Dec 31 18:00:00 EST 1999
  ```
- Add an offset to a date or timestamp:
  ```sh
  unix '1983-01-01 13:37:11' -a 8h10m
  # 410305631

  unix 410276231 -a -13h
  # Sat Jan  1 00:37:11 UTC 1983
  ```
