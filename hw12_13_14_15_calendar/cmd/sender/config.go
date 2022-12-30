package main

import (
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/cmd"
	"log"
	"os"
)

type Config struct {
	RabbitInfo RabbitConfig     `yaml:"optQueue"`
	Logger     cmd.LoggerConfig `yaml:"logger"`
}

type RabbitConfig struct {
	Url       string `yaml:"url"`
	Name      string `yaml:"name"`
	NoWait    bool   `yaml:"noWait"`
	NoLocal   bool   `yaml:"noLocal"`
	AutoAck   bool   `yaml:"autoAck"`
	Exclusive bool   `yaml:"exclusive"`
}

func (rc *RabbitConfig) Set() {
	exclusive := os.Getenv("EXCLUSIVE")
	autoAck := os.Getenv("AUTO_ACK")
	noLocal := os.Getenv("NO_LOCAL")
	noWait := os.Getenv("NO_WAIT")

	rc.Url = os.Getenv("URL")
	rc.Name = os.Getenv("NAME_Q")

	envs, err := cmd.ParseBool(exclusive, autoAck, noLocal, noWait)
	if err != nil {
		log.Fatal(err)
	}

	rc.Exclusive = envs[exclusive]
	rc.AutoAck = envs[autoAck]
	rc.NoLocal = envs[noLocal]
	rc.NoWait = envs[noWait]
}

func NewConfig() (config Config) {
	config.Logger.Set()
	config.RabbitInfo.Set()

	return config
}
