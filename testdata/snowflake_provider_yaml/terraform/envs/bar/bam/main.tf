locals {
  foo = data.terraform_remote_state.foo
}

resource random_number foo {
  name = local.foo
}
