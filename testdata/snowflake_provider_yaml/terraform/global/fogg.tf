# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

provider "snowflake" {
  account = "foo"
  role    = "bar"
  region  = "us-west-2"
}

provider "assert" {}
terraform {
  required_version = "=1.1.1"

  backend "s3" {

    bucket = "bucket"

    key     = "terraform/foo/global.tfstate"
    encrypt = true
    region  = "region"
    profile = "foo"


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

    snowflake = {
      source = "Snowflake-Labs/snowflake"

      version = "0.55.1"

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
  default = ""
}
# tflint-ignore: terraform_unused_declarations
variable "project" {
  type    = string
  default = "foo"
}
# tflint-ignore: terraform_unused_declarations
variable "component" {
  type    = string
  default = "global"
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
# tflint-ignore: terraform_unused_declarations
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project                = "foo"
    env                    = ""
    service                = "global"
    owner                  = "foo@example.com"
    terraformLastApplyTime = timestamp()
    terraformWorkspaceDir  = "/terraform/global"
    gitRepository          = "https://github.com/chanzuckerberg/fogg"
    managedBy              = "terraform"
  }
}
# tflint-ignore: terraform_unused_declarations
variable "aws_accounts" {
  type = map(string)
  default = {

  }
}
