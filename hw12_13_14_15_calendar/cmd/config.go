package cmd

import (
	"os"
	"strconv"
)

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

func (l LoggerConfig) Set() {
	l.Level = os.Getenv("LOGGER_LEVEL")
}

func ParseBool(envs ...string) (result map[string]bool, err error) {
	result = make(map[string]bool)
	var opt bool
	for _, env := range envs {
		opt, err = strconv.ParseBool(env)
		if err != nil {
			return nil, err
		}

		result[env] = opt
	}

	return result, err
}
