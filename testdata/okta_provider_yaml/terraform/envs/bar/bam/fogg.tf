# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.










// https://github.com/articulate/terraform-provider-okta
provider "okta" {
  version  = "~>aversion"
  org_name = "orgname"
}









terraform {
  required_version = "~>1.1.1"


  backend "s3" {
    bucket = "bucket"

    key     = "terraform/foofoo/envs/bar/components/bam.tfstate"
    encrypt = true
    region  = "region"
    profile = "foofoo"
  }

}

variable "env" {
  type    = string
  default = "bar"
}

variable "project" {
  type    = string
  default = "foofoo"
}



variable "component" {
  type    = string
  default = "bam"
}




variable "owner" {
  type    = string
  default = "foo@example.com"
}

variable "tags" {
  type = map(string)
  default = {
    project   = "foofoo"
    env       = "bar"
    service   = "bam"
    owner     = "foo@example.com"
    managedBy = "terraform"
  }
}



data "terraform_remote_state" "global" {
  backend = "s3"

  config = {
    bucket = "bucket"

    key     = "terraform/foofoo/global.tfstate"
    region  = "region"
    profile = "foofoo"
  }
}




# remote state for accounts

data "terraform_remote_state" "foo" {
  backend = "s3"

  config = {
    bucket = "bucket"

    key     = "terraform/foofoo/accounts/foo.tfstate"
    region  = "region"
    profile = "foofoo"
  }
}


# map of aws_accounts
variable "aws_accounts" {
  type = map
  default = {



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

provider null {
  version = "~> 2.1"
}

provider local {
  version = "~> 1.4"
}

provider tls {
  version = "~> 2.1"
}
