# CDKTF Fogg Helpers

This CDKTF library provides Helpers to integrate CDKTF applications with Fogg generated component configuration.

[npmjs](https://www.npmjs.com/org/vincenthsh)

## Usage

### FoggStack

FoggStack exposes Fogg component configuration to CDKTF configurations.

> [!NOTE]
> Refer to [unit tests](https://github.com/vincenthsh/fogg/blob/feat-multi-module-components/packages/cdktf-fogg-constructs/test/fogg-stack.test.ts)
> for supported functionality.

```typescript
import { FoggStack } from "@vincenthsh/cdktf-fogg-helpers"
import { Construct } from "constructs"

export class ComponentStack extends FoggStack {
  constructor(scope: Construct, id: string) {
    super(scope, id, {
      // Optionally force Fogg generated Remote State references to be loaded:
      //   By default Fogg generates remote states for all other components in the
      //   same environment due to backward compatibility. FoggStack ignores these
      //   unless forceRemoteStates is enabled.
      forceRemoteStates: false
    })

    // Example setting variables for a single `module_source` backed Component
    this.setMainModuleVariables({
      // The AWS region (required by fogg) is exposed directly
      region: this.region ?? "us-east-1",
      foo: "bar",
      baz: "qux"
    })
  }
}
```

Get typed Remote State outputs from non CDKTF components:

```typescript
    interface NetworkOutputs {
      vpc_id: string
    }
    const outputs = stack.remoteState<NetworkOutputs>("network")
    // Access network output simply through property getters
    // for example `vpc_id` to get the output named "vpc_id"
    const vpc = outputs.vpc_id
```

See [Output Schema Utilities](#output-schema-utilities) for more details on managing Remote State output types.

> Fogg generates a [type declaration](https://github.com/vincenthsh/fogg/blob/feat-multi-module-components/testdata/v2_cdktf_components/terraform/envs/test/network-old/index.d.ts)
for every HCL component.
>
> Outputs can be loaded through Component dependencies:
>
> ```typescript
> import { type NetworkOldOutputs } from "test-network-old"
> ```

Set variables for any TF Module configured for the Component, by "Module Name":

```typescript
    // Configure module variables by module name
    this.setModuleVariables("my-module-name", {
      foo: "bar",
      baz: "qux",
      tags: {
        // Also access other Fogg component Attributes
        Project: this.foggComp.project,
        Owner: this.foggComp.owner,
        Environment: this.foggComp.env,
      }
    })
```

Define new AWS Resources directly in a component

```typescript
// Add import statement
import { dataAwsAvailabilityZones } from "@cdktf/provider-aws"

// Instantiated in Stack body (using scope: `this`)
    const azs = new dataAwsAvailabilityZones.DataAwsAvailabilityZones(
      this,
      "azs",
      {}
    )
```

Add any custom CDKTF Constructs within the Repository [Workspace](https://pnpm.io/workspaces) or from the NPM registry.

```typescript
import { MyConstrct } from "@vincenthsh/my-construct"

// instantiate in Stack body
    new MyConstruct(this, "my-construct", {
      foo: "bar",
    })
```

### FoggTerraStack

FoggTerraStack exposes Fogg component configuration for advanced [TerraConstructs.dev](https://terraconstructs.dev) AwsStacks.

> [!NOTE]
> Refer to [unit tests](https://github.com/vincenthsh/fogg/blob/feat-multi-module-components/packages/cdktf-fogg-constructs/test/fogg-terrastack.test.ts)
> for supported functionality.

```typescript
import { FoggTerraStack } from "@vincenthsh/fogg"
import { Construct } from "constructs"
import { Duration } from "terraconstructs"
import { Bucket } from "terraconstructs/lib/aws/storage"
import { Queue } from "terraconstructs/lib/aws/notify"

export class ComponentStack extends FoggTerraStack {
  constructor(scope: Construct, id: string) {
    super(scope, id, {
      // optionally force Fogg generated Remote State references to be loaded
      forceRemoteStates: false,
      // Prefix for resource UUID, should never change for component lifecycle
      gridUUID: "123-456",
      // This is used as a tag and may change over time
      environmentName: "development"
    })

    // Use Advanced https://terraconstructs.dev L2 Constructs.
    new Bucket(this, "MyEventBridgeBucket", {
      forceDestroy: true,
      eventBridgeEnabled: true,
      enforceSSL: true
    })

    new Queue(stack, "Queue", {
      namePrefix: "queue.fifo",
      messageRetentionSeconds: Duration.days(14).toSeconds(),
      visibilityTimeoutSeconds: Duration.minutes(15).toSeconds()
    })
  }
}
```

### Output Schema Utilities

The Fogg helpers library also provides Utility types to work with Terraform Outputs across CDKTF Stacks.

#### Fixed Outputs

The type of Outputs is usually known in advance, for example a `Network` Construct
may provide known outputs for:

- `vpc_id`: The VPC ID for the network
- `app_subnet_ids`: A list of Application Subnets
- `endpoints`: A Map of Service Endpoints provisioned within the VPC

In this case, the `Network` Construct authors may define a `NetworkOutputs` Interface for
the Terraform Outputs created by their Construct as follows:

```typescript
import type { MapOutput, ListOutput } from '@vincenthsh/cdktf-fogg-helpers'

export interface ServiceEndpoints {
  apigw: string,
  s3: string
}
export interface NetworkOutputs {
  vpc_id: string
  app_subnet_ids: ListOutput<string>
  endpoints: MapOutput<ServiceEndpoints>
}
```

Because TypeScript Type information is not available during CDKTF synthesis, Construct authors should also
provide the schema as a value available at runtime, using the `InferOutputSchema` Utility type:

```typescript
import type { InferOutputSchema } from '@vincenthsh/cdktf-fogg-helpers'
export const networkOutputsSchema: InferOutputSchema<NetworkOutputs> = {
  vpc_id: 'string',
  app_subnet_ids: 'list', // Compilation will fail if not set to 'list'
  endpoints: 'map'        // Compilation will fail if not set to 'map'
}
```

> [!NOTE]
> Frameworks such as [zod](https://zod.dev/) are under evaluation to provide alternative runtime validations with less duplication.

Once Construct authors provide this information, remote state can be accessed across states
with Type Checks as shown in the following example.

In this example `my-network-state` is a workspaced CDKTF project using the `Network` Construct and
re-exporting its `NetworkOutputs` Interface as well as the `networkOutputsSchema` runtime information.

Any Fogg Component depending on `my-network-state`, may acces these Outputs with strong type checks as follows:

```typescript
import { NetworkOutputs, networkOutputsSchema } from 'my-network-state'

//.. in class MyComponent extends FoggStack
  const networkOutputs = this.remoteState<NetworkOutputs>(
    'network',              // the remote state reference name
    networkOutputsSchema    // the remote state output schema
  )

  // Output property access -> synthesized Terraform
  networkOutputs.vpc_id
  // -> ${data.terraform_remote_state.network.outputs.vpc_id}

  networkOutputs.app_subnet_ids.element(0)
  // -> ${element(data.terraform_remote_state.network.outputs.app_subnet_ids, 0)}

  networkOutputs.endpoints.lookup('apigw') // Compilation error unless one of ['apigw', 's3']
  // -> ${data.terraform_remote_state.network.outputs.endpoints.apigw}
```

#### Dynamic Outputs

In some cases, the shape of Terraform Outputs is dynamic and their Interface
is not known until instantiated (i.e. depending on certain input variables).

In this example the number of Network Endpoints are dynamically created depending
on the `Network` Construct `endpoints` input, consumers may need to know the
private domain name for endpoints to keep network traffic internal:

```typescript
// terraform/envs/staging/landing-zone
import { Network } from "landing-zone-constructs"

export { NetworkOutput, networkOutputSchema } from "landing-zone-constructs"
// Must be an imutable array (... as const)
export const serviceEndpoints = ['apigw', 's3'] as const

// ... class StagingStack extends FoggStack

  new Network(this, "network", {
    // convert to dynamic string[] with Array.from()
    endpoints: Array.from(serviceEndpoints)
  })
```

A Consumer may depend on the `staging-landing-zone` Component and have strongly typed output access as follows:

```typescript
import { NetworkOutput, networkOutputSchema, serviceEndpoints } from "staging-landing-zone"

// ideally these are defined in staging-landing-zone for all dependants
type Endpoints = {
  // immutable string array of endpoint names
  [K in (typeof serviceEndpoints)[number]]: string
}
interface NetworkOutputs {
  vpc_id: string
  app_subnet_ids: ListOutput<string>
  endpoints: MapOutput<Endpoints>
}

//.. in class MyStack extends FoggStack
  const networkOutputs = this.remoteState<NetworkOutputs>(
    'landing-zone',         // the remote state alias for 'staging-landing-zone'
    networkOutputsSchema    // the remote state output schema
  )
  // TypeScript validates if staging-landing-zone has 'apigw' serviceEndpoints input
  networkOutputs.endpoints.lookup('apigw');
  // -> ${data.terraform_remote_state.network.outputs.endpoints.apigw}

  // NOTE: nested attributes are also accessible, but are not type-checked
  networkOutputs.endpoints.lookupNested('apigw', 'domain');
  // -> ${data.terraform_remote_state.network.outputs.endpoints.apigw.domain}

  // WARNING: lookupNested() is NOT type-checked,
  // Typescript allows this, but `terraform plan` will fail
  networkOutputs.endpoints.lookupNested('apiGateway', 'domain')
  // -> ${data.terraform_remote_state.network.outputs.endpoints.apiGateway.domain}
```

#### Unit Testing Output Interfaces

To ensure the Output Interface is kept in-sync with the actual Terraform Outputs in the synthesized stack,
the `cdktf-fogg-constructs` library provides [Vitest](https://vitest.dev/) Helpers which are used as follows:

```typescript
// test/outputs.test.ts
import * as path from 'node:path'
import {
  SchemaVitestMatchers, // Vitest Schema Assertions
  toMatchSchema         // Vitest Schema Matcher
} from '@vincenthsh/cdktf-fogg-helpers'
import { Testing } from 'cdktf'
import { describe, it, expect, beforeEach } from 'vitest'
import { MyConstruct } from "../src"

// Extends Vitest's expect types with `toMatchSchema` matcher
// See https://vitest.dev/guide/extending-matchers.html
expect.extend({
    toMatchSchema
})
// Ambient module declaration to satisify Typescript
declare module 'vitest' {
    interface Assertion<T = any> extends SchemaVitestMatchers<T> {}
}

describe('MyConstructOutputs', () => {
    let outputs: any
    // MyConstructsOutput TypeScript Interface is defined in this source File
    const sourceFilePath = path.resolve(__dirname, '../src/outputs.ts')

    beforeEach(() => {
        const template = Testing.synthScope((stack) => {
            new MyConstruct(stack, 'Default', {
              // MyConstruct Construct Properties
            })
        });
        // Read Terraform Outputs from synthesized template
        terraformOutputs = JSON.parse(template).output
    })

    it('should match MyConstruct actual Terraform Outputs', () => {
        expect(terraformOutputs).toBeDefined()
        const ignoredOutputs: string[] = ['fooProperty', 'barProperty']
        expect(terraformOutputs).toMatchSchema(
          sourceFilePath,       // the source File to load TypesScript interface
          'MyConstructOutputs', // the name of the Interface to load
          ignoredOutputs        // any terraformOutputs to ignore (Dynamic)
        )
    })
})
```

> [!NOTE]
> Refer to [unit tests](https://github.com/vincenthsh/fogg/blob/feat-multi-module-components/packages/cdktf-fogg-constructs/test/matchers/schema.test.ts)
> for supported functionality.

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
