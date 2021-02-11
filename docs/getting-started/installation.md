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

1. Visit the fogg releases page: <https://github.com/chanzuckerberg/fogg/releases> and download the
   latest version for your architecture.
2. Run `curl -s https://raw.githubusercontent.com/chanzuckerberg/fogg/master/download.sh | bash -s -- -b FOGG_PATH VERSION`
   1. FOGG_PATH is the directory where you want to install fogg
   2. VERSION is the release you want
3. To verify you installed the desired version, you can run `fogg version`.
