# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
terraform {
  required_version = "=0.100.0"

  backend s3 {

    bucket = "buck"

    key     = "terraform/proj/global.tfstate"
    encrypt = true
    region  = "us-west-2"
    profile = "profile"


  }
  required_providers {

    archive = {
      source = "hashicorp/archive"

      version = "~> 2.0"

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

    template = {
      source = "hashicorp/template"

      version = "~> 2.2"

    }

    tls = {
      source = "hashicorp/tls"

      version = "~> 3.0"

    }

  }
  variable env {
    type    = string
    default = ""
  }
  variable project {
    type    = string
    default = "proj"
  }
  variable component {
    type    = string
    default = "global"
  }
  variable owner {
    type    = string
    default = "foo@example.com"
  }
  variable tags {
    type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
    default = {
      project   = "proj"
      env       = ""
      service   = "global"
      owner     = "foo@example.com"
      managedBy = "terraform"
    }
  }
  variable foo {
    type    = string
    default = "bar1"
  }
  variable aws_accounts {
    type = map
    default = {

    }
  }
