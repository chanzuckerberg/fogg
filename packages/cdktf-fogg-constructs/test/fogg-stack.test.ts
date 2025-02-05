import { beforeEach, describe, expect, it, vi } from "vitest";
import "cdktf/lib/testing/adapters/jest";
import { App, Fn, TerraformOutput, Testing } from "cdktf";
import merge from "deepmerge";
import { FoggStack } from "../src/fogg-stack";
import type { Component } from "../src/imports/fogg-types.generated";
import { type OutputSchema } from "../src/util/types";
import { Template } from "./assertions";

// hoisted mock, can't capture variable for vi.fn() mock
vi.mock(import("../src/util/load-component-config"), async (importOriginal) => {
  const original = await importOriginal();
  return {
    ...original,
    loadComponentConfig: vi.fn(),
  };
});
// get Mock handle through module re-import with mock in effect
import { loadComponentConfig as mockedLoadCompCfg } from "../src/util/load-component-config";
const loadComponentConfig = mockedLoadCompCfg as ReturnType<typeof vi.fn>;

let app: App;
describe("FoggStack", () => {
  beforeEach(() => {
    app = Testing.app();
    loadComponentConfig.mockReset();
  });

  it("translates fogg backend config to tf backend config", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        backend: {
          kind: "s3",
          s3: {
            bucket: "bucket-1",
            key_path: "fake-keypath",
            region: "us-fake-1",
          },
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack");
    // THEN
    Template.fromStack(stack).toMatchObject({
      terraform: {
        backend: {
          s3: {
            bucket: "bucket-1",
            key: "fake-keypath",
            region: "us-fake-1",
          },
        },
      },
    });
  });

  it("throws error on unsupported backend", () => {
    // WHEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        backend: { kind: "unknown" },
      })
    );
    // THEN
    expect(() => new FoggStack(app, "MyStack")).toThrow(
      /Unsupported backend configuration/
    );
  });

  it("skips creating remote states if forceRemoteStates is false", () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        component_backends: {
          network: {
            kind: "s3",
            s3: { bucket: "some-bucket" },
          },
        },
        component_backends_filtered: false,
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack", {
      forceRemoteStates: false,
    });
    // THEN
    // no remote states expected
    expect(() => stack.remoteState("network")).toThrow(
      /Remote state network not found/
    );
  });

  describe("creates remote states if forceRemoteStates is true", () => {
    let stack: FoggStack;
    beforeEach(() => {
      loadComponentConfig.mockReturnValue(
        merge(getDefaultComponentConfig(), {
          component_backends: {
            network: {
              kind: "s3",
              s3: { bucket: "some-bucket" },
            },
          },
        })
      );
      stack = new FoggStack(app, "MyStack", {
        forceRemoteStates: true,
      });
    });
    // with output Typing
    interface NetworkOutputs {
      vpc_id: string;
      app_subnet_ids: string[];
    }
    const networkOutputsSchema: OutputSchema<NetworkOutputs> = {
      vpc_id: "string",
      app_subnet_ids: "list",
    };
    it("as data terraform remote state", () => {
      Template.fromStack(stack).toMatchObject({
        data: {
          terraform_remote_state: {
            network: {
              backend: "s3",
              config: {
                bucket: "some-bucket",
              },
            },
          },
        },
      });
    });
    it("with strongly typed access to simple outputs", () => {
      const outputs = stack.remoteState<NetworkOutputs>("network");
      new TerraformOutput(stack, "vpc_id", {
        value: outputs.vpc_id,
        staticId: true,
      });
      Template.fromStack(stack).toMatchObject({
        output: {
          vpc_id: {
            value: "${data.terraform_remote_state.network.outputs.vpc_id}",
          },
        },
      });
    });
    // TODO: This isn't throwing an error...
    it.todo("throws on incorrectly accessing list items", () => {
      const outputs = stack.remoteState<NetworkOutputs>("network");
      new TerraformOutput(stack, "app_subnet1", {
        value: outputs.app_subnet_ids[0],
        staticId: true,
      });
      expect(() => Testing.synth(stack)).toThrow(
        /Found an encoded list token string in a scalar string context/
      );
    });
    it("throws on incorrectly accessing list items", () => {
      const outputs = stack.remoteState<NetworkOutputs>(
        "network",
        networkOutputsSchema
      );
      new TerraformOutput(stack, "app_subnet1", {
        // should be Fn.element(outputs.app_subnet_ids, 0)
        value: outputs.app_subnet_ids[0],
        staticId: true,
      });
      expect(() => Testing.synth(stack)).toThrow(
        /Found an encoded list token string in a scalar string context/
      );
    });
    it("with strongly typed access to complex outputs", () => {
      const outputs = stack.remoteState<NetworkOutputs>(
        "network",
        networkOutputsSchema
      );
      new TerraformOutput(stack, "vpc_id", {
        value: outputs.vpc_id,
        staticId: true,
      });
      new TerraformOutput(stack, "app_subnet1", {
        value: Fn.element(outputs.app_subnet_ids, 0),
        staticId: true,
      });
      new TerraformOutput(stack, "app_subnet2", {
        value: Fn.element(outputs.app_subnet_ids, 1),
        staticId: true,
      });
      Template.fromStack(stack, { snapshot: false }).toMatchObject({
        output: {
          app_subnet1: {
            value:
              "${element(data.terraform_remote_state.network.outputs.app_subnet_ids, 0)}",
          },
          app_subnet2: {
            value:
              "${element(data.terraform_remote_state.network.outputs.app_subnet_ids, 1)}",
          },
          vpc_id: {
            value: "${data.terraform_remote_state.network.outputs.vpc_id}",
          },
        },
      });
    });
  });

  it("exposes file dependencies through locals", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        locals_block: {
          foo_foo_foo: 'yamldecode(file("../../../foo-fooFoo.yaml"))',
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack");
    // THEN
    // can access locals in CDKTF code
    expect(stack.getLocal("foo_foo_foo")).toBeDefined();
    // locals block is defined in the template
    Template.fromStack(stack).toMatchObject({
      locals: {
        foo_foo_foo: '${yamldecode(file("../../../foo-fooFoo.yaml"))}',
      },
    });
  });

  it("allows configuring modules", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        modules: [
          {
            name: "network",
            source: "terraform-aws-modules/vpc/aws",
            version: "5.12.0",
            // TODO: should we validate config with this?
            variables: ["name", "cidr", "azs"],
          },
        ],
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack");
    const networkInputs = {
      name: "my-vpc",
      cidr: "10.0.0.0/16",
      azs: ["us-west-2a", "us-west-2b", "us-west-2c"],
    };
    stack.setModuleVariables("network", networkInputs);
    // THEN
    Template.fromStack(stack).toMatchObject({
      module: {
        network: {
          ...networkInputs,
          source: "terraform-aws-modules/vpc/aws",
          version: "5.12.0",
        },
      },
    });
  });

  it("allows configuring main module", async () => {
    // GIVEN
    const module_source = "terraform/modules/foo";
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        module_source,
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack");
    const moduleInputs = {
      foo: "bar",
      baz: "qux",
    };
    stack.setMainModuleVariables(moduleInputs);
    // THEN
    Template.fromStack(stack).toMatchObject({
      module: {
        main: {
          ...moduleInputs,
          source: module_source,
        },
      },
    });
  });

  it("throws error setting variable without module_source", () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(getDefaultComponentConfig());
    // THEN
    const stack = new FoggStack(app, "MyStack");
    expect(() => stack.setMainModuleVariables({})).toThrow(
      /Module main not found/
    );
  });

  it("exposes default Aws provider", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        providers_configuration: {
          aws: {
            account_id: "123456789012",
            region: "us-west-2",
            role_arn: "arn:aws:iam::123456789012:role/role",
          },
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack");

    // THEN
    expect(stack.defaultAwsProvider).toBeDefined();
  });

  it("exposes aliased providers", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        providers_configuration: {
          aws: {
            account_id: "123456789012",
            region: "us-west-2",
            role_arn: "arn:aws:iam::123456789012:role/role",
          },
          aws_regional_providers: [
            {
              account_id: "210987654321",
              alias: "shared_services",
              region: "ap-southeast-1",
              role_arn: "arn:aws:iam::210987654321:role/role",
            },
          ],
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack");
    // THEN
    Template.fromStack(stack).toMatchObject({
      provider: {
        aws: [
          {
            assume_role: [
              {
                role_arn: "arn:aws:iam::123456789012:role/role",
              },
            ],
            region: "us-west-2",
          },
          {
            alias: "shared_services",
            assume_role: [
              {
                role_arn: "arn:aws:iam::210987654321:role/role",
              },
            ],
            region: "ap-southeast-1",
          },
        ],
      },
    });
  });

  it("sets foggComponent default tags", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        providers_configuration: {
          aws: {
            account_id: "123456789012",
            region: "us-west-2",
            default_tags: {
              enabled: true,
              tags: {
                env: "test",
                owner: "me",
              },
            },
          },
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack");
    // THEN
    Template.fromStack(stack).toMatchObject({
      provider: {
        aws: [
          {
            region: "us-west-2",
            default_tags: [
              {
                tags: {
                  env: "test",
                  owner: "me",
                  managedBy: "terraform",
                  service: "fake-name",
                },
              },
            ],
          },
        ],
      },
    });
  });

  it("sets foggComponent ignore tags", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        providers_configuration: {
          aws: {
            account_id: "123456789012",
            region: "us-west-2",
            ignore_tags: {
              enabled: true,
              keys: ["foo", "bar"],
            },
          },
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack");
    // THEN
    Template.fromStack(stack).toMatchObject({
      provider: {
        aws: [
          {
            region: "us-west-2",
            ignore_tags: [
              {
                keys: ["foo", "bar"],
              },
            ],
          },
        ],
      },
    });
  });

  it("sets foggComponent default and constructor default tags", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        providers_configuration: {
          aws: {
            account_id: "123456789012",
            region: "us-west-2",
            default_tags: {
              enabled: true,
              tags: {
                env: "test",
                owner: "me",
              },
            },
          },
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack", {
      defaultTags: {
        owner: "me2",
        service: "my-service",
        purpose: "my-purpose",
      },
    });
    // THEN
    Template.fromStack(stack).toMatchObject({
      provider: {
        aws: [
          {
            region: "us-west-2",
            default_tags: [
              {
                tags: {
                  env: "test",
                  owner: "me2",
                  managedBy: "terraform",
                  service: "my-service",
                  purpose: "my-purpose",
                },
              },
            ],
          },
        ],
      },
    });
  });

  it("sets foggComponent ignore and constructor ignore tags", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        providers_configuration: {
          aws: {
            account_id: "123456789012",
            region: "us-west-2",
            ignore_tags: {
              enabled: true,
              keys: ["foo", "bar"],
            },
          },
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack", {
      ignoreTags: {
        keys: ["foo", "bar", "baz"],
      },
    });
    // THEN
    Template.fromStack(stack).toMatchObject({
      provider: {
        aws: [
          {
            region: "us-west-2",
            ignore_tags: [
              {
                keys: ["foo", "bar", "baz"],
              },
            ],
          },
        ],
      },
    });
  });

  it("sets constructor ignore tags", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        providers_configuration: {
          aws: {
            account_id: "123456789012",
            region: "us-west-2",
          },
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack", {
      ignoreTags: {
        keys: ["foo", "bar", "baz"],
      },
    });
    // THEN
    Template.fromStack(stack).toMatchObject({
      provider: {
        aws: [
          {
            region: "us-west-2",
            ignore_tags: [
              {
                keys: ["foo", "bar", "baz"],
              },
            ],
          },
        ],
      },
    });
  });

  it("exposes fogg component configuration values", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(getDefaultComponentConfig());
    // WHEN
    const stack = new FoggStack(app, "MyStack");
    // THEN
    expect(stack.foggComp.project).toBeDefined();
    expect(stack.foggComp.owner).toBeDefined();
    expect(stack.foggComp.env).toBeDefined();
  });

  it("parses bundled and generic provider configurations", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        providers_configuration: {
          datadog: {
            enabled: true,
            version: "2.0.0",
          },
        },
        required_providers: {
          cloudflare: {
            enabled: true,
            source: "cloudflare/cloudflare",
            version: "~> 4.0",
            config: {
              // fogg does not generate CDKTF's expected `apiToken`!
              api_token: "bar",
            },
          },
        },
      })
    );
    // WHEN
    const stack = new FoggStack(app, "MyStack");
    // THEN
    expect(stack.region).toBe("us-fake-1");
    Template.fromStack(stack).toMatchObject({
      provider: {
        aws: [
          {
            region: "us-fake-1",
          },
        ],
        cloudflare: [
          {
            api_token: "bar",
          },
        ],
        datadog: [{}],
      },
      terraform: {
        required_providers: {
          aws: {
            source: "aws",
          },
          cloudflare: {
            source: "cloudflare/cloudflare",
          },
          datadog: {
            source: "DataDog/datadog",
          },
        },
      },
    });
  });

  it("throws for unsupported generic providers", async () => {
    // GIVEN
    loadComponentConfig.mockReturnValue(
      merge(getDefaultComponentConfig(), {
        required_providers: {
          bar: {
            enabled: true,
            source: "foo/bar",
            version: "~> 4.0",
            config: {
              api_token: "bar",
            },
          },
        },
      })
    );
    // THEN
    expect(() => new FoggStack(app, "MyStack")).toThrow(
      /Unsupported provider foo\/bar/
    );
  });
});

/**
 * Helper function to generate a mock component configuration
 */
function getDefaultComponentConfig(): Component {
  const defaultCiComponent = {
    Enabled: false,
    Buildevents: false,
    AWSProfileName: "",
    AWSRoleName: "",
    AWSRegion: "",
    AWSAccountID: "",
    Command: "",
  };
  return {
    path_to_repo_root: ".",
    terraform_version: "1.7.5",
    account_backends: {},
    all_accounts: {},
    backend: {
      kind: "s3",
      s3: {
        account_name: "fake-account",
        bucket: "fake-bucket",
        key_path: "fake-keypath",
        region: "us-fake-1",
      },
    },
    package_fields: {},
    component_backends: {},
    autoplan_relative_globs: [],
    autoplan_files: [],
    locals_block: {},
    component_backends_filtered: false,
    env: "test-env-1",
    extra_vars: {},
    name: "fake-name",
    owner: "fake-owner",
    project: "fake-project",
    providers_configuration: {
      aws: {
        account_id: "123456789012",
        region: "us-fake-1",
      },
      aws_regional_providers: [],
    },
    required_providers: {},
    provider_versions: {},
    cdktf_dependencies: {},
    cdktf_dev_dependencies: {},
    tf_lint: {
      enabled: false,
    },
    TravisCI: defaultCiComponent,
    CircleCI: {
      ...defaultCiComponent,
      SSHFingerprints: [],
    },
    GitHubActionsCI: defaultCiComponent,
    eks: {
      cluster_name: "fake-cluster",
    },
    providers: {},
    variables: [],
    outputs: [],
    modules: [],
  };
}
