# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

provider "assert" {}

provider "okta" {
  org_name = "orgname"
  base_url = "oktapreview.com"
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

    local = {
      source = "hashicorp/local"

      version = "~> 2.0"

    }

    null = {
      source = "hashicorp/null"

      version = "3.1.1"

    }

    okta = {
      source = "oktadeveloper/okta"

      version = "aversion"

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
# tflint-ignore: terraform_unused_declarations
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project               = "foofoo"
    env                   = "accounts"
    service               = "foo"
    owner                 = "foo@example.com"
    date                  = "2023-09-13"
    terraformWorkspaceDir = "/terraform/accounts/foo"
    foggUser              = ""
    gitRepository         = "git@github.com:chanzuckerberg/fogg"
    commitHash            = ""
    managedBy             = "terraform"
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
