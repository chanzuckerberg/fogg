# CDKTF Fogg Helpers

This CDKTF library provides Helpers to integrate CDKTF applications with Fogg generated component configuration.

## Compile

```console
    pnpm run get           Import/update Terraform providers and modules (you should check-in this directory)
    pnpm run go:generate   Import/update Typescript types for Fogg golang structs
    pnpm run watch         Watch for changes and compile typescript in the background
    pnpm run build         Compile typescript
```

## Test

```console
    pnpm run test        Runs unit tests (edit __tests__/main-test.ts to add your own tests)
    pnpm run test:watch  Watches the tests and reruns them on change
```

> [!IMPORTANT]
> To update the snapshopts, run `pnpm run test --update`

## Usage

## FoggStack

FoggStack exposes Fogg component configuration to CDKTF configurations.

> [!NOTE]
> Refer to [unit tests](./test/fogg-stack.test.ts) for supported functionality.

```typescript
import { FoggStack } from "@vincenthsh/fogg";
import { Construct } from "constructs";
import { }

export class ComponentStack extends FoggStack {
  constructor(scope: Construct, id: string) {
    super(scope, id, {
      forceRemoteStates: false,
    });

    // Example setting variables for `module_source` backed Component
    this.setMainModuleVariables({
      foo: "bar",
      baz: "qux",
    });
  }
}
```
