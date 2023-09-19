# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
variable "TFC_RUN_ID" {
  type    = string
  default = "unknown"
}
variable "TFC_WORKSPACE_NAME" {
  type    = string
  default = "unknown"
}
variable "TFC_WORKSPACE_SLUG" {
  type    = string
  default = "unknown"
}
variable "TFC_CONFIGURATION_VERSION_GIT_BRANCH" {
  type    = string
  default = "unknown"
}
variable "TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA" {
  type    = string
  default = "unknown"
}
variable "TFC_CONFIGURATION_VERSION_GIT_TAG" {
  type    = string
  default = "unknown"
}
variable "TFC_PROJECT_NAME" {
  type    = string
  default = "unknown"
}
provider "aws" {

  region = "us-west-2"


  assume_role {
    role_arn = "arn:aws:iam::626314663667:role/tfe-si"
  }

  default_tags {
    tags = {
      TFC_RUN_ID                               = var.TFC_RUN_ID
      TFC_WORKSPACE_NAME                       = var.TFC_WORKSPACE_NAME
      TFC_WORKSPACE_SLUG                       = var.TFC_WORKSPACE_SLUG
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = var.TFC_CONFIGURATION_VERSION_GIT_BRANCH
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA
      TFC_CONFIGURATION_VERSION_GIT_TAG        = var.TFC_CONFIGURATION_VERSION_GIT_TAG
      TFC_PROJECT_NAME                         = var.TFC_PROJECT_NAME
    }
  }
  allowed_account_ids = ["626314663667"]
}
# Aliased Providers (for doing things in every region).


provider "assert" {}
terraform {
  required_version = "=1.1.1"

  backend "remote" {

    hostname     = "si.prod.tfe.czi.technology"
    organization = "shared-infra"
    workspaces {
      name = "accounts-account"
    }

  }
  required_providers {

    archive = {
      source = "hashicorp/archive"

      version = "~> 2.0"

    }

    assert = {
      source = "bwoznicki/assert"

      version = "0.0.1"

    }

    aws = {
      source = "hashicorp/aws"

      version = "3.63.0"

    }

    local = {
      source = "hashicorp/local"

      version = "~> 2.0"

    }

    null = {
      source = "hashicorp/null"

      version = "3.1.1"

    }

    okta-head = {
      source = "okta/okta"

      version = "~> 3.30"

    }

    random = {
      source = "hashicorp/random"

      version = "~> 3.4"

    }

    tls = {
      source = "hashicorp/tls"

      version = "~> 3.0"

    }

  }
}
# tflint-ignore: terraform_unused_declarations
variable "env" {
  type    = string
  default = "accounts"
}
# tflint-ignore: terraform_unused_declarations
variable "project" {
  type    = string
  default = "foo"
}
# tflint-ignore: terraform_unused_declarations
variable "region" {
  type    = string
  default = "us-west-2"
}
# tflint-ignore: terraform_unused_declarations
variable "component" {
  type    = string
  default = "account"
}
# tflint-ignore: terraform_unused_declarations
variable "account" {
  type    = string
  default = "account"
}
# tflint-ignore: terraform_unused_declarations
variable "owner" {
  type    = string
  default = "foo@example.com"
}
# tflint-ignore: terraform_unused_declarations
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project   = "foo"
    env       = "accounts"
    service   = "account"
    owner     = "foo@example.com"
    managedBy = "terraform"
  }
}
# tflint-ignore: terraform_unused_declarations
data "terraform_remote_state" "global" {
  backend = "remote"
  config = {


    hostname     = "si.prod.tfe.czi.technology"
    organization = "shared-infra"
    workspaces = {
      name = "global"
    }

  }
}
data "terraform_remote_state" "account" {
  backend = "remote"
  config = {


    hostname     = "si.prod.tfe.czi.technology"
    organization = "shared-infra"
    workspaces = {
      name = "accounts-account"
    }

  }
}
# tflint-ignore: terraform_unused_declarations
variable "aws_accounts" {
  type = map(string)
  default = {

  }
}
