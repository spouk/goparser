package lib

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"log"
	"io/ioutil"
)

const (
	LOG_CONFIG_PREFIX = "[config] "
	LOG_CONFIG_FLAGS  = log.Ldate | log.Ltime | log.Lshortfile
)

type Config struct {
	//база данных
	DB string `yaml:"db"`
	//логфайл
	LogfileName string `yaml:"logfilename"`
	LogFile     io.Writer
	//количество работающих воркеров
	CountWorkerRequest int `yaml:"workrequest"`
	CountWorkerFile    int `yaml:"workfile"`
	//имя файла где указаны запросы
	FileRequest string `yaml:"filerequest"`
	configfile  io.Writer
	//logger
	log *log.Logger
	//request file
	RequestFile string `yaml:"requestfile"`
}

func NewConfig(configFileName string, logout io.Writer) *Config {
	//создаю инстанс
	c := &Config{}
	if logout != nil {
		c.log = log.New(logout, LOG_CONFIG_PREFIX, LOG_CONFIG_FLAGS)
	}
	//открываю файл с конфигом для чтения
	f, err := os.Open(configFileName)
	if err != nil {
		c.log.Printf(err.Error())
		panic(err)
	}
	c.configfile = f
	b, err := ioutil.ReadAll(f)
	if err != nil {
		c.log.Printf(err.Error())
		panic(err)
	}

	//парсю файл
	err = yaml.Unmarshal(b, c)
	if err != nil {
		c.log.Printf(err.Error())
		panic(err)
	}
	//если указан логфайл в конфиге то открываю файл лога
	if c.LogfileName != "" {
		f, err = os.OpenFile(c.LogfileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			c.log.Printf(err.Error())
		}
		c.LogFile = f
		//изменяю вывод логирования в файл
		c.log = log.New(c.LogFile, LOG_CONFIG_PREFIX, LOG_CONFIG_FLAGS)
	}
	return c
}
