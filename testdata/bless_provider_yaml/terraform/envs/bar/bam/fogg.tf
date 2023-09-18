# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

provider "assert" {}

provider "bless" {
  region  = "bar"
  profile = "foofoofoo"
}
provider "bless" {
  alias  = "a"
  region = "a"

  profile = "foofoofoo"
}
provider "bless" {
  alias  = "b"
  region = "b"

  profile = "foofoofoo"
}
terraform {
  required_version = "=1.1.1"

  backend "s3" {

    bucket = "bucket"

    key     = "terraform/foofoo/envs/bar/components/bam.tfstate"
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

    bless = {
      source = "chanzuckerberg/bless"

      version = "0.0.0"

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
  default = "bar"
}
# tflint-ignore: terraform_unused_declarations
variable "project" {
  type    = string
  default = "foofoo"
}
# tflint-ignore: terraform_unused_declarations
variable "component" {
  type    = string
  default = "bam"
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
# tflint-ignore: terraform_unused_declarations
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project                = "foofoo"
    env                    = "bar"
    service                = "bam"
    owner                  = "foo@example.com"
    terraformLastApplyTime = timestamp()
    terraformWorkspaceDir  = "/terraform/envs/bar/bam"
    gitRepository          = "git@github.com:chanzuckerberg/fogg"
    gitSHA                 = data.external.git_sha.result.sha
    gitUser                = data.external.git_user.result.name
    gitEmail               = data.external.git_email.result.email
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
