# CDKTF Fogg Helpers

This CDKTF library provides Helpers to integrate CDKTF applications with Fogg generated component configuration.

## Usage

### FoggStack

FoggStack exposes Fogg component configuration to CDKTF configurations.

> [!NOTE]
> Refer to [unit tests](https://github.com/vincenthsh/fogg/blob/feat-multi-module-components/packages/cdktf-fogg-constructs/test/fogg-stack.test.ts)
> for supported functionality.

```typescript
import { FoggStack } from "@vincenthsh/cdktf-fogg-helpers";
import { Construct } from "constructs";

export class ComponentStack extends FoggStack {
  constructor(scope: Construct, id: string) {
    super(scope, id, {
      // optionally force Fogg generated Remote State references to be loaded
      forceRemoteStates: false,
    });

    // Example setting variables for the single `module_source` backed Component
    this.setMainModuleVariables({
      // The component aws region is exposed directly
      region: this.region ?? "us-east-1",
      foo: "bar",
      baz: "qux",
    });
  }
}
```

Get typed Remote State outputs from non CDKTF components:

```typescript
    interface NetworkOutputs {
      vpc_id: string;
      app_subnet_ids: string[];
    }
    const outputs = stack.remoteState<NetworkOutputs>("network");
    new TerraformOutput(stack, "vpc_id", {
      // network outputs are access through property getters
      // for example `vpc_id` to get the output named "vpc_id"
      value: outputs.vpc_id,
      staticId: true,
    });
```

> Fogg generates a [type declaration](https://github.com/vincenthsh/fogg/blob/feat-multi-module-components/testdata/v2_cdktf_components/terraform/envs/test/network-old/index.d.ts)
for every HCL component.

```typescript
import { type Outputs as NetworkOutputs } from "../../network";
```

Set variables for any TF Module configured for the Component, by `name`:

```typescript
    // Configure fogg module variables here
    this.setModuleVariables("my-module-name", {
      foo: "bar",
      baz: "qux",
      tags: {
        // Use fogg component configuration values
        Project: this.foggComp.project,
        Owner: this.foggComp.owner,
        Environment: this.foggComp.env,
      },
    });
```

Define new AWS Resources directly in component

```typescript
// Add import statement
import { dataAwsAvailabilityZones } from "@cdktf/provider-aws";

// Instantiate in Stack body
    const azs = new dataAwsAvailabilityZones.DataAwsAvailabilityZones(
      this,
      "azs",
      {}
    );
```

Add any custom CDKTF Constructs within [Workspace](https://pnpm.io/workspaces) or from NPM registry.

```typescript
import { MyConstrct } from "@vincenthsh/my-construct"

// instantiate in Stack body
    new MyConstruct(this, "my-construct", {
      foo: "bar",
    });
```

### FoggTerraStack

FoggTerraStack exposes Fogg component configuration for [TerraConstructs.dev](https://terraconstructs.dev) AwsStacks.

> [!NOTE]
> Refer to [unit tests](https://github.com/vincenthsh/fogg/blob/feat-multi-module-components/packages/cdktf-fogg-constructs/test/fogg-terrastack.test.ts)
> for supported functionality.

```typescript
import { FoggTerraStack } from "@vincenthsh/fogg";
import { Construct } from "constructs";
import { Duration } from "terraconstructs";
import { Bucket } from "terraconstructs/lib/aws/storage";
import { Queue } from "terraconstructs/lib/aws/notify";

export class ComponentStack extends FoggTerraStack {
  constructor(scope: Construct, id: string) {
    super(scope, id, {
      // optionally force Fogg generated Remote State references to be loaded
      forceRemoteStates: false,
      // Prefix for resource UUID, should never change for component lifecycle
      gridUUID: "123-456",
      // This is used as a tag and may change over time
      environmentName: "development",
    });

    // Use Advanced https://terraconstructs.dev L2 Constructs.
    new Bucket(this, "MyEventBridgeBucket", {
      forceDestroy: true,
      eventBridgeEnabled: true,
      enforceSSL: true,
    });

    new Queue(stack, "Queue", {
      namePrefix: "queue.fifo",
      messageRetentionSeconds: Duration.days(14).toSeconds(),
      visibilityTimeoutSeconds: Duration.minutes(15).toSeconds(),
    });
  }
}
```

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
