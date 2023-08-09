package validation

import "github.com/KusoKaihatsuSha/startup/internal/helpers"

var (
	Valids []Valid
)

type Valid interface {
	Valid(string, string) (any, bool)
}

// AddValidation - add custom validation. Example in interface docs.
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

// Implements default validations
func (o tmpFileValid) Valid(key, value string) (any, bool) {
	if key == string(o) {
		return helpers.ValidTempFile(value), true
	}
	return nil, false
}

// Implements default validations
func (o fileValid) Valid(key, value string) (any, bool) {
	if key == string(o) {
		return helpers.ValidFile(value), true
	}
	return nil, false
}

// Implements default validations
func (o urlValid) Valid(key, value string) (any, bool) {
	if key == string(o) {
		return helpers.ValidURL(value), true
	}
	return nil, false
}

// Implements default validations
func (o boolValid) Valid(key, value string) (any, bool) {
	if key == string(o) {
		return helpers.ValidBool(value), true
	}
	return nil, false
}

// Implements default validations
func (o intValid) Valid(key, value string) (any, bool) {
	if key == string(o) {
		return helpers.ValidInt(value), true
	}
	return nil, false
}

// Implements default validations
func (o floatValid) Valid(key, value string) (any, bool) {
	if key == string(o) {
		return helpers.ValidFloat(value), true
	}
	return nil, false
}

// Implements default validations
func (o durationValid) Valid(key, value string) (any, bool) {
	if key == string(o) {
		return helpers.ValidTimer(value), true
	}
	return nil, false
}

// Implements default validations
func (o uuidValid) Valid(key, value string) (any, bool) {
	if key == string(o) {
		return helpers.ValidUUID(value), true
	}
	return nil, false
}
