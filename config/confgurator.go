package config

import (
	"fmt"
	"os"
	"os/exec"
)

const configTemplate = `
{
    "watchFolder": "Folder to store screenshots (e. g. ~/Screenshots)",
    "s3": {
		"key":        "S3 access key",
		"secret":     "S3 secret",
		"endpoint":   "URL of your S3-compatible server",
		"region":     "S3 region",
		"bucket":     "S3 bucket",
		"publicURIs": false,
		"duration":   "if publicURIs is false, this is the duration of the presigned URL. e.g. 24h",
		"cdn":        "custom domain for sharing screenshots from your S3"
	},
	"screenshots": {
		"jpegQuality": 999,
		"removeOriginals": true 
	}
}`

// forceConfig creates config file if it doesn't exist
func forceConfig(p string) error {
	if _, err := os.Stat(p); err == nil {
		return nil
	}

	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("error creating config file %s: %w", p, err)
	}

	defer f.Close()

	_, err = f.WriteString(configTemplate)
	if err != nil {
		return fmt.Errorf("error writing config file %s: %w", p, err)
	}

	return nil
}

// openEditor opens editor for given file and gives it control of the terminal
func openEditor(p string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmdStr := fmt.Sprintf("%s %s", editor, p)

	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error opening editor %s: %w", editor, err)
	}

	return nil
}

func configure(configPath string) error {
	if err := forceConfig(configPath); err != nil {
		return err
	}

	return openEditor(configPath)
}

// RunConfigure saves config to home folder
func RunConfigure() error {
	configDir := expandHomeFolder("~/.config/foxyshot")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return configure(configDir + "/config.json")
}
