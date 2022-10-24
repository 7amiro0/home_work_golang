package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBInfo  DBConfig     `yaml:"dbInfo"`
	HTTP    HTTPConfig   `yaml:"http"`
	GRPC    GRPCConfig   `yaml:"grpc"`
	Logger  LoggerConfig `yaml:"logger"`
	Storage string       `yaml:"storage"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

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

func NewConfig(path string) (config Config) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	config.DBInfo.SetEnv()

	return config
}

func (db *DBConfig) SetEnv() {
	if err := os.Setenv("DATABASE_USER", db.User); err != nil {
		fmt.Println(err)
	} else if err := os.Setenv("DATABASE_PORT", db.Port); err != nil {
		fmt.Println(err)
	} else if err := os.Setenv("DATABASE_HOST", db.Host); err != nil {
		fmt.Println(err)
	} else if err := os.Setenv("DATABASE_PASSWORD", db.Password); err != nil {
		fmt.Println(err)
	} else if err := os.Setenv("DATABASE_NAME", db.Name); err != nil {
		fmt.Println(err)
	}
}
