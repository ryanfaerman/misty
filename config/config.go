package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	ApplicationName = "misty"
	ApplicationDesc = "misty interacts with a misty robot"
)

var (
	ConfigFile string

	Config config
	Logger = log.New()
)

func init() {
	Logger.SetLevel(log.InfoLevel)
}

type config struct {
	Addr string `yaml:"addr"`
}

func Addr() string { return Config.Addr }

// Load the config at the given path
func Load(path string) error {
	if path != "" {
		ConfigFile = path
	}

	b, err := ioutil.ReadFile(ConfigFile)

	if err != nil {
		return errors.Wrap(err, "cannot read config")
	}

	c := config{}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return errors.Wrap(err, "cannot unmarshal config")
	}

	Config = c

	return nil
}

// Reload the previously loaded config file
func Reload() error {
	return Load(ConfigFile)
}
