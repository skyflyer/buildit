package conf

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const filename = "buildit.yml"

// DefaultWorkingDirectory is where git file repo is stored
const DefaultWorkingDirectory = "buildit-repo"

// LastHeadMarker is the name of the file that stores last commit hash
const LastHeadMarker = ".lasthead"

// Conf struct
type Conf struct {
	Repo             string    `yaml:"repo"`
	Auth             *AuthConf `yaml:"auth"`
	WorkingDirectory string    `yaml:"workdir"`
	Branch           string    `yaml:"branch"`
	Steps            []string  `yaml:"steps,flow"`
}

// AuthConf struct
type AuthConf struct {
	Key         string `yaml:"key"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	UseSSHAgent bool   `yaml:"use_ssh_agent"`
}

func parseYaml(content []byte) (Conf, error) {
	c := Conf{}
	err := yaml.Unmarshal(content, &c)
	if err != nil {
		return c, err
	}

	if c.Repo == "" {
		return c, errors.New("Repository must be defined")
	}

	if len(c.Steps) == 0 {
		return c, errors.New("At least one step must be defined")
	}

	if len(c.WorkingDirectory) == 0 {
		c.WorkingDirectory = DefaultWorkingDirectory
	}

	return c, nil
}

// ReadConfig reads config from filename
func ReadConfig() (Conf, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return Conf{}, err
	}

	return parseYaml(content)
}
