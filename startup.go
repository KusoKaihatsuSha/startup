package startup

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/KusoKaihatsuSha/startup/internal/helpers"
	"github.com/KusoKaihatsuSha/startup/internal/order"

	tags "github.com/KusoKaihatsuSha/startup/internal/tag"
	"github.com/KusoKaihatsuSha/startup/internal/validation"
)

var (
	this      any
	onceFlags = sync.Once{}
)

var DEBUG = false

const (
	// FLAG Get data from flags
	FLAG = order.FLAG

	// FILE Get data from json file
	FILE = order.FILE

	// ENV Get data from environments
	ENV = order.ENV

	// PreloadConfigEnvThenFlag - Get filepath config from environments and then flags
	// Default
	PreloadConfigEnvThenFlag = order.PreloadConfigEnvThenFlag

	// PreloadConfigFlagThenEnv - Get filepath config from flags and then environments
	PreloadConfigFlagThenEnv = order.PreloadConfigFlagThenEnv

	// PreloadConfigFlag - Get filepath config from flags
	PreloadConfigFlag = order.PreloadConfigFlag

	// PreloadConfigEnv - Get filepath config from environments
	PreloadConfigEnv = order.PreloadConfigEnv

	// NoPreloadConfig - Get filepath config file only from ordered stages
	NoPreloadConfig = order.NoPreloadConfig
)

// Concat structs
type temp[T any] struct {
	Stages []order.Stages
	tags.Tags
	CustomerConfiguration T
	Configuration         configuration
}

/*
Configuration consists of settings that are filled in at startup.
Default fields:
  - "Config" - filepath for config file
*/
type configuration struct {
	Config string `json:"startup_configuration_file" default:"config.ini" flag:"config" env:"CONFIG" help:"Configuration settings file" valid:"default_configuration_file"`
}

/*
AddValidation using for add custom validation

Example:

		// Custom struct. Struct will be implement in program with selected 'Stages' variable.
		type Test struct {
			NewValid []string `json:"new-valid" default:"new valid is default" flag:"valid" text:"-" valid:"test"`
		}
		// Custom type.
		type testValid string
		// Custom validation.
		var testValidation testValid = "test"
		// Custom method.
		func (o testValid) Valid(stringValue string, value any) (any, bool) {
			return []string{stringValue + "+++"}, true
		}

	    // add custom validation
		func MyFunc() {
			...
			startup.AddValidation(testValidation)
			...
			// Implement all types of configs (Json file -> Environment -> Flags).
			configurations := startup.Get[Test](startup.FILE, startup.ENV, startup.FLAG)
			// Test print.
			fmt.Println(configurations)
		}

Default validations:
  - `tmp_file` - Check exist inside Temp folder and create if not exist  (string in struct)
  - `file` - Check exist the filepath and create if not exist (string in struct)
  - `url` - Check url is correct (string in struct)
  - `bool` - Parse Bool (bool in struct)
  - `int` - Parse int (int64 in struct)
  - `float` - Parse float (float64 in struct)
  - `duration` - Parse duration (time.Duration in struct)
  - `uuid` - Check uuid. Return new if not exist (string in struct)

Caution:
flags are reserved:
  - config
*/
func AddValidation(value ...validation.Valid) {
	validation.Add(value...)
}

/*
GetForce will initialize scan the flags, environment and config-file with the right order:
  - order.FLAG - flag
  - order.FILE - config file
  - order.ENV - environment

Caution! flags are reserved:
  - config
*/
func GetForce[T any](stages ...order.Stages) T {
	// ---debug---
	if DEBUG {
		helpers.PrintDebug(stages...)
	}
	// ---debug---

	return get[T](stages...).CustomerConfiguration
}

func (t *temp[T]) prepare(config tags.Tags) *temp[T] {
	elements := reflect.ValueOf(&t.CustomerConfiguration).Elem()
	t.Tags = make(tags.Tags, elements.NumField())
	for ii := 0; ii < elements.NumField(); ii++ {
		name := elements.Type().Field(ii).Name
		t.Tags[name] = tags.Fill[T](name, t.Stages...)
	}
	for configTagName, configTagData := range config {
		t.Tags[configTagName] = configTagData
	}
	return t
}

func (t *temp[T]) preparePreload() *temp[T] {
	elements := reflect.ValueOf(&t.Configuration).Elem()
	t.Tags = make(tags.Tags, elements.NumField())
	for ii := 0; ii < elements.NumField(); ii++ {
		name := elements.Type().Field(ii).Name
		t.Tags[name] = tags.Fill[configuration](name, t.Stages...)
	}
	return t
}

func (t *temp[T]) fillPreload(fileExist bool) *temp[T] {
	for _, v := range t.Stages {
		switch v {
		case order.FLAG:
			// ---debug---
			if DEBUG && fileExist {
				printing := ""
				for _, arg := range os.Args {
					if strings.Contains(arg, "-config=") {
						printing = "\tinfo about config file:\tGet filepath from flag '-config'"
						break
					} else {
						printing = "\tinfo about config file:\tFlag '-config' with filepath not set"
					}
				}
				fmt.Println(printing)
			}
			// ---debug---

			t.flagNoParse()
		case order.FILE:
			t.conf()
		case order.ENV:
			// ---debug---
			if DEBUG && fileExist {
				if os.Getenv("CONFIG") != "" {
					fmt.Println("\tinfo about config file:\tGet filepath from environment 'CONFIG'")
				} else {
					fmt.Println("\tinfo about config file:\tEnvironment 'CONFIG' with filepath not set")
				}
			}
			// ---debug---

			t.env()
		}
	}
	return t
}

func (t *temp[T]) fill() *temp[T] {
	for _, v := range t.Stages {
		switch v {
		case order.FLAG:
			t.flag()
		case order.FILE:
			t.conf()
		case order.ENV:
			t.env()
		}
	}
	return t
}

func get[T any](stages ...order.Stages) temp[T] {
	fileExistInStages := helpers.FileConfExistInStages(stages...)
	preload := (&temp[T]{
		Stages:                helpers.PresetPreload(stages...),
		CustomerConfiguration: *new(T),
		Configuration:         configuration{},
	}).
		preparePreload().
		fillPreload(fileExistInStages).
		conf()

	// ---debug---
	if DEBUG && fileExistInStages {
		cfg := ""
		switch {
		case preload.Configuration.Config == "config.ini":
			if helpers.FileExist(preload.Configuration.Config) {
				cfg = "default config.ini"
			} else {
				cfg = "skipped default config.ini(not exist)"
			}
		case len(preload.Configuration.Config) > 0:
			if helpers.FileExist(preload.Configuration.Config) {
				// custom config file
				cfg = filepath.Base(preload.Configuration.Config)
			} else {
				cfg = "skipped config file(not exist)"
			}
		default:
			cfg = "not any config file"
		}
		fmt.Printf("FILE => %s\n", cfg)
	}
	// ---debug---

	load := (&temp[T]{
		Stages:                stages,
		CustomerConfiguration: *new(T),
		Configuration:         preload.Configuration,
	}).prepare(preload.Tags)
	load.
		fill().
		valid()

	// ---debug---
	if DEBUG {
		fmt.Printf("DATA => %v\n\n", load.CustomerConfiguration)
	}
	// ---debug---

	return *load
}

/*
Get will initialize scan the flags(one time), environment and config-file with the right order:
  - order.FLAG - flag
  - order.FILE - config file
  - order.ENV - environment

Caution! flags are reserved:
  - config
*/
func Get[T any](stages ...order.Stages) T {
	onceFlags.Do(
		func() {
			this = get[T](stages...)
		})
	return this.(temp[T]).CustomerConfiguration
}

func (t *temp[T]) dummy() *temp[T] {
	for _, v := range t.Tags {
		v.DummyFlags()
	}
	flag.Parse()
	return t
}

func (t *temp[T]) flag() *temp[T] {
	for _, v := range t.Tags {
		v.Flag()
	}
	return t.dummy()
}

func (t *temp[T]) flagNoParse() *temp[T] {
	for _, v := range t.Tags {
		v.Flag()
	}
	return t
}

func (t *temp[T]) env() *temp[T] {
	for _, v := range t.Tags {
		v.Env()
	}
	return t
}

// valid check info and make some correcting
func (t *temp[T]) valid() *temp[T] {
	for k, v := range t.Tags {
		field := reflect.ValueOf(&t.CustomerConfiguration).Elem().FieldByName(k)
		if field.CanSet() {
			field.Set(reflect.ValueOf(v.Valid()))
		}
	}
	return t
}

// conf get info from the configuration file
func (t *temp[T]) conf() *temp[T] {
	confFile := ""
	for _, f := range t.Tags["Config"].Flags {
		if f.Value.String() != "" {
			confFile = f.Value.String()
		}
	}
	if confFile != "" {
		reflect.ValueOf(&t.Configuration).Elem().FieldByName("Config").Set(reflect.ValueOf(t.Tags["Config"].Valid()))
		tmpConfig := helpers.SettingsFile(t.Configuration.Config)
		for _, v := range t.Tags {
			v.ConfigFile(tmpConfig)
		}
	}
	return t
}
