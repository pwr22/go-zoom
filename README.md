# zoom

[![Build Status](https://travis-ci.com/pwr22/zoom.svg?branch=master)](https://travis-ci.com/pwr22/zoom)
[![Build status](https://ci.appveyor.com/api/projects/status/cuptxx2040f6f9sa/branch/master?svg=true)](https://ci.appveyor.com/project/pwr22/zoom/branch/master)
[![codecov](https://codecov.io/gh/pwr22/zoom/branch/master/graph/badge.svg)](https://codecov.io/gh/pwr22/zoom)
[![Go Report Card](https://goreportcard.com/badge/github.com/pwr22/zoom)](https://goreportcard.com/report/github.com/pwr22/zoom)
[![Downloads](https://img.shields.io/github/downloads/pwr22/zoom/total.svg)](https://github.com/pwr22/zoom/releases)
[![](https://tokei.rs/b1/github/pwr22/zoom)](https://github.com/pwr22/zoom).
[![AUR version](https://img.shields.io/aur/version/zoom-parallel.svg)](https://aur.archlinux.org/packages/zoom-parallel/)

Parallel command executor with a focus on simplicity and good cross-platform behaviour 

## Usage

    cat args.txt | zoom [options] [command] 
    zoom [options] [command] [::: arg1 arg2 arg3 ...] [:::: argfile1.txt argfile2.txt ...] ...

    --dry-run           print the commands that would be run instead of running them
    -j, --jobs int      number of jobs to run at once or 0 for as many as possible (default 4)
    -k, --keep-order    print output in the order jobs were run instead of the order they finish
    -V, --version       print version information


For detailed usage and lots of examples take a look [here](USAGE.md)

## Installation

### Arch Linux

You can install from the AUR with `yay -S zoom-parallel`

### Anything Else

Head over to the [releases](https://github.com/pwr22/zoom/releases) page, download the binary for your operating system and put it somewhere in your `$PATH`

## Why

`zoom` is inspired by [rush](https://github.com/shenwei356/rush) but I needed different behaviour on command failure and found the codebase a little challenging as a newcomer so felt I couldn't add to it easily

I was also interested in using Go more, and its a good language for this sort of problem so I decided to write `zoom`

## See also

- [rush](https://github.com/shenwei356/rush)
- [gargs](https://github.com/brentp/gargs)
- [GNU Parallel](https://www.gnu.org/software/parallel/)