package main

import (
	"io/ioutil"
	"os"

	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/Xuanwo/beancollect/bean"
	"github.com/Xuanwo/beancollect/collect"
	"github.com/Xuanwo/beancollect/constants"
	"github.com/Xuanwo/beancollect/types"
)

func run(c *cli.Context) (err error) {
	log.SetLevel(log.DebugLevel)

	if c.NArg() == 0 {
		println("Please specify the file to collect")
		os.Exit(1)
	}

	fi, err := ioutil.ReadDir("collect")
	if err != nil {
		println("Please collect from your bean main directory")
		os.Exit(1)
	}

	schema := c.String("schema")
	if schema == "" {
		println("Please specify the schema to collect")
		os.Exit(1)
	}

	globalConfig, schemaConfig := &types.Config{}, &types.Config{}

	for _, v := range fi {
		if v.Name() == "global.yaml" {
			cfgContent, err := ioutil.ReadFile("collect/global.yaml")
			if err != nil {
				log.Errorf("Open config failed for %v", err)
				return err
			}
			err = yaml.Unmarshal(cfgContent, globalConfig)
			if err != nil {
				log.Errorf("Load config failed for %v", err)
				return err
			}
			continue
		}
		if v.Name() == schema+".yaml" {
			cfgContent, err := ioutil.ReadFile("collect/" + schema + ".yaml")
			if err != nil {
				log.Errorf("Open config failed for %v", err)
				return err
			}
			err = yaml.Unmarshal(cfgContent, schemaConfig)
			if err != nil {
				log.Errorf("Load config failed for %v", err)
				return err
			}
			continue
		}
		continue
	}

	if err := mergo.Merge(schemaConfig, globalConfig, mergo.WithOverride); err != nil {
		log.Errorf("Config merge failed for %v", err)
		return err
	}

	var t types.Transactions

	f, err := os.Open(c.Args().Get(0))
	if err != nil {
		log.Errorf("Open file failed for %v", err)
		return err
	}
	defer f.Close()

	t, err = collect.NewCollector(schema).Parse(schemaConfig, f)
	if err != nil {
		log.Errorf("Parse failed for %v", err)
		return err
	}

	bean.Generate(schemaConfig, &t)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = constants.Name
	app.Usage = constants.Usage
	app.Version = constants.Version
	app.Action = run

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "schema, s",
			Usage: "schema for the collect",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
