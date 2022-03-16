# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
terraform {
  required_version = "=1.1.1"

  backend "remote" {

    hostname     = "tfe.example.com"
    organization = "test-org"
    workspaces {
      name = "accounts-acct1"
    }

  }
  required_providers {

    archive = {
      source = "hashicorp/archive"

      version = "~> 2.0"

    }

    local = {
      source = "hashicorp/local"

      version = "~> 2.0"

    }

    null = {
      source = "hashicorp/null"

      version = "~> 3.0"

    }

    random = {
      source = "hashicorp/random"

      version = "~> 2.2"

    }

    template = {
      source = "hashicorp/template"

      version = "~> 2.2"

    }

    tls = {
      source = "hashicorp/tls"

      version = "~> 3.0"

    }

  }
}
variable "env" {
  type    = string
  default = "accounts"
}
variable "project" {
  type    = string
  default = "foo"
}
variable "component" {
  type    = string
  default = "acct1"
}
variable "account" {
  type    = string
  default = "acct1"
}
variable "owner" {
  type    = string
  default = "foo@example.com"
}
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string, repo : string, folderPath : string })
  default = {
    project    = "foo"
    env        = "accounts"
    service    = "acct1"
    owner      = "foo@example.com"
    repo       = ""
    folderPath = ""
    managedBy  = "terraform"
  }
}
data "terraform_remote_state" "global" {
  backend = "remote"
  config = {


    hostname     = "tfe.example.com"
    organization = "test-org"
    workspaces = {
      name = "global"
    }

  }
}
data "terraform_remote_state" "acct1" {
  backend = "remote"
  config = {


    hostname     = "tfe.example.com"
    organization = "test-org"
    workspaces = {
      name = "accounts-acct1"
    }

  }
}
variable "aws_accounts" {
  type = map(string)
  default = {

  }
}
