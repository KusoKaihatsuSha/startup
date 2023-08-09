package startup

import (
	"flag"
	"os"
	"reflect"
	"sync"

	"github.com/KusoKaihatsuSha/startup/internal/helpers"
	tags "github.com/KusoKaihatsuSha/startup/internal/tag"
	"github.com/KusoKaihatsuSha/startup/internal/validation"
)

var (
	this      any
	onceFlags sync.Once = sync.Once{}
)

type stages int

const (
	// only flags
	Flag stages = iota
	// only config file (json)
	File
	// only environment
	Env
	// config file -> environment
	FileEnv
	// config file -> flags
	FileFlag
	// environment -> flags
	EnvFlag
	// environment -> config file
	EnvFile
	// flags -> environment
	FlagEnv
	// flags -> config file
	FlagFile
	// config file -> environment -> flags
	FileEnvFlag
	// config file -> flags -> environment
	FileFlagEnv
	Default
	// environment -> flags -> config file
	EnvFlagFile
	// environment -> config file -> flags
	EnvFileFlag
	// flags -> environment -> config file
	FlagEnvFile
	// flags -> environment -> config file
	FlagFileEnv
)

// Concat structs
type temp[T any] struct {
	tags.TagsInfo
	CustomerConfigurationOfUnknownStruct31415926535 T
	CustomerConfigurationFromFile31415926535        configuration
}

/*
Configuration consists of settings that are filled in at startup.
Default fields:
  - "Config" - filepath for config file
  - "Help" - print flags info
*/
type configuration struct {
	Config string `default:"" flag:"config" env:"CONFIG" text:"Config file" valid:"file"`
	// Usage -> 'bin.app -h=true' because all flags reading as String
	Help bool `default:"false" flag:"h,help" text:"show help info" valid:"bool"`
}

// Stages are parameters of Init function, sequence read type of settings on startup

/*
AddValidation using for add custom validation

Example:

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
			// Implement all types of configs (Environment -> Flags -> Json file).
			configurations := startup.InitForce[Test](startup.EnvFlagFile)
			// Test print.
			fmt.Println(configurations)
		}

Default validations:
  - 'tmp_file' - As 'file', but if empty returm file from Temp folder  (string in struct)
  - 'file' - Check exist the filepath (string in struct)
  - 'url' - Check url is correct (string in struct)
  - 'bool' - Parse Bool (bool in struct)
  - 'int' - Parse int (int64 in struct)
  - 'float' - Parse float (float64 in struct)
  - 'duration' - Parse duration (time.Duration in struct)
  - 'uuid' - Check uuid. Return new if not exist (string in struct)

Caution:
flags are reserved:
  - h
  - help
  - config
*/
func AddValidation(value ...validation.Valid) {
	validation.Add(value...)
}

// GetForce will initialize scan the flags(every time), environment and config-file with the right order.
func GetForce[T any](stages stages) T {
	onceFlags = sync.Once{}
	return Get[T](stages)
}

/*
Get will initialize scan the flags(one time), environment and config-file with the right order:
  - Flag - only flag
  - File - only config file
  - Env - only environment
  - FileEnv
  - FileFlag
  - EnvFlag
  - EnvFile
  - FlagEnv
  - FlagFile
  - FileEnvFlag
  - FileFlagEnv = Default
  - EnvFlagFile
  - EnvFileFlag
  - FlagEnvFile
  - FlagFileEnv

Caution! flags are reserved:
  - h
  - help
  - config
*/
func Get[T any](stages stages) T {
	onceFlags.Do(
		func() {
			t := temp[T]{}
			elements := reflect.ValueOf(&t).Elem()
			t.TagsInfo = make(tags.TagsInfo, elements.NumField())
			t.CustomerConfigurationOfUnknownStruct31415926535 = *new(T)
			t.CustomerConfigurationFromFile31415926535 = configuration{}
			for i := 0; i < elements.NumField(); i++ {
				name := elements.Type().Field(i).Name
				switch name {
				case "CustomerConfigurationOfUnknownStruct31415926535":
					elements := reflect.ValueOf(&t.CustomerConfigurationOfUnknownStruct31415926535).Elem()
					for ii := 0; ii < elements.NumField(); ii++ {
						name := elements.Type().Field(ii).Name
						t.TagsInfo[name] = tags.Fill[T](name)
					}
				case "CustomerConfigurationFromFile31415926535":
					elements := reflect.ValueOf(&t.CustomerConfigurationFromFile31415926535).Elem()
					for ii := 0; ii < elements.NumField(); ii++ {
						name := elements.Type().Field(ii).Name
						t.TagsInfo[name] = tags.Fill[configuration](name)
					}
				}
			}
			t.dummy()

			t.flag()
			t.env()

			switch stages {
			case Flag:
				t.flag()
			case File:
				t.conf()
			case Env:
				t.env()
			case FileEnv:
				t.conf()
				t.env()
			case FileFlag:
				t.conf()
				t.flag()
			case EnvFlag:
				t.env()
				t.flag()
			case EnvFile:
				t.env()
				t.conf()
			case FlagEnv:
				t.flag()
				t.env()
			case FlagFile:
				t.flag()
				t.conf()
			case FileEnvFlag:
				t.conf()
				t.env()
				t.flag()
			case FileFlagEnv, Default:
				t.conf()
				t.flag()
				t.env()
			case EnvFlagFile:
				t.env()
				t.flag()
				t.conf()
			case EnvFileFlag:
				t.env()
				t.conf()
				t.flag()
			case FlagEnvFile:
				t.flag()
				t.env()
				t.conf()
			case FlagFileEnv:
				t.flag()
				t.conf()
				t.env()
			}
			t.valid()
			this = t
			if t.CustomerConfigurationFromFile31415926535.Help {
				flag.PrintDefaults()
				os.Exit(0)
			}
		})
	return this.(temp[T]).CustomerConfigurationOfUnknownStruct31415926535
}

func (c *temp[T]) dummy() {
	for _, v := range c.TagsInfo {
		v.DummyFlags()
	}
	flag.Parse()
}

func (c *temp[T]) flag() {
	for _, v := range c.TagsInfo {
		v.Flag()
	}
}

func (c *temp[T]) env() {
	for _, v := range c.TagsInfo {
		v.Env()
	}
}

// Valid check info and make some correcting
func (c *temp[T]) valid() {
	for k, v := range c.TagsInfo {
		field := reflect.ValueOf(&c.CustomerConfigurationOfUnknownStruct31415926535).Elem().FieldByName(k)
		if field.CanSet() {
			field.Set(reflect.ValueOf(v.Valid()))
		}
	}
}

// Conf get info from the configuration file
func (c *temp[T]) conf() {
	confFile := ""
	for _, f := range c.TagsInfo["Config"].Flags {
		if f.Value.String() != "" {
			confFile = f.Value.String()
		}
	}
	if confFile != "" {
		reflect.ValueOf(&c.CustomerConfigurationFromFile31415926535).Elem().FieldByName("Config").Set(reflect.ValueOf(c.TagsInfo["Config"].Valid()))
		tmpConfig := helpers.SettingsFile(c.CustomerConfigurationFromFile31415926535.Config)
		for _, v := range c.TagsInfo {
			v.ConfigFile(tmpConfig)
		}
	}
	help := reflect.ValueOf(&c.CustomerConfigurationFromFile31415926535).Elem().FieldByName("Help")
	if help.CanSet() {
		help.Set(reflect.ValueOf(c.TagsInfo["Help"].Valid()))
	}
}
