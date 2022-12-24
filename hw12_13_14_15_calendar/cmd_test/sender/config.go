package main

import (
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/cmd"
	"log"
	"os"
	"strconv"
)

type Config struct {
	RabbitInfo RabbitConfig
	Logger     cmd.LoggerConfig
}

type RabbitConfig struct {
	Url        string
	Name       string
	NameTest   string
	Exchange   string
	Mandatory  bool
	Immediate  bool
	Exclusive  bool
	AutoAck    bool
	NoLocal    bool
	NoWait     bool
	AutoDelete bool
	Durable    bool
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
	autoDelete := os.Getenv("AUTO_DELETE")
	durable := os.Getenv("DURABLE")
	mandatory := os.Getenv("MANDATORY")
	immediate := os.Getenv("IMMEDIATE")

	rc.Exchange = os.Getenv("EXCHANGE")
	rc.Url = os.Getenv("URL")
	rc.Name = os.Getenv("NAME_Q")
	rc.NameTest = os.Getenv("NAME_TEST_Q")

	envs, err := parseBool(exclusive, autoAck, noLocal, noWait, exclusive, autoDelete, durable)

	if err != nil {
		log.Fatal(err)
	}

	rc.Mandatory = envs[mandatory]
	rc.Immediate = envs[immediate]
	rc.Exclusive = envs[exclusive]
	rc.AutoAck = envs[autoAck]
	rc.NoLocal = envs[noLocal]
	rc.NoWait = envs[noWait]
	rc.AutoDelete = envs[autoDelete]
	rc.Durable = envs[durable]
}

func NewConfig() (config Config) {
	config.Logger.Set()
	config.RabbitInfo.Set()

	return config
}
