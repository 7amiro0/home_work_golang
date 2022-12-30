package cmd

import "os"

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
