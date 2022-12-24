package main

import (
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/cmd"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	year  = 365 * 24 * time.Hour
	mouth = 30 * 24 * time.Hour
	day   = 24 * time.Hour
)

type Config struct {
	RabbitInfo RabbitConfig     `yaml:"optQueue"`
	OptAdd     OptAddInQueue    `yaml:"optAdd"`
	Logger     cmd.LoggerConfig `yaml:"logger"`
	DBInfo     cmd.DBConfig     `yaml:"dbInfo"`
	Ticker     Ticker           `yaml:"ticker"`
}

type Ticker struct {
	Add   time.Duration `yaml:"add"`
	Clear time.Duration `yaml:"clear"`
}

func convertToDuration(per []string) time.Duration {
	yInt, _ := strconv.Atoi(strings.Split(per[0], "y")[0])
	mInt, _ := strconv.Atoi(strings.Split(per[1], "m")[0])
	dInt, _ := strconv.Atoi(strings.Split(per[2], "d")[0])
	hInt, _ := strconv.Atoi(strings.Split(per[3], "h")[0])
	minInt, _ := strconv.Atoi(strings.Split(per[4], "min")[0])
	sInt, _ := strconv.Atoi(strings.Split(per[5], "s")[0])

	y := time.Duration(yInt) * year
	m := time.Duration(mInt) * mouth
	d := time.Duration(dInt) * day
	h := time.Duration(hInt) * time.Hour
	min := time.Duration(minInt) * time.Minute
	s := time.Duration(sInt) * time.Second

	return y + m + d + h + min + s
}

func (t *Ticker) Set() {
	addSlice := strings.Split(os.Getenv("ADD"), " ")
	clearSlice := strings.Split(os.Getenv("CLEAR"), " ")
	t.Add = convertToDuration(addSlice)
	t.Clear = convertToDuration(clearSlice)
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
	config.RabbitInfo.Set()
	config.Ticker.Set()

	return config
}
