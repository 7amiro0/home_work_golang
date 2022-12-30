package main

import (
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/cmd"
	"os"
)

type Config struct {
	DBInfo  cmd.DBConfig     `yaml:"dbInfo"`
	HTTP    HTTPConfig       `yaml:"http"`
	GRPC    GRPCConfig       `yaml:"grpc"`
	Logger  cmd.LoggerConfig `yaml:"logger"`
	Storage string           `yaml:"storage"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (h *HTTPConfig) Set() {
	h.Host = os.Getenv("HTTP_HOST")
	h.Port = os.Getenv("HTTP_PORT")
}

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (c *GRPCConfig) Set() {
	c.Host = os.Getenv("GRPC_HOST")
	c.Port = os.Getenv("GRPC_PORT")
}

func NewConfig() (config Config) {
	config.HTTP.Set()
	config.GRPC.Set()
	config.Logger.Set()

	return config
}
