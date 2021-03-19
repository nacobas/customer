package customer

import (
	"reflect"
	"regexp"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("person-name", ValidatePersonName)
	v.RegisterValidation("org-name", ValidateOrgName)
	v.RegisterValidation("before", ValidateBeforeNow)
	v.RegisterCustomTypeFunc(ValidateDate, date.Date{})
	return v
}

// single unicode Letter or 2 to many unicode Letters, allowing - in the middle
var personNameRegexp *regexp.Regexp = regexp.MustCompile(`^([\p{L}])$|^([\p{L}])([\p{L}-])*([\p{L}])$`)

func ValidatePersonName(fl validator.FieldLevel) bool {
	return personNameRegexp.MatchString(fl.Field().String())
}

// single unicode Letter,digit OR 2 to many unicode Letters, digits, allowing - & and space in the middle
var orgNameRegexp *regexp.Regexp = regexp.MustCompile(`^([\p{L}\d])$|^([\p{L}\d])([\p{L}\d-& ])*([\p{L}\d])$`)

func ValidateOrgName(fl validator.FieldLevel) bool {
	return orgNameRegexp.MatchString(fl.Field().String())
}

func ValidateBeforeNow(fl validator.FieldLevel) bool {

	return fl.Field().Interface().(time.Time).Before(time.Now())
}

func ValidateDate(field reflect.Value) interface{} {

	if date, ok := field.Interface().(date.Date); ok {
		return date.ToTime()
	}

	return nil
}
