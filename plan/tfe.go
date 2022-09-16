package plan

type TFEConfig struct {
	ReadTeams                      []string
	Branch                         string
	GithubOrg                      string
	GithubRepo                     string
	TFEOrg                         string
	SSHKeyName                     string
	ExcludedGithubRequiredChecks   []string
	AdditionalGithubRequiredChecks []string
}
