# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
provider "aws" {

  region  = "us-west-2"
  profile = "profile"

  allowed_account_ids = ["00456"]
}
# Aliased Providers (for doing things in every region).


provider "assert" {}
terraform {
  required_version = "=0.100.0"

  backend "s3" {

    bucket = "buck"

    key     = "terraform/proj/envs/prod/components/vpc.tfstate"
    encrypt = true
    region  = "us-west-2"
    profile = "profile"


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

      version = "0.12.0"

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
  default = "prod"
}
# tflint-ignore: terraform_unused_declarations
variable "project" {
  type    = string
  default = "proj"
}
# tflint-ignore: terraform_unused_declarations
variable "region" {
  type    = string
  default = "us-west-2"
}
# tflint-ignore: terraform_unused_declarations
variable "component" {
  type    = string
  default = "vpc"
}
# tflint-ignore: terraform_unused_declarations
variable "aws_profile" {
  type    = string
  default = "profile"
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
    project                = "proj"
    env                    = "prod"
    service                = "vpc"
    owner                  = "foo@example.com"
    terraformLastApplyTime = timestamp()
    terraformWorkspaceDir  = "/terraform/envs/prod/vpc"
    gitRepository          = "git@github.com:chanzuckerberg/fogg"
    gitSHA                 = data.external.git_sha.result.sha
    gitUser                = data.external.git_user.result.name
    gitEmail               = data.external.git_email.result.email
    gitBranch              = data.external.git_branch.result.branch
    managedBy              = "terraform"
  }
}
# tflint-ignore: terraform_unused_declarations
variable "foo" {
  type    = string
  default = "bar1"
}
# tflint-ignore: terraform_unused_declarations
data "terraform_remote_state" "global" {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/global.tfstate"
    region  = "us-west-2"
    profile = "profile"


  }
}
data "terraform_remote_state" "datadog" {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/envs/prod/components/datadog.tfstate"
    region  = "us-west-2"
    profile = "profile"


  }
}
data "terraform_remote_state" "hero" {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/envs/prod/components/hero.tfstate"
    region  = "us-west-2"
    profile = "profile"


  }
}
data "terraform_remote_state" "okta" {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/envs/prod/components/okta.tfstate"
    region  = "us-west-2"
    profile = "profile"


  }
}
data "terraform_remote_state" "sentry" {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/envs/prod/components/sentry.tfstate"
    region  = "us-west-2"
    profile = "profile"


  }
}
data "terraform_remote_state" "bar" {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/accounts/bar.tfstate"
    region  = "us-west-2"
    profile = "profile"


  }
}
data "terraform_remote_state" "foo" {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/accounts/foo.tfstate"
    region  = "us-west-2"
    profile = "profile"


  }
}
# tflint-ignore: terraform_unused_declarations
variable "aws_accounts" {
  type = map(string)
  default = {

    bar = "00456"

    foo = "123"

  }
}
