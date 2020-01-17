# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.













terraform {
  required_version = "=0.100.0"

  backend "s3" {
    bucket = "buck"

    key     = "terraform/proj/accounts/foo.tfstate"
    encrypt = true
    region  = "us-west-2"
    profile = "profile"
  }
}

variable "project" {
  type    = string
  default = "proj"
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





provider random {
  version = "~> 2.2"
}

provider template {
  version = "~> 2.1"
}

provider archive {
  version = "~> 1.3"
}
