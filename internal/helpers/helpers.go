// Package helpers working with logging
// and other non-main/other help-function.
package helpers

import (
	"encoding/json"
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

// Logger output types
const (
	LogOutNullTS = iota
	LogOutFastTS
	LogOutHumanTS
)

// Logger Error output types
const (
	LogErrNullTS = iota + 100
	LogErrFastTS
	LogErrHumanTS
)

// LogNull Logger null
const (
	LogNull = 1000
)

// second as default
const (
	defaultTimePostfix = "s"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

// ToLog splits notifications into Error and Debug(text).
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

// ToLogWithType splits notifications into Error and Debug(text). Using timestamped.
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

// logger return Main logger with timestamp.
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

// DeleteTmp delete file in temp folder. Actually means delete any file by path
func DeleteTmp(path string) {
	ToLog(os.RemoveAll(path))
}

// FileExist checking if file exists by path
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Print - clear dbl spaces
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
func SettingsFile(filename string) (compare map[string]any) {
	if runtime.GOOS != "windows" {
		filename = fmt.Sprintf("/%s", filename)
	}
	f, err := os.ReadFile(filename)
	ToLog(err)
	err = json.Unmarshal(f, &compare)
	ToLog(err)
	return
}

// ValidUUID - validation on the UUID
func ValidUUID(v string) string {
	if tmp, err := uuid.Parse(v); err == nil {
		return tmp.String()
	}
	return uuid.New().String()
}

// ValidInt - validation on the int
func ValidInt(v string) int64 {
	if tmp, err := strconv.ParseInt(v, 10, 64); err == nil {
		return tmp
	}
	return 0
}

// ValidUint - validation on the uint
func ValidUint(v string) uint64 {
	if tmp, err := strconv.ParseUint(v, 10, 64); err == nil {
		return tmp
	}
	return 0
}

// ValidFloat - validation on the float
func ValidFloat(v string) float64 {
	if tmp, err := strconv.ParseFloat(v, 64); err == nil {
		return tmp
	}
	return 0
}

// ValidBool - validation on the boolean
func ValidBool(v string) bool {
	if tmp, err := strconv.ParseBool(v); err == nil {
		return tmp
	}
	return false
}

// CloseFile close file. Using with defer
func CloseFile(file *os.File) {
	if errClose := file.Close(); errClose != nil {
		ToLog(errClose)
	}
}

// FilepathElements - split filename using separate symbols
func FilepathElements(path string) []string {
	return strings.FieldsFunc(path, func(symbols rune) bool {
		return strings.ContainsAny(string(symbols), `\/`)
	})
}

// CreateFile create file or temp file
func CreateFile(path string) string {
	if path == "" {
		file, err := os.CreateTemp("", "*.tmpgo")
		defer CloseFile(file)
		if err != nil {
			ToLog(fmt.Sprintf("file '%s' not created", path))
		}
		return file.Name()
	} else {
		file, err := os.Create(path)
		defer CloseFile(file)
		if err != nil {
			ToLog(fmt.Sprintf("file '%s' not created", path))
		}
		return file.Name()
	}
}

// ValidTempFile - validation type.
// create same name file in Temp folder or create random temp file
func ValidTempFile(filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return CreateFile("")
	}
	elementsOfPath := FilepathElements(filename)
	newFilepath := filepath.Join(os.TempDir(), elementsOfPath[len(elementsOfPath)-1])
	if !FileExist(newFilepath) {
		CreateFile(newFilepath)
	}
	return newFilepath
}

// ValidFile - validation type.
// Create file in not exist.
func ValidFile(filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return CreateFile("")
	}
	elementsOfPath := FilepathElements(filename)
	newFilepath := filepath.Join(elementsOfPath...)
	if !FileExist(newFilepath) {
		CreateFile(newFilepath)
	}
	return newFilepath
}

// ValidDuration - validation on the duration
func ValidDuration(v string) time.Duration {
	if v[0] == '-' || v[0] == '+' {
		v = v[1:]
	}
	l := len(v) - 1
	if '0' <= v[l] && v[l] <= '9' {
		v += defaultTimePostfix
	}
	tmp, err := time.ParseDuration(v)
	ToLog(err)
	return tmp
}

// ValidURL - validation on the url
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
