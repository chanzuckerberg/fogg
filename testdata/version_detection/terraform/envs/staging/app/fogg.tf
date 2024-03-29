# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.


provider "aws" {
  version             = "~> 1.27.0"
  region              = "us-west-2"
  profile             = "acme@infra.exmaple"
  allowed_account_ids = [000]
}

# Aliased Providers (for doing things in every region).











terraform {
  required_version = "~>0.12.5"

  backend "s3" {
    bucket         = "acme"
    dynamodb_table = "acme-auth"

    key = "terraform/acme/envs/staging/components/app.tfstate"


    encrypt = true
    region  = "us-west-2"
    profile = "acme@infra.exmaple"
  }
}

variable "env" {
  type    = string
  default = "staging"
}

variable "project" {
  type    = string
  default = "acme"
}


variable "region" {
  type    = string
  default = "us-west-2"
}


variable "component" {
  type    = string
  default = "app"
}


variable "aws_profile" {
  type    = string
  default = "acme@infra.exmaple"
}



variable "owner" {
  type    = string
  default = "acme-infra"
}

variable "tags" {
  type = map(string)
  default = {
    project   = "acme"
    env       = "staging"
    service   = "app"
    owner     = "acme-infra"
    managedBy = "terraform"
  }
}



data "terraform_remote_state" "global" {
  backend = "s3"

  config = {
    bucket         = "acme"
    dynamodb_table = "acme-auth"
    key            = "terraform/acme/global.tfstate"
    region         = "us-west-2"
    profile        = "acme@infra.exmaple"
  }
}




# remote state for accounts


# map of aws_accounts
variable "aws_accounts" {
  type = map(string)
  default = {

  }
}
