# sub

sub is a command-line tool with a clean interface that runs a command in
multiple sub-directories in parallel. It was originally created to speed up
developing involving multiple git repositories.

```
$ sub -h
Usage: sub [options] [--] <command>...

Options:
  -c string
    	Configuration file
  -j int
    	Number of concurrent subprocesses (default -1)
  -v	Print out the version and exit
```

## Usage

sub obtains the list of sub-directories to run in from the configuration file,
but if one isn't provided, then it just runs in all non-hidden sub-directories,
a sensible default.

The format of a configuration file is a JSON formatted list of strings, each
representing a sub-directory. It can technically be the path to any
directory:

```json
[
  "foo",
  "bar",
  "baz",
  "../im/complicated"
]
```

The positional arguments to sub is the command itself. If a command requires
flags, then you must prefix the command with `--`:

```
$ sub -- git branch -vv
```

By default, sub runs all the commands at once. If you have to throttle the
number of simultaneous commands, pass a value  for `-j`:

```
$ sub -j 4 git fetch
```
