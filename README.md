[![godoc](https://godoc.org/github.com/KusoKaihatsuSha/startup?status.svg)](https://godoc.org/github.com/KusoKaihatsuSha/startup) [![Go Report Card](https://goreportcard.com/badge/github.com/KusoKaihatsuSha/startup)](https://goreportcard.com/report/github.com/KusoKaihatsuSha/startup) [![go test](https://github.com/KusoKaihatsuSha/startup/actions/workflows/test.yml/badge.svg)](https://github.com/KusoKaihatsuSha/startup/actions/workflows/test.yml)


# Startup

> Use the package's functionality for filling custom golang 'struct' from environments, flags or configuration file (JSON) with the required order. 

### **Usage**

```
go get github.com/KusoKaihatsuSha/startup
```

### **Example**

```golang
// Custom struct. Struct will be implement in program with selected 'Stages' variable.
type Test struct {
    NewValid []string `json:"new-valid" default:"new valid is default" flag:"valid" text:"-" valid:"test"`
}

// Custom type.
type testValid string

// Custom variable.
var testValidation testValid = "test"

// Custom method.
func (o testValid) Valid(key, value string) (any, bool) {
        return []string{value + "+++"}, true
    }

// add custom validation
func MyFunc() {
    ...
    startup.AddValidation(testValidation)
    ...
    // Implement all types of configs (ORDER: Environment -> Flags -> Json file).
    configurations := startup.InitForce[Test](startup.EnvFlagFile)
    // Test print.
    fmt.Println(configurations)
}
```

### **Default validations** (in tag 'valid' inside annotation)

  - 'tmp_file' - As 'file', but if empty returm file from Temp folder  (string in struct)
  - 'file' - Check exist the filepath (string in struct)
  - 'url' - Check url is correct (string in struct)
  - 'bool' - Parse Bool (bool in struct)
  - 'int' - Parse int (int64 in struct)
  - 'float' - Parse float (float64 in struct)
  - 'duration' - Parse duration (time.Duration in struct)
  - 'uuid' - Check uuid. Return new if not exist (string in struct)

### **Caution**

flags are reserved:
  - h
  - help
  - config