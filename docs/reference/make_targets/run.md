---
parent: Make Targets
grand_parent: Reference
layout: default
title: run
---

A wrapper around the `terraform` cli that supports running arbitrary terraform commands. Accepts a CMD argument

## Usage

```shell
make run CMD="apply -refresh=false"
```
