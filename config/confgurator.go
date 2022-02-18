package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var configFields = []struct {
	key      string
	question string
}{
	{"watchFolder", "Folder to store screenshots"},
	{"creds.endpoint", "URL of your S3-compatible server"},
	{"creds.key", "S3 access key"},
	{"creds.secret", "S3 secret"},
	{"creds.region", "S3 region"},
	{"creds.bucket", "S3 bucket"},
}

// configure asks the user to enter data needed for config
// TODO add tests, remove duplication
func configure(v *viper.Viper, p string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Setup foxyshot")
	fmt.Println("---------------------")

	for _, cf := range configFields {
		fmt.Printf("%s: ", cf.question)
		val, _ := reader.ReadString('\n')
		val = strings.Replace(val, "\n", "", -1)
		v.Set(cf.key, val)
	}

	fmt.Println("Done! Saving to", p)

	err := v.SafeWriteConfigAs(p)
	if errors.Is(err, viper.ConfigFileAlreadyExistsError(p)) {
		fmt.Println("Config file already exists. Do you really want to overwrite the existing config? [y/n]")
		val, _ := reader.ReadString('\n')
		if val == "y\n" || val == "yes\n" {
			fmt.Printf("Done!")
			return v.WriteConfigAs(p)
		}

		fmt.Println("No changes were applied!")
		return nil
	}
	return err
}

// RunConfigure saves config to home folder
func RunConfigure() error {
	configDir := expandHomeFolder("~/.config/foxyshot")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return configure(viper.GetViper(), configDir+"/config.json")
}
