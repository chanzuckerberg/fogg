# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.


# Default Provider
provider "aws" {
  version             = "~> 0.12.0"
  region              = "us-west-foo1"
  profile             = "czi-foo"
  allowed_account_ids = [2]
}

# Aliased Providers (for doing things in every region).

provider "aws" {
  alias               = "us-east-foo2"
  version             = "~> 0.12.0"
  region              = "us-east-foo2"
  profile             = "czi-foo"
  allowed_account_ids = [2]
}













terraform {
  required_version = "=0.12.0"

  backend "s3" {
    bucket         = "foo-bucket"
    dynamodb_table = "foo-table"
    key            = "terraform/foo-project/accounts/foo.tfstate"
    encrypt        = true
    region         = "us-west-foo1"
    profile        = "czi-foo"
  }
}

variable "project" {
  type    = string
  default = "foo-project"
}


variable "region" {
  type    = string
  default = "us-west-foo1"
}



variable "aws_profile" {
  type    = string
  default = "czi-foo"
}


variable "owner" {
  type    = string
  default = "foo@example.com"
}

variable "aws_accounts" {
  type = map
  default = {


    bar = 3



    foo = 2


  }
}


variable "foo" {
  type    = string
  default = "foo"
}


data "terraform_remote_state" "global" {
  backend = "s3"

  config = {
    bucket         = "the-bucket"
    dynamodb_table = "the-table"
    key            = "terraform/test-project/global.tfstate"
    region         = "us-west-2"
    profile        = "czi"
  }
}




data "terraform_remote_state" "bar" {
  backend = "s3"

  config = {
    bucket         = "foo-bucket"
    dynamodb_table = "foo-table"
    key            = "terraform/foo-project/accounts/bar.tfstate"
    region         = "us-west-foo1"
    profile        = "czi-foo"
  }
}





provider random {
  version = "~> 2.2"
}

provider template {
  version = "~> 2.1"
}

provider archive {
  version = "~> 1.3"
}
