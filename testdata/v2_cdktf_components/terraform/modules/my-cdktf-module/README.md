# CDKTF Module

```console
make synth
```

## Working with existing Terraform Modules

> [!NOTE]
> To import TF modules from Handshakes Org:
> Authenticate with Handshakes AWS Org (using `aws-vault`).

[Add Modules to cdktf.json](https://developer.hashicorp.com/terraform/cdktf/concepts/modules#add-module-to-cdktf-json)

[Generate module bindings](https://developer.hashicorp.com/terraform/cdktf/concepts/modules#generate-module-bindings)

```console
    pnpm run get           Import/update Terraform providers and modules (you should check-in this directory)
```

## Test

```console
    pnpm run test        Runs unit tests (edit __tests__/main-test.ts to add your own tests)
    pnpm run test:watch  Watches the tests and reruns them on change
```

> [!IMPORTANT]
> To update the snapshops, run `pnpm run test --update`

## Usage

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
