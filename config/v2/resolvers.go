package v2

import (
	"encoding/json"

	"github.com/chanzuckerberg/fogg/util"
)

// lastNonNil, despite its name can return nil if all results are nil
func lastNonNil(getter func(Common) *string, commons ...Common) *string {
	var s *string
	for _, c := range commons {
		t := getter(c)
		if t != nil {
			s = t
		}
	}
	return s
}

// lastNonNilInt, despite its name can return nil if all results are nil
func lastNonNilInt64(getter func(Common) *int64, commons ...Common) *int64 {
	var s *int64
	for _, c := range commons {
		t := getter(c)
		if t != nil {
			s = t
		}
	}
	return s
}

// lastNonNilStringSlice, despite its name can return nil if all results are nil
func lastNonNilStringSlice(getter func(Common) []string, commons ...Common) []string {
	var s []string
	for _, c := range commons {
		t := getter(c)
		if t != nil {
			s = t
		}
	}
	return s
}

// ResolveRequiredString will resolve the value and panic if it is nil. Only to be used after validations are run.
func ResolveRequiredString(getter func(Common) *string, commons ...Common) string {
	return *lastNonNil(getter, commons...)
}

// ResolveRequiredInt will resolve the value and panic if it is nil. Only to be used after validations are run.
func ResolveRequiredInt64(getter func(Common) *int64, commons ...Common) int64 {
	return *lastNonNilInt64(getter, commons...)
}

func ResolveOptionalString(getter func(Common) *string, commons ...Common) *string {
	return lastNonNil(getter, commons...)
}

func ResolveOptionalStringSlice(getter func(Common) []string, commons ...Common) []string {
	return lastNonNilStringSlice(getter, commons...)
}

func ResolveStringArray(def []string, override []string) []string {
	if override != nil {
		return override
	}
	return def
}

func ResolveStringMap(getter func(Common) map[string]string, commons ...Common) map[string]string {
	resolved := map[string]string{}

	for _, c := range commons {
		m := getter(c)
		for k, v := range m {
			resolved[k] = v
		}
	}
	return resolved
}

func ResolveAuth0Provider(commons ...Common) *Auth0Provider {
	var domain, version, source *string
	for _, c := range commons {
		if c.Providers == nil || c.Providers.Auth0 == nil {
			continue
		}
		if c.Providers.Auth0.Version != nil {
			version = c.Providers.Auth0.Version
		}

		if c.Providers.Auth0.Domain != nil {
			domain = c.Providers.Auth0.Domain
		}

		if c.Providers.Auth0.Source != nil {
			source = c.Providers.Auth0.Source
		}
	}

	if domain != nil && version != nil {
		return &Auth0Provider{Version: version, Domain: domain, Source: source}
	}
	return nil
}

func ResolveAssertProvider(commons ...Common) *AssertProvider {
	var version *string
	for _, c := range commons {
		if c.Providers == nil || c.Providers.Assert == nil {
			continue
		}
		if c.Providers.Assert.Version != nil {
			version = c.Providers.Assert.Version
		}
	}

	if version != nil {
		return &AssertProvider{Version: version}
	}
	return nil
}

// ResolveAWSProvider will return an AWSProvder if one of the required fields is set somewhere in
// the set of Common config objects passed in. Otherwise it will return nil.
func ResolveAWSProvider(commons ...Common) *AWSProvider {
	var profile, region, role, version *string
	var accountID *json.Number
	var additionalRegions []string
	var additionalProviders map[string]*AWSProvider

	for _, c := range commons {
		if c.Providers != nil && c.Providers.AWS != nil {
			p := c.Providers.AWS

			// Profile and Role are mutually exclusive, so if one is set then we set the other to
			// nil Our validations in validateAWSProviderAuth will assure that they are not
			// both set in the same stuct.
			if p.Profile != nil {
				profile = p.Profile
				role = nil
			} else if p.Role != nil {
				role = p.Role
				profile = nil
			}

			if p.Region != nil {
				region = p.Region
			}

			if p.Version != nil {
				version = p.Version
			}

			if p.AccountID != nil {
				accountID = p.AccountID
			}

			if p.AdditionalRegions != nil {
				additionalRegions = p.AdditionalRegions
			}

			if p.AdditionalProviders != nil {
				additionalProviders = p.AdditionalProviders
			}
		}
	}

	if profile != nil || role != nil || region != nil || version != nil {
		return &AWSProvider{
			Profile: profile,
			Region:  region,
			Role:    role,
			Version: version,

			// optional fields
			AccountID:           accountID,
			AdditionalRegions:   additionalRegions,
			AdditionalProviders: additionalProviders,
		}
	}
	return nil
}

// ResolveBackend returns the Backend configuration for a given component, after applying all inheritance rules
func ResolveBackend(commons ...Common) *Backend {
	var ret *Backend
	for _, c := range commons {
		if c.Backend != nil {
			if ret == nil {
				ret = &Backend{Kind: util.Ptr("s3")}
			}
			b := c.Backend
			if b.Kind != nil {
				ret.Kind = b.Kind
			}

			if b.AccountID != nil {
				ret.AccountID = b.AccountID
			}

			if b.Bucket != nil {
				ret.Bucket = b.Bucket
			}

			if b.DynamoTable != nil {
				ret.DynamoTable = b.DynamoTable
			}

			if b.Region != nil {
				ret.Region = b.Region
			}

			// Profile and Role are mutually exclusive, so if one is set then we set the other to
			// nil Our validations in validateBackend will assure that they are not both set or missing in the
			// same stuct.
			if b.Profile != nil {
				ret.Profile = b.Profile
				ret.Role = nil
			} else if b.Role != nil {
				ret.Role = b.Role
				ret.Profile = nil
			}

			if b.HostName != nil {
				ret.HostName = b.HostName
			}

			if b.Organization != nil {
				ret.Organization = b.Organization
			}
		}
	}

	return ret
}

// ResolveGithubProvider will return an GithubProvder iff one of the required fields is set somewhere in the set of Common
// config objects passed in. Otherwise it will return nil.
func ResolveGithubProvider(commons ...Common) *GithubProvider {
	org := lastNonNil(GithubProviderOrganizationGetter, commons...)

	if org == nil {
		return nil
	}

	return &GithubProvider{
		Organization: org,

		// optional fields
		BaseURL: lastNonNil(GithubProviderBaseURLGetter, commons...),
		Version: lastNonNil(GithubProviderVersionGetter, commons...),
	}
}

func ResolveSnowflakeProvider(commons ...Common) *SnowflakeProvider {
	account := lastNonNil(SnowflakeProviderAccountGetter, commons...)
	role := lastNonNil(SnowflakeProviderRoleGetter, commons...)
	region := lastNonNil(SnowflakeProviderRegionGetter, commons...)
	version := lastNonNil(SnowflakeProviderVersionGetter, commons...)

	if account != nil || role != nil || region != nil {
		return &SnowflakeProvider{
			Account: account,
			Role:    role,
			Region:  region,
			Version: version,
		}
	}
	return nil
}

func ResolveOktaProvider(commons ...Common) *OktaProvider {
	orgName := lastNonNil(OktaProviderOrgNameGetter, commons...)
	baseURL := lastNonNil(OktaProviderBaseURLGetter, commons...)
	registryNamespace := lastNonNil(OktaProviderRegistryNamespaceGetter, commons...)

	// required fields
	if orgName == nil {
		return nil
	}

	return &OktaProvider{
		OrgName:           orgName,
		Version:           lastNonNil(OktaProviderVersionGetter, commons...),
		BaseURL:           baseURL,
		RegistryNamespace: registryNamespace,
	}
}

func ResolveBlessProvider(commons ...Common) *BlessProvider {
	profile := lastNonNil(BlessProviderProfileGetter, commons...)
	roleArn := lastNonNil(BlessProviderRoleArnGetter, commons...)
	region := lastNonNil(BlessProviderRegionGetter, commons...)

	// required fields
	if (profile == nil && roleArn == nil) || region == nil {
		return nil
	}

	return &BlessProvider{
		AWSProfile: profile,
		AWSRegion:  region,
		RoleArn:    roleArn,

		Version:           lastNonNil(BlessProviderVersionGetter, commons...),
		AdditionalRegions: ResolveOptionalStringSlice(BlessProviderAdditionalRegionsGetter, commons...),
	}
}

func ResolveHerokuProvider(commons ...Common) *HerokuProvider {
	var p *HerokuProvider
	for _, c := range commons {
		if c.Providers == nil || c.Providers.Heroku == nil {
			continue
		}
		p = c.Providers.Heroku
	}

	version := lastNonNil(HerokuProviderVersionGetter, commons...)

	if version != nil {
		return &HerokuProvider{
			Version: version,
		}
	}
	return p
}

func ResolveDatadogProvider(commons ...Common) *DatadogProvider {
	var p *DatadogProvider
	for _, c := range commons {
		if c.Providers == nil || c.Providers.Datadog == nil {
			continue
		}
		p = c.Providers.Datadog
	}

	version := lastNonNil(DatadogProviderVersionGetter, commons...)

	if version != nil {
		return &DatadogProvider{
			Version: version,
		}
	}
	return p
}

func ResolvePagerdutyProvider(commons ...Common) *PagerdutyProvider {
	var p *PagerdutyProvider
	for _, c := range commons {
		if c.Providers == nil || c.Providers.Pagerduty == nil {
			continue
		}
		p = c.Providers.Pagerduty
	}

	version := lastNonNil(PagerdutyProviderVersionGetter, commons...)

	if version != nil {
		return &PagerdutyProvider{
			Version: version,
		}
	}
	return p
}

func ResolveOpsGenieProvider(commons ...Common) *OpsGenieProvider {
	var p *OpsGenieProvider
	for _, c := range commons {
		if c.Providers == nil || c.Providers.OpsGenie == nil {
			continue
		}
		p = c.Providers.OpsGenie
	}

	version := lastNonNil(OpsGenieProviderVersionGetter, commons...)

	if version != nil {
		return &OpsGenieProvider{
			Version: version,
		}
	}
	return p
}

func ResolveDatabricksProvider(commons ...Common) *DatabricksProvider {
	var p *DatabricksProvider
	for _, c := range commons {
		if c.Providers == nil || c.Providers.Databricks == nil {
			continue
		}
		p = c.Providers.Databricks
	}

	version := lastNonNil(DatabricksProviderVersionGetter, commons...)

	if version != nil {
		return &DatabricksProvider{
			Version: version,
		}
	}
	return p
}

func ResolveSentryProvider(commons ...Common) *SentryProvider {
	var p *SentryProvider
	for _, c := range commons {
		if c.Providers == nil || c.Providers.Datadog == nil {
			continue
		}
		p = c.Providers.Sentry
	}

	version := lastNonNil(SentryProviderVersionGetter, commons...)
	baseURL := lastNonNil(SentryProviderBaseURLGetter, commons...)

	if version != nil {
		return &SentryProvider{
			Version: version,
			BaseURL: baseURL,
		}
	}
	return p
}

func ResolveTfeProvider(commons ...Common) *TfeProvider {
	var version *string
	var enabled *bool
	var hostname *string

	for _, c := range commons {
		if c.Providers != nil && c.Providers.Tfe != nil {
			t := c.Providers.Tfe

			if t.Enabled != nil {
				enabled = t.Enabled
			}

			if t.Version != nil {
				version = t.Version
			}

			if t.Hostname != nil {
				hostname = t.Hostname
			}
		}
	}

	return &TfeProvider{
		CommonProvider: CommonProvider{
			Enabled: enabled,
			Version: version,
		},
		Hostname: hostname,
	}
}

func ResolveSopsProvider(commons ...Common) *SopsProvider {
	var version *string
	var enabled *bool

	for _, c := range commons {
		if c.Providers != nil && c.Providers.Sops != nil {
			t := c.Providers.Sops

			if t.Enabled != nil {
				enabled = t.Enabled
			}

			if t.Version != nil {
				version = t.Version
			}
		}
	}

	return &SopsProvider{
		CommonProvider: CommonProvider{
			Enabled: enabled,
			Version: version,
		},
	}
}

func ResolveKubernetesProvider(commons ...Common) *KubernetesProvider {
	var version *string
	var enabled *bool

	for _, c := range commons {
		if c.Providers != nil && c.Providers.Kubernetes != nil {
			t := c.Providers.Kubernetes

			if t.Enabled != nil {
				enabled = t.Enabled
			}

			if t.Version != nil {
				version = t.Version
			}
		}
	}

	return &KubernetesProvider{
		CommonProvider: CommonProvider{
			Enabled: enabled,
			Version: version,
		},
	}
}

func ResolveGrafanaProvider(commons ...Common) *GrafanaProvider {
	var version *string
	var enabled *bool

	for _, c := range commons {
		if c.Providers != nil && c.Providers.Grafana != nil {
			t := c.Providers.Grafana

			if t.Enabled != nil {
				enabled = t.Enabled
			}

			if t.Version != nil {
				version = t.Version
			}
		}
	}

	return &GrafanaProvider{
		CommonProvider: CommonProvider{
			Enabled: enabled,
			Version: version,
		},
	}
}

func ResolveTfLint(commons ...Common) TfLint {
	enabled := false
	for _, c := range commons {
		if c.Tools != nil && c.Tools.TfLint != nil && c.Tools.TfLint.Enabled != nil {
			enabled = *c.Tools.TfLint.Enabled
		}
	}

	return TfLint{
		Enabled: &enabled,
	}
}

func ResolveTravis(commons ...Common) *TravisCI {
	enabled := false
	buildevents := false
	testCommand := "check"

	for _, c := range commons {
		if c.Tools != nil && c.Tools.TravisCI != nil && c.Tools.TravisCI.Enabled != nil {
			enabled = *c.Tools.TravisCI.Enabled
		}

		if c.Tools != nil && c.Tools.TravisCI != nil && c.Tools.TravisCI.Command != nil {
			testCommand = *c.Tools.TravisCI.Command
		}
		if c.Tools != nil && c.Tools.TravisCI != nil && c.Tools.TravisCI.Buildevents != nil {
			buildevents = *c.Tools.TravisCI.Buildevents
		}
	}

	roleName := lastNonNil(TravisRoleNameGetter, commons...)

	return &TravisCI{
		CommonCI: CommonCI{
			Enabled:        &enabled,
			Buildevents:    &buildevents,
			AWSIAMRoleName: roleName,
			Command:        &testCommand,
		},
	}
}

func ResolveGitHubActionsCI(commons ...Common) *GitHubActionsCI {
	enabled := false
	buildevents := false
	testCommand := "check"

	for _, c := range commons {
		if c.Tools != nil && c.Tools.GitHubActionsCI != nil {
			if c.Tools.GitHubActionsCI.Enabled != nil {
				enabled = *c.Tools.GitHubActionsCI.Enabled
			}
			if c.Tools.GitHubActionsCI.Command != nil {
				testCommand = *c.Tools.GitHubActionsCI.Command
			}
			if c.Tools.GitHubActionsCI.Buildevents != nil {
				buildevents = *c.Tools.GitHubActionsCI.Buildevents
			}
		}
	}

	roleName := lastNonNil(GitHubActionsRoleNameGetter, commons...)
	region := lastNonNil(GitHubActionsRegionGetter, commons...)
	return &GitHubActionsCI{
		CommonCI: CommonCI{
			Enabled:        &enabled,
			Buildevents:    &buildevents,
			AWSIAMRoleName: roleName,
			AWSRegion:      region,
			Command:        &testCommand,
		},
	}
}

func ResolveCircleCI(commons ...Common) *CircleCI {
	enabled := false
	buildevents := false
	testCommand := "check"
	var providers map[string]CIProviderConfig

	for _, c := range commons {
		if c.Tools != nil && c.Tools.CircleCI != nil && c.Tools.CircleCI.Enabled != nil {
			enabled = *c.Tools.CircleCI.Enabled
		}

		if c.Tools != nil && c.Tools.CircleCI != nil && c.Tools.CircleCI.Command != nil {
			testCommand = *c.Tools.CircleCI.Command
		}

		if c.Tools != nil && c.Tools.CircleCI != nil && c.Tools.CircleCI.Buildevents != nil {
			buildevents = *c.Tools.CircleCI.Buildevents
		}

		if c.Tools != nil && c.Tools.CircleCI != nil {
			providers = c.Tools.CircleCI.Providers
		}
	}

	sshFingerprints := ResolveOptionalStringSlice(CircleCISSHFingerprintsGetter, commons...)
	roleName := lastNonNil(CircleCIRoleNameGetter, commons...)

	return &CircleCI{
		CommonCI: CommonCI{
			Enabled:        &enabled,
			Buildevents:    &buildevents,
			AWSIAMRoleName: roleName,
			Command:        &testCommand,
			Providers:      providers,
		},
		SSHKeyFingerprints: sshFingerprints,
	}
}

func OwnerGetter(comm Common) *string {
	return comm.Owner
}

func ProjectGetter(comm Common) *string {
	return comm.Project
}

func TerraformVersionGetter(comm Common) *string {
	return comm.TerraformVersion
}

func BackendBucketGetter(comm Common) *string {
	if comm.Backend != nil {
		return comm.Backend.Bucket
	}
	return nil
}

func BackendRegionGetter(comm Common) *string {
	if comm.Backend != nil {
		return comm.Backend.Region
	}
	return nil
}

func BackendDynamoTableGetter(comm Common) *string {
	if comm.Backend != nil {
		return comm.Backend.DynamoTable
	}
	return nil
}

func BackendProfileGetter(comm Common) *string {
	if comm.Backend != nil {
		return comm.Backend.Profile
	}
	return nil
}

// BackendKindGetter retrieves the Kind for the current common object
func BackendKindGetter(comm Common) *string {
	if comm.Backend == nil {
		return nil
	}
	return comm.Backend.Kind
}

func BackendAccountIDGetter(comm Common) *string {
	if comm.Backend == nil {
		return nil
	}

	return comm.Backend.AccountID
}

func BackendHostNameGetter(comm Common) *string {
	if comm.Backend == nil {
		return nil
	}

	return comm.Backend.HostName
}

func BackendOrganizationGetter(comm Common) *string {
	if comm.Backend == nil {
		return nil
	}

	return comm.Backend.Organization
}

func GithubProviderOrganizationGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.Github != nil {
		return comm.Providers.Github.Organization
	}
	return nil
}

func GithubProviderBaseURLGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.Github != nil {
		return comm.Providers.Github.BaseURL
	}
	return nil
}

func GithubProviderVersionGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.Github != nil {
		return comm.Providers.Github.Version
	}
	return nil
}

func ExtraVarsGetter(comm Common) map[string]string {
	if comm.ExtraVars != nil {
		return comm.ExtraVars
	}
	return map[string]string{}
}

func ResolveModuleTerraformVersion(def Defaults, module Module) *string {
	if module.TerraformVersion != nil {
		return module.TerraformVersion
	}
	return def.TerraformVersion
}

func SnowflakeProviderAccountGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.Snowflake != nil {
		return comm.Providers.Snowflake.Account
	}
	return nil
}
func SnowflakeProviderRoleGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.Snowflake != nil {
		return comm.Providers.Snowflake.Role
	}
	return nil
}
func SnowflakeProviderRegionGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.Snowflake != nil {
		return comm.Providers.Snowflake.Region
	}
	return nil
}
func SnowflakeProviderVersionGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.Snowflake != nil {
		return comm.Providers.Snowflake.Version
	}
	return nil
}

func BlessProviderProfileGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Bless == nil {
		return nil
	}
	return comm.Providers.Bless.AWSProfile
}

func BlessProviderRoleArnGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Bless == nil {
		return nil
	}
	return comm.Providers.Bless.RoleArn
}

func BlessProviderRegionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Bless == nil {
		return nil
	}
	return comm.Providers.Bless.AWSRegion
}
func BlessProviderVersionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Bless == nil {
		return nil
	}
	return comm.Providers.Bless.Version
}
func BlessProviderAdditionalRegionsGetter(comm Common) []string {
	if comm.Providers == nil || comm.Providers.Bless == nil {
		return nil
	}
	return comm.Providers.Bless.AdditionalRegions
}

func HerokuProviderVersionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Heroku == nil {
		return nil
	}
	return comm.Providers.Heroku.Version
}

func DatadogProviderVersionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Datadog == nil {
		return nil
	}
	return comm.Providers.Datadog.Version
}

func PagerdutyProviderVersionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Pagerduty == nil {
		return nil
	}
	return comm.Providers.Pagerduty.Version
}

func OpsGenieProviderVersionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.OpsGenie == nil {
		return nil
	}
	return comm.Providers.OpsGenie.Version
}

func DatabricksProviderVersionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Databricks == nil {
		return nil
	}
	return comm.Providers.Databricks.Version
}

func SentryProviderVersionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Sentry == nil {
		return nil
	}
	return comm.Providers.Sentry.Version
}

func SentryProviderBaseURLGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Sentry == nil {
		return nil
	}
	return comm.Providers.Sentry.BaseURL
}

func OktaProviderVersionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Okta == nil {
		return nil
	}
	return comm.Providers.Okta.Version
}

func OktaProviderBaseURLGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Okta == nil {
		return nil
	}
	return comm.Providers.Okta.BaseURL
}

func OktaProviderRegistryNamespaceGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Okta == nil {
		return nil
	}
	return comm.Providers.Okta.RegistryNamespace
}

func OktaProviderRegistryNamespacegetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Okta == nil {
		return nil
	}
	return comm.Providers.Okta.RegistryNamespace
}

func OktaProviderOrgNameGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Okta == nil {
		return nil
	}
	return comm.Providers.Okta.OrgName
}

func TravisRoleNameGetter(comm Common) *string {
	if comm.Tools == nil || comm.Tools.TravisCI == nil {
		return nil
	}
	return comm.Tools.TravisCI.AWSIAMRoleName
}

func GitHubActionsRoleNameGetter(comm Common) *string {
	if comm.Tools == nil || comm.Tools.GitHubActionsCI == nil {
		return nil
	}
	return comm.Tools.GitHubActionsCI.AWSIAMRoleName
}

func GitHubActionsRegionGetter(comm Common) *string {
	if comm.Tools == nil || comm.Tools.GitHubActionsCI == nil {
		return nil
	}
	return comm.Tools.GitHubActionsCI.AWSRegion
}

func CircleCIRoleNameGetter(comm Common) *string {
	if comm.Tools == nil || comm.Tools.CircleCI == nil {
		return nil
	}
	return comm.Tools.CircleCI.AWSIAMRoleName
}

func CircleCISSHFingerprintsGetter(comm Common) []string {
	if comm.Tools == nil || comm.Tools.CircleCI == nil {
		return nil
	}
	return comm.Tools.CircleCI.SSHKeyFingerprints
}

func DependsOnAccountsGetter(comm Common) []string {
	if comm.DependsOn == nil {
		return nil
	}
	return comm.DependsOn.Accounts
}

func DependsOnComponentsGetter(comm Common) []string {
	if comm.DependsOn == nil {
		return nil
	}
	return comm.DependsOn.Components
}
