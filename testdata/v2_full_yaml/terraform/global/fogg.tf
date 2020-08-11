# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
provider aws {

  version = "~> 0.12.0"
  region  = "us-west-2"
  profile = "profile"

  allowed_account_ids = [456]
}
# Aliased Providers (for doing things in every region).

terraform {
  required_version = "~>0.100.0"

  backend s3 {

    bucket = "buck"

    key     = "terraform/proj/global.tfstate"
    encrypt = true
    region  = "us-west-2"
    profile = "profile"


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
variable region {
  type    = string
  default = "us-west-2"
}
variable component {
  type    = string
  default = "global"
}
variable aws_profile {
  type    = string
  default = "profile"
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
