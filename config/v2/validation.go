package v2

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/chanzuckerberg/fogg/errs"
	multierror "github.com/hashicorp/go-multierror"
	validator "gopkg.in/go-playground/validator.v9"
)

// Validate validates the config
func (c *Config) Validate() error {

	v := validator.New()
	// This func gives us the ability to get the full path for a field deeply
	// nested in the structure
	// https://github.com/go-playground/validator/issues/323#issuecomment-343670840
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})
	err := v.Struct(c)

	if err != nil {
		return err
	}

	err = c.validateExtraVars()

	if err != nil {
		return err
	}

	return nil
}

func validOwner(owner string) bool {
	return len(owner) > 0
}

func (c *Config) validateOwners() *multierror.Error {
	var err *multierror.Error

	// For each account, we need a valid owner to be set in either the defaults or account
	for acctName, acct := range c.Accounts {
		if !(validOwner(c.Defaults.Owner) || validOwner(acct.Owner)) {
			err = multierror.Append(err, fmt.Errorf("account %s must have a valid owner set at either the account or defaults level", acctName))
		}
	}

	// For each component, we need a valid owner to be set at one of defaults, env or component
	for envName, env := range c.Envs {
		for componentName, component := range env.Components {
			if !(validOwner(c.Defaults.Owner) || validOwner(env.Owner) || validOwner(component.Owner)) {
				err = multierror.Append(err, fmt.Errorf("componnent %s/%s must have a valid owner", envName, componentName))
			}
		}
	}

	return err
}

// validateBackend
// validateProject
// validateTerraformVerion
// validateAWSProviders

// validateExtraVars make sure users don't specify reserved variables
func (c *Config) validateExtraVars() error {
	var err *multierror.Error
	validate := func(extraVars map[string]string) {
		for extraVar := range extraVars {
			if _, ok := reservedVariableNames[extraVar]; ok {
				err = multierror.Append(err, fmt.Errorf("extra_var[%s] is a fogg reserved variable name", extraVar))
			}
		}
	}
	extraVars := []map[string]string{}
	extraVars = append(extraVars, c.Defaults.ExtraVars)
	for _, env := range c.Envs {
		extraVars = append(extraVars, env.ExtraVars)
		for _, component := range env.Components {
			extraVars = append(extraVars, component.ExtraVars)
		}
	}

	for _, acct := range c.Accounts {
		extraVars = append(extraVars, acct.ExtraVars)
	}

	for _, extraVar := range extraVars {
		validate(extraVar)
	}

	if err.ErrorOrNil() != nil {
		return errs.WrapUser(err.ErrorOrNil(), "extra_vars contains reserved variable names")
	}
	return nil
}
