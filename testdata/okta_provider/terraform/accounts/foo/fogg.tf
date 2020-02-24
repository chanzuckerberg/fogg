# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.










// https://github.com/articulate/terraform-provider-okta
provider "okta" {
  version  = "~>aversion"
  org_name = "orgname"
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
}

variable "project" {
  type    = string
  default = "foofoo"
}





variable "owner" {
  type    = string
  default = "foo@example.com"
}

variable "aws_accounts" {
  type = map
  default = {



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
