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
	errs = multierror.Append(errs, c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("project", ProjectGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("terraform version", TerraformVersionGetter, nonEmptyString))

	errs = multierror.Append(errs, c.validateInheritedStringField("backend bucket", BackendBucketGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("backend region", BackendRegionGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("backend profile", BackendProfileGetter, nonEmptyString))

	errs = multierror.Append(errs, c.validateInheritedStringField("provider region", AWSProviderRegionGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("provider profile", AWSProviderProfileGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedStringField("provider version", AWSProviderVersionGetter, nonEmptyString))
	errs = multierror.Append(errs, c.validateInheritedIntField("provider account id", AWSProviderAccountIdGetter, nonZeroInt64))

	errs = multierror.Append(errs, c.validateModules())

	return errs.ErrorOrNil()
}

func nonEmptyString(s string) bool {
	return len(s) > 0
}

func OwnerGetter(comm Common) *string {
	return comm.Owner
}

func ProjectGetter(comm Common) *string {
	return comm.Project
}

func TerraformVersionGetter(comm Common) *string {
	return comm.TerraformVersion
}

func BackendBucketGetter(comm Common) *string {
	if comm.Backend != nil {
		return comm.Backend.Bucket
	}
	return nil
}
func BackendRegionGetter(comm Common) *string {
	if comm.Backend != nil {
		return comm.Backend.Region
	}
	return nil
}

func BackendDynamoTableGetter(comm Common) *string {
	if comm.Backend != nil {
		return comm.Backend.DynamoTable
	}
	return nil
}

func BackendProfileGetter(comm Common) *string {
	if comm.Backend != nil {
		return comm.Backend.Profile
	}
	return nil
}

func AWSProviderRegionGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.Region

	}
	return nil
}

func AWSProviderVersionGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.Version
	}
	return nil
}

func AWSProviderProfileGetter(comm Common) *string {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.Profile
	}
	return nil
}

func AWSProviderAccountIdGetter(comm Common) *int64 {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.AccountID
	}
	return nil
}

func AWSProviderAdditionalRegionsGetter(comm Common) []string {
	if comm.Providers != nil && comm.Providers.AWS != nil {
		return comm.Providers.AWS.AdditionalRegions
	}
	return nil
}

func nonZeroInt64(i int64) bool {
	return true
}

// lastNonNil, despite its name can return nil if all results are nil
func lastNonNil(getter func(Common) *string, commons ...Common) *string {
	var s *string
	for _, c := range commons {
		t := getter(c)
		if t != nil {
			s = t
		}
	}
	return s
}

// lastNonNilInt, despite its name can return nil if all results are nil
func lastNonNilInt64(getter func(Common) *int64, commons ...Common) *int64 {
	var s *int64
	for _, c := range commons {
		t := getter(c)
		if t != nil {
			s = t
		}
	}
	return s
}

// lastNonNilStringSlice, despite its name can return nil if all results are nil
func lastNonNilStringSlice(getter func(Common) []string, commons ...Common) []string {
	var s []string
	for _, c := range commons {
		t := getter(c)
		if t != nil {
			s = t
		}
	}
	return s
}

// ResolveRequiredString will resolve the value and panic if it is nil. Only to be used after validations are run.
func ResolveRequiredString(getter func(Common) *string, commons ...Common) string {
	return *lastNonNil(getter, commons...)
}

// ResolveRequiredInt will resolve the value and panic if it is nil. Only to be used after validations are run.
func ResolveRequiredInt64(getter func(Common) *int64, commons ...Common) int64 {
	return *lastNonNilInt64(getter, commons...)
}

func ResolveOptionalString(getter func(Common) *string, commons ...Common) *string {
	return lastNonNil(getter, commons...)
}

func ResolveOptionalStringSlice(getter func(Common) []string, commons ...Common) []string {
	return lastNonNilStringSlice(getter, commons...)
}

func ResolveStringArray(def []string, override []string) []string {
	if override != nil {
		return override
	}
	return def
}

// validateInheritedStringField will walk all accounts and components and ensure that a given field is valid at at least
// one level of the inheritance hierarchy. We should eventually distinuish between not present and invalid because
// if the value is present but invalid we should probably mark it as such, rather than papering over it.
func (c *Config) validateInheritedStringField(fieldName string, getter func(Common) *string, validator func(string) bool) *multierror.Error {
	var err *multierror.Error

	// For each account, we need the field to be valid in either the defaults or account
	for acctName, acct := range c.Accounts {
		v := lastNonNil(getter, c.Defaults.Common, acct.Common)
		if v == nil {
			err = multierror.Append(err, fmt.Errorf("account %s must have field %s at either the account or defaults level", acctName, fieldName))
		} else if !validator(*v) {
			err = multierror.Append(err, fmt.Errorf("account %s must have a valid %s set at either the account or defaults level", acctName, fieldName))
		}
	}

	// global
	v := lastNonNil(getter, c.Defaults.Common, c.Global.Common)
	if v == nil {
		err = multierror.Append(err, fmt.Errorf("global must have a %s set at either the global or defaults level", fieldName))
	} else if !validator(*v) {
		err = multierror.Append(err, fmt.Errorf("global must have a valid %s set at either the global or defaults level", fieldName))
	}

	// For each component, we need the field to be valid at one of defaults, env or component
	for envName, env := range c.Envs {
		for componentName, component := range env.Components {
			v := lastNonNil(getter, c.Defaults.Common, env.Common, component.Common)

			if v == nil {
				err = multierror.Append(err, fmt.Errorf("componnent %s/%s must have a %s", envName, componentName, fieldName))
			} else if !validator(*v) {
				err = multierror.Append(err, fmt.Errorf("componnent %s/%s must have a valid %s", envName, componentName, fieldName))
			}
		}
	}

	return err
}

// validateInheritedIntField will walk all accounts and components and ensure that a given field is valid at at least
// one level of the inheritance hierarchy. We should eventually distinuish between not present and invalid because
// if the value is present but invalid we should probably mark it as such, rather than papering over it.
func (c *Config) validateInheritedIntField(fieldName string, getter func(Common) *int64, validator func(int64) bool) *multierror.Error {
	// TODO figure out how to refactor to consolidate with above func
	var err *multierror.Error

	// For each account, we need the field to be valid in either the defaults or account
	for acctName, acct := range c.Accounts {
		v := lastNonNilInt64(getter, c.Defaults.Common, acct.Common)
		if v == nil {
			err = multierror.Append(err, fmt.Errorf("account %s must have field %s at either the account or defaults level", acctName, fieldName))
		} else if !validator(*v) {
			err = multierror.Append(err, fmt.Errorf("account %s must have a valid %s set at either the account or defaults level", acctName, fieldName))
		}
	}

	// global
	v := lastNonNilInt64(getter, c.Defaults.Common, c.Global.Common)
	if v == nil {
		err = multierror.Append(err, fmt.Errorf("global must have a %s set at either the global or defaults level", fieldName))
	} else if !validator(*v) {
		err = multierror.Append(err, fmt.Errorf("global must have a valid %s set at either the global or defaults level", fieldName))
	}

	// For each component, we need the field to be valid at one of defaults, env or component
	for envName, env := range c.Envs {
		for componentName, component := range env.Components {
			v := lastNonNilInt64(getter, c.Defaults.Common, env.Common, component.Common)

			if v == nil {
				err = multierror.Append(err, fmt.Errorf("componnent %s/%s must have a %s", envName, componentName, fieldName))
			} else if !validator(*v) {
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

func ResolveModuleTerraformVersion(def Defaults, module v1.Module) *string {
	if module.TerraformVersion != nil {
		return module.TerraformVersion
	}
	return def.TerraformVersion
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
