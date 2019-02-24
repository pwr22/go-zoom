# zoom

[![Build Status](https://travis-ci.com/pwr22/zoom.svg?branch=master)](https://travis-ci.com/pwr22/zoom)
[![Build status](https://ci.appveyor.com/api/projects/status/cuptxx2040f6f9sa/branch/master?svg=true)](https://ci.appveyor.com/project/pwr22/zoom/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pwr22/zoom)](https://goreportcard.com/report/github.com/pwr22/zoom)
[![Downloads](https://img.shields.io/github/downloads/pwr22/zoom/total.svg)](https://github.com/pwr22/zoom/releases)
[![AUR version](https://img.shields.io/aur/version/zoom-parallel.svg)](https://aur.archlinux.org/packages/zoom-parallel/)

Parallel command executor with a focus on simplicity and good cross-platform behaviour 

## Usage

    cat args.txt | zoom [optional command] 
    zoom [optional command] [::: arg1 arg2 arg3 ...] [:::: argfile1.txt argfile2.txt ...] ...

There are two main modes of operation

- Args from standard input
- Args from command line and / or files (can still read from stdin by passing `-`)

In either mode the optional command may contain a placeholder `{}` which will be replaced with each argument provided, otherwise arguments will be appended to the end of the command. If there isn't a command then each argument is a full command itself

For each of these zoom will invoke a `$SHELL` so you can use things like `&&`, `||` and other goodness. Watch out quote them properly if you pass them on the commandline but you shouldn't need to worry if loading them from a file 

An example with arguments from standard input

    $ cat args.txt

    8.8.8.8
    8.8.4.4

    $ cat args.txt | zoom ping -c1

An example with commands from standard input

    $ cat commands.txt

    ping -c1 8.8.8.8
    ping -c1 8.8.4.4

    $ cat commands.txt | zoom

An example with arguments on the command line and using a placeholder

    $ zoom ping {} -c1 ::: 8.8.8.8 8.8.4.4

An example with commands on the command line

    $ zoom ::: "ping -c1 8.8.8.8" "ping -c1 8.8.4.4"

An example taking arguments from a file

    $ cat args.txt

    8.8.8.8
    8.8.4.4

    $ zoom ping -c1 :::: args.txt

An example taking commands from a file

    $ cat commands.txt

    ping -c1 8.8.8.8
    ping -c1 8.8.4.4

    $ zoom :::: commands.txt

`:::` and `::::` can be used multiple times and intermixed as needed. Each set of arguments will be permuted together when building commands to run so they can be used to replace loop functionality

An example using `:::`

    $ zoom echo ::: a b c ::: 1 2 3

    a 3
    a 2
    a 1
    b 1
    b 2
    b 3
    c 2
    c 1
    c 3

An example using `::::`

    $ cat letters.txt

    a
    b
    c

    $ cat numbers.txt

    1
    2
    3

    $ zoom echo :::: letters.txt numbers.txt

    a 2
    a 3
    b 1
    a 1
    b 2
    c 1
    b 3
    c 2
    c 3

An example using both `:::` and `::::`

    $ cat numbers.txt

    1
    2
    3

    $ zoom echo ::: a b c :::: numbers.txt

    a 3
    a 1
    a 2
    b 1
    b 2
    b 3
    c 1
    c 2
    c 3

The arguments given after `:::` are all taken as a single set to permute but `::::` gives each file a set of its own. That means these all permute

    $ zoom echo ::: a b c ::: 1 2 3

    a 2
    a 3
    b 1
    a 1
    b 2
    c 1
    c 2
    b 3
    c 3

    $ cat letters.txt

    a
    b
    c

    $ cat numbers.txt

    1
    2
    3

    $ zoom echo :::: letters.txt numbers.txt

    a 1
    a 3
    b 1
    a 2
    b 2
    b 3
    c 1
    c 2
    c 3

    $ zoom echo :::: letters.txt :::: numbers.txt

    a 1
    a 2
    a 3
    b 1
    b 2
    b 3
    c 1
    c 2
    c 3

But this does not permute

    $ zoom echo ::: a b c 1 2 3

    a
    b
    2
    3
    1
    c

This behaviour is how [GNU Parallel](https://www.gnu.org/software/parallel/) behaves. [MIT/Rust Parallel](https://github.com/mmstick/parallel) implements `::::` differently in that it takes all arguments given in files and concatenates them into a single set like this

    $ cat letters.txt 

    a
    b
    c

    $ cat numbers.txt 

    1
    2
    3
    
    $ rust-parallel echo :::: letters.txt numbers.txt

    a
    b
    c
    1
    2
    3

I find this to be more consistent across `:::` and `::::` so the zoom semantics may change before v1.0.0 to match this. If so I'll likely provide a flag for parallel compatibility

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