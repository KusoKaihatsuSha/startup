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
	"github.com/KusoKaihatsuSha/startup/internal/order"
)

var defArgs = os.Args

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

func runNoPreload() {
	startup.DEBUG = true

	startup.GetForce[Configuration](order.NoPreloadConfig)

	startup.GetForce[Configuration](
		order.FILE,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.FLAG,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.ENV,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.NoPreloadConfig,
		order.FILE,
		order.FLAG,
	)

	startup.GetForce[Configuration](
		order.FLAG,
		order.FILE,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.NoPreloadConfig,
		order.FILE,
		order.ENV,
	)

	startup.GetForce[Configuration](
		order.ENV,
		order.FILE,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.FLAG,
		order.ENV,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.ENV,
		order.FLAG,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.FILE,
		order.FLAG,
		order.ENV,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.FLAG,
		order.FILE,
		order.ENV,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.FILE,
		order.ENV,
		order.FLAG,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.ENV,
		order.FILE,
		order.FLAG,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.FLAG,
		order.ENV,
		order.FILE,
		order.NoPreloadConfig,
	)

	startup.GetForce[Configuration](
		order.ENV,
		order.FLAG,
		order.FILE,
		order.NoPreloadConfig,
	)
}

func run() {
	startup.DEBUG = true

	startup.GetForce[Configuration]()

	startup.GetForce[Configuration](
		order.FILE,
	)

	startup.GetForce[Configuration](
		order.FLAG,
	)

	startup.GetForce[Configuration](
		order.ENV,
	)

	startup.GetForce[Configuration](
		order.FILE,
		order.FLAG,
	)

	startup.GetForce[Configuration](
		order.FLAG,
		order.FILE,
	)

	startup.GetForce[Configuration](
		order.FILE,
		order.ENV,
	)

	startup.GetForce[Configuration](
		order.ENV,
		order.FILE,
	)

	startup.GetForce[Configuration](
		order.FLAG,
		order.ENV,
	)

	startup.GetForce[Configuration](
		order.ENV,
		order.FLAG,
	)

	startup.GetForce[Configuration](
		order.FILE,
		order.FLAG,
		order.ENV,
	)

	startup.GetForce[Configuration](
		order.FLAG,
		order.FILE,
		order.ENV,
	)

	startup.GetForce[Configuration](
		order.FILE,
		order.ENV,
		order.FLAG,
	)

	startup.GetForce[Configuration](
		order.ENV,
		order.FILE,
		order.FLAG,
	)

	startup.GetForce[Configuration](
		order.FLAG,
		order.ENV,
		order.FILE,
	)

	startup.GetForce[Configuration](
		order.ENV,
		order.FLAG,
		order.FILE,
	)
	startup.DEBUG = false
}

func Example_configOrderDef() {
	os.Args = defArgs
	run()

	// Output:
	// EMPTY - only defaults
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags]
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments]
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments]
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags]
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags] ↣ [Environments]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File] ↣ [Environments]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments] ↣ [Flags]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File] ↣ [Flags]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}

}

func Example_configOrderFlagEnv() {
	os.Args = defArgs
	err := os.Setenv("TEST_ORDER", "http://env")
	if err != nil {
		fmt.Println(err)
	}

	flagVal := "http://flag"
	os.Args = append(
		os.Args,
		"-to="+"http://flag-skip",
		"-test="+flagVal,
	)

	fmt.Printf(
		`
PRESETS
┌──────────────────────────────────────────────────────────────────────────
│                         flag : %s
│                  environment : %s
└──────────────────────────────────────────────────────────────────────────
`,
		flagVal,
		os.Getenv("TEST_ORDER"),
	)

	run()

	// Output:
	// PRESETS
	// ┌──────────────────────────────────────────────────────────────────────────
	// │                         flag : http://flag
	// │                  environment : http://env
	// └──────────────────────────────────────────────────────────────────────────
	// EMPTY - only defaults
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags]
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments]
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments]
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags]
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags] ↣ [Environments]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File] ↣ [Environments]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments] ↣ [Flags]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File] ↣ [Flags]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}

}

func Example_configOrderConfFlagEnv() {
	os.Args = defArgs
	testDataDef := `{
	             "test-order": "http://config.ini"
	     }`

	fileDef := helpers.CreateFile("config.ini")
	defer helpers.DeleteFile(fileDef)
	err := os.WriteFile(fileDef, []byte(testDataDef), 0755)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Setenv("TEST_ORDER", "http://env")
	if err != nil {
		fmt.Println(err)
	}

	flagVal := "http://flag"
	os.Args = append(
		os.Args,
		"-to="+"http://flag-skip",
		"-test="+flagVal,
	)

	fmt.Printf(
		`
PRESETS
┌──────────────────────────────────────────────────────────────────────────
│                         flag : %s
│                  environment : %s
│                  config file : %s
└──────────────────────────────────────────────────────────────────────────
`,
		flagVal,
		os.Getenv("TEST_ORDER"),
		"config.ini",
	)

	run()

	// Output:
	// PRESETS
	// ┌──────────────────────────────────────────────────────────────────────────
	// │                         flag : http://flag
	// │                  environment : http://env
	// │                  config file : config.ini
	// └──────────────────────────────────────────────────────────────────────────
	// EMPTY - only defaults
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {config.ini:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags]
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments]
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {config.ini:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {config.ini:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments]
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags]
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags] ↣ [Environments]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File] ↣ [Environments]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments] ↣ [Flags]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File] ↣ [Flags]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {config.ini:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags] ↣ [JSON File]
	//	info about config file:	Environment 'CONFIG' with filepath not set
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => default config.ini
	// DATA => {config.ini:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}

}

func Example_configOrderManyConfFlagEnv() {
	os.Args = defArgs
	testDataDef := `{
	             "test-order": "http://config.ini"
	     }`
	testDataFlag := `{
                 "test-order": "http://file-flag"
         }`
	testDataEnv := `{
                 "test-order": "http://file-env"
         }`
	fileDef := helpers.CreateFile("config.ini")
	defer helpers.DeleteFile(fileDef)
	err := os.WriteFile(fileDef, []byte(testDataDef), 0755)
	if err != nil {
		fmt.Println(err)
	}

	fileFlag := helpers.ValidTempFile("file-flag.confile")
	defer helpers.DeleteFile(fileFlag)
	err = os.WriteFile(fileFlag, []byte(testDataFlag), 0755)
	if err != nil {
		fmt.Println(err)
	}

	fileEnv := helpers.ValidTempFile("file-env.confile")
	defer helpers.DeleteFile(fileEnv)
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
PRESETS
┌──────────────────────────────────────────────────────────────────────────
│                         flag : %s
│                  environment : %s
│   config-file in environment : %s
│          config-file in flag : %s
│              def config-file : %s
└──────────────────────────────────────────────────────────────────────────
`,
		flagVal,
		os.Getenv("TEST_ORDER"),
		strings.ReplaceAll(strings.ReplaceAll(testDataEnv, " ", ""), string([]rune{10}), ""),
		strings.ReplaceAll(strings.ReplaceAll(testDataFlag, " ", ""), string([]rune{10}), ""),
		"config.ini",
	)

	run()

	// Output:
	// PRESETS
	// ┌──────────────────────────────────────────────────────────────────────────
	// │                         flag : http://flag
	// │                  environment : http://env
	// │   config-file in environment : {"test-order":"http://file-env"}
	// │          config-file in flag : {"test-order":"http://file-flag"}
	// │              def config-file : config.ini
	// └──────────────────────────────────────────────────────────────────────────
	// EMPTY - only defaults
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags]
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments]
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments]
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags]
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments] ↣ [Flags]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File] ↣ [Flags]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}

}

func Example_configOrderOtherConfFlagEnv() {
	os.Args = defArgs
	testDataFlag := `{
                 "test-order": "http://file-flag"
         }`
	testDataEnv := `{
                 "test-order": "http://file-env"
         }`

	fileFlag := helpers.ValidTempFile("file-flag.confile")
	defer helpers.DeleteFile(fileFlag)
	err := os.WriteFile(fileFlag, []byte(testDataFlag), 0755)
	if err != nil {
		fmt.Println(err)
	}

	fileEnv := helpers.ValidTempFile("file-env.confile")
	defer helpers.DeleteFile(fileEnv)
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
PRESETS
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

	run()

	// Output:
	// PRESETS
	// ┌──────────────────────────────────────────────────────────────────────────
	// │                         flag : http://flag
	// │                  environment : http://env
	// │   config-file in environment : {"test-order":"http://file-env"}
	// │          config-file in flag : {"test-order":"http://file-flag"}
	// └──────────────────────────────────────────────────────────────────────────
	// EMPTY - only defaults
	// DATA => {def:81 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags]
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments]
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments]
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags]
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments] ↣ [Flags]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File] ↣ [Flags]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-env:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Get filepath from flag '-config'
	// FILE => file-flag.confile
	// DATA => {file-flag:80 email [] 11 1s true 1 111 127.0.0.1 { } {default_001 default_002}}

}

func Example_configOrigin() {
	os.Args = defArgs

	startup.AddValidation(testValidation)
	startup.AddValidation(emailValidation)
	startup.AddValidation(maxValidation)

	err := os.Setenv("TEST_SLICE", "999")
	if err != nil {
		fmt.Println(err)
	}

	os.Args = append(
		os.Args,
		"-test-email="+"flag@email.post",
		"-test-int="+"100",
		"-test-slice="+"100,200,300",
		"-test-json="+"{\"param1\":\"new_003\",\"param2\":\"new_004\"}",
	)

	testDataEnv := `{
                 "test-email": "fileenv@mail.com",
                 "test-slice": "18,19,20"
         }`

	fileEnv := helpers.ValidTempFile("test.confile")
	defer helpers.DeleteFile(fileEnv)
	err = os.WriteFile(fileEnv, []byte(testDataEnv), 0755)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Setenv("CONFIG", fileEnv)
	if err != nil {
		fmt.Println(err)
	}

	run()
	runNoPreload()

	// Output:
	// EMPTY - only defaults
	// DATA => {def:81 email@example.com [1 2 3 4 5 6] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {def:81 fileenv@mail.com [18 19 20] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags]
	// DATA => {def:81 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments]
	// DATA => {env:80 email@example.com [999] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {def:81 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {def:81 fileenv@mail.com [18 19 20] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [999] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [18 19 20] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments]
	// DATA => {env:80 flag@email.post [999] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags]
	// DATA => {env:80 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Flags] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 flag@email.post [999] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [JSON File] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [999] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[JSON File] ↣ [Environments] ↣ [Flags]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [JSON File] ↣ [Flags]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Flags] ↣ [Environments] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [18 19 20] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// PreloadConfigEnvThenFlag/Default - Preload find config in Env then Flag
	// Structure filling order:
	//	[Environments] ↣ [Flags] ↣ [JSON File]
	//	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [18 19 20] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	// DATA => {def:81 email@example.com [1 2 3 4 5 6] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[JSON File] ↣ FILE => not any config file
	// DATA => {def:81 email@example.com [1 2 3 4 5 6] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Flags] ↣ DATA => {def:81 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Environments] ↣ DATA => {env:80 email@example.com [999] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	// [JSON File] ↣ [Flags]
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Flags] ↣ [JSON File] ↣ 	info about config file:	Flag '-config' with filepath not set
	// FILE => not any config file
	// DATA => {def:81 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	// [JSON File] ↣ [Environments]
	//	info about config file:	Get filepath from environment 'CONFIG'
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [999] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Environments] ↣ [JSON File] ↣ 	info about config file:	Get filepath from environment 'CONFIG'
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [18 19 20] 10 1s true 1 111 127.0.0.1 { } {default_001 default_002}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Flags] ↣ [Environments] ↣ DATA => {env:80 flag@email.post [999] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Environments] ↣ [Flags] ↣ DATA => {env:80 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[JSON File] ↣ [Flags] ↣ [Environments] ↣ 	info about config file:	Flag '-config' with filepath not set
	//	info about config file:	Get filepath from environment 'CONFIG'
	// FILE => test.confile
	// DATA => {env:80 flag@email.post [999] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Flags] ↣ [JSON File] ↣ [Environments] ↣ 	info about config file:	Flag '-config' with filepath not set
	//	info about config file:	Get filepath from environment 'CONFIG'
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [999] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[JSON File] ↣ [Environments] ↣ [Flags] ↣ 	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Environments] ↣ [JSON File] ↣ [Flags] ↣ 	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 flag@email.post [100 200 300] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Flags] ↣ [Environments] ↣ [JSON File] ↣ 	info about config file:	Flag '-config' with filepath not set
	//	info about config file:	Get filepath from environment 'CONFIG'
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [18 19 20] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}
	//
	// NoPreloadConfig - Disable find config in other places. Only in list
	// Structure filling order:
	//	[Environments] ↣ [Flags] ↣ [JSON File] ↣ 	info about config file:	Get filepath from environment 'CONFIG'
	//	info about config file:	Flag '-config' with filepath not set
	// FILE => test.confile
	// DATA => {env:80 fileenv@mail.com [18 19 20] 10 1s true 1 111 127.0.0.1 { } {new_003 new_004}}

}
