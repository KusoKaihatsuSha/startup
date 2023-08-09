// Package helpers working with logging
// and other non-main/other help-function.
package helpers

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogOutNullTS = iota
	LogOutFastTS
	LogOutHumanTS
)

// Error output types
const (
	LogErrNullTS = iota + 100
	LogErrFastTS
	LogErrHumanTS
)

// LogNull Error null
const (
	LogNull = 1000
)

const (
	defTimePostfix = "s"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

// ToLog splitting notifications to Err and Debug.
func ToLog(err any) {
	if err == nil || err == "" {
		return
	}
	switch val := err.(type) {
	case error:
		if val.Error() != "" {
			log.Error().Msgf("%v", val)
		}
	default:
		log.Debug().Msgf("%v", val)
	}
}

// ToLogWithType splitting notifications to Err and Debug. With the type of timestamps.
func ToLogWithType(err any, typ int) {
	if err == nil || err == "" {
		return
	}
	if typ == LogNull {
		std := logger(typ)
		std.Debug().Msgf("%v", err)
		return
	}
	switch val := err.(type) {
	case error:
		if val.Error() != "" {
			std := logger(typ + LogErrNullTS)
			std.Error().Msgf("%v", val)
		}
	default:
		std := logger(typ)
		std.Debug().Msgf("%v", val)
	}
}

// logger return logger with TS.
func logger(typ int) zerolog.Logger {
	switch typ {
	case LogOutNullTS:
		return zerolog.New(os.Stdout).With().Logger()
	case LogErrNullTS:
		return zerolog.New(os.Stderr).With().Logger()
	case LogOutFastTS:
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	case LogErrFastTS:
		return zerolog.New(os.Stderr).With().Timestamp().Logger()
	case LogOutHumanTS:
		return zerolog.New(os.Stdout).With().Str("time", time.Now().Format("200601021504")).Logger()
	case LogErrHumanTS:
		return zerolog.New(os.Stderr).With().Str("time", time.Now().Format("200601021504")).Logger()
	case LogNull:
		return zerolog.New(io.Discard).With().Logger()
	default:
		return zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}

// CreateTmp create file in Temp folder
func CreateTmp() string {
	fileEnv, err := os.CreateTemp("", "tmp_golang_")
	ToLog(err)
	defer func(path string) {
		ToLog(fileEnv.Close())
	}(fileEnv.Name())
	if runtime.GOOS != "windows" {
		return string(os.PathSeparator) + fileEnv.Name()
	}
	return fileEnv.Name()
}

// DeleteTmp delete file in Temp folder. Actually imply delete any file by path
func DeleteTmp(path string) {
	ToLog(os.RemoveAll(path))
}

// FileExist check exist
func FileExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

// Print - clear spaces
func Print(v any) {
	fmt.Println(
		strings.ReplaceAll(
			fmt.Sprint(v),
			"  ",
			" ",
		),
	)
}

// SettingsFile return map[string]string from the setting file
func SettingsFile(filename string) (compare map[string]string) {
	if runtime.GOOS != "windows" {
		filename = fmt.Sprintf("/%s", filename)
	}
	f, err := os.ReadFile(filename)
	ToLog(err)
	err = json.Unmarshal(f, &compare)
	ToLog(err)
	return
}

func ValidUUID(v string) string {
	if tmp, err := uuid.Parse(v); err == nil {
		return tmp.String()
	}
	return uuid.New().String()
}

func ValidInt(v string) int64 {
	if tmp, err := strconv.ParseInt(v, 10, 64); err == nil {
		return tmp
	}
	return 0
}

func ValidFloat(v string) float64 {
	if tmp, err := strconv.ParseFloat(v, 64); err == nil {
		return tmp
	}
	return 0
}

func ValidBool(v string) bool {
	if tmp, err := strconv.ParseBool(v); err == nil {
		return tmp
	}
	return false
}

func ValidTempFile(v string) string {
	if strings.TrimSpace(v) == "" {
		return ""
	}
	var file string
	if runtime.GOOS == "windows" {
		check := strings.FieldsFunc(v, func(ss rune) bool {
			return strings.ContainsAny(string(ss), `\/`)
		})
		tmpString := filepath.Join(os.TempDir(), check[len(check)-1])
		file = tmpString
	} else {
		file = "/" + v
	}
	return file
}

func ValidFile(v string) string {
	check := strings.FieldsFunc(v, func(ss rune) bool {
		return strings.ContainsAny(string(ss), `\/`)
	})
	if strings.TrimSpace(v) == "" {
		return ""
	}
	if !FileExist(v) {
		ToLog(fmt.Sprintf("file '%s' not found", v))
	}
	return strings.Join(check, string(os.PathSeparator))
}

func ValidTimer(v string) time.Duration {
	if v[0] == '-' || v[0] == '+' {
		v = v[1:]
	}
	l := len(v) - 1
	if '0' <= v[l] && v[l] <= '9' {
		v += defTimePostfix
	}
	tmp, err := time.ParseDuration(v)
	ToLog(err)
	return tmp
}

func ValidURL(v string) string {
	// trim prefix
	re := regexp.MustCompile(`^.*(://|^)[^/]+`)
	trimPrefix := re.FindString(v)
	re = regexp.MustCompile(`^.*(://|^)`)
	fullAddress := re.ReplaceAllString(trimPrefix, "")
	// trim port
	re = regexp.MustCompile(`^[^/:$]+`)
	address := re.FindString(fullAddress)
	// fill address
	if strings.TrimSpace(address) == "" {
		return ""
	}
	// check ip
	isIP := false
	re = regexp.MustCompile(`\d+`)
	isIPTest := re.ReplaceAllString(address, "")
	isIPTest = strings.ReplaceAll(isIPTest, ".", "")
	if strings.TrimSpace(isIPTest) == "" {
		isIP = true
	}
	// correct IP
	if isIP {
		re = regexp.MustCompile(`\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}`)
		addressIP := re.FindString(address)
		if strings.TrimSpace(addressIP) == "" {
			return ""
		}
	}
	// check and correct port
	re = regexp.MustCompile(`:.*`)
	correctPort := re.FindString(fullAddress)
	correctPort = strings.Replace(correctPort, ":", "", 1)
	re = regexp.MustCompile(`\D`)
	correctPort = re.ReplaceAllString(correctPort, "")
	correctPort = strings.Replace(correctPort, ":", "", 1)
	if strings.TrimSpace(correctPort) == "" {
		return address + ":80"
	}
	return address + ":" + correctPort
}

// pointerFlag return the flag pointer
func PointerFlag(name string, fs *flag.FlagSet) *flag.Flag {
	var current *flag.Flag
	fs.VisitAll(func(f *flag.Flag) {
		if f.Name == name {
			current = f
		}
	})
	return current
}
