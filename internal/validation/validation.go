package validation

import "github.com/KusoKaihatsuSha/startup/internal/helpers"

var (
	Valids []Valid
)

type Valid interface {
	Valid(string, any) (any, bool)
}

// Add - add custom validation. Example in interface docs.
func Add(value ...Valid) {
	Valids = append(Valids, value...)
}

var (
	tmpFileValidation  tmpFileValid  = "tmp_file"
	fileValidation     fileValid     = "file"
	urlValidation      urlValid      = "url"
	boolValidation     boolValid     = "bool"
	intValidation      intValid      = "int"
	floatValidation    floatValid    = "float"
	durationValidation durationValid = "duration"
	uuidValidation     uuidValid     = "uuid"
)

// Default validations
type (
	tmpFileValid  string
	fileValid     string
	urlValid      string
	boolValid     string
	intValid      string
	floatValid    string
	durationValid string
	uuidValid     string
)

func init() {
	// Will add the default validation checks in the handle of the struct
	Add(
		tmpFileValidation,
		fileValidation,
		urlValidation,
		boolValidation,
		intValidation,
		floatValidation,
		durationValidation,
		uuidValidation,
	)
}

// Valid Implements default validations
func (o tmpFileValid) Valid(def string, value any) (any, bool) {
	return helpers.ValidTempFile(value.(string)), true
}

// Valid Implements default validations
func (o fileValid) Valid(key string, value any) (any, bool) {
	return helpers.ValidFile(value.(string)), true
}

// Valid Implements default validations
func (o urlValid) Valid(key string, value any) (any, bool) {
	return helpers.ValidURL(value.(string)), true
}

// Valid Implements default validations
func (o boolValid) Valid(key string, value any) (any, bool) {
	return helpers.ValidBool(value.(string)), true
}

// Valid Implements default validations
func (o intValid) Valid(key string, value any) (any, bool) {
	return helpers.ValidInt(value.(string)), true
}

// Valid Implements default validations
func (o floatValid) Valid(key string, value any) (any, bool) {
	return helpers.ValidFloat(value.(string)), true
}

// Valid Implements default validations
func (o durationValid) Valid(key string, value any) (any, bool) {
	return helpers.ValidDuration(value.(string)), true
}

// Valid Implements default validations
func (o uuidValid) Valid(key string, value any) (any, bool) {
	return helpers.ValidUUID(value.(string)), true
}
