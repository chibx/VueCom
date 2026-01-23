package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New(validator.WithRequiredStructEnabled())

func TagNameFunc(fld reflect.StructField) string {
	name := fld.Tag.Get("name")
	if name == "" {
		return fld.Name
	}
	return name
}

// func UserStructLevelValidation(sl validator.StructLevel) {

// 	user := sl.Current().Interface().(User)

// 	if len(user.FirstName) == 0 && len(user.LastName) == 0 {
// 		sl.ReportError(user.FirstName, "fname", "FirstName", "fnameorlname", "")
// 		sl.ReportError(user.LastName, "lname", "LastName", "fnameorlname", "")
// 	}

// 	// plus can do more, even with different tag than "fnameorlname"
// }
