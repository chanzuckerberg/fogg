package v2

import (
	"encoding/json"

	v1 "github.com/chanzuckerberg/fogg/config/v1"
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

// ResolveGithubProvider will return an GithubProvder iff one of the required fields is set somewhere in the set of Common
// config objects passed in. Otherwise it will return nil.
func ResolveGithubProvider(commons ...Common) *GithubProvider {
	org := lastNonNil(GithubProviderOrganizationGetter, commons...)

	if org != nil {
		return &GithubProvider{
			Organization: org,

			// optional fields
			BaseURL: lastNonNil(GithubProviderBaseURLGetter, commons...),
		}
	}
	return nil
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

func ResolveTfLint(commons ...Common) v1.TfLint {
	enabled := false
	for _, c := range commons {
		if c.Tools != nil && c.Tools.TfLint != nil && c.Tools.TfLint.Enabled != nil {
			enabled = *c.Tools.TfLint.Enabled
		}
	}

	return v1.TfLint{
		Enabled: &enabled,
	}
}

func ResolveAtlantis(commons ...Common) *Atlantis {
	enabled := false
	for _, c := range commons {
		if c.Tools != nil && c.Tools.Atlantis != nil && c.Tools.Atlantis.Enabled != nil {
			enabled = *c.Tools.Atlantis.Enabled
		}
	}

	roleName := lastNonNil(AtlantisRoleNameGetter, commons...)
	rolePath := lastNonNil(AtlantisRolePathGetter, commons...)

	return &Atlantis{
		Enabled:  &enabled,
		RoleName: roleName,
		RolePath: rolePath,
	}
}

func ResolveTravis(commons ...Common) *v1.TravisCI {
	enabled := false
	testCommand := "check"
	for _, c := range commons {
		if c.Tools != nil && c.Tools.TravisCI != nil && c.Tools.TravisCI.Enabled != nil {
			enabled = *c.Tools.TravisCI.Enabled
		}

		if c.Tools != nil && c.Tools.TravisCI != nil && c.Tools.TravisCI.Command != nil {
			testCommand = *c.Tools.TravisCI.Command
		}
	}

	roleName := lastNonNil(TravisRoleNameGetter, commons...)

	return &v1.TravisCI{
		Enabled:        &enabled,
		AWSIAMRoleName: roleName,
		Command:        &testCommand,
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

func ExtraVarsGetter(comm Common) map[string]string {
	if comm.ExtraVars != nil {
		return comm.ExtraVars
	}
	return map[string]string{}
}

func ResolveModuleTerraformVersion(def Defaults, module v1.Module) *string {
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

func AtlantisRolePathGetter(comm Common) *string {
	if comm.Tools == nil || comm.Tools.Atlantis == nil {
		return nil
	}
	return comm.Tools.Atlantis.RolePath
}

func AtlantisRoleNameGetter(comm Common) *string {
	if comm.Tools == nil || comm.Tools.Atlantis == nil {
		return nil
	}
	return comm.Tools.Atlantis.RoleName
}

func TravisRoleNameGetter(comm Common) *string {
	if comm.Tools == nil || comm.Tools.TravisCI == nil {
		return nil
	}
	return comm.Tools.TravisCI.AWSIAMRoleName
}

func TravisTestCommandGetter(comm Common) *string {
	if comm.Tools == nil || comm.Tools.TravisCI == nil {
		return nil
	}
	return comm.Tools.TravisCI.Command
}
