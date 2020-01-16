# Usage Manual

    cat args.txt | go-zoom [options] [command] 
    go-zoom [options] [command] [::: arg1 arg2 arg3 ...] [:::: argfile1.txt argfile2.txt ...] ...

go-zoom has two ways of taking arguments:

- Standard input.
- Inline on the command line and / or from files (standard input as well if you give `-` as the name of a file).

In either mode the command can contain a placeholder `{}` and any instances of this are replaced with arguments. If there's no placeholder then arguments are appended to the end of the command. If there isn't a command then each argument is itself a command.

For each command zoom will invoke a `$SHELL` on *NIX or `%COMSPEC%` on Windows so you can use things like `&&`, `||` and other shellisms. Watch out to quote them properly if passing them on the commandline but you shouldn't need to worry if loading them from a file.

## Simple Examples

### Arguments from standard input

    $ cat args.txt

    8.8.8.8
    8.8.4.4

    $ cat args.txt | go-zoom ping -c1

### Commands from standard input

    $ cat commands.txt

    ping -c1 8.8.8.8
    ping -c1 8.8.4.4

    $ cat commands.txt | go-zoom

### Commands on the command line

    $ go-zoom ::: "ping -c1 8.8.8.8" "ping -c1 8.8.4.4"

### Using a placeholder

    $ go-zoom ping {} -c1 ::: 8.8.8.8 8.8.4.4

### Arguments from a file

    $ cat args.txt

    8.8.8.8
    8.8.4.4

    $ go-zoom ping -c1 {} :::: args.txt

## Examples with multiple argument sources

`:::` and `::::` can be used multiple times and intermixed as needed. Each set of arguments will be permuted together when building commands to run so they can be used to replace loop functionality

### Arguments from the command line

    $ go-zoom echo ::: a b c ::: 1 2 3

    a 3
    a 2
    a 1
    b 1
    b 2
    b 3
    c 2
    c 1
    c 3

### Arguments from files

    $ cat letters.txt

    a
    b
    c

    $ cat numbers.txt

    1
    2
    3

    $ go-zoom echo :::: letters.txt numbers.txt

    a 2
    a 3
    b 1
    a 1
    b 2
    c 1
    b 3
    c 2
    c 3

### Arguments from the command line and a file

    $ cat numbers.txt

    1
    2
    3

    $ go-zoom echo ::: a b c :::: numbers.txt

    a 3
    a 1
    a 2
    b 1
    b 2
    b 3
    c 1
    c 2Using both `:::` and `::::`
    c 3

## Differences between `:::` and `::::`

The arguments given after `:::` are all taken as a single set to permute but `::::` gives each file a set of its own. That means these all permute

    $ go-zoom echo ::: a b c ::: 1 2 3

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

    $ go-zoom echo :::: letters.txt numbers.txt

    a 1
    a 3
    b 1
    a 2
    b 2
    b 3
    c 1
    c 2
    c 3

    $ go-zoom echo :::: letters.txt :::: numbers.txt

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

    $ go-zoom echo ::: a b c 1 2 3

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

I find this to be more consistent across `:::` and `::::` so the go-zoom semantics may change before v1.0.0 to match this. If so I'll likely provide a flag for parallel compatibility