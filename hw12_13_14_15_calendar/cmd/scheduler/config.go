package main

import (
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/cmd"
	"log"
	"os"
	"strconv"
)

type Config struct {
	RabbetInfo RabbitConfig     `yaml:"optQueue"`
	OptAdd     OptAddInQueue    `yaml:"optAdd"`
	Logger     cmd.LoggerConfig `yaml:"logger"`
	DBInfo     cmd.DBConfig     `yaml:"dbInfo"`
}

type RabbitConfig struct {
	Url        string `yaml:"url"`
	Name       string `yaml:"name"`
	NoWait     bool   `yaml:"noWait"`
	Durable    bool   `yaml:"durable"`
	Exclusive  bool   `yaml:"exclusive"`
	AutoDelete bool   `yaml:"autoDelete"`
}

func parseBool(envs ...string) (result map[string]bool, err error) {
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

func (rc *RabbitConfig) Set() {
	exclusive := os.Getenv("EXCLUSIVE")
	autoDelete := os.Getenv("AUTO_DELETE")
	durable := os.Getenv("DURABLE")
	noWait := os.Getenv("NO_WAIT")

	rc.Url = os.Getenv("URL")
	rc.Name = os.Getenv("NAME_Q")

	envs, err := parseBool(exclusive, autoDelete, durable, noWait)

	if err != nil {
		log.Fatal(err)
	}

	rc.Exclusive = envs[exclusive]
	rc.AutoDelete = envs[autoDelete]
	rc.Durable = envs[durable]
	rc.NoWait = envs[noWait]
}

type OptAddInQueue struct {
	exchange  string `yaml:"exchange"`
	mandatory bool   `yaml:"mandatory"`
	immediate bool   `yaml:"immediate"`
}

func (opt *OptAddInQueue) Set() {
	opt.exchange = os.Getenv("EXCHANGE")

	mandatory := os.Getenv("MANDATORY")
	immediate := os.Getenv("IMMEDIATE")

	envs, err := parseBool(mandatory, immediate)

	if err != nil {
		log.Fatal(err)
	}

	opt.mandatory = envs[mandatory]
	opt.immediate = envs[immediate]
}

func NewConfig() (config Config) {
	config.Logger.Set()
	config.OptAdd.Set()
	config.RabbetInfo.Set()

	return config
}
