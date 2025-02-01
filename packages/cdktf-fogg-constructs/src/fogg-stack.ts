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
  TerraformStack,
} from "cdktf";
import { Construct, IConstruct } from "constructs";
import {
  AWSProvider,
  Backend,
  Component,
  DatadogProvider,
  GenericProvider,
} from "./imports/fogg-types.generated";
import { loadComponentConfig } from "./util/load-component-config";
import type { Mutable, OutputSchema } from "./util/types";

const DEFAULT_AWS_PROVIDER_ID = "DefaultAwsProvider";

export interface FoggStackProps {
  /**
   * Force remote state configuration
   * @default false - only configure remote states if `component_backends_filtered` is true
   */
  forceRemoteStates?: boolean;

  /**
   * Default tags to apply to resources on top of the default tags set by the Fogg component configuration
   *
   * @default - No additional tags
   */
  defaultTags?: Record<string, string>;
  /**
   * Tags to ignore on top of the ignoreTags set by the Fogg component configuration
   */
  ignoreTags?: awsProvider.AwsProviderIgnoreTags;
}

export interface IFoggStack extends IConstruct {
  /**
   * The Raw Fogg component configuration
   */
  foggComp: Component;
  /**
   * A Map of module name to TerraformHclModule Constructs parsed from the fogg component configuration
   */
  modules: Record<string, TerraformHclModule>;
  /**
   * A Map of local name to TerraformLocal Constructs parsed from the fogg component configuration
   */
  locals: Record<string, TerraformLocal>;
  /**
   * The AWS region defined in the aws providers configuration (if provided)
   * or undefined if not set
   */
  region?: string;

  /**
   * The default AWS provider
   */
  defaultAwsProvider: awsProvider.AwsProvider;

  /**
   * Set variables for the main module defined in the fogg component configuration.
   *
   * @param variables - The variables to set for the module
   */
  setMainModuleVariables(variables: Record<string, any>): void;
  /**
   * Return a remote state defined in the fogg component configuration.
   * @param name - The name of the remote state to get
   * @returns the DataTerraformRemoteState object
   * @throws if the remote state is not found
   */
  remoteState<T>(name: string, schema?: OutputSchema<T>): T;
  /**
   * Set variables for a module included in the fogg component modules[] configuration.
   *
   * @param name - The module name as defined in the fogg component configuration
   * @param variables - The variables to set for the module
   */
  setModuleVariables(name: string, variables: Record<string, any>): void;
  /**
   * Get a local defined in the fogg component configuration.
   *
   * @param name the name of the local to get
   * @returns the TerraformLocal object
   */
  getLocal(name: string): TerraformLocal;
}

/**
 * Helper stack to wrap Fogg component configuration and set up configured providers and backends.
 */
export class FoggStack extends TerraformStack implements IFoggStack {
  public readonly foggComp: Component;
  public readonly modules: Record<string, TerraformHclModule> = {};
  public readonly locals: Record<string, TerraformLocal> = {};
  public readonly region?: string;
  private readonly _providers: Record<string, awsProvider.AwsProvider> = {};
  private readonly _remoteStates: Record<string, DataTerraformRemoteState> = {};
  private readonly _defaultTags?: Record<string, string>;
  private readonly _ignoreTags?: awsProvider.AwsProviderIgnoreTags;

  constructor(scope: Construct, id: string, props: FoggStackProps = {}) {
    super(scope, id);
    this.foggComp = loadComponentConfig();
    this._defaultTags = props.defaultTags;
    this._ignoreTags = props.ignoreTags;
    this.region = this.foggComp.providers_configuration.aws?.region;

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

  public setMainModuleVariables(variables: Record<string, any>): void {
    this.foggComp.module_name = this.foggComp.module_name ?? "main";
    this.setModuleVariables(this.foggComp.module_name, variables);
  }

  public remoteState<T>(name: string, schema?: OutputSchema<T>): T {
    if (!this._remoteStates[name]) {
      throw new Error(`Remote state ${name} not found`);
    }
    return new Proxy(this._remoteStates[name], {
      get(target, prop: string | symbol, _receiver) {
        if (typeof prop !== "string") return undefined;
        if (!schema) {
          return target.getString(prop);
        }
        const type = schema[prop as keyof T];
        switch (type) {
          case "string":
            return target.getString(prop);
          case "list":
            return target.getList(prop);
          case "number":
            return target.getNumber(prop);
          case "boolean":
            return target.getBoolean(prop);
          default:
            return target.get(prop);
        }
      },
    }) as unknown as T;
  }

  /**
   * Return a provider defined in the fogg component configuration.
   * @param alias - The alias of the provider to get
   * @returns the AwsProvider object
   * @throws if the provider is not found
   */
  public awsProvider(alias: string): awsProvider.AwsProvider {
    if (!this._providers[alias]) {
      throw new Error(`Provider ${alias} not found`);
    }
    return this._providers[alias];
  }

  public get defaultAwsProvider(): awsProvider.AwsProvider {
    if (!this._providers[DEFAULT_AWS_PROVIDER_ID]) {
      throw new Error(`Default AWS Provider not defined`);
    }
    return this.awsProvider(DEFAULT_AWS_PROVIDER_ID);
  }

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
    if (providers.aws) {
      this.parseAwsProviderConfig(providers.aws);
    }
    if (providers.aws_regional_providers) {
      for (let i = 0; i < providers.aws_regional_providers.length; i++) {
        this.parseAwsProviderConfig(
          providers.aws_regional_providers[i],
          `aws-${i}`
        );
      }
    }
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

  private parseAwsProviderConfig(
    config: AWSProvider,
    id: string = DEFAULT_AWS_PROVIDER_ID
  ): void {
    const c: Mutable<awsProvider.AwsProviderConfig> = {
      region: config.region,
      alias: config.alias,
    };
    if (config.default_tags && config.default_tags.enabled) {
      c.defaultTags = [
        {
          tags: {
            env: this.foggComp.env,
            owner: this.foggComp.owner,
            project: this.foggComp.project,
            managedBy: "terraform",
            service: this.foggComp.name,
            ...(this.foggComp.backend.s3?.key_path && {
              tfstateKey: this.foggComp.backend.s3?.key_path,
            }),
            ...(this.foggComp.providers_configuration?.aws?.default_tags
              ?.enabled &&
              (this.foggComp.providers_configuration?.aws?.default_tags?.tags ??
                {})),
            ...(this._defaultTags ?? {}),
          },
        },
      ];
    }
    let ignoreTags: awsProvider.AwsProviderIgnoreTags | undefined;
    if (config.ignore_tags && config.ignore_tags.enabled) {
      ignoreTags = {
        keys: config.ignore_tags.keys,
        keyPrefixes: config.ignore_tags.key_prefixes,
      };
    }
    if (this._ignoreTags) {
      const keys = new Set([
        ...(ignoreTags?.keys ?? []),
        ...(this._ignoreTags.keys ?? []),
      ]);
      const keyPrefixes = new Set([
        ...(ignoreTags?.keyPrefixes ?? []),
        ...(this._ignoreTags.keyPrefixes ?? []),
      ]);
      ignoreTags = {
        keys: keys.size > 0 ? Array.from(keys) : undefined,
        keyPrefixes: keyPrefixes.size > 0 ? Array.from(keyPrefixes) : undefined,
      };
    }
    if (ignoreTags) {
      c.ignoreTags = [ignoreTags];
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
    this._providers[c.alias ?? DEFAULT_AWS_PROVIDER_ID] =
      new awsProvider.AwsProvider(this, id, c);
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
