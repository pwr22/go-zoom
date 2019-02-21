# zoom

[![Build Status](https://travis-ci.com/pwr22/zoom.svg?branch=master)](https://travis-ci.com/pwr22/zoom)
[![Build status](https://ci.appveyor.com/api/projects/status/cuptxx2040f6f9sa/branch/master?svg=true)](https://ci.appveyor.com/project/pwr22/zoom/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pwr22/zoom)](https://goreportcard.com/report/github.com/pwr22/zoom)
[![Downloads](https://img.shields.io/github/downloads/pwr22/zoom/total.svg)](https://github.com/pwr22/zoom/releases)

Parallel command executor with a focus on simplicity and good cross-platform behaviour 

## Usage

    cat args.txt | zoom [optional command] 

The file can be arguments for the command, or if none was given, full commands. In either case it's one per line

An example with arguments

    $ cat args.txt

    8.8.8.8
    8.8.4.4

    $ cat args.txt | zoom ping

An example with commands

    $ cat commands.txt

    ping 8.8.8.8
    ping 8.8.4.4

    $ cat commands.txt | zoom

`zoom` will build jobs by taking each argument, prefixing it with the command if you gave one and then run those jobs for you in parallel. It will invoke a `$SHELL` for each command so you can use things like `&&`, `||` and other goodness 

## Installation

Head over to the [releases](https://github.com/pwr22/zoom/releases) page, download the binary for your operating system and put it somewhere in your `$PATH`

## Why

`zoom` is inspired by [rush](https://github.com/shenwei356/rush) but I needed different behaviour on command failure and found the codebase a little challenging as a newcomer so felt I couldn't add to it easily

I was also interested in using Go more, and its a good language for this sort of problem so I decided to write `zoom`

## See also

- [rush](https://github.com/shenwei356/rush)
- [gargs](https://github.com/brentp/gargs)