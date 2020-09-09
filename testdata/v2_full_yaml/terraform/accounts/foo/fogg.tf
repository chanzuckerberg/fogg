# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
provider aws {

  version = "~> 0.12.0"
  region  = "us-west-2"


  assume_role {
    role_arn = "arn:aws:iam::123:role/roll"
  }

  allowed_account_ids = [123]
}
# Aliased Providers (for doing things in every region).


provider bless {
  version = "~>0.4.2"
  region  = "us-west-2"
  profile = "prof"
}
terraform {
  required_version = "=0.100.0"
  backend s3 {

    bucket = "buck"

    key     = "terraform/proj/accounts/foo.tfstate"
    encrypt = true
    region  = "us-west-2"
    profile = "profile"


  }
}
variable project {
  type    = string
  default = "proj"
}
variable region {
  type    = string
  default = "us-west-2"
}
variable owner {
  type    = string
  default = "foo@example.com"
}
variable account {
  type    = string
  default = "foo"
}
variable tags {
  type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
  default = {
    project   = "proj"
    env       = "accounts"
    service   = "foo"
    owner     = "foo@example.com"
    managedBy = "terraform"
  }
}
# map of aws_accounts
variable aws_accounts {
  type = map
  default = {


    bar = 456



    foo = 123


  }
}
variable foo {
  type    = string
  default = "bar1"
}
data terraform_remote_state global {
  backend = "s3"
  config = {


    bucket = "buck"

    key     = "terraform/proj/global.tfstate"
    region  = "us-west-2"
    profile = "profile"


  }
}
data terraform_remote_state bar {
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
provider null {
  version = "~> 2.1"
}
provider local {
  version = "~> 1.4"
}
provider tls {
  version = "~> 2.1"
}
