---
parent: Reference
nav_order: 1
layout: default
title: Configuration file
has_toc: true
---

## Configuration file

{: .no_toc }

Fogg reads its configuration from a fogg.yml file in the root of your terraform repository.

1. TOC
{:toc}

## Example fogg.yml

```yaml
version: 2
accounts:
  logging-account:
    providers:
      aws:
        account_id: '123456789123'
        additional_regions: &id001
          - us-east-1
          - us-west-2
        role: terraform-role
  main-account:
    providers:
      aws:
        account_id: '987654321321'
        additional_regions: *id001
defaults:
  backend:
    account_id: '987654321321'
    bucket: terraform-state-bucket
    dynamodb_table: terraform-statelocks
    profile: terraform-profile
    region: us-west-2
  extra_vars:
    globalvariable: 'value'
  owner: infra-eng@mydomain.com
  project: shared-infra
  providers:
    aws:
      account_id: '987654321321'
      region: us-west-2
      version: 2.65.0
      role: terraform-role
  terraform_version: 0.12.25
  tools:
    circle_ci:
      command: lint
      enabled: true
    tflint:
      enabled: true
envs:
  development:
    components:
      webserver:
        module_source: terraform/modules/webserver
      database:
        providers:
          datadog:
            version: v2.9.0
    extra_vars:
      database_username: dev_website
  production:
    components:
      webserver:
        module_source: terraform/modules/webserver
      database:
        providers:
          datadog:
            version: v2.9.0
    extra_vars:
      database_username: prod_website
modules:
  webserver: {}
```

## Top Level Arguments

- `version` - Required, the current fogg config version is 2
- `accounts` - Specify a map of workspaces to manage in `terraform/accounts` whose state outputs are
  available to all other workspaces.
- `defaults` - Default workspace arguments that are applied to every workspace unless overridden by that workspace.
- `envs` - Create separate environments that don't share config
- `modules` - Manage [modules](https://www.terraform.io/docs/modules/index.html) in
  `terraform/modules` that can be invoked in any workspace for better code reuse.

## Common workspace arguments

The following config paths all accept a standard `workspace` configuration spec:

- `accounts.{accountname}`
- `envs.{envname}.{component}`
- `defaults`

### Arguments

- `backend` - Configure the terraform remote state backend
- `extra_vars` - Any extra terraform variables to add to the workspace
- `owner` - Set the `var.owner` terraform variable to this string (email address recommended)
- `project` - Set the `var.project` terraform variable to this string
- `providers` - Preconfigure any terraform providers
- `terraform_version` - Which terraform version to use for this workspace.
- `tools` - En/Disable CI for this workspace

### `backend`

Defines terraform remote state storage for a workspace

For S3 backends:

- `account_id` - AWS Account ID
- `bucket` - Name of the S3 bucket for state storage
- `dynamodb_table` - Name of the DynamoDB table for state locks
- `profile` - AWS profile to use for authentication
- `role` - AWS role to assume for state storage
- `region` - AWS region

For Terraform Enterprise / Terraform Cloud backends

- `host_name` - Hostname of the Terraform Cloud/Enterprise instance.
- `organization` - Organization this repository belongs to.

### `providers`

Providers is a map of provider names to provider configuration values.

The currently supported providers are:

- `aws`
  - `account_id` - AWS Account ID
  - `additional_regions` - List of regions to generate provider aliases for
  - `profile` - AWS profile to use
  - `region` - Primary AWS region
  - `version` - Version of the provider to use
- `bless`
  - `aws_profile` - AWS profile to use
  - `aws_region` - AWS region to use
  - `additional_regions` - List of regions to generate provider aliases for
  - `version` - Version of the provider to use
- `datadog`
  - `version` - Version of the provider to use
- `github`
  - `organization` - GitHub Organization
  - `base_url` - URL for on-premise GitHub installations
  - `version` - Version of the provider to use
- `heroku`
  - `version` - Version of the provider to use
- `okta`
  - `org_name` - Okta organization name
  - `version` - Version of the provider to use
- `snowflake`
  - `account` - Snowflake account
  - `role` - Role to use
  - `region` - Region
  - `version` - Version of the provider to use
- `tfe` - Terraform Enterprise
  - `hostname` - Hostname for the Terraform Enterprise instance
  - `version` - Version of the provider to use

### `tools`

Define CI integrations.

The currently supported CI tools are:

- `travis_ci`
  - `aws_iam_role_name` - Configure the CI toole to assume this role before running checks
  - `test_buckets` - Whether to run checks
  - `command` - Which makefile target to use to validate the workspace (default is "check")
  - `enabled` - boolean, enable this CI tool
- `circle_ci`
  - `aws_iam_role_name` -
  - `test_buckets` -
  - `command` -
  - `ssh_key_fingerprints` -
  - `enabled` - boolean, enable this CI tool
- `github_actions_ci`
  - `aws_iam_role_name` -
  - `test_buckets` -
  - `command` -
  - `enabled` - boolean, enable this CI tool
  - `ssh_key_secrets` - list of strings, github actions secrets containing ssh keys for pulling code from private repos
- `tflint`
  - `enabled` - boolean, enable this CI tool

## envs

Fogg manages workspaces within environments (think "dev", "staging" and "prod" for example) that
don't share state between them.

Each workspace within an environment is called a `component`

### env arguments

- `components` - This is a map of workspace names to [workspace configurations](#common-workspace-arguments)

### component arguments

In addition to the common workspace arguments, components also support the following arguments:

- `eks` - Object that contains a `cluster_name` key to point terraform to a kubectl context
- `kind` - String,  only `terraform` is supported (DEPRECATED)
- `module_source` - String, path to a terraform module to use for this component. If this is a
  relative path such as `terraform/modules/webserver`, fogg will generate a main.tf in this
  component with a full invocation to this module.
- `module_name` - If module_source is supplied, this is the name to use for the module invocation resource.
