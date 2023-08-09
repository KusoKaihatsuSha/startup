package startup_test

import (
	"fmt"
	"os"

	"github.com/KusoKaihatsuSha/startup"
	"github.com/KusoKaihatsuSha/startup/internal/helpers"
)

type Test struct {
	Address  string   `json:"target" default:"123.123.123.123:3389" flag:"t,target" env:"PROXY_TARGET" text:"Target address and port" valid:"url"`
	NewValid []string `json:"new-valid" default:"new valid is default" flag:"valid" text:"-" valid:"test"`
}

type testValid string

var testValidation testValid = "test"

func (o testValid) Valid(key, value string) (any, bool) {
	if key == string(o) {
		return []string{value + "+++"}, true
	}
	return nil, false
}

func Example_configOrder() {
	startup.AddValidation(testValidation)

	testDataFlag := `{
                 "target": "http://file-flag"
         }`
	testDataEnv := `{
                 "target": "http://file-env"
         }`
	fileFlag := helpers.CreateTmp()
	err := os.WriteFile(fileFlag, []byte(testDataFlag), 0755)
	if err != nil {
		fmt.Println(err)
	}
	fileEnv := helpers.CreateTmp()
	err = os.WriteFile(fileEnv, []byte(testDataEnv), 0755)

	if err != nil {
		fmt.Println(err)
	}

	err = os.Setenv("PROXY_TARGET", "http://env")
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
		"-t="+flagVal,
		"-valid="+flagVal,
	)

	fmt.Println("PRESET:")
	fmt.Println("flag:", flagVal)
	fmt.Println("environment:", os.Getenv("PROXY_TARGET"))
	fmt.Println("config-file in environment:")
	fmt.Println("", testDataEnv)
	fmt.Println("config-file in flag:")
	fmt.Println("", testDataFlag)

	fmt.Println("last 'file' is 'file'", startup.GetForce[Test](startup.File).Address)

	fmt.Println("last 'flag' is 'flag'", startup.GetForce[Test](startup.Flag).Address)

	fmt.Println("last 'env' is 'env'", startup.GetForce[Test](startup.Env).Address)

	fmt.Println("last 'file --> flag' is 'flag'", startup.GetForce[Test](startup.FileFlag).Address)

	fmt.Println("last 'file --> env' is 'env'", startup.GetForce[Test](startup.FileEnv).Address)

	fmt.Println("last 'flag --> env' is 'env'", startup.GetForce[Test](startup.FlagEnv).Address)

	fmt.Println("last 'file --> flag --> env' is 'env'", startup.GetForce[Test](startup.FileFlagEnv).Address)

	fmt.Println("last 'flag --> file --> env' is 'env'", startup.GetForce[Test](startup.FlagFileEnv).Address)

	fmt.Println("last 'file --> env --> flag' is 'flag'", startup.GetForce[Test](startup.FileEnvFlag).Address)

	fmt.Println("last 'env --> file --> flag' is 'flag'", startup.GetForce[Test](startup.EnvFileFlag).Address)

	fmt.Println("last 'flag --> env --> file' is 'file'", startup.GetForce[Test](startup.FlagEnvFile).Address)

	fmt.Println("last 'env --> flag --> file' is 'file'", startup.GetForce[Test](startup.EnvFlagFile).Address)

	fmt.Println(startup.GetForce[Test](startup.EnvFlagFile).NewValid)

	// Output:
	// PRESET:
	// flag: http://flag
	// environment: http://env
	// config-file in environment:
	//  {
	//                  "target": "http://file-env"
	//          }
	// config-file in flag:
	//  {
	//                  "target": "http://file-flag"
	//          }
	// last 'file' is 'file' file-env:80
	// last 'flag' is 'flag' flag:80
	// last 'env' is 'env' env:80
	// last 'file --> flag' is 'flag' flag:80
	// last 'file --> env' is 'env' env:80
	// last 'flag --> env' is 'env' env:80
	// last 'file --> flag --> env' is 'env' env:80
	// last 'flag --> file --> env' is 'env' env:80
	// last 'file --> env --> flag' is 'flag' flag:80
	// last 'env --> file --> flag' is 'flag' flag:80
	// last 'flag --> env --> file' is 'file' file-env:80
	// last 'env --> flag --> file' is 'file' file-flag:80
	// [http://flag+++]

}
