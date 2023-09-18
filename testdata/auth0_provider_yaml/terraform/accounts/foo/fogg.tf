# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

provider "assert" {}

provider "auth0" {
  domain = "adomain"
}
terraform {
  required_version = "=1.1.1"

  backend "s3" {

    bucket = "bucket"

    key     = "terraform/foofoo/accounts/foo.tfstate"
    encrypt = true
    region  = "region"
    profile = "foofoo"


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

    auth0 = {
      source = "blah/blah"

      version = "aversion"

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
  default = "foofoo"
}
# tflint-ignore: terraform_unused_declarations
variable "component" {
  type    = string
  default = "foo"
}
# tflint-ignore: terraform_unused_declarations
variable "account" {
  type    = string
  default = "foo"
}
# tflint-ignore: terraform_unused_declarations
variable "owner" {
  type    = string
  default = "foo@example.com"
}
data "external" "git_sha" {
  program = [
    "make",
    "soft_git_log",
  ]
}
data "external" "git_user" {
  program = [
    "make",
    "soft_git_user"
  ]
}
data "external" "git_email" {
  program = [
    "make",
    "soft_git_email"
  ]
}
data "external" "git_branch" {
  program = [
    "make",
    "soft_git_branch"
  ]
}
data "external" "git_authors" {
  program = [
    "make",
    "soft_git_authors"
  ]
}
# tflint-ignore: terraform_unused_declarations
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project                = "foofoo"
    env                    = "accounts"
    service                = "foo"
    owner                  = "foo@example.com"
    terraformLastApplyTime = timestamp()
    terraformWorkspaceDir  = "/terraform/accounts/foo"
    gitRepository          = "git@github.com:chanzuckerberg/fogg"
    gitSHA                 = data.external.git_sha.result.sha
    gitUser                = data.external.git_user.result.name
    gitEmail               = data.external.git_email.result.email
    gitBranch              = data.external.git_branch.result.branch
    gitAuthors             = data.external.git_authors.result.authors
    managedBy              = "terraform"
  }
}
# tflint-ignore: terraform_unused_declarations
data "terraform_remote_state" "global" {
  backend = "s3"
  config = {


    bucket = "bucket"

    key     = "terraform/foofoo/global.tfstate"
    region  = "region"
    profile = "foofoo"


  }
}
data "terraform_remote_state" "foo" {
  backend = "s3"
  config = {


    bucket = "bucket"

    key     = "terraform/foofoo/accounts/foo.tfstate"
    region  = "region"
    profile = "foofoo"


  }
}
# tflint-ignore: terraform_unused_declarations
variable "aws_accounts" {
  type = map(string)
  default = {

  }
}
