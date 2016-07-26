# sub

sub is a command-line tool with a clean interface that runs a command in
multiple sub-directories in parallel. It was originally created to speed up
developing involving multiple git repositories.

```
Usage: sub [options] [--] <command>...

Options:
  -d value
        Run only in this directory (default [])
  -j int
        Number of concurrent subprocesses (default 4)
  -v    Print out the version and exit
```

## Usage

sub obtains the list of sub-directories in this order:

1. If one or more `-d` flags are provided, the command is run in each of them.
2. If no `-d` flags are provided and `./sub.cnf` exists, the command is run in
the sub-directories specified in that file.
3. If no `-d` flags are provided and `./sub.cnf` does not exist, the command is
run in all non-hidden sub-directories.

In this sense, `-d` flags override the `./sub.cnf` file, and the default is
just whatever non-hidden sub-directories are present.

A non-hidden sub-directory is any sub-directory whose name does not start with
a period.

The format of a `./sub.cnf` is simple: place each directory on its own line:

```
$ cat sub.cnf
foo-dir
bar-dir
baz-dir
```

Command output is printed to stdout and stderr in the order the commands are
finished. This provides a discontinuous, but legible and informative
interaction.

## Examples

Let's start off with three empty directories: `mkdir go build tools`.

Initialize a git repository in each:

```
$ sub git init
go: done
Initialized existing Git repository in /tmp/go/.git/
tools: done
Initialized existing Git repository in /tmp/tools/.git/
build: done
Initialized existing Git repository in /tmp/build/.git/
```

Add remote repositories for each:

```
$ sub git remote add origin https://github.com/golang/{}.git
go: done
build: done
tools: done
```

Network operations. Now this is where parallelism comes in handy. Let's fetch:

```
$ sub git fetch
build: done
From https://github.com/golang/build
 * branch            master     -> FETCH_HEAD
 * [new branch]      master     -> origin/master
tools: done
From https://github.com/golang/tools
 * branch            master     -> FETCH_HEAD
 * [new branch]      master     -> origin/master
go: done
From https://github.com/golang/go
 * branch            master     -> FETCH_HEAD
 * [new branch]      master     -> origin/master
```
