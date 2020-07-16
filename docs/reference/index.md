---
layout: default
has_children: true
title: Reference
has_toc: true
nav_order: 2
---

# Reference

The `fogg` cli reads the contents of `fogg.yml` to generate a `Makefile` and `fogg.tf` scaffolding that configures providers, state storage, and dependencies between fogg-managed components.

Fogg has two primary CLI interfaces: the `fogg` command line utility itself, and the `make` targets supplied in the `Makefile`s that fogg generates.

