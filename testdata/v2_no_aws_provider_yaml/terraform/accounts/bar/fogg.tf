# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

provider "assert" {}
terraform {
  required_version = "=0.100.0"

  backend "s3" {

    bucket = "buck"

    key     = "terraform/proj/accounts/bar.tfstate"
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
  default = "accounts"
}
variable "project" {
  type    = string
  default = "proj"
}
variable "component" {
  type    = string
  default = "bar"
}
variable "account" {
  type    = string
  default = "bar"
}
variable "owner" {
  type    = string
  default = "foo@example.com"
}
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project   = "proj"
    env       = "accounts"
    service   = "bar"
    owner     = "foo@example.com"
    managedBy = "terraform"
  }
}
variable "foo" {
  type    = string
  default = "bar1"
}
data "terraform_remote_state" "global" {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/global.tfstate"
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
variable "aws_accounts" {
  type = map(string)
  default = {

  }
}
