# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.
# tflint-ignore: terraform_unused_declarations
data "terraform_remote_state" "global" {
  backend = "remote"
  config = {


    hostname     = "tfe.example.com"
    organization = "test-org"
    workspaces = {
      name = "global"
    }

  }
}
# tflint-ignore: terraform_unused_declarations
data "terraform_remote_state" "acct1" {
  backend = "remote"
  config = {


    hostname     = "tfe.example.com"
    organization = "test-org"
    workspaces = {
      name = "accounts-acct1"
    }

  }
}