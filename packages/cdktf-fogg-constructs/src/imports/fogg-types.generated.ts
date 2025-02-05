/* Do not change, this code is generated from Golang structs */


export interface IntegrationRegistryMap {
    format?: string;
    drop_component?: boolean;
    path?: string;
    for_each?: boolean;
    path_for_each?: string;
}
export interface ModuleIntegrationConfig {
    mode?: string;
    format?: string;
    drop_prefix?: boolean;
    drop_component?: boolean;
    path_infix?: string;
    providers?: string[];
    outputs_map?: {[key: string]: IntegrationRegistryMap};
}
export interface ComponentModule {
    source?: string;
    version?: string;
    name?: string;
    prefix?: string;
    variables?: string[];
    outputs?: string[];
    integration?: ModuleIntegrationConfig;
    providers?: {[key: string]: string};
    for_each?: string;
    depends_on?: string[];
}
export interface EKSConfig {
    cluster_name: string;
}
export interface GitHubActionsComponent {
    Enabled: boolean;
    Buildevents: boolean;
    AWSProfileName: string;
    AWSRoleName: string;
    AWSRegion: string;
    AWSAccountID: string;
    Command: string;
}
export interface CircleCIComponent {
    Enabled: boolean;
    Buildevents: boolean;
    AWSProfileName: string;
    AWSRoleName: string;
    AWSRegion: string;
    AWSAccountID: string;
    Command: string;
    SSHFingerprints: string[];
}
export interface TravisCIComponent {
    Enabled: boolean;
    Buildevents: boolean;
    AWSProfileName: string;
    AWSRoleName: string;
    AWSRegion: string;
    AWSAccountID: string;
    Command: string;
}
export interface TfLint {
    enabled: boolean;
}
export interface ProviderVersion {
    source: string;
    version?: string;
}
export interface GenericProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
    source: string;
    config: {[key: string]: any};
}
export interface SopsProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
}
export interface TfeProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
    hostname?: string;
}
export interface SnowflakeProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
    account?: string;
    role?: string;
    region?: string;
}
export interface SentryProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
    base_url?: string;
}
export interface OktaProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
    org_name?: string;
    base_url?: string;
}
export interface KubernetesProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
}
export interface HerokuProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
}
export interface GrafanaProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
}
export interface GithubProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
    organization: string;
    base_url?: string;
}
export interface DatadogProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
}
export interface BlessProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
    additional_regions?: string[];
    aws_profile?: string;
    aws_region?: string;
    role_arn?: string;
}
export interface AWSProviderIgnoreTags {
    enabled?: boolean;
    keys?: string[];
    key_prefixes?: string[];
}
export interface AWSProviderDefaultTags {
    enabled?: boolean;
    tags?: {[key: string]: string};
}
export interface AWSProvider {
    account_id: string;
    alias?: string;
    profile?: string;
    region: string;
    role_arn?: string;
    default_tags?: AWSProviderDefaultTags;
    ignore_tags?: AWSProviderIgnoreTags;
}
export interface Auth0Provider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
    domain?: string;
}
export interface AssertProvider {
    custom_provider?: boolean;
    enabled?: boolean;
    version?: string;
}
export interface ProviderConfiguration {
    assert?: AssertProvider;
    auth0?: Auth0Provider;
    aws?: AWSProvider;
    aws_regional_providers: AWSProvider[];
    bless?: BlessProvider;
    datadog?: DatadogProvider;
    github?: GithubProvider;
    grafana?: GrafanaProvider;
    heroku?: HerokuProvider;
    kubernetes?: KubernetesProvider;
    okta?: OktaProvider;
    sentry?: SentryProvider;
    snowflake?: SnowflakeProvider;
    tfe?: TfeProvider;
    sops?: SopsProvider;
}
export interface Number {

}
export interface RemoteBackend {
    host_name: string;
    organization: string;
    workspace: string;
}
export interface S3Backend {
    account_id?: string;
    account_name: string;
    bucket: string;
    dynamo_table?: string;
    key_path: string;
    profile?: string;
    region: string;
    role_arn?: string;
}
export interface Backend {
    kind: string;
    s3?: S3Backend;
    remote?: RemoteBackend;
}
export interface Component {
    path_to_repo_root: string;
    terraform_version: string;
    account_backends: {[key: string]: Backend};
    all_accounts: {[key: string]: Number};
    backend: Backend;
    component_backends: {[key: string]: Backend};
    autoplan_relative_globs: string[];
    autoplan_files: string[];
    locals_block: {[key: string]: any};
    component_backends_filtered: boolean;
    env: string;
    extra_vars: {[key: string]: string};
    name: string;
    owner: string;
    project: string;
    providers_configuration: ProviderConfiguration;
    required_providers: {[key: string]: GenericProvider};
    provider_versions: {[key: string]: ProviderVersion};
    integration_registry?: string;
    cdktf_dependencies: {[key: string]: string};
    cdktf_dev_dependencies: {[key: string]: string};
    package_fields: {[key: string]: any};
    tf_lint: TfLint;
    TravisCI: TravisCIComponent;
    CircleCI: CircleCIComponent;
    GitHubActionsCI: GitHubActionsComponent;
    eks?: EKSConfig;
    kind?: string;
    module_source?: string;
    module_name?: string;
    module_for_each?: string;
    providers: {[key: string]: string};
    variables: string[];
    outputs: string[];
    modules: ComponentModule[];
    global?: Component;
}