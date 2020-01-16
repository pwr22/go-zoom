# go-zoom

[![Build Status](https://travis-ci.com/pwr22/zoom.svg?branch=master)](https://travis-ci.com/pwr22/zoom)
[![Build status](https://ci.appveyor.com/api/projects/status/cuptxx2040f6f9sa/branch/master?svg=true)](https://ci.appveyor.com/project/pwr22/zoom/branch/master)
[![codecov](https://codecov.io/gh/pwr22/zoom/branch/master/graph/badge.svg)](https://codecov.io/gh/pwr22/zoom)
[![Go Report Card](https://goreportcard.com/badge/github.com/pwr22/zoom)](https://goreportcard.com/report/github.com/pwr22/zoom)
[![Downloads](https://img.shields.io/github/downloads/pwr22/zoom/total.svg)](https://github.com/pwr22/zoom/releases)
[![](https://tokei.rs/b1/github/pwr22/zoom)](https://github.com/pwr22/zoom).
[![AUR version](https://img.shields.io/aur/version/zoom-parallel.svg)](https://aur.archlinux.org/packages/zoom-parallel/)

Parallel command executor focussed on simplicity and good cross-platform behaviour.

## Usage

    cat args.txt | go-zoom [options] [command] 
    go-zoom [options] [command] [::: arg1 arg2 arg3 ...] [:::: argfile1.txt argfile2.txt ...] ...

    --dry-run           Print the commands that would be run with out doing so.
    -j, --jobs int      How many jobs to run at once. Give 0 to run as many as possible. Defaults to the number of CPUs available.
    -k, --keep-order    Print output from jobs in the order they were started instead of the order they finish.
    -V, --version       Print version and licensing info.


For detailed usage and lots of examples take a look [here](USAGE.md).

## Installation

### Arch Linux

You can install from the AUR with `yay -S go-zoom`.

### Anything Else

Head over to the [releases](https://github.com/pwr22/go-zoom/releases) page, download the binary for your operating system and put it somewhere in your `$PATH`.

## Why

`go-zoom` is inspired by [rush](https://github.com/shenwei356/rush) but I needed different behaviour on failure and found the codebase a little challenging as a newcomer so felt I couldn't add to it easily.

I was also wanted to use Go more, and its a good language for this sort of problem so I decided to write `go-zoom`.

## See also

- [rush](https://github.com/shenwei356/rush)
- [gargs](https://github.com/brentp/gargs)
- [MIT/Rust Parallel](https://github.com/mmstick/parallel)
- [GNU Parallel](https://www.gnu.org/software/parallel/)