package conf

import (
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var (
	// P : Configuration properties
	P *Properties
)

func init() {
	props, err := mergeProperties()
	if err != nil {
		fmt.Printf("Unable to load config: %v\n", err)
		os.Exit(2)
	}
	P = props
}

func mergeProperties() (*Properties, error) {
	props := &Properties{}
	err := fromFile(props)
	if err != nil {
		return props, err
	}
	err = envconfig.Process("", props)
	if err != nil {
		return props, err
	}
	return props, nil
}

func fromFile(props *Properties) error {
	fh, err := os.Open("conf/defaults.yaml")
	if err != nil {
		return err
	}
	defer fh.Close()
	decoder := yaml.NewDecoder(fh)
	err = decoder.Decode(props)
	if err != nil {
		return err
	}
	return nil
}

func handleError(err error) {
	log.Fatalln(err)
	os.Exit(2)
}
