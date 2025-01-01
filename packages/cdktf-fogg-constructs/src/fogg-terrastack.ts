import { provider as awsProvider } from "@cdktf/provider-aws";
import { provider as cloudflareProvider } from "@cdktf/provider-cloudflare";
import { provider as dataDogProvider } from "@cdktf/provider-datadog";
import {
  DataTerraformRemoteState,
  DataTerraformRemoteStateS3,
  DataTerraformRemoteStateS3Config,
  S3Backend,
  S3BackendConfig,
  TerraformHclModule,
  TerraformLocal,
} from "cdktf";
import { Construct } from "constructs";
import { AwsStack, AwsStackProps } from "terraconstructs/lib/aws";
import {
  AWSProvider,
  Backend,
  Component,
  DatadogProvider,
  GenericProvider,
} from "./imports/fogg-types.generated";
import { loadComponentConfig } from "./util/load-component-config";

export interface FoggTerraStackProps
  extends Omit<AwsStackProps, "providerConfig"> {
  /**
   * Force remote state configuration
   * @default false - only configure remote states if `component_backends_filtered` is true
   */
  forceRemoteStates?: boolean;
}

// TODO: Convert Fogg component configuration handlers to component class

/**
 * Helper AwsStack to wrap Fogg component configuration and set up configured providers and backends.
 */
export class FoggTerraStack extends AwsStack {
  public readonly foggComp: Component;
  public readonly modules: Record<string, TerraformHclModule> = {};
  public readonly locals: Record<string, TerraformLocal> = {};
  private readonly _providers: Record<string, awsProvider.AwsProvider> = {};
  private readonly _remoteStates: Record<string, DataTerraformRemoteState> = {};

  constructor(scope: Construct, id: string, props: FoggTerraStackProps) {
    const foggComp = loadComponentConfig();
    const awsProviderConfig = parseAwsProviderConfig(foggComp);
    super(scope, id, {
      ...props,
      providerConfig: awsProviderConfig,
    });
    this.foggComp = foggComp;

    this.parseBackendConfig();
    this.parseBundledProviderConfig();
    for (const p of Object.values(this.foggComp.required_providers)) {
      if (p.enabled) this.parseGenericProviderConfig(p);
    }
    // parse remote backends
    const forceRemoteBackend = props.forceRemoteStates ?? false;
    if (forceRemoteBackend || this.foggComp.component_backends_filtered) {
      for (const [name, remoteStateConfig] of Object.entries(
        this.foggComp.component_backends
      )) {
        this.parseRemoteState(name, remoteStateConfig);
      }
    }
    this.parseLocalsBlock();
    this.parseModules();
  }

  /**
   * Set variables for the main module defined in the fogg component configuration.
   *
   * @param variables - The variables to set for the module
   */
  public setMainModuleVariables(variables: Record<string, any>): void {
    const id = (this.foggComp.module_name =
      this.foggComp.module_name ?? "main");
    this.setModuleVariables(id, variables);
  }

  /**
   * Return a remote state defined in the fogg component configuration.
   * @param name - The name of the remote state to get
   * @returns the DataTerraformRemoteState object
   * @throws if the remote state is not found
   */
  public remoteState(name: string): DataTerraformRemoteState {
    if (!this._remoteStates[name]) {
      throw new Error(`Remote state ${name} not found`);
    }
    return this._remoteStates[name];
  }

  /**
   * Set variables for a module included in the fogg component modules[] configuration.
   *
   * @param name - The module name as defined in the fogg component configuration
   * @param variables - The variables to set for the module
   */
  public setModuleVariables(
    name: string,
    variables: Record<string, any>
  ): void {
    if (!this.modules[name]) {
      throw new Error(`Module ${name} not found`);
    }
    for (const [key, value] of Object.entries(variables)) {
      this.modules[name].set(key, value);
    }
  }

  /**
   * Get a local defined in the fogg component configuration.
   *
   * @param name the name of the local to get
   * @returns the TerraformLocal object
   */
  public getLocal(name: string): TerraformLocal {
    if (!this.locals[name]) {
      throw new Error(`Local ${name} not found`);
    }
    return this.locals[name];
  }

  private parseBackendConfig(): void {
    if (this.foggComp.backend.kind === "s3" && this.foggComp.backend.s3) {
      const s3Config = this.foggComp.backend.s3;
      const s3BackendConfig: Mutable<S3BackendConfig> = {
        bucket: s3Config.bucket,
        dynamodbTable: s3Config.dynamo_table,
        key: s3Config.key_path,
        region: s3Config.region,
        encrypt: true,
      };
      if (s3Config.profile) {
        s3BackendConfig.profile = s3Config.profile;
      } else if (s3Config.role_arn) {
        s3BackendConfig.assumeRole = {
          roleArn: s3Config.role_arn,
        };
      }
      // console.log(
      //   `Setting S3 backend Config ${JSON.stringify(s3BackendConfig, null, 2)}`
      // );
      new S3Backend(this, s3BackendConfig);
    } else {
      throw new Error(
        `Unsupported backend configuration ${this.foggComp.backend.kind}`
      );
    }
  }

  private parseRemoteState(id: string, remoteConfig: Backend): void {
    if (remoteConfig.kind === "s3" && remoteConfig.s3) {
      const s3Config = remoteConfig.s3;
      const remoteStateConfig: Mutable<DataTerraformRemoteStateS3Config> = {
        bucket: s3Config.bucket,
        dynamodbTable: s3Config.dynamo_table,
        key: s3Config.key_path,
        region: s3Config.region,
        encrypt: true,
      };
      if (s3Config.profile) {
        remoteStateConfig.profile = s3Config.profile;
      } else if (s3Config.role_arn) {
        remoteStateConfig.assumeRole = {
          roleArn: s3Config.role_arn,
        };
      }
      this._remoteStates[id] = new DataTerraformRemoteStateS3(
        this,
        id,
        remoteStateConfig
      );
    } else {
      throw new Error(`Unsupported backend configuration ${remoteConfig.kind}`);
    }
  }

  private parseBundledProviderConfig(): void {
    const providers = this.foggComp.providers_configuration;
    if (providers.datadog) {
      this.parseDataDogProviderConfig(providers.datadog);
    }
  }

  /**
   * Instantiate supported generic providers.
   */
  private parseGenericProviderConfig(provider: GenericProvider): void {
    switch (provider.source) {
      case "cloudflare/cloudflare":
        this.parseCloudflareProviderConfig(provider);
        break;
      default:
        throw new Error(`Unsupported provider ${provider.source}`);
    }
  }

  private parseLocalsBlock() {
    if (this.foggComp.locals_block) {
      for (const [key, value] of Object.entries(this.foggComp.locals_block)) {
        this.locals[key] = new TerraformLocal(this, key, `\${${value}}`);
      }
    }
  }

  private parseModules() {
    if (this.foggComp.module_source) {
      const id = (this.foggComp.module_name =
        this.foggComp.module_name ?? "main");
      this.modules[id] = new TerraformHclModule(this, id, {
        source: this.foggComp.module_source,
      });
    }

    for (let i = 0; i < this.foggComp.modules.length; i++) {
      const moduleConfig = this.foggComp.modules[i];
      const id = moduleConfig.name ?? `module_${i}`;
      if (!moduleConfig.source) {
        console.warn(`Module ${id} does not have a source, skipping`);
        continue;
      }
      if (this.modules[id]) {
        throw new Error(`Module ${id} already exists`);
      }
      if (!moduleConfig.name) {
        console.log(
          `Module ${moduleConfig.source} does not have a name, using ${id}`
        );
      }
      this.modules[id] = new TerraformHclModule(this, id, {
        source: moduleConfig.source,
        version: moduleConfig.version,
      });
      // TODO: Add validation for module variables
      // TODO: Export module outputs
    }
  }

  private parseDataDogProviderConfig(_provider: DatadogProvider): void {
    // TODO: There's no datadog provider config?
    new dataDogProvider.DatadogProvider(this, "datadog", {});
  }

  private parseCloudflareProviderConfig(provider: GenericProvider): void {
    // The fogg provided provider.config is not `CloudflareProviderConfig` must use overrides
    const cf = new cloudflareProvider.CloudflareProvider(this, "cloudflare");
    for (const [key, value] of Object.entries(provider.config)) {
      cf.addOverride(key, value);
    }
  }
}

// Parse AWS Provider config from Fogg component configuration
function parseAwsProviderConfig(
  foggComp: Component
): awsProvider.AwsProviderConfig {
  const providers = foggComp.providers_configuration;
  if (
    providers.aws_regional_providers &&
    providers.aws_regional_providers.length > 0
  ) {
    throw new Error(
      "AWS regional providers are not supported by terraconstruct components"
    );
  }
  if (providers.aws) {
    return getAwsProviderConfig(foggComp, providers.aws);
  }
  throw new Error("AWS provider configuration not found");
}

function getAwsProviderConfig(
  foggComp: Component,
  config: AWSProvider
): awsProvider.AwsProviderConfig {
  const c: Mutable<awsProvider.AwsProviderConfig> = {
    region: config.region,
    alias: config.alias,
  };
  if (config.default_tags && config.default_tags.enabled) {
    c.defaultTags = [
      {
        tags: {
          env: foggComp.env,
          owner: foggComp.owner,
          project: foggComp.project,
          managedBy: "terraform",
          service: foggComp.name,
          ...(foggComp.backend.s3?.key_path && {
            tfstateKey: foggComp.backend.s3?.key_path,
          }),
          ...(foggComp.providers_configuration?.aws?.default_tags?.enabled &&
            (foggComp.providers_configuration?.aws?.default_tags?.tags ?? {})),
        },
      },
    ];
  }
  if (config.profile) {
    c.profile = config.profile;
  } else if (config.role_arn) {
    c.assumeRole = [
      {
        roleArn: config.role_arn,
      },
    ];
  }
  return c;
}

// helper type to make readonly interface properties mutable
type Mutable<T> = {
  -readonly [P in keyof T]: T[P];
};
