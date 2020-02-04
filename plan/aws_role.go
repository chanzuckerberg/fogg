package plan

type AWSRole struct {
	AccountID string `yaml:"account_id"`
	RolePath  string `yaml:"role_path"`
	RoleName  string `yaml:"role_name"`
}
