package migrate

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"

	"github.com/antzucaro/matchr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	prompt "github.com/segmentio/go-prompt"
	"github.com/sirupsen/logrus"
)

func generatePlan(planPath string) error {
	cmd := exec.Command("make", "init")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Could not run make init")
	}

	cmd = exec.Command("make", "run")
	cmd.Env = append(cmd.Env, fmt.Sprintf("plan -out %s", planPath))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	return errors.Wrap(err, "Could not run terraform plan")
}

// parsePlan
func parsePlan(planPath string) error {
	f, err := os.Open(planPath)
	if err != nil {
		return errors.Wrapf(err, "Could not read plan at %s", planPath)
	}
	defer f.Close()
	// TODO: also remove the plan?
	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return errors.Wrapf(err, "Terraform could not parse plan at %s", planPath)
	}
	if plan.Diff == nil {
		logrus.Debug("nil diff")
		return nil
	}

	deletions := map[string]bool{}
	additions := map[string]bool{}

	for _, module := range plan.Diff.Modules {
		moduleName := strings.TrimPrefix(strings.Join(module.Path, "."), "root.")
		for name, instance := range module.Resources {
			fullName := fmt.Sprintf("%s.%s", moduleName, name)
			if instance.Destroy {
				deletions[fullName] = true
			} else {
				additions[fullName] = true
			}
		}
	}

	for addition := range additions {
		currScore := math.MaxInt64
		var replace *string

		for deletion, ok := range deletions {
			if !ok {
				continue
			}
			score := matchr.DamerauLevenshtein(addition, deletion)
			if score < currScore {
				currScore = score
				replace = aws.String(deletion)
			}
		}

		if replace == nil {
			continue
		}

		if !prompt.Confirm("Would you like us to move %s to %s", *replace, addition) {
			continue
		}

		deletions[*replace] = false
		cmd := exec.Command("make", "run")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Env = append(cmd.Env, fmt.Sprintf("CMD=state mv %s %s", *replace, addition))

		err = cmd.Run()
		if err != nil {
			return errors.Wrapf(err, "Could not move %s to %s", *replace, addition)
		}
	}
	return nil
}

// Migrate migrates
func Migrate(planPath string) error {
	defer os.Remove(planPath)
	err := generatePlan(planPath)
	if err != nil {
		return err
	}
	return parsePlan(planPath)
}
