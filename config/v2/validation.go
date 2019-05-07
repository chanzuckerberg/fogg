package v2

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	multierror "github.com/hashicorp/go-multierror"
	validator "gopkg.in/go-playground/validator.v9"
)

// Validate validates the config
func (c *Config) Validate() error {
	if c == nil {
		return errs.NewInternal("config is nil")
	}
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
	errs = multierror.Append(errs, c.validateProviderRegion())
	errs = multierror.Append(errs, c.validateProviderProfile())
	errs = multierror.Append(errs, c.validateProviderVersion())

	return errs.ErrorOrNil()
}

func nonEmptyString(s string) bool {
	return len(s) > 0
}

func (c *Config) validateOwners() *multierror.Error {
	var getter = func(comm Common) string {
		return comm.Owner
	}
	return c.validateInheritedStringField("owner", getter, nonEmptyString)
}

func (c *Config) validateProjects() *multierror.Error {
	var getter = func(comm Common) string {
		return comm.Project
	}
	return c.validateInheritedStringField("project", getter, nonEmptyString)
}

func (c *Config) validateTerraformVerion() *multierror.Error {
	var getter = func(comm Common) string {
		return comm.TerraformVersion
	}
	return c.validateInheritedStringField("terraform version", getter, nonEmptyString)
}

func (c *Config) validateBackendBucket() *multierror.Error {
	var getter = func(comm Common) string {
		return comm.Backend.Bucket
	}
	return c.validateInheritedStringField("backend bucket", getter, nonEmptyString)
}

func (c *Config) validateBackendRegion() *multierror.Error {
	var getter = func(comm Common) string {
		return comm.Backend.Region
	}
	return c.validateInheritedStringField("backend region", getter, nonEmptyString)
}

func (c *Config) validateProviderRegion() *multierror.Error {
	var getter = func(comm Common) string {
		if comm.Providers.AWS != nil && comm.Providers.AWS.Region != nil {
			return *comm.Providers.AWS.Region

		}
		// once we reconcile the use of nils across the v2 config this should get simpler
		return ""

	}
	return c.validateInheritedStringField("provider region", getter, nonEmptyString)
}

func (c *Config) validateProviderProfile() *multierror.Error {
	var getter = func(comm Common) string {
		if comm.Providers.AWS != nil && comm.Providers.AWS.Profile != nil {
			return *comm.Providers.AWS.Profile

		}
		// once we reconcile the use of nils across the v2 config this should get simpler
		return ""

	}
	return c.validateInheritedStringField("provider profile", getter, nonEmptyString)
}

func (c *Config) validateProviderVersion() *multierror.Error {
	var getter = func(comm Common) string {
		if comm.Providers.AWS != nil && comm.Providers.AWS.Version != nil {
			return *comm.Providers.AWS.Version

		}
		// once we reconcile the use of nils across the v2 config this should get simpler
		return ""

	}
	return c.validateInheritedStringField("provider version", getter, nonEmptyString)
}

// validateInheritedStringField will walk all accounts and components and ensure that a given field is valid at at least
// one level of the inheritance hierarchy. We should eventually distinuish between not present and invalid because
// if the value is present but invalid we should probably mark it as such, rather than papering over it.
func (c *Config) validateInheritedStringField(fieldName string, getter func(Common) string, validator func(string) bool) *multierror.Error {
	var err *multierror.Error

	// For each account, we need the field to be valid in either the defaults or account
	for acctName, acct := range c.Accounts {
		if !(validator(getter(c.Defaults.Common)) || validator(getter(acct.Common))) {
			err = multierror.Append(err, fmt.Errorf("account %s must have a valid %s set at either the account or defaults level", acctName, fieldName))
		}
	}

	// global
	if !(validator(getter(c.Defaults.Common)) || validator(getter(c.Global.Common))) {
		err = multierror.Append(err, fmt.Errorf("global must have a valid %s set at either the global or defaults level", fieldName))
	}

	// For each component, we need the field to be valid at one of defaults, env or component
	for envName, env := range c.Envs {
		for componentName, component := range env.Components {
			if !(validator(getter(c.Defaults.Common)) || validator(getter(env.Common)) || validator(getter(component.Common))) {
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
			if _, ok := v1.ReservedVariableNames[extraVar]; ok {
				err = multierror.Append(err, fmt.Errorf("extra_var[%s] is a fogg reserved variable name", extraVar))
			}
		}
	}
	extraVars := []map[string]string{}
	if c.Defaults.ExtraVars != nil {
		extraVars = append(extraVars, c.Defaults.ExtraVars)
	}
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
