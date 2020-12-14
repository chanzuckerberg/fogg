resource null_resource foo {
  triggers = {
    foo = data.terraform_remote_state.one.outputs.foo
  }
}

locals {
  this = data.terraform_remote_state.two.outputs.foo
}

output three {
  value = data.terraform_remote_state.three.outputs.foo
}

provider aws {
  profile = data.terraform_remote_state.four.outputs.foo
}

data tempate_file foo {
  tempalte = "..."

  vars = {
    it = data.terraform_remote_state.five.outputs.foo
  }
}

module foo {
  source = "../"

  foo = data.terraform_remote_state.six.outputs.foo
}
