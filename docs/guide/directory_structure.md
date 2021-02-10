---
parent: Guide
nav_order: 1
layout: default
title: Directory Structure
has_toc: true
---

Given the following minimal `fogg.yml` file:

```yaml
version: 2
defaults:
  project: foo
  owner: foo@foobar.com
  terraform_version: 0.12.25
  backend:
    bucket: my-terraform-state-bucket
    dynamodb_table: terraform-statelocks
    profile: default
    region: us-west-2
accounts:
  logging-account: {}
  main-account: {}
envs:
  development:
    components:
      database: {}
      webserver: {}
modules:
  rds_instance: {}
```

`fogg apply` will generate the following directory structure.

All of the *leaf* directories that Fogg creates are terraform workspaces. Fogg installs `fogg.tf`
files in each workspace to manage remote state data sources and state storage, as well as a
`Makefile` for running fogg-provided make targets.

```bash
├── scripts
└── terraform
    ├── accounts
    │   ├── logging-account
    │   └── main-account
    ├── envs
    │   └── development
    │       ├── database
    │       └── webserver
    ├── global
    └── modules
        └── rds_instance
```

Fogg will also write an empty `main.tf`, `variables.tf` and `outputs.tf` file to each workspace to
encourage best practices for workspace layout:

- Any workspaces that need to expose output values to other workspaces should define those outputs in `outputs.tf`
- `main.tf` can invoke terraform modules and define terraform-controlled resources
- Values that are local to a workspace go in a `locals{}` block in `variables.tf`
