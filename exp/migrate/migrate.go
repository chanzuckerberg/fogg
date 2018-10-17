package migrate

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func generatePlan(planPath string) error {
	cmd := exec.Command("make", "run")
	cmd.Env = append(cmd.Env, fmt.Sprintf("plan -out %s", planPath))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
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
		log.Debug("nil diff")
		return nil
	}

	spew.Dump(plan.Diff.Modules)
	return nil
}

// Migrate migrates
func Migrate(planPath string) error {
	err := generatePlan(planPath)
	if err != nil {
		return err
	}
	err = parsePlan(planPath)
	return err
}
