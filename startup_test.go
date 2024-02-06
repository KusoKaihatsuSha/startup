package startup_test

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/KusoKaihatsuSha/startup"
	"github.com/KusoKaihatsuSha/startup/internal/helpers"
)

type Configuration struct {
	TestOrder         string            `json:"test-order"    default:"http://def:81" flag:"to,test,order" env:"TEST_ORDER"    help:"order" valid:"url"`
	TestEmail         string            `json:"test-email"    default:"email"         flag:"test-email"    env:"TEST_EMAIL"    help:"email" valid:"email"`
	TestSlice         []string          `json:"test-slice"    default:"1,2,3,4;5;6"   flag:"test-slice"    env:"TEST_SLICE"    help:"slice" valid:"test"`
	TestInt           int64             `json:"test-int"      default:"11"            flag:"test-int"      env:"TEST_INT"      help:"int"   valid:"max10"`
	TestDuration      time.Duration     `json:"test-duration" default:"1s"            flag:"test-duration" env:"TEST_DURATION" help:"duration"`
	TestBool          bool              `json:"test-bool"     default:"true"          flag:"test-bool"     env:"TEST_BOOL"     help:"bool"`
	TestFloat         float64           `json:"test-float"    default:"1"             flag:"test-float"    env:"TEST_FLOAT"    help:"float"`
	TestUint          uint64            `json:"test-uint"     default:"111"           flag:"test-uint"     env:"TEST_UINT"     help:"uint"`
	TestIP            net.IP            `json:"test-ip"       default:"127.0.0.1"     flag:"test-ip"       env:"TEST_IP"       help:"ip"`
	TestNonMethodJSON TestNonMethodJSON `json:"test-nmj"      default:""              flag:"test-nmj"      env:"TEST_NMJ"      help:"non method json"`
	TestJSON          TestJSON          `json:"test-json"     default:"{\"param1\":\"default_001\",\"param2\":\"default_002\"}" flag:"test-json" env:"TEST_JSON" help:"json"`
}

type TestNonMethodJSON struct {
	P1 string `json:"param1"`
	P2 string `json:"param2"`
}

type TestJSON struct {
	P1 string `json:"param1"`
	P2 string `json:"param2"`
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *TestJSON) UnmarshalText(text []byte) error {
	type Tmp *TestJSON
	if err := json.Unmarshal(text, Tmp(t)); err != nil {
		return err
	}
	return nil
}

// custom email validation
type emailValid string

var emailValidation emailValid = "email"

func (o emailValid) Valid(stringValue string, value any) (any, bool) {
	matchString, err := regexp.MatchString("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])", stringValue)
	if matchString && err == nil {
		return value, true
	}
	return stringValue + "@example.com", true
}

// custom slice validation
type testValid string

var testValidation testValid = "test"

func (o testValid) Valid(stringValue string, value any) (any, bool) {
	// 'value' can be used for Type check
	newValue := strings.ReplaceAll(stringValue, ";", ",")
	return strings.Split(newValue, ","), true
}

// custom max10 validation
type maxValid string

var maxValidation maxValid = "max10"

func (o maxValid) Valid(stringValue string, value any) (any, bool) {
	if value.(int64) > 10 {
		return int64(10), true
	}
	return value, true
}

func Example_configOrder() {
	testDataFlag := `{
                 "test-order": "http://file-flag"
         }`
	testDataEnv := `{
                 "test-order": "http://file-env"
         }`
	fileFlag := helpers.CreateFile("")
	defer helpers.DeleteFile(fileFlag)
	err := os.WriteFile(fileFlag, []byte(testDataFlag), 0755)
	if err != nil {
		fmt.Println(err)
	}
	fileEnv := helpers.CreateFile("")
	defer helpers.DeleteFile(fileFlag)
	err = os.WriteFile(fileEnv, []byte(testDataEnv), 0755)

	if err != nil {
		fmt.Println(err)
	}

	err = os.Setenv("TEST_ORDER", "http://env")
	if err != nil {
		fmt.Println(err)
	}
	err = os.Setenv("CONFIG", fileEnv)
	if err != nil {
		fmt.Println(err)
	}
	flagVal := "http://flag"
	os.Args = append(
		os.Args,
		"-config="+fileFlag,
		"-to="+"http://flag-skip",
		"-test="+flagVal,
	)

	fmt.Printf(
		`
┌──────────────────────────────────────────────────────────────────────────
│                         flag : %s
│                  environment : %s
│   config-file in environment : %s
│          config-file in flag : %s
└──────────────────────────────────────────────────────────────────────────
`,
		flagVal,
		os.Getenv("TEST_ORDER"),
		strings.ReplaceAll(strings.ReplaceAll(testDataEnv, " ", ""), string([]rune{10}), ""),
		strings.ReplaceAll(strings.ReplaceAll(testDataFlag, " ", ""), string([]rune{10}), ""),
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"file ↣ "+
			""+
			"",
		startup.GetForce[Configuration](
			startup.File,
		).TestOrder,
		"| file config.ini not exist and default only",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"flag ↣ "+
			""+
			"",
		startup.GetForce[Configuration](
			startup.Flag,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"env ↣ "+
			""+
			"",
		startup.GetForce[Configuration](
			startup.Env,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"file ↣ "+
			"flag ↣ "+
			"",
		startup.GetForce[Configuration](
			startup.File,
			startup.Flag,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"flag ↣ "+
			"file ↣ "+
			"",
		startup.GetForce[Configuration](
			startup.Flag,
			startup.File,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"file ↣ "+
			"env ↣"+
			"",
		startup.GetForce[Configuration](
			startup.File,
			startup.Env,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"env ↣ "+
			"file ↣ "+
			"",
		startup.GetForce[Configuration](
			startup.Env,
			startup.File,
		).TestOrder,
		"|",
	)

	errEnv := os.Unsetenv("CONFIG")
	helpers.ToLog(errEnv)
	fmt.Printf(
		"ORDER %s %s %s\n",
		"env ↣ "+
			"file ↣ "+
			"",
		startup.GetForce[Configuration](
			startup.Env,
			startup.File,
		).TestOrder,
		"| not filepath in env CONFIG",
	)
	errEnv = os.Setenv("CONFIG", fileEnv)
	helpers.ToLog(errEnv)
	fmt.Printf(
		"ORDER %s %s %s\n",
		"flag ↣ "+
			"env ↣ "+
			"",
		startup.GetForce[Configuration](
			startup.Flag,
			startup.Env,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"env ↣ "+
			"flag ↣ "+
			"",
		startup.GetForce[Configuration](
			startup.Env,
			startup.Flag,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"file ↣ "+
			"flag ↣ "+
			"env ↣ ",
		startup.GetForce[Configuration](
			startup.File,
			startup.Flag,
			startup.Env,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"flag ↣ "+
			"file ↣ "+
			"env ↣ ",
		startup.GetForce[Configuration](
			startup.Flag,
			startup.File,
			startup.Env,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"file ↣ "+
			"env ↣ "+
			"flag ↣ ",
		startup.GetForce[Configuration](
			startup.File,
			startup.Env,
			startup.Flag,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"env ↣ "+
			"file ↣ "+
			"flag ↣ ",
		startup.GetForce[Configuration](
			startup.Env,
			startup.File,
			startup.Flag,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"flag ↣ "+
			"env ↣ "+
			"file ↣ ",
		startup.GetForce[Configuration](
			startup.Flag,
			startup.Env,
			startup.File,
		).TestOrder,
		"|",
	)

	fmt.Printf(
		"ORDER %s %s %s\n",
		"env ↣ "+
			"flag ↣ "+
			"file ↣ ",
		startup.GetForce[Configuration](
			startup.Env,
			startup.Flag,
			startup.File,
		).TestOrder,
		"|",
	)

	// Output:
	// ┌──────────────────────────────────────────────────────────────────────────
	// │                         flag : http://flag
	// │                  environment : http://env
	// │   config-file in environment : {"test-order":"http://file-env"}
	// │          config-file in flag : {"test-order":"http://file-flag"}
	// └──────────────────────────────────────────────────────────────────────────
	// ORDER file ↣  def:81 | file config.ini not exist and default only
	// ORDER flag ↣  flag:80 |
	// ORDER env ↣  env:80 |
	// ORDER file ↣ flag ↣  flag:80 |
	// ORDER flag ↣ file ↣  file-flag:80 |
	// ORDER file ↣ env ↣ env:80 |
	// ORDER env ↣ file ↣  file-env:80 |
	// ORDER env ↣ file ↣  env:80 | not filepath in env CONFIG
	// ORDER flag ↣ env ↣  env:80 |
	// ORDER env ↣ flag ↣  flag:80 |
	// ORDER file ↣ flag ↣ env ↣  env:80 |
	// ORDER flag ↣ file ↣ env ↣  env:80 |
	// ORDER file ↣ env ↣ flag ↣  flag:80 |
	// ORDER env ↣ file ↣ flag ↣  flag:80 |
	// ORDER flag ↣ env ↣ file ↣  file-env:80 |
	// ORDER env ↣ flag ↣ file ↣  file-flag:80 |

}

func Example_configTypes() {
	startup.AddValidation(testValidation)
	startup.AddValidation(emailValidation)
	startup.AddValidation(maxValidation)

	os.Args = append(
		os.Args,
		"-test-email="+"my@email.post",
		"-test-int="+"100",
		"-test-json="+"{\"param1\":\"default_003\",\"param2\":\"default_004\"}",
	)

	get := startup.GetForce[Configuration](
		startup.File,
		startup.Env,
		startup.Flag,
	)

	fmt.Printf(
		"custom email ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
		get.TestEmail,
	)

	fmt.Printf(
		"custom slice ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
		get.TestSlice,
	)

	fmt.Printf(
		"custom integer max 10 ↣ TYPE [%[1]T] VALUE [%[1]v]\n",
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

	fmt.Printf(
		"json non method ↣ TYPE [%[1]T] VALUE [%#[1]v]\n",
		get.TestNonMethodJSON,
	)

	// Output:
	// custom email ↣ TYPE [string] VALUE [my@email.post]
	// custom slice ↣ TYPE [[]string] VALUE [[1 2 3 4 5 6]]
	// custom integer max 10 ↣ TYPE [int64] VALUE [10]
	// duration ↣ TYPE [time.Duration] VALUE [1s]
	// boolean ↣ TYPE [bool] VALUE [true]
	// float ↣ TYPE [float64] VALUE [1]
	// uint ↣ TYPE [uint64] VALUE [111]
	// ip ↣ TYPE [net.IP] VALUE [127.0.0.1]
	// json unmarshal ↣ TYPE [startup_test.TestJSON] VALUE [startup_test.TestJSON{P1:"default_003", P2:"default_004"}]
	// json non method ↣ TYPE [startup_test.TestNonMethodJSON] VALUE [startup_test.TestNonMethodJSON{P1:"", P2:""}]

}
