# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
provider "aws" {

  region  = "us-west-2"
  profile = "profile"

  allowed_account_ids = ["00456"]
}
# Aliased Providers (for doing things in every region).


provider "assert" {}
data "aws_eks_cluster" "cluster" {
  name = data.terraform_remote_state.comp1.outputs.cluster_id
}
data "aws_eks_cluster_auth" "cluster" {
  name = data.terraform_remote_state.comp1.outputs.cluster_id
}
provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}
terraform {
  required_version = "=0.100.0"

  backend "remote" {

    hostname     = "example.com"
    organization = "foo"
    workspaces {
      name = "staging-k8s-comp"
    }

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

    kubernetes = {
      source = "hashicorp/kubernetes"

      version = "2.19.0"

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
  default = "staging"
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
  default = "k8s-comp"
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
# tflint-ignore: terraform_unused_declarations
variable "tags" {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project       = "proj"
    env           = "staging"
    service       = "k8s-comp"
    owner         = "foo@example.com"
    date          = "2023-09-13"
    filePath      = "BasePathFs"
    foggUser      = ""
    gitRepository = "git@github.com:chanzuckerberg/fogg"
    commitHash    = ""
    managedBy     = "terraform"
  }
}
# tflint-ignore: terraform_unused_declarations
variable "foo" {
  type    = string
  default = "bar2"
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
data "terraform_remote_state" "comp1" {
  backend = "remote"
  config = {


    hostname     = "example.com"
    organization = "foo"
    workspaces = {
      name = "staging-comp1"
    }

  }
}
data "terraform_remote_state" "comp2" {
  backend = "remote"
  config = {


    hostname     = "example.com"
    organization = "foo"
    workspaces = {
      name = "staging-comp2"
    }

  }
}
data "terraform_remote_state" "vpc" {
  backend = "remote"
  config = {


    hostname     = "example.com"
    organization = "foo"
    workspaces = {
      name = "staging-vpc"
    }

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
