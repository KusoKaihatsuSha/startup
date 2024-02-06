package tags

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/fatih/color"

	"github.com/KusoKaihatsuSha/startup/internal/helpers"
	"github.com/KusoKaihatsuSha/startup/internal/order"
	"github.com/KusoKaihatsuSha/startup/internal/validation"
)

const (
	defaultTag     = "default"
	flagTag        = "flag"
	EnvironmentTag = "env"
	helpTextTag    = "help"
	validationTag  = "valid"
	jsonTag        = "json"

	testTrigger = "-test."
)

var durationTypesMap = map[string]string{
	"ns": "Nanosecond",
	"us": "Microsecond",
	"ms": "Millisecond",
	"s":  "Second",
	"m":  "Minute",
	"h":  "Hour",
}

// Tags consist information when reading/valid configs
type Tags map[string]Tag

// Tag of TagInfo store tags Config struct
type Tag struct {
	ConfigFile func(map[string]any) Tag
	Valid      func() any
	DummyFlags func() Tag
	Env        func() Tag
	Flag       func() Tag
	FlagSet    *flag.FlagSet
	Flags      map[string]*flag.Flag
	Name       string
	Annotation
}

// Annotation - store general values
type Annotation struct {
	valid string
	json  string
	desc  string
	env   string
	def   string
}

type storage struct {
	Flag        bool
	Store       any
	StoreString string
	Type        reflect.StructField
	Default     string
	Desc        string
	Env         string
	JSON        string
	Name        string
}

// Set 'flag' interface implementation
func (s *storage) Set(value string) error {
	var err error
	s.Store, err = comparatorStringType(nil, s.Type, value, s.Flag)
	s.StoreString = value
	return err
}

// String - Stringer interface implementation
func (s *storage) String() string {
	return fmt.Sprint(s.Store)
}

// Fill - filling the 'tag'
func Fill[T any](field string, order ...order.Stages) Tag {
	var t T
	tagData := Tag{}

	tagData.Flags = make(map[string]*flag.Flag)

	tagData.Name = field
	tagData.FlagSet = flag.NewFlagSet("", flag.ContinueOnError)
	fieldByName, exist := reflect.TypeOf(t).FieldByName(tagData.Name)
	if !exist {
		return Tag{}
	}
	if v, ok := fieldByName.Tag.Lookup(helpTextTag); ok {
		tagData.desc = v
	}
	if v, ok := fieldByName.Tag.Lookup(defaultTag); ok {
		tagData.def = v
	}
	if v, ok := fieldByName.Tag.Lookup(EnvironmentTag); ok {
		tagData.env = strings.ToUpper(v)
	}
	if v, ok := fieldByName.Tag.Lookup(jsonTag); ok {
		tagData.json = v
	}
	if v, ok := fieldByName.Tag.Lookup(flagTag); ok {
		fv := new(storage)
		fv.Type = fieldByName
		fv.Flag = true
		fv.Default = tagData.def
		fv.Env = tagData.env
		fv.JSON = tagData.json
		fv.Desc = tagData.desc
		fv.Name = v
		err := fv.Set(tagData.def)
		helpers.ToLog("", fmt.Sprintf("set flag data '%s' %v", tagData.def, err)) // skip info and error parse
		for _, fTag := range strings.Split(v, ",") {
			tagData.FlagSet.Var(fv, fTag, tagData.desc)
			fl := tagData.FlagSet.Lookup(fTag)
			tagData.Flags[fl.Name] = fl
		}
	}

	if v, ok := fieldByName.Tag.Lookup(validationTag); ok {
		tagData.valid = v
	}

	tagData.ConfigFile = func(m map[string]any) Tag {
		for k, v := range m {
			if tagData.json == k {
				for _, f := range tagData.Flags {
					f.DefValue = fmt.Sprintf("%v", v)
					err := f.Value.Set(f.DefValue)
					helpers.ToLog(err, fmt.Sprintf("set flag data '%s' error", f.DefValue))
				}
			}
		}
		return tagData
	}
	tagData.DummyFlags = func() Tag {
		for _, v := range os.Args {
			if strings.Contains(v, testTrigger) {
				var _ = func() bool {
					testing.Init()
					return true
				}()
				break
			}
		}
		for _, f := range tagData.Flags {
			if ff := flag.Lookup(f.Name); ff == nil {
				flag.Var(f.Value, f.Name, f.Usage)
			}
		}
		return tagData
	}
	tagData.Env = func() Tag {
		env, ok := os.LookupEnv(tagData.env)
		if ok {
			for _, f := range tagData.Flags {
				err := f.Value.Set(env)
				helpers.ToLog(err, fmt.Sprintf("set flag data '%s' error", env))
				break
			}
		}
		return tagData
	}

	tagData.Flag = func() Tag {
		def := tagData.FlagSet.Output()
		tagData.FlagSet.SetOutput(io.Discard)
		for _, arg := range os.Args {
			err := tagData.FlagSet.Parse([]string{arg})
			helpers.ToLogWithType(err, helpers.LogNull)
		}
		tagData.FlagSet.SetOutput(def)
		return tagData
	}

	flag.Usage = func() {
		PrintDefaults(flag.CommandLine, order...)
	}

	tagData.Valid = func() any {
		var stringValueType any
		var stringValueTypeString string
		for _, f := range tagData.Flags {
			stringValueType = f.Value.(*storage).Store
			stringValueTypeString = f.Value.(*storage).StoreString
			break
		}
		for _, v := range validation.Valids {
			if tagData.valid == fmt.Sprint(v) {
				if ret, ok := v.Valid(stringValueTypeString, stringValueType); ok {
					return ret
				}
			}

		}
		return stringValueType
	}
	return tagData
}

// PrintDefaults - printing help
func PrintDefaults(f *flag.FlagSet, o ...order.Stages) {
	yellow := color.New(color.FgYellow).SprintfFunc()
	red := color.New(color.FgRed).SprintfFunc()
	cyan := color.New(color.FgCyan).SprintfFunc()

	var def strings.Builder
	def.WriteString("Order of priority for settings (low -> high): \n")
	for k, v := range o {
		if k > 0 {
			def.WriteString(" --> ")
		}
		switch v {
		case order.Flag:
			def.WriteString("Flags")
		case order.File:
			def.WriteString("Config file (JSON)")
		case order.Env:
			def.WriteString("Environment")
		}
	}
	_, err := fmt.Fprint(f.Output(), yellow("%s \n\n", def.String()))
	helpers.ToLog(err, fmt.Sprintf("print data '%s' error", def.String()))
	f.VisitAll(func(lf *flag.Flag) {
		switch lf.Value.(type) {
		case *storage:
			// dummy
		default:
			return
		}
		var b strings.Builder
		_, errPrintf := fmt.Fprintf(&b, "  %s%s", red("%s", "-"), red("%s", lf.Name))
		helpers.ToLog(errPrintf, "print data error")
		// Two spaces before -; see next two comments.
		name, usage := flag.UnquoteUsage(lf)
		if len(name) > 0 {
			b.WriteString(" ")
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if b.Len() <= 4 { // space, space, '-', 'x'.
			b.WriteString("\t")
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			b.WriteString("\n    \t")
		}
		vvv := comparatorInfo(lf.Value.(*storage), o...)
		b.WriteString(strings.ReplaceAll(cyan("%s", usage), "\n", "\n    \t"))
		b.WriteString("\n    \t")
		b.WriteString(fmt.Sprintf("Default value: %v\n    \t", lf.DefValue))
		b.WriteString(strings.ReplaceAll(yellow("%s", vvv), "\n", "\n    \t"))
		_, errPrint := fmt.Fprint(f.Output(), b.String(), "\n")
		helpers.ToLog(errPrint, "print data error")
	})
}

func comparatorStringType(flagType flag.Value, reflectType reflect.StructField, stringValue any, onlyForMarshaller bool) (any, error) {
	reflectTypeName := strings.ToLower(strings.TrimSpace(reflectType.Type.Name()))
	switch reflectTypeName {
	case "string":
		return stringValue, nil
	case "bool":
		return helpers.ValidBool(fmt.Sprintf("%v", stringValue)), nil
	case "duration":
		return helpers.ValidDuration(fmt.Sprintf("%v", stringValue)), nil
	case "int8", "int16", "int32", "int64", "rune":
		return helpers.ValidInt(fmt.Sprintf("%v", stringValue)), nil
	case "int":
		return int(helpers.ValidInt(fmt.Sprintf("%v", stringValue))), nil
	case "uint8", "uint16", "uint32", "uint64":
		return helpers.ValidUint(fmt.Sprintf("%v", stringValue)), nil
	case "uint":
		return uint(helpers.ValidUint(fmt.Sprintf("%v", stringValue))), nil
	case "float32", "float64":
		return helpers.ValidFloat(fmt.Sprintf("%v", stringValue)), nil
	default:
		if method, ok := reflect.PointerTo(reflectType.Type).MethodByName("UnmarshalText"); ok {
			in := make([]reflect.Value, method.Type.NumIn())
			yyy := reflect.New(reflectType.Type).Interface()
			in[0] = reflect.ValueOf(yyy)
			in[1] = reflect.ValueOf([]byte(fmt.Sprintf("%v", stringValue)))
			method.Func.Call(in)[0].Interface()
			return in[0].Elem().Interface(), nil
		}
		yyy := reflect.New(reflectType.Type).Interface()
		return reflect.ValueOf(yyy).Elem().Interface(), errors.New("parse error")
	}
}

func comparatorInfo(t *storage, o ...order.Stages) string {
	tt := t.Type.Type.Name()
	fileName, err := os.Executable()
	if err != nil {
		fileName = "appImageBinary"
	} else {
		fileName = filepath.Base(fileName)
	}
	ret := ""
	reflectTypeName := strings.ToLower(strings.TrimSpace(tt))
	for _, v := range o {
		switch reflectTypeName {
		case "string":
			switch v {
			case order.Flag:
				ret += sample(fileName, t.Name, t.Default)
			case order.File:
				ret += sampleJson(t.JSON, t.Store)
			case order.Env:
				ret += sampleEnv(t.Env, t.Default)
			}
		case "bool":
			switch v {
			case order.Flag:
				ret += sampleBool(fileName, t.Name)
			case order.File:
				ret += sampleJson(t.JSON, t.Store)
			case order.Env:
				ret += sampleEnv(t.Env, t.Default)
			}
		case "duration":
			switch v {
			case order.Flag:
				ret += durationSample(fileName, t.Name, t.Default)
			case order.File:
				ret += sampleJson(t.JSON, t.Store)
			case order.Env:
				ret += sampleEnv(t.Env, t.Default)
			}
		case "int8", "int16", "int32", "int64", "rune":
			switch v {
			case order.Flag:
				ret += sample(fileName, t.Name, t.Default)
			case order.File:
				ret += sampleJson(t.JSON, t.Store)
			case order.Env:
				ret += sampleEnv(t.Env, t.Default)
			}
		case "int":
			switch v {
			case order.Flag:
				ret += sample(fileName, t.Name, t.Default)
			case order.File:
				ret += sampleJson(t.JSON, t.Store)
			case order.Env:
				ret += sampleEnv(t.Env, t.Default)
			}
		case "uint8", "uint16", "uint32", "uint64":
			switch v {
			case order.Flag:
				ret += sample(fileName, t.Name, t.Default)
			case order.File:
				ret += sampleJson(t.JSON, t.Store)
			case order.Env:
				ret += sampleEnv(t.Env, t.Default)
			}
		case "uint":
			switch v {
			case order.Flag:
				ret += sample(fileName, t.Name, t.Default)
			case order.File:
				ret += sampleJson(t.JSON, t.Store)
			case order.Env:
				ret += sampleEnv(t.Env, t.Default)
			}
		case "float32", "float64":
			switch v {
			case order.Flag:
				ret += sample(fileName, t.Name, fmt.Sprintf("%f", helpers.ValidFloat(t.Default)))
			case order.File:
				ret += sampleJson(t.JSON, t.Store)
			case order.Env:
				ret += sampleEnv(t.Env, t.Default)
			}
		default:
			switch v {
			case order.Flag:
				// ret += sample(fileName, flagName, def, flagOrder)
			case order.File:
				ret += sampleJson(t.JSON, t.Store)
			case order.Env:
				ret += sampleEnv(t.Env, t.Default)
			}
		}
	}
	return ret
}

func sampleEnv(envValue, def string) string {
	return fmt.Sprintf("Sample environment:\t%s=%s\n", strings.ToUpper(envValue), def)
}

func sampleJson(jsonValue string, def any) string {
	m := make(map[string]any, 1)
	m[jsonValue] = def
	out, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return ""
	}
	return fmt.Sprintf("Sample JSON config:\t\n%s\n", string(out))
}

func sample(fileName, flagName, def string) string {
	return fmt.Sprintf("Sample flag value:\t%s -%s=%s\n", fileName, flagName, def)
}

func sampleBool(fileName, flagName string) string {
	return fmt.Sprintf("Sample(TRUE):\t%s -%s\n", fileName, flagName) +
		fmt.Sprintf("Sample(FALSE):\t%s\n", fileName) +
		extSample("TRUE", fileName, flagName, "true") +
		extSample("FALSE", fileName, flagName, "false") +
		extSample("TRUE", fileName, flagName, "1") +
		extSample("FALSE", fileName, flagName, "0") +
		extSample("TRUE", fileName, flagName, "t") +
		extSample("FALSE", fileName, flagName, "f")
}

func extSample(value, fileName, flagName, def string) string {
	return fmt.Sprintf("Sample(%s):\t%s -%s=%s\n", value, fileName, flagName, def)
}

func durationSample(fileName, flagName, def string) string {
	return durationSample1(fileName, flagName) + durationSample2(fileName, flagName) + durationSample3(fileName, flagName)
}

func durationSample1(fileName, flagName string) (print string) {
	for k, v := range durationTypesMap {
		print += extSample(v, fileName, flagName, "1"+k)
	}
	return
}

func durationSample2(fileName, flagName string) string {
	return extSample("1 Hour 2 Minutes and 3 Seconds", fileName, flagName, "1h2m3s")
}

func durationSample3(fileName, flagName string) string {
	return extSample("111 Seconds", fileName, flagName, "111")
}
