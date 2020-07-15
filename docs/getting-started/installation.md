---
parent: Getting Started
nav_order: 1
layout: default
title: Installing fogg
---

* Installing
{:toc}

## OSX

Homebrew is the recommended installation method for fogg:

```shell
brew tap chanzuckerberg/tap
brew install fogg
```

To upgrade fogg:

```shell
brew upgrade fogg
```

## Linux

Binaries are available on the releases page. Download one for your architecture, put it in your path and make it executable.

Instructions on downloading the binary:
1. Visit the fogg releases page: https://github.com/chanzuckerberg/fogg/releases and download the latest version for your architecture.
2. Run `rm -r WHICH_FOGG_PATH; curl -s https://raw.githubusercontent.com/chanzuckerberg/fogg/master/download.sh | bash -s -- -b FOGG_PATH VERSION` where FOGG_PATH is where you want to install fogg and VERSION is the specific release version you want to install (format is vx.yy.z). To find the path of your current fogg, you can run `which fogg`. Then use the path that is outputted as WHICH_FOGG_PATH in the command. The FOGG_PATH is the folder in which fogg will be installed.
3. To verify you installed the desired version, you can run `fogg version`.

