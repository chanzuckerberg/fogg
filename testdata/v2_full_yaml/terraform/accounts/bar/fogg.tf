# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
variable "TFC_RUN_ID" {
  type    = string
  default = "unknown"
}
variable "TFC_WORKSPACE_NAME" {
  type    = string
  default = "unknown"
}
variable "TFC_WORKSPACE_SLUG" {
  type    = string
  default = "unknown"
}
variable "TFC_CONFIGURATION_VERSION_GIT_BRANCH" {
  type    = string
  default = "unknown"
}
variable "TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA" {
  type    = string
  default = "unknown"
}
variable "TFC_CONFIGURATION_VERSION_GIT_TAG" {
  type    = string
  default = "unknown"
}
variable "TFC_PROJECT_NAME" {
  type    = string
  default = "unknown"
}
provider "aws" {

  region = "us-west-2"


  assume_role {
    role_arn = "arn:aws:iam::00456:role/foo"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["00456"]
}
# Aliased Providers (for doing things in every region).


provider "aws" {
  alias  = "another_account_same_role"
  region = "us-west-2"


  assume_role {
    role_arn = "arn:aws:iam::different_account:role/foo"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["different_account"]
}


provider "aws" {
  alias  = "another_account_same_role-us-east-2"
  region = "us-east-2"


  assume_role {
    role_arn = "arn:aws:iam::different_account:role/foo"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["different_account"]
}


provider "aws" {
  alias  = "another_account_same_role-us-east-1"
  region = "us-east-1"


  assume_role {
    role_arn = "arn:aws:iam::different_account:role/foo"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["different_account"]
}


provider "aws" {
  alias  = "another_account_more_regions"
  region = "us-west-2"


  assume_role {
    role_arn = "arn:aws:iam::another_different_account:role/foo"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["another_different_account"]
}


provider "aws" {
  alias  = "another_account_more_regions-eu-west-2"
  region = "eu-west-2"


  assume_role {
    role_arn = "arn:aws:iam::another_different_account:role/foo"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["another_different_account"]
}


provider "aws" {
  alias  = "another_account_different_role"
  region = "us-west-2"


  assume_role {
    role_arn = "arn:aws:iam::789:role/different_role"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["789"]
}


provider "aws" {
  alias  = "another_account_different_role-us-east-2"
  region = "us-east-2"


  assume_role {
    role_arn = "arn:aws:iam::789:role/different_role"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["789"]
}


provider "aws" {
  alias  = "another_account_different_role-us-east-1"
  region = "us-east-1"


  assume_role {
    role_arn = "arn:aws:iam::789:role/different_role"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["789"]
}


provider "aws" {
  alias  = "us-east-2"
  region = "us-east-2"


  assume_role {
    role_arn = "arn:aws:iam::00456:role/foo"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["00456"]
}


provider "aws" {
  alias  = "us-east-1"
  region = "us-east-1"


  assume_role {
    role_arn = "arn:aws:iam::00456:role/foo"
  }

  # this is the new way of injecting AWS tags to all AWS resources
  # var.tags should be considered deprecated
  default_tags {
    tags = {
      TFC_RUN_ID                               = coalesce(var.TFC_RUN_ID, "unknown")
      TFC_WORKSPACE_NAME                       = coalesce(var.TFC_WORKSPACE_NAME, "unknown")
      TFC_WORKSPACE_SLUG                       = coalesce(var.TFC_WORKSPACE_SLUG, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_BRANCH     = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_BRANCH, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA, "unknown")
      TFC_CONFIGURATION_VERSION_GIT_TAG        = coalesce(var.TFC_CONFIGURATION_VERSION_GIT_TAG, "unknown")
      TFC_PROJECT_NAME                         = coalesce(var.TFC_PROJECT_NAME, "unknown")
      project                                  = coalesce(var.tags.project, "unknown")
      env                                      = coalesce(var.tags.env, "unknown")
      service                                  = coalesce(var.tags.service, "unknown")
      owner                                    = coalesce(var.tags.owner, "unknown")
      managedBy                                = "terraform"
    }
  }
  allowed_account_ids = ["00456"]
}


provider "assert" {}

provider "bless" {
  region   = "us-west-2"
  role_arn = "arn:aws:iam::1234567890:role/roll"
}
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

    aws = {
      source = "hashicorp/aws"

      version = "0.12.0"

    }

    bless = {
      source = "chanzuckerberg/bless"

      version = "0.4.2"

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
  default = "bar"
}
# tflint-ignore: terraform_unused_declarations
variable "account" {
  type    = string
  default = "bar"
}
# tflint-ignore: terraform_unused_declarations
variable "owner" {
  type    = string
  default = "foo@example.com"
}
# tflint-ignore: terraform_unused_declarations
# DEPRECATED: this field is deprecated in favor or 
# AWS provider default tags.
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
