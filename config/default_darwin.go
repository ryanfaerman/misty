package config

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		Logger.Fatal("cannot find home directory")
	}
	ConfigFile = filepath.Join(home, ".config/misty/misty.yml")
}
