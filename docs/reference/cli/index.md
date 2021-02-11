---
parent: Reference
nav_order: 2
layout: default
has_children: true
title: CLI commands
---

Most fogg commands can be invoked with a -h or --help to get detailed help.

To view a list of the available commands at any time, just run fogg with no arguments:

```
$ fogg
Usage:
  fogg [command]

Available Commands:
  apply       Apply model defined in fogg.yml to the current tree.
  completion
  exp         Experimental commands
  help        Help about any command
  init        Initialize a new repo for use with fogg
  migrate     Runs all possible fogg migrations
  plan        Run a plan
  setup       Setup dependencies for curent working directory
  version     Print the version number of fogg

Flags:
  -p, --cpuprofile string   activate cpu profiling via pprof and write to file
      --debug               enable verbose output
  -h, --help                help for fogg
  -q, --quiet               do not output to console; use return code to determine success/failure

Use "fogg [command] --help" for more information about a command.
```
