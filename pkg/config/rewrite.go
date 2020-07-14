package config

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Rewrite struct {
	Rules []*RewriteRule `validate:"dive"`
}

type RewriteRule struct {
	Pattern         string         `validate:"required"`
	CompiledPattern *regexp.Regexp `yaml:"-"`
	Replace         string         `validate:"required"`
}

func init() {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		rule := sl.Current().Interface().(RewriteRule)
		_, err := regexp.Compile(rule.Pattern)
		if err != nil {
			sl.ReportError(rule.Pattern, "pattern", "Pattern", "", "")
		}
	}, RewriteRule{})
}
