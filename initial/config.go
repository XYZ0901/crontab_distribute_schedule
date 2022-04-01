package initial

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func configInit() {
	f, err := os.Open("./config/config.yaml")
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Cfg)
	if err != nil {
		log.Fatalln(err)
	}
}
