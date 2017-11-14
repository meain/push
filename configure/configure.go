package configure

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Conf has the settings configured by the user
type Conf struct {
	Token  string `yaml:"token"`
	Device string `yaml:"device"`
}

func (c *Conf) writeConfFile(file string) error {
	file = expandFilePath(file)
	yamlFile, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, yamlFile, 0644)
	return err
}

// GetConf parses the config file and get the user config
func (c *Conf) GetConf(file string) error {
	file = expandFilePath(file)
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conf) SetDefaultDevice(file string, iden string) error {
	c.Device = iden
	err := c.writeConfFile(file)
	return err
}

func (c *Conf) SetToken(file string, token string) error {
	c.Token = token
	err := c.writeConfFile(file)
	return err
}
