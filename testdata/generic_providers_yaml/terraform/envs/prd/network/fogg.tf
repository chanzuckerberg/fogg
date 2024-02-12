# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

provider "assert" {}

provider "sops" {}
provider "bar" {
}
provider "baz" {
  baz_token = "prod_token_arn"
}
provider "foo" {
  foo_host = "prod"
  foo_tls  = true
}
terraform {
  required_version = "=1.1.1"

  backend "s3" {

    bucket = "bucket"

    key     = "terraform/foo/envs/prd/components/network.tfstate"
    encrypt = true
    region  = "region"
    profile = "foo"


  }
  required_providers {
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.0"
    }
    assert = {
      source  = "bwoznicki/assert"
      version = "0.0.1"
    }
    bar = {
      source  = "czi/bar"
      version = "~> 0.1.0"
    }
    baz = {
      source  = "czi/baz"
      version = "~> 0.1.0"
    }
    foo = {
      source  = "czi/foo"
      version = "~> 0.2"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
    okta-head = {
      source  = "okta/okta"
      version = "~> 3.30"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.4"
    }
    sops = {
      source  = "carlpett/sops"
      version = "0.7.2"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 3.0"
    }
  }
}
# tflint-ignore: terraform_unused_declarations
variable "env" {
  type    = string
  default = "prd"
}
# tflint-ignore: terraform_unused_declarations
variable "project" {
  type    = string
  default = "foo"
}
# tflint-ignore: terraform_unused_declarations
# tflint-ignore: terraform_unused_declarations
variable "component" {
  type    = string
  default = "network"
}
# tflint-ignore: terraform_unused_declarations
variable "owner" {
  type    = string
  default = "foo@example.com"
}
# tflint-ignore: terraform_unused_declarations
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string, tfstateKey : string })
  default = {
    project    = "foo"
    env        = "prd"
    service    = "network"
    owner      = "foo@example.com"
    tfstateKey = "terraform/foo/envs/prd/components/network.tfstate"

    managedBy = "terraform"
  }
}
# tflint-ignore: terraform_unused_declarations
variable "aws_accounts" {
  type = map(string)
  default = {

  }
}