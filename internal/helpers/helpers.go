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
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/KusoKaihatsuSha/startup/internal/order"
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
func ToLog(err any, text ...string) {
	if err == nil || err == "" {
		return
	}
	switch val := err.(type) {
	case error:
		if val.Error() != "" {
			log.Error().Msgf("%s |-> %v", text, val)
		}
	default:
		log.Debug().Msgf("%s |-> %v", text, val)
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

// DeleteFile delete file in temp folder. Actually means delete any file by path
func DeleteFile(filename string) {
	filename = separatorCorrect(filename)
	ToLog(os.RemoveAll(filename))
}

// FileExist checking if file exists by path
func FileExist(filename string) bool {
	if filename == "." {
		return false
	}
	_, err := os.Stat(filename)
	return err == nil
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
	if err := file.Close(); err != nil {
		ToLog(err, fmt.Sprintf("file '%s' close error", file.Name()))
	}
}

// SettingsFile return map[string]string from the setting file
func SettingsFile(filename string) (compare map[string]any) {
	filename = separatorCorrect(filename)
	if FileExist(filename) {
		f, err := os.ReadFile(filename)
		ToLog(err, fmt.Sprintf("settings file '%s' is not access", filename))
		err = json.Unmarshal(f, &compare)
		ToLog(err, fmt.Sprintf("settings file '%s' parse error", filename))
	}
	return
}

// CreateFile create file or temp file
func CreateFile(filename string) string {
	filename = separatorCorrect(filename)
	if filename == "" || filename == "." {
		file, err := os.CreateTemp("", "*.tmpgo")
		defer CloseFile(file)
		ToLog(err, fmt.Sprintf("file '%s' not created", filename))
		return file.Name()
	} else {
		file, err := os.Create(filename)
		defer CloseFile(file)
		ToLog(err, fmt.Sprintf("file '%s' not created", filename))
		return file.Name()
	}
}

// ValidTempFile - validation type.
// create same name file in Temp folder or create random temp file
func ValidTempFile(filename string) string {
	filename = separatorCorrect(filename)
	_, file := filepath.Split(filename)

	if filename == "" || filename == "." {
		return CreateFile(filename)
	}

	filename = filepath.Join(os.TempDir(), file)
	if !FileExist(filename) {
		return CreateFile(filename)
	}

	return filename
}

// ValidConfigFile - validation type.
// create same name file in Temp folder or create random temp file
func ValidConfigFile(filename string) string {
	filename = separatorCorrect(filename)

	if !FileExist(filename) || filename == "" || filename == "." {
		return ""
	}

	return filename
}

// ValidFile - validation type.
// Create file in not exist.
func ValidFile(filename string) string {
	filename = separatorCorrect(filename)

	if filename == "" || filename == "." {
		return CreateFile(filename)
	}

	if !FileExist(filename) {
		CreateFile(filename)
	}
	return filename
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
	ToLog(err, fmt.Sprintf("duration '%s' parse error", v))
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

func separatorCorrect(filename string) string {
	r := strings.NewReplacer("/", string(os.PathSeparator), "\\", string(os.PathSeparator))
	return filepath.Clean(r.Replace(strings.TrimSpace(filename)))
}

func FileConfExistInStages(stages ...order.Stages) bool {
	for _, stage := range stages {
		if stage == order.FILE {
			return true
		}
	}
	return false
}

func printDebug(stages ...order.Stages) {
	fmt.Println("Structure filling order:")
	for k, stage := range stages {
		prefix := ""
		suf := "\n"
		if k != len(stages)-1 {
			suf = " ↣ "
		}
		if k == 0 {
			prefix = "\t"
		}
		switch stage {
		case order.FLAG:
			fmt.Printf("%s%s%s", prefix, "[Flags]", suf)
		case order.FILE:
			fmt.Printf("%s%s%s", prefix, "[JSON File]", suf)
		case order.ENV:
			fmt.Printf("%s%s%s", prefix, "[Environments]", suf)
		}
	}
}

func PrintDebug(stages ...order.Stages) {
	for _, stage := range stages {
		switch stage {
		case order.NoPreloadConfig:
			fmt.Printf("%s - Disable find config in other places. Only in list\n", "NoPreloadConfig")
			printDebug(stages...)
			return
		case order.PreloadConfigFlagThenEnv:
			fmt.Printf("%s - Preload find config in Flag then Env\n", "PreloadConfigFlagThenEnv")
			printDebug(stages...)
			return
		case order.PreloadConfigFlag:
			fmt.Printf("%s - Preload find config in Flag\n", "PreloadConfigFlag")
			printDebug(stages...)
			return
		case order.PreloadConfigEnv:
			fmt.Printf("%s - Preload find config in Env\n", "PreloadConfigEnv")
			printDebug(stages...)
			return
		}
	}
	if stages == nil {
		fmt.Printf("%s - only defaults\n", "EMPTY")
	} else {
		fmt.Printf("%s - Preload find config in Env then Flag\n", "PreloadConfigEnvThenFlag/Default")
		printDebug(stages...)
	}

}

func PresetPreload(stages ...order.Stages) []order.Stages {
	for _, v := range stages {
		switch v {
		case order.NoPreloadConfig:
			return stages
		case order.PreloadConfigEnvThenFlag:
			return []order.Stages{order.ENV, order.FLAG}
		case order.PreloadConfigFlagThenEnv:
			return []order.Stages{order.FLAG, order.ENV}
		case order.PreloadConfigFlag:
			return []order.Stages{order.FLAG}
		case order.PreloadConfigEnv:
			return []order.Stages{order.ENV}
		}
	}
	return []order.Stages{order.ENV, order.FLAG}
}
