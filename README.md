# prefix
prefix: execute commands in repl

# Usage
```
Usage: prefix <binary> [<args> ...]

prefix: execute commands in repl. Use ctrl+C or ":q" to quit

Arguments:
  <binary>        binary to run
  [<args> ...]    initial arguments

Flags:
  -h, --help     Show context-sensitive help.
  -d, --debug    print full command to stderr before executing
```

# Examples
```
prefix ls -la
>
... current dir content ...
> app/
... app/ dir content ...
> --color cmd/
... colored content of cmd/ dir ...
```

```
prefix -d git
> status
running: git status
... git status output ...
> add README.md
running: git add README.md
... git add output ...
```