# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
provider "aws" {

  region  = "us-west-2"
  profile = "profile"

  allowed_account_ids = ["00456"]
}
# Aliased Providers (for doing things in every region).


provider "sentry" {
}
terraform {
  required_version = "=0.100.0"

  backend "s3" {

    bucket = "buck"

    key     = "terraform/proj/envs/prod/components/sentry.tfstate"
    encrypt = true
    region  = "us-west-2"
    profile = "profile"


  }
  required_providers {

    archive = {
      source = "hashicorp/archive"

      version = "~> 2.0"

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

      version = "~> 3.0"

    }

    random = {
      source = "hashicorp/random"

      version = "~> 2.2"

    }

    sentry = {
      source = "jianyuan/sentry"

      version = "1.2.3"

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
  default = "prod"
}
variable "project" {
  type    = string
  default = "proj"
}
variable "region" {
  type    = string
  default = "us-west-2"
}
variable "component" {
  type    = string
  default = "sentry"
}
variable "aws_profile" {
  type    = string
  default = "profile"
}
variable "owner" {
  type    = string
  default = "foo@example.com"
}
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project   = "proj"
    env       = "prod"
    service   = "sentry"
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
data "terraform_remote_state" "hero" {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/envs/prod/components/hero.tfstate"
    region  = "us-west-2"
    profile = "profile"


  }
}
variable "aws_accounts" {
  type = map(string)
  default = {

    bar = "00456"

    foo = "123"

  }
}
