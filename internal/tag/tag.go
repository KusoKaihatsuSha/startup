package tags

import (
	"flag"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/KusoKaihatsuSha/startup/internal/helpers"
	"github.com/KusoKaihatsuSha/startup/internal/validation"
)

const (
	defaultTag     = "default"
	flagTag        = "flag"
	EnvironmentTag = "env"
	helpTextTag    = "text"
	validationTag  = "valid"
	jsonTag        = "json"

	testTrigger = "-test."
)

// Tmp consist information when reading/valid configs
type TagsInfo map[string]tagInfo

// TagInfo store tags Config struct
type tagInfo struct {
	ConfigFile func(map[string]string) tagInfo
	Valid      func() any
	DummyFlags func() tagInfo
	Env        func() tagInfo
	Flag       func() tagInfo
	FlagSet    *flag.FlagSet
	Flags      []*flag.Flag
	Name       string
	Annotation
}

// Tag store general values
type Annotation struct {
	valid string
	json  string
	desc  string
	env   string
	def   string
	flag  []string
}

// tags filling the 'tag'
func Fill[T any](field string) tagInfo {
	var structAny T
	ret := tagInfo{}
	ret.Name = field
	ret.FlagSet = flag.NewFlagSet("", flag.ContinueOnError)
	tmp, exist := reflect.TypeOf(structAny).FieldByName(ret.Name)
	if !exist {
		return tagInfo{}
	}
	if v, ok := tmp.Tag.Lookup(defaultTag); ok {
		ret.def = v
	}
	if v, ok := tmp.Tag.Lookup(flagTag); ok {
		ret.flag = strings.Split(v, ",")
		var all string
		for _, flagTag := range ret.flag {
			ret.FlagSet.StringVar(&all, flagTag, ret.def, ret.desc)
			ret.Flags = append(ret.Flags, helpers.PointerFlag(flagTag, ret.FlagSet))
		}
	}
	if v, ok := tmp.Tag.Lookup(EnvironmentTag); ok {
		ret.env = v
	}
	if v, ok := tmp.Tag.Lookup(helpTextTag); ok {
		ret.desc = v
	}
	if v, ok := tmp.Tag.Lookup(validationTag); ok {
		ret.valid = v
	}
	if v, ok := tmp.Tag.Lookup(jsonTag); ok {
		ret.json = v
	}
	ret.ConfigFile = func(m map[string]string) tagInfo {
		for k, v := range m {
			if ret.json == k {
				for _, f := range ret.Flags {
					f.DefValue = v
					helpers.ToLog(
						f.Value.Set(f.DefValue),
					)
				}
			}
		}
		return ret
	}
	ret.DummyFlags = func() tagInfo {
		for _, flagTag := range ret.flag {
			for _, v := range os.Args {
				if strings.Contains(v, testTrigger) {
					var _ = func() bool {
						testing.Init()
						return true
					}()
					break
				}
			}
			if flag.Lookup(flagTag) == nil {
				flag.StringVar(new(string), flagTag, ret.def, ret.desc)
			}
		}
		return ret
	}
	ret.Flag = func() tagInfo {
		def := ret.FlagSet.Output()
		ret.FlagSet.SetOutput(io.Discard)
		for _, arg := range os.Args {
			err := ret.FlagSet.Parse([]string{arg})
			helpers.ToLogWithType(err, helpers.LogNull)
		}
		ret.FlagSet.SetOutput(def)
		return ret
	}
	ret.Env = func() tagInfo {
		env, ok := os.LookupEnv(ret.env)
		if ok {
			for _, f := range ret.Flags {
				err := f.Value.Set(env)
				helpers.ToLog(err)
				break
			}
		}
		return ret
	}
	ret.Valid = func() any {
		value := ""
		for _, f := range ret.Flags {
			value = f.Value.String()
			break
		}
		for _, v := range validation.Valids {
			if ret, ok := v.Valid(ret.valid, value); ok {
				return ret
			}
		}
		return value
	}
	return ret
}
