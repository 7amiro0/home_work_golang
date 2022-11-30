package main

import (
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/cmd"
	"log"
	"os"
	"strconv"
)

type Config struct {
	RabbitInfo RabbitConfig     `yaml:"optQueue"`
	Logger     cmd.LoggerConfig `yaml:"logger"`
}

type RabbitConfig struct {
	Url       string `yaml:"url"`
	Name      string `yaml:"name"`
	Exclusive bool   `yaml:"exclusive"`
	AutoAck   bool   `yaml:"autoAck"`
	NoLocal   bool   `yaml:"noLocal"`
	NoWait    bool   `yaml:"noWait"`
}

func parseBool(envs ...string) (result map[string]bool, err error) {
	var opt bool
	result = make(map[string]bool)

	for _, env := range envs {
		opt, err = strconv.ParseBool(env)
		if err != nil {
			return nil, err
		}

		result[env] = opt
	}

	return result, err
}

func (rc *RabbitConfig) Set() {
	exclusive := os.Getenv("EXCLUSIVE")
	autoAck := os.Getenv("AUTO_ACK")
	noLocal := os.Getenv("NO_LOCAL")
	noWait := os.Getenv("NO_WAIT")

	rc.Url = os.Getenv("URL")
	rc.Name = os.Getenv("NAME_Q")

	envs, err := parseBool(exclusive, autoAck, noLocal, noWait)

	if err != nil {
		log.Fatal(err)
	}

	rc.Exclusive = envs[exclusive]
	rc.AutoAck = envs[autoAck]
	rc.NoLocal = envs[noLocal]
	rc.NoWait = envs[noWait]
}

func NewConfig() (config Config) {
	config.RabbitInfo.Set()

	return config
}
