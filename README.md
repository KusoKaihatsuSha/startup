[![godoc](https://godoc.org/github.com/KusoKaihatsuSha/startup?status.svg)](https://godoc.org/github.com/KusoKaihatsuSha/startup) [![Go Report Card](https://goreportcard.com/badge/github.com/KusoKaihatsuSha/startup)](https://goreportcard.com/report/github.com/KusoKaihatsuSha/startup) [![go test](https://github.com/KusoKaihatsuSha/startup/actions/workflows/test.yml/badge.svg)](https://github.com/KusoKaihatsuSha/startup/actions/workflows/test.yml)


# Package `startup`

> Use the package's functionality to simplify the use of flags, environments and config file when starting the application.  
> Fill custom Golang `struct` from environments, flags or configuration file (JSON) with desired order.

### Install

```shell
go get github.com/KusoKaihatsuSha/startup
```

### Example

```golang
package main

import (
    "encoding/json"
    "fmt"
    "net"
    "strings"
    "time"
    
    "github.com/KusoKaihatsuSha/startup"
)

// Custom struct
type Configuration struct {
    TestEmail    string        `json:"test-email"    default:"a@b.c"         flag:"test-email"    env:"TEST_EMAIL"    help:"email"`
    TestSlice    []string      `json:"test-slice"    default:"1,2,3,4;5;6"   flag:"test-slice"    env:"TEST_SLICE"    help:"slice" valid:"slice"`
    TestInt      int64         `json:"test-int"      default:"11"            flag:"test-int"      env:"TEST_INT"      help:"int"`
    TestDuration time.Duration `json:"test-duration" default:"1s"            flag:"test-duration" env:"TEST_DURATION" help:"duration"`
    TestBool     bool          `json:"test-bool"     default:"true"          flag:"test-bool"     env:"TEST_BOOL"     help:"bool"`
    TestFloat    float64       `json:"test-float"    default:"1"             flag:"test-float"    env:"TEST_FLOAT"    help:"float"`
    TestUint     uint64        `json:"test-uint"     default:"111"           flag:"test-uint"     env:"TEST_UINT"     help:"uint"`
    TestIP       net.IP        `json:"test-ip"       default:"127.0.0.1"     flag:"test-ip"       env:"TEST_IP"       help:"ip"`
    TestJSON     TestJSON      `json:"test-json"     default:"{\"param1\":\"default_001\",\"param2\":\"default_002\"}" flag:"test-json" env:"TEST_JSON" help:"json"`
}

// JSON field and unmashalText
type TestJSON struct {
    P1 string `json:"param1"`
    P2 string `json:"param2"`
}

func (t *TestJSON) UnmarshalText(text []byte) error {
    type Tmp TestJSON
    tmp := (*Tmp)(t)
    if err := json.Unmarshal(text, tmp); err != nil {
        return err
    }
    t = (*TestJSON)(tmp)
    return nil
}

// Custom slice validation
type sliceValid string

var sliceValidation sliceValid = "slice"

func (o sliceValid) Valid(def string, value any) (any, bool) {
    newValue := strings.ReplaceAll(def, ";", ",")
    return strings.Split(newValue, ","), true
}

func main() {
    
    // Add custom validation
    startup.AddValidation(sliceValidation)
    
	// Implement all types of configs (ORDER: Json file -> Environment -> Flags ).
    get := startup.GetForce[Configuration](
        startup.File,
        startup.Env,
        startup.Flag,
    )
    
    //---PRINT SECTION---
    fmt.Printf(
        "string ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
        get.TestEmail,
    )
    
    fmt.Printf(
        "custom slice ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
        get.TestSlice,
    )
    
    fmt.Printf(
        "integer max 10 ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
        get.TestInt,
    )
    
    fmt.Printf(
        "duration ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
        get.TestDuration,
    )
    
    fmt.Printf(
        "boolean ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
        get.TestBool,
    )
    
    fmt.Printf(
        "float ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
        get.TestFloat,
    )
    
    fmt.Printf(
        "uint ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
        get.TestUint,
    )
    
    fmt.Printf(
        "ip ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
        get.TestIP,
    )
    
    fmt.Printf(
        "json unmarshal ↣ TYPE [%[1]T] VALUE [%#[1]v]\n",
        get.TestJSON,
    )
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
  - config