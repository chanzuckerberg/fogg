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

// lastNonNilJsonNumber, despite its name can return nil if all results are nil
func lastNonNilJsonNumber(getter func(Common) *json.Number, commons ...Common) *json.Number {
	var jsonNumber *json.Number
	for _, c := range commons {
		j := getter(c)
		if j != nil {
			jsonNumber = j
		}
	}
	return jsonNumber
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

// ResolveAWSProvider will return an AWSProvder iff one of the required fields is set somewhere in the set of Common
// config objects passed in. Otherwise it will return nil.
func ResolveAWSProvider(commons ...Common) *AWSProvider {

	// we may in the future want invert this implementation and walk the structs first
	profile := lastNonNil(AWSProviderProfileGetter, commons...)
	region := lastNonNil(AWSProviderRegionGetter, commons...)
	version := lastNonNil(AWSProviderVersionGetter, commons...)

	if profile != nil || region != nil || version != nil {
		return &AWSProvider{
			Profile: profile,
			Region:  region,
			Version: version,

			// optional fields
			AccountID:         lastNonNilJsonNumber(AWSProviderAccountIdGetter, commons...),
			AdditionalRegions: ResolveOptionalStringSlice(AWSProviderAdditionalRegionsGetter, commons...),
		}
	}
	return nil
}

// ResolveBackend returns the Backend configuration for a given component, after applying all inheritance rules
func ResolveBackend(commons ...Common) *Backend {
	kind := lastNonNil(BackendKindGetter, commons...)

	// This feels like a somewhat hacky way to do this, but not sure of a better place yet
	if kind == nil {
		kind = util.StrPtr("s3")
	}

	accountID := lastNonNil(BackendAccountIDGetter, commons...)
	bucket := lastNonNil(BackendBucketGetter, commons...)
	dynamoTable := lastNonNil(BackendDynamoTableGetter, commons...)
	profile := lastNonNil(BackendProfileGetter, commons...)
	region := lastNonNil(BackendRegionGetter, commons...)
	hostName := lastNonNil(BackendHostNameGetter, commons...)
	organization := lastNonNil(BackendOrganizationGetter, commons...)

	return &Backend{
		Kind: kind,

		AccountID:    accountID,
		Bucket:       bucket,
		DynamoTable:  dynamoTable,
		Profile:      profile,
		Region:       region,
		HostName:     hostName,
		Organization: organization,
	}
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

	// required fields
	if orgName == nil {
		return nil
	}

	return &OktaProvider{
		OrgName: orgName,
		Version: lastNonNil(OktaProviderVersionGetter, commons...),
	}
}

func ResolveBlessProvider(commons ...Common) *BlessProvider {
	profile := lastNonNil(BlessProviderProfileGetter, commons...)
	region := lastNonNil(BlessProviderRegionGetter, commons...)

	// required fields
	if profile == nil || region == nil {
		return nil
	}

	return &BlessProvider{
		AWSProfile: profile,
		AWSRegion:  region,

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
	} else {
		return p
	}
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
	} else {
		return p
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
		if c.Tools != nil && c.Tools.GitHubActionsCI != nil && c.Tools.GitHubActionsCI.Enabled != nil {
			enabled = *c.Tools.GitHubActionsCI.Enabled
		}
		if c.Tools != nil && c.Tools.GitHubActionsCI != nil && c.Tools.GitHubActionsCI.Command != nil {
			testCommand = *c.Tools.GitHubActionsCI.Command
		}
		if c.Tools != nil && c.Tools.GitHubActionsCI != nil && c.Tools.GitHubActionsCI.Buildevents != nil {
			buildevents = *c.Tools.GitHubActionsCI.Buildevents
		}
	}

	roleName := lastNonNil(GitHubActionsRoleNameGetter, commons...)
	return &GitHubActionsCI{
		CommonCI: CommonCI{
			Enabled:        &enabled,
			Buildevents:    &buildevents,
			AWSIAMRoleName: roleName,
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

func BackendAccountIdGetter(comm Common) *string {
	if comm.Backend != nil {
		return comm.Backend.AccountID
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

func AWSProviderRegionGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.Region
	}
	return nil
}

func AWSProviderVersionGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.Version
	}
	return nil
}

func AWSProviderProfileGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.Profile
	}
	return nil
}

func AWSProviderAccountIdGetter(comm Common) *json.Number {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.AccountID
	}
	return nil
}

func AWSProviderAdditionalRegionsGetter(comm Common) []string {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.AdditionalRegions
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

func OktaProviderVersionGetter(comm Common) *string {
	if comm.Providers == nil || comm.Providers.Okta == nil {
		return nil
	}
	return comm.Providers.Okta.Version
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
