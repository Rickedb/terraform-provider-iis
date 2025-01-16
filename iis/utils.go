package iis

import (
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateAllowedValues(allowedValues []string) schema.SchemaValidateDiagFunc {
	return func(val interface{}, path cty.Path) diag.Diagnostics {
		v := val.(string)
		for _, allowed := range allowedValues {
			if v == allowed {
				return nil
			}
		}

		return diag.Errorf("%q must be one of %v, got: %q", path, allowedValues, v)
	}
}

func greaterOrEqualThan(minValue int) schema.SchemaValidateDiagFunc {
	return func(val interface{}, path cty.Path) diag.Diagnostics {
		v := val.(int)
		if v >= minValue {
			return nil
		}

		return diag.Errorf("%q must be greater than %v", path, v)
	}
}

func isInBetweenValues(minValue int, maxValue int) schema.SchemaValidateDiagFunc {
	return func(val interface{}, path cty.Path) diag.Diagnostics {
		v := val.(int)
		if v >= minValue && v <= maxValue {
			return nil
		}

		return diag.Errorf("%q must be between %v and %v", path, minValue, maxValue)
	}
}

func validateAllowedIntValues(allowedValues []int) schema.SchemaValidateDiagFunc {
	return func(val interface{}, path cty.Path) diag.Diagnostics {
		v := val.(int)
		for _, allowed := range allowedValues {
			if v == allowed {
				return nil
			}
		}

		return diag.Errorf("%q must be one of %v, got: %q", path, allowedValues, v)
	}
}

func isValid(regex *regexp.Regexp) schema.SchemaValidateDiagFunc {
	return func(val interface{}, path cty.Path) diag.Diagnostics {
		v := val.(string)
		if regex.MatchString(v) {
			return nil
		}

		return diag.Errorf("%q is not valid", path)
	}
}

func isValidPath(onlyBackslashOnly bool) schema.SchemaValidateDiagFunc {
	pattern := `^[a-zA-Z]:[\\|^/](?:[^/:*?"<>|\\]+[\\|/]?)*$`
	if onlyBackslashOnly {
		pattern = `^[a-zA-Z]:[\\](?:[^/:*?"<>|\\]+[\\]?)*$`
	}

	return isValid(regexp.MustCompile(pattern))
}
