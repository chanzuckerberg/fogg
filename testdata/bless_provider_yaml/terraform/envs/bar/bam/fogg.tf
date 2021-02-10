# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

provider bless {
  region  = "bar"
  profile = "foofoofoo"
}
provider bless {
  alias  = "a"
  region = "a"

  profile = "foofoofoo"
}
provider bless {
  alias  = "b"
  region = "b"

  profile = "foofoofoo"
}
terraform {
  required_version = "=1.1.1"

  backend s3 {

    bucket = "bucket"

    key     = "terraform/foofoo/envs/bar/components/bam.tfstate"
    encrypt = true
    region  = "region"
    profile = "foofoo"


  }
  required_providers {

    archive = {
      source = "hashicorp/archive"

      version = "~> 2.0"

    }

    bless = {
      source = "chanzuckerberg/bless"

      version = "0.0.0"

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
    default = "bar"
  }
  variable project {
    type    = string
    default = "foofoo"
  }
  variable component {
    type    = string
    default = "bam"
  }
  variable owner {
    type    = string
    default = "foo@example.com"
  }
  variable tags {
    type = object({ project : string, env : string, service : string, owner : string, managedBy : string })
    default = {
      project   = "foofoo"
      env       = "bar"
      service   = "bam"
      owner     = "foo@example.com"
      managedBy = "terraform"
    }
  }
  data terraform_remote_state global {
    backend = "s3"
    config = {


      bucket = "bucket"

      key     = "terraform/foofoo/global.tfstate"
      region  = "region"
      profile = "foofoo"


    }
  }
  data terraform_remote_state foo {
    backend = "s3"
    config = {


      bucket = "bucket"

      key     = "terraform/foofoo/accounts/foo.tfstate"
      region  = "region"
      profile = "foofoo"


    }
  }
  variable aws_accounts {
    type = map
    default = {

    }
  }
