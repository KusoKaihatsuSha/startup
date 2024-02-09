[![godoc](https://godoc.org/github.com/KusoKaihatsuSha/startup?status.svg)](https://godoc.org/github.com/KusoKaihatsuSha/startup) [![Go Report Card](https://goreportcard.com/badge/github.com/KusoKaihatsuSha/startup)](https://goreportcard.com/report/github.com/KusoKaihatsuSha/startup) [![go test](https://github.com/KusoKaihatsuSha/startup/actions/workflows/test.yml/badge.svg)](https://github.com/KusoKaihatsuSha/startup/actions/workflows/test.yml)
[![Go Coverage](https://github.com/KusoKaihatsuSha/startup/wiki/coverage.svg)](https://raw.githack.com/wiki/KusoKaihatsuSha/startup/coverage.html)

# Package `startup`
> Use the package's functionality to simplify the use of flags, environments and config file when starting the application.  
> Fill custom Golang `struct` from environments, flags or configuration file (JSON) with desired order.  
> Convenient for setting up an application for a `container` (`Docker`, `Podman`)

### Install
```shell
go get github.com/KusoKaihatsuSha/startup
```

### Update
```shell
go get -u
go mod tidy
```

### Usage
![img_3.png](image/usage.png)

```go
type CustomConf struct {
    VarStr string `json:"t-str" default:"abcd" flag:"t-str" env:"T_STR" help:"description"`
}
```
Order (`low` -> `high`) [`FILE` -> `ENV` -> `FLAG`]:  
```go
filledCustomConf := startup.Get[CustomConf](
	startup.FILE,
	startup.ENV,
	startup.FLAG,
	)
```

### Example 01
```go
// Some struct
type Configuration struct {
    TestString   string        `json:"test-string"   default:"abcd"          flag:"test-string"   env:"TEST_STR"      help:"string"`
    TestInt      int64         `json:"test-int"      default:"11"            flag:"test-int"      env:"TEST_INT"      help:"int"`
    TestDuration time.Duration `json:"test-duration" default:"1s"            flag:"test-duration" env:"TEST_DURATION" help:"duration"`
    TestBool     bool          `json:"test-bool"     default:"false"         flag:"test-bool"     env:"TEST_BOOL"     help:"bool"`
    TestFloat    float64       `json:"test-float"    default:"1"             flag:"test-float"    env:"TEST_FLOAT"    help:"float"`
    TestUint     uint64        `json:"test-uint"     default:"111"           flag:"test-uint"     env:"TEST_UINT"     help:"uint"`
}

func main() {
    // Emulated settings at startup
    os.Setenv("TEST_STR", "dcba")
    os.Setenv("TEST_UINT", "999")
    os.Setenv("TEST_DURATION", "11h")
    os.Args = append(
        os.Args,
        "-test-bool",
        "-test-uint=741",
    )
    
    // Get filled configs (ORDER: Environment -> Flags ).
    config := startup.GetForce[Configuration](
        startup.ENV,
        startup.FLAG,
    )
    
    // Print result => {dcba 11 11h0m0s true 1 741}
    fmt.Printf("%v\n", config)
}
```

### Example 02
```go
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
    get := startup.Get[Configuration](
        startup.FILE,
        startup.ENV,
        startup.FLAG,
    )
    
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

### Default validations (in tag `valid` inside annotation)
  - `tmp_file` - Check exist inside Temp folder and create if not exist  (string in struct)
  - `file` - Check exist the filepath and create if not exist (string in struct)
  - `url` - Check url is correct (string in struct)
  - `bool` - Parse Bool (bool in struct)
  - `int` - Parse int (int64 in struct)
  - `float` - Parse float (float64 in struct)
  - `duration` - Parse duration (time.Duration in struct)
  - `uuid` - Check uuid. Return new if not exist (string in struct)

### Caution
Default config filename:
  - `config.ini`

Flags are reserved:
  - `config`

Environments are reserved:
  - `CONFIG`

### Print `-h` or `-help` tag Example
```
Order of priority for settings (low -> high):
Config file (JSON) --> Environment --> Flags

  -config
        Configuration settings file
        Default value: config.ini
        Sample JSON config:
        {
          "startup_configuration_file": "config.ini"
        }
        Sample environment:     CONFIG=config.ini
        Sample flag value:      testee.exe -config=config.ini

  -test-bool
        bool
        Default value: true
        Sample JSON config:
        {
          "test-bool": true
        }
        Sample environment:     TEST_BOOL=true
        Sample(TRUE):     testee.exe -test-bool
        Sample(default):  testee.exe
        Sample(TRUE):   testee.exe -test-bool=true
        Sample(FALSE):  testee.exe -test-bool=false
        Sample(TRUE):   testee.exe -test-bool=1
        Sample(FALSE):  testee.exe -test-bool=0
        Sample(TRUE):   testee.exe -test-bool=t
        Sample(FALSE):  testee.exe -test-bool=f

  -test-duration
        duration
        Default value: 1s
        Sample JSON config:
        {
          "test-duration": 1000000000
        }
        Sample environment:     TEST_DURATION=1s
        Sample(Millisecond):    testee.exe -test-duration=1ms
        Sample(Second): testee.exe -test-duration=1s
        Sample(Minute): testee.exe -test-duration=1m
        Sample(Hour):   testee.exe -test-duration=1h
        Sample(Nanosecond):     testee.exe -test-duration=1ns
        Sample(Microsecond):    testee.exe -test-duration=1us
        Sample(1 Hour 2 Minutes and 3 Seconds): testee.exe -test-duration=1h2m3s
        Sample(111 Seconds):    testee.exe -test-duration=111

  -test-email
        email
        Default value: a@b.c
        Sample JSON config:
        {
          "test-email": "a@b.c"
        }
        Sample environment:     TEST_EMAIL=a@b.c
        Sample flag value:      testee.exe -test-email=a@b.c

  -test-float
        float
        Default value: 1
        Sample JSON config:
        {
          "test-float": 1
        }
        Sample environment:     TEST_FLOAT=1
        Sample flag value:      testee.exe -test-float=1.000000

  -test-int
        int
        Default value: 11
        Sample JSON config:
        {
          "test-int": 11
        }
        Sample environment:     TEST_INT=11
        Sample flag value:      testee.exe -test-int=11

  -test-ip
        ip
        Default value: 127.0.0.1
        Sample JSON config:
        {
          "test-ip": "127.0.0.1"
        }
        Sample environment:     TEST_IP=127.0.0.1

  -test-json
        json
        Default value: {default_001 default_002}
        Sample JSON config:
        {
          "test-json": {
            "param1": "default_001",
            "param2": "default_002"
          }
        }
        Sample environment:     TEST_JSON={"param1":"default_001","param2":"default_002"}

  -test-slice
        slice
        Default value: []
        Sample JSON config:
        {
          "test-slice": null
        }
        Sample environment:     TEST_SLICE=1,2,3,4;5;6

  -test-uint
        uint
        Default value: 111
        Sample JSON config:
        {
          "test-uint": 111
        }
        Sample environment:     TEST_UINT=111
        Sample flag value:      testee.exe -test-uint=111
```
