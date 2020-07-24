---
parent: Guide
nav_order: 1
layout: default
title: Using Remote State
has_toc: true
---

# Using Remote State
Fogg manages [remote state data sources](https://www.terraform.io/docs/providers/terraform/d/remote_state.html) for all workspaces in a repository.

Given the following `fogg.yml` snippet, the `webserver` workspace can refer to outputs of the `database` and `main_accounts` workspace (see [the fogg docs homepage]({% link index.md %}#fogg-concepts) for more detail on these relationships)

```yaml
<snip>
accounts:
  main_account: {}
envs:
  development:
    components:
      webserver: {}
      database: {}
```

Often enough, our web service is dependent on some values from our database service!

If the `database` workspace defines some useful outputs in `outputs.tf`:
```hcl
output database_uri {
  value = aws_rds_cluster.db.endpoint
}
```

Terraform code in the `webserver` workspace can refer to it and use it directly in resource definitions or module invocations:
```hcl
module my_web_service {
  source       = "../../../modules/webserver"
  database_uri = data.terraform_remote_state.database.outputs.database_uri
}
```
