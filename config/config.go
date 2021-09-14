package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Name     string   `yaml:"name"`
		Host     string   `yaml:"host"`
		Port     string   `yaml:"port"`
		Mode     string   `json:",default=pro,options=dev|test|rt|pre|pro" yaml:"mode"`
		Log      Log      `yaml:"log"`
		Mysql    Mysql    `yaml:"mysql"`
		Throttle Throttle `yaml:"throttle"`
		Mail     Mail     `yaml:"mail"`
	}

	Mysql struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"db_name"`
		DumpPath string `yaml:"dump_path"`

		Prefix                 string `json:",default="`
		SingularTable          bool   `json:",default=false"`
		SetConnMaxLifetime     int    `json:",default=1000"`
		SetConnMaxIdleTime     int    `json:",default=1000"`
		SetMaxOpenConn         int    `json:",default=100"`
		SetMaxIdleConn         int    `json:",default=200"`
		SkipDefaultTransaction bool   `json:",default=true"`
	}

	Log struct {
		File string `yaml:"file"`
	}

	Throttle struct {
		KeyPrefix string `yaml:"key_prefix"`
		Seconds   int    `yaml:"seconds"`
		Quota     int    `yaml:"quota"`
	}

	Mail struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	}
)

type Schedule struct {
	Spec string
	Func string
}

type Schedules struct {
	Data []*Schedule
}

func ParseJsonFile(file string, v interface{}) error {
	fileData, err := readFile(file)

	if err != nil {
		return err
	}

	//解析json
	err = json.Unmarshal(fileData, v)
	if err != nil {
		return err
	}
	return nil
}

func ParseYamlFile(file string, v interface{}) error {
	fileData, err := readFile(file)
	if err != nil {
		return err
	}
	err = yaml.UnmarshalStrict(fileData, v)

	if err != nil {
		return err
	}
	//fmt.Printf("11%#v", v)
	return nil
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	fd, err := ioutil.ReadAll(file)

	return fd, nil
}
