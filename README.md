# zoom

[![Go Report Card](https://goreportcard.com/badge/github.com/pwr22/zoom)](https://goreportcard.com/report/github.com/pwr22/zoom)
[![Downloads](https://img.shields.io/github/downloads/pwr22/zoom/total.svg)](https://github.com/pwr22/zoom/releases)

Parallel command executor with a focus on simplicity and good cross-platform behaviour 

## Usage

    zoom <file_containing_commands>

The file should contain a list of commands to be executed, one per line. For example

    ping 8.8.8.8
    ping 8.8.4.4

`zoom` will spawn a `$SHELL` for each command so you can use things like `&&` and `||` 

## Installation

Head over to the [releases](https://github.com/pwr22/zoom/releases) page, download the binary for your operating system and put it somewhere in your `$PATH`

## Why

`zoom` is inspired by [rush](https://github.com/shenwei356/rush) but I needed different behaviour on command failure and found the codebase a little challenging as a newcomer so felt I couldn't add to it easily

I was also interested in using Go more, and its a good language for this sort of problem so I decided to write `zoom`

## See also

- [rush](https://github.com/shenwei356/rush)
- [gargs](https://github.com/brentp/gargs)