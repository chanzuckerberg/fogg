# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.











terraform {
  required_version = "~>1.1.1"

  backend "s3" {
    bucket = "bucket"


    key = "terraform/foo/global.tfstate"

    encrypt = true
    region  = "region"
    profile = "foo"
  }
}

variable "env" {
  type    = "string"
  default = ""
}

variable "project" {
  type    = "string"
  default = "foo"
}



variable "component" {
  type    = "string"
  default = "global"
}



variable "owner" {
  type    = "string"
  default = "foo@example.com"
}

variable "tags" {
  type = "map"
  default = {
    project   = "foo"
    env       = ""
    service   = "global"
    owner     = "foo@example.com"
    managedBy = "terraform"
  }
}





