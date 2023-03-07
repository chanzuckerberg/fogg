# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

provider "assert" {}
terraform {
  required_version = "=0.100.0"

  backend "remote" {

    hostname     = "si.prod.tfe.czi.technology"
    organization = "k8s-test-app-infra"
    workspaces {
      name = "staging-comp2"
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
variable "env" {
  type    = string
  default = "staging"
}
variable "project" {
  type    = string
  default = "proj"
}
variable "component" {
  type    = string
  default = "comp2"
}
variable "owner" {
  type    = string
  default = "foo@example.com"
}
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project   = "proj"
    env       = "staging"
    service   = "comp2"
    owner     = "foo@example.com"
    managedBy = "terraform"
  }
}
variable "foo" {
  type    = string
  default = "bar2"
}
data "terraform_remote_state" "global" {
  backend = "remote"
  config = {


    hostname     = "si.prod.tfe.czi.technology"
    organization = "k8s-test-app-infra"
    workspaces = {
      name = "global"
    }

  }
}
data "terraform_remote_state" "comp1" {
  backend = "remote"
  config = {


    hostname     = "si.prod.tfe.czi.technology"
    organization = "k8s-test-app-infra"
    workspaces = {
      name = "staging-comp1"
    }

  }
}
data "terraform_remote_state" "vpc" {
  backend = "remote"
  config = {


    hostname     = "si.prod.tfe.czi.technology"
    organization = "k8s-test-app-infra"
    workspaces = {
      name = "staging-vpc"
    }

  }
}
data "terraform_remote_state" "bar" {
  backend = "remote"
  config = {


    hostname     = "si.prod.tfe.czi.technology"
    organization = "k8s-test-app-infra"
    workspaces = {
      name = "accounts-bar"
    }

  }
}
data "terraform_remote_state" "foo" {
  backend = "remote"
  config = {


    hostname     = "si.prod.tfe.czi.technology"
    organization = "k8s-test-app-infra"
    workspaces = {
      name = "accounts-foo"
    }

  }
}
variable "aws_accounts" {
  type = map(string)
  default = {

  }
}
