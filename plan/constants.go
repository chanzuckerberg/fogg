package plan

var reservedVariableNames = map[string]bool{
	"project":      true,
	"region":       true,
	"aws_profile":  true,
	"owner":        true,
	"aws_accounts": true,
}
