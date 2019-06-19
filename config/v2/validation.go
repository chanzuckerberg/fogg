package v2

import (
	"fmt"
	"reflect"
	"strings"

	v1 "github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	multierror "github.com/hashicorp/go-multierror"
	validator "gopkg.in/go-playground/validator.v9"
)

// Validate validates the config
func (c *Config) Validate() ([]string, error) {
	if c == nil {
		return nil, errs.NewInternal("config is nil")
	}

	var errs *multierror.Error
	var warnings []string

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
	errs = multierror.Append(errs, c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("project", ProjectGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("terraform version", TerraformVersionGetter, nonEmptyString))

	errs = multierror.Append(errs, c.validateInheritedStringField("backend bucket", BackendBucketGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("backend region", BackendRegionGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("backend profile", BackendProfileGetter, nonEmptyString))

	errs = multierror.Append(errs, c.ValidateAWSProviders())
	errs = multierror.Append(errs, c.ValidateSnowflakeProviders())
	errs = multierror.Append(errs, c.ValidateBlessProviders())
	errs = multierror.Append(errs, c.ValidateOktaProviders())
	errs = multierror.Append(errs, c.validateModules())

	// refactor to make it easier to manage these
	w, e := c.ValidateToolsTravis()
	warnings = append(warnings, w...)
	errs = multierror.Append(errs, e)
	w, e = c.ValidateToolsTfLint()
	warnings = append(warnings, w...)
	errs = multierror.Append(errs, e)

	if c.Docker {
		warnings = append(warnings, "Docker support is no longer supported and the config entry is ignored.")
	}

	return warnings, errs.ErrorOrNil()
}

func merge(warnings []string, err *multierror.Error, w []string, e error) ([]string, *multierror.Error) {
	return append(warnings, w...), multierror.Append(err, e)
}

func ValidateAWSProvider(p *AWSProvider, component string) error {
	var errs *multierror.Error
	if p == nil {
		return nil // nothing to validate
	}

	if p.Region == nil {
		errs = multierror.Append(errs, fmt.Errorf("aws provider region for %s", component))
	}

	if p.Profile == nil {
		errs = multierror.Append(errs, fmt.Errorf("aws provider profile %s ", component))
	}

	if p.Version == nil {
		errs = multierror.Append(errs, fmt.Errorf("aws provider version for %s ", component))
	}

	if p.AccountID == nil || *p.AccountID == "" {
		errs = multierror.Append(errs, fmt.Errorf("aws provider account id for %s", component))
	}
	return errs
}

func (c *Config) ValidateAWSProviders() error {
	var errs *multierror.Error

	c.WalkComponents(func(component string, comms ...Common) {
		v := ResolveAWSProvider(comms...)
		if e := ValidateAWSProvider(v, component); e != nil {
			errs = multierror.Append(errs, e)
		}
	})

	return errs.ErrorOrNil()
}

func (p *BlessProvider) Validate(component string) error {
	var errs *multierror.Error
	if p == nil {
		return nil // nothing to do
	}
	if p.AWSProfile == nil {
		errs = multierror.Append(errs, fmt.Errorf("bless provider aws_profile required in %s", component))
	}
	if p.AWSRegion == nil {
		errs = multierror.Append(errs, fmt.Errorf("bless provider aws_region required in %s", component))
	}
	return errs
}

func (p *SnowflakeProvider) Validate(component string) error {
	var errs *multierror.Error
	if p == nil {
		return nil // nothing to do
	}

	if p.Account == nil {
		errs = multierror.Append(errs, fmt.Errorf("snowflake provider account must be set in %s", component))
	}

	if p.Role == nil {
		errs = multierror.Append(errs, fmt.Errorf("snowflake provider role must be set in %s", component))
	}

	if p.Region == nil {
		errs = multierror.Append(errs, fmt.Errorf("snowflake provider region must be set in %s", component))
	}

	return errs
}

func (o *OktaProvider) Validate(component string) error {
	var errs *multierror.Error
	if o == nil {
		return nil
	}
	if o.OrgName == nil {
		errs = multierror.Append(errs, fmt.Errorf("okta provider org_name required in %s", component))
	}
	return errs
}

func (c *Config) ValidateSnowflakeProviders() error {
	var errs *multierror.Error
	c.WalkComponents(func(component string, comms ...Common) {
		p := ResolveSnowflakeProvider(comms...)
		if e := p.Validate(component); e != nil {
			errs = multierror.Append(errs, e)
		}
	})
	return errs.ErrorOrNil()
}

func (c *Config) ValidateOktaProviders() error {
	var errs *multierror.Error
	c.WalkComponents(func(component string, comms ...Common) {
		p := ResolveOktaProvider(comms...)
		if err := p.Validate(component); err != nil {
			errs = multierror.Append(errs, err)
		}
	})
	return errs
}

func (c *Config) ValidateBlessProviders() error {
	var errs *multierror.Error
	c.WalkComponents(func(component string, comms ...Common) {
		p := ResolveBlessProvider(comms...)
		if err := p.Validate(component); err != nil {
			errs = multierror.Append(errs, err)
		}
	})
	return errs
}

func (c *Config) ValidateToolsTravis() ([]string, error) {
	var warns []string
	var errs *multierror.Error
	c.WalkComponents(func(component string, comms ...Common) {
		c := comms[len(comms)-1]
		if c.Tools != nil && c.Tools.TravisCI != nil {
			warns = append(warns, fmt.Sprintf("per-component travisci config is not implemented, ignoring config in %s", component))
		}
	})

	return warns, errs
}

func (c *Config) ValidateToolsTfLint() ([]string, error) {
	var warns []string
	var errs *multierror.Error
	c.WalkComponents(func(component string, comms ...Common) {
		c := comms[len(comms)-1]
		if c.Tools != nil && c.Tools.TfLint != nil {
			warns = append(warns, fmt.Sprintf("per-component tflint config is not implemented, ignoring config in %s", component))
		}
	})

	return warns, errs
}
func (c *Config) WalkComponents(f func(component string, commons ...Common)) {
	for name, acct := range c.Accounts {
		f(fmt.Sprintf("accounts/%s", name), c.Defaults.Common, acct.Common)
	}

	f("global", c.Defaults.Common, c.Global.Common)

	for envName, env := range c.Envs {
		for componentName, component := range env.Components {
			f(fmt.Sprintf("%s/%s", envName, componentName), c.Defaults.Common, env.Common, component.Common)
		}
	}
}

// validateInheritedStringField will walk all accounts and components and ensure that a given field is valid at at least
// one level of the inheritance hierarchy. We should eventually distinuish between not present and invalid because
// if the value is present but invalid we should probably mark it as such, rather than papering over it.
func (c *Config) validateInheritedStringField(fieldName string, getter func(Common) *string, validator func(*string) bool) *multierror.Error {
	var err *multierror.Error

	// For each account, we need the field to be valid in either the defaults or account
	for acctName, acct := range c.Accounts {
		v := lastNonNil(getter, c.Defaults.Common, acct.Common)
		if !validator(v) {
			err = multierror.Append(err, fmt.Errorf("account %s must have a valid %s set at either the account or defaults level", acctName, fieldName))
		}
	}

	// global
	v := lastNonNil(getter, c.Defaults.Common, c.Global.Common)
	if !validator(v) {
		err = multierror.Append(err, fmt.Errorf("global must have a valid %s set at either the global or defaults level", fieldName))
	}

	// For each component, we need the field to be valid at one of defaults, env or component
	for envName, env := range c.Envs {
		for componentName, component := range env.Components {
			v := lastNonNil(getter, c.Defaults.Common, env.Common, component.Common)
			if !validator(v) {
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

func (c *Config) validateModules() error {
	for name, module := range c.Modules {
		version := ResolveModuleTerraformVersion(c.Defaults, module)
		if version == nil {
			return fmt.Errorf("must set terrform version for module %s at either defaults or module level", name)
		}
	}
	return nil
}

func nonEmptyString(s *string) bool {
	return s != nil && len(*s) > 0
}
