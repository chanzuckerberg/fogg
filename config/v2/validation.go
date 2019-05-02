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
	var errs *multierror.Error

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

	errs = multierror.Append(errs, v.Struct(c))
	errs = multierror.Append(errs, c.validateExtraVars())
	errs = multierror.Append(errs, c.validateOwners())
	errs = multierror.Append(errs, c.validateProjects())
	errs = multierror.Append(errs, c.validateTerraformVerion())
	errs = multierror.Append(errs, c.validateBackendBucket())
	errs = multierror.Append(errs, c.validateBackendRegion())

	return errs.ErrorOrNil()
}

func nonEmptyString(s string) bool {
	return len(s) > 0
}

func (c *Config) validateOwners() *multierror.Error {
	var getter = func(comm common) string {
		return comm.Owner
	}
	return c.validateInheritedStringField("owner", getter, nonEmptyString)
}

func (c *Config) validateProjects() *multierror.Error {
	var getter = func(comm common) string {
		return comm.Project
	}
	return c.validateInheritedStringField("project", getter, nonEmptyString)
}

func (c *Config) validateTerraformVerion() *multierror.Error {
	var getter = func(comm common) string {
		return comm.Project
	}
	return c.validateInheritedStringField("project", getter, nonEmptyString)
}

func (c *Config) validateBackendBucket() *multierror.Error {
	var getter = func(comm common) string {
		return comm.Backend.Bucket
	}
	return c.validateInheritedStringField("backend bucket", getter, nonEmptyString)
}

func (c *Config) validateBackendRegion() *multierror.Error {
	var getter = func(comm common) string {
		return comm.Backend.Region
	}
	return c.validateInheritedStringField("backend region", getter, nonEmptyString)
}

// TODO validateAWSProviders

func (c *Config) validateInheritedStringField(fieldName string, getter func(common) string, validator func(string) bool) *multierror.Error {
	var err *multierror.Error

	// For each account, we need a valid owner to be set in either the defaults or account
	for acctName, acct := range c.Accounts {
		if !(validator(getter(c.Defaults.common)) || validator(getter(acct.common))) {
			err = multierror.Append(err, fmt.Errorf("account %s must have a valid %s set at either the account or defaults level", acctName, fieldName))
		}
	}

	// For each component, we need a valid owner to be set at one of defaults, env or component
	for envName, env := range c.Envs {
		for componentName, component := range env.Components {
			if !(validator(getter(c.Defaults.common)) || validator(getter(env.common)) || validator(getter(component.common))) {
				err = multierror.Append(err, fmt.Errorf("componnent %s/%s must have a valid %s", envName, componentName, fieldName))
			}
		}
	}

	return err
}

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
