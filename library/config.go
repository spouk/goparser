package library

import (
	"io/ioutil"
	"os"
	"time"
	"github.com/go-yaml/yaml"
)

type (
	ConfigTable struct {
		LogfileName string        `yaml:"logfilename"`
		WorkRequest int           `yaml:"workrequest"`
		WorkPause   time.Duration `yaml:"workpause"`
		FileRequest string        `yaml:"filerequest"`
		RequestFile string        `yaml:"requestfile"`
		Database    string        `yaml:"database"`
		Pathsave    string        `yaml:"pathsave"`
	}
)

func (p *Parser) newConfigRead(filenameConfig string) (error, *ConfigTable) {

	//open file for read config
	f, err := os.Open(filenameConfig)

	//if error open return err
	if err != nil {
		return err, nil
	}
	//close  handler open file after exit function
	defer f.Close()

	//read file data
	b, err := ioutil.ReadAll(f)

	//return if error read file
	if err != nil {
		return err, nil
	}
	//make new instanse ConfigTable
	ct := new(ConfigTable)

	//parse config file
	err = yaml.Unmarshal(b, ct)

	//if error return error
	if err != nil {
		return err, nil
	}

	//return success make ConfigTable and nil error
	return nil, ct

}
