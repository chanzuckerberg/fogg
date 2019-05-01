package v1

var reservedVariableNames = map[string]bool{
	"aws_accounts": true,
	"aws_profile":  true,
	"env":          true,
	"owner":        true,
	"project":      true,
	"region":       true,
	"tags":         true,
}
