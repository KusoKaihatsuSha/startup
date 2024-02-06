package startup

import (
	"flag"
	"reflect"
	"sync"

	"github.com/KusoKaihatsuSha/startup/internal/helpers"
	"github.com/KusoKaihatsuSha/startup/internal/order"

	tags "github.com/KusoKaihatsuSha/startup/internal/tag"
	"github.com/KusoKaihatsuSha/startup/internal/validation"
)

var (
	this      any
	onceFlags sync.Once = sync.Once{}
)

const (
	// flags
	Flag = order.Flag
	// config file (json)
	File = order.File
	// environment
	Env = order.Env
)

// Concat structs
type temp[T any] struct {
	tags.Tags
	CustomerConfigurationOfUnknownStruct31415926535 T
	CustomerConfigurationFromFile31415926535        configuration
}

/*
Configuration consists of settings that are filled in at startup.
Default fields:
  - "Config" - filepath for config file
*/
type configuration struct {
	Config string `json:"startup_configuration_file" default:"config.ini" flag:"config" env:"CONFIG" help:"Configuration settings file" valid:"default_configuration_file"`
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
func GetForce[T any](stages ...order.Stages) T {
	return get[T](stages...).CustomerConfigurationOfUnknownStruct31415926535
}

func get[T any](stages ...order.Stages) temp[T] {
	load := temp[T]{}
	preload := temp[T]{}

	elementsPreload := reflect.ValueOf(&preload).Elem()
	preload.Tags = make(tags.Tags, elementsPreload.NumField())
	preload.CustomerConfigurationFromFile31415926535 = configuration{}
	for i := 0; i < elementsPreload.NumField(); i++ {
		name := elementsPreload.Type().Field(i).Name
		switch name {
		case "CustomerConfigurationFromFile31415926535":
			elements := reflect.ValueOf(&preload.CustomerConfigurationFromFile31415926535).Elem()
			for ii := 0; ii < elements.NumField(); ii++ {
				name := elements.Type().Field(ii).Name
				preload.Tags[name] = tags.Fill[configuration](name, stages...)
			}
		}
	}

	elements := reflect.ValueOf(&load).Elem()
	load.Tags = make(tags.Tags, elements.NumField())
	load.CustomerConfigurationOfUnknownStruct31415926535 = *new(T)
	load.CustomerConfigurationFromFile31415926535 = configuration{}
	for i := 0; i < elements.NumField(); i++ {
		name := elements.Type().Field(i).Name
		switch name {
		case "CustomerConfigurationOfUnknownStruct31415926535":
			elements := reflect.ValueOf(&load.CustomerConfigurationOfUnknownStruct31415926535).Elem()
			for ii := 0; ii < elements.NumField(); ii++ {
				name := elements.Type().Field(ii).Name
				load.Tags[name] = tags.Fill[T](name, stages...)
			}
		case "CustomerConfigurationFromFile31415926535":
			elements := reflect.ValueOf(&load.CustomerConfigurationFromFile31415926535).Elem()
			for ii := 0; ii < elements.NumField(); ii++ {
				name := elements.Type().Field(ii).Name
				load.Tags[name] = tags.Fill[configuration](name, stages...)
			}
		}
	}

	preload.flagNoParse()
	preload.env()
	load.CustomerConfigurationFromFile31415926535 = preload.CustomerConfigurationFromFile31415926535

	for _, v := range stages {
		switch v {
		case Flag:
			load.flag()
		case File:
			load.conf()
		case Env:
			load.env()
		}
	}
	load.valid()
	return load
}

/*
Get will initialize scan the flags(one time), environment and config-file with the right order:
  - order.Flag - only flag
  - order.File - only config file
  - order.Env - only environment

Caution! flags are reserved:
  - h
  - help
  - config
*/
func Get[T any](stages ...order.Stages) T {
	onceFlags.Do(
		func() {
			this = get[T](stages...)
		})
	return this.(temp[T]).CustomerConfigurationOfUnknownStruct31415926535
}

func (c *temp[T]) dummy() *temp[T] {
	for _, v := range c.Tags {
		v.DummyFlags()
	}
	flag.Parse()
	return c
}

func (c *temp[T]) flag() *temp[T] {
	for _, v := range c.Tags {
		v.Flag()
	}
	return c.dummy()
}

func (c *temp[T]) flagNoParse() *temp[T] {
	for _, v := range c.Tags {
		v.Flag()
	}
	return c
}

func (c *temp[T]) env() *temp[T] {
	for _, v := range c.Tags {
		v.Env()
	}
	return c
}

// valid check info and make some correcting
func (c *temp[T]) valid() *temp[T] {
	for k, v := range c.Tags {
		field := reflect.ValueOf(&c.CustomerConfigurationOfUnknownStruct31415926535).Elem().FieldByName(k)
		if field.CanSet() {
			field.Set(reflect.ValueOf(v.Valid()))
		}
	}
	return c
}

// conf get info from the configuration file
func (c *temp[T]) conf() *temp[T] {
	confFile := ""
	for _, f := range c.Tags["Config"].Flags {
		if f.Value.String() != "" {
			confFile = f.Value.String()
		}
	}
	if confFile != "" {
		reflect.ValueOf(&c.CustomerConfigurationFromFile31415926535).Elem().FieldByName("Config").Set(reflect.ValueOf(c.Tags["Config"].Valid()))
		tmpConfig := helpers.SettingsFile(c.CustomerConfigurationFromFile31415926535.Config)
		for _, v := range c.Tags {
			v.ConfigFile(tmpConfig)
		}
	}
	return c
}
