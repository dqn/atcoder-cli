package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/dqn/atcoder-cli"
	"gopkg.in/yaml.v2"
)

func readConfig(path string) (*atcoder.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var c atcoder.Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func run() error {
	c, err := readConfig("./atcoderrc.yml")
	if err != nil {
		return err
	}
	client := atcoder.New(c)

	flag.Parse()
	switch flag.Arg(0) {
	case "init":
		if flag.NArg() != 2 {
			flag.Usage()
			os.Exit(2)
		}
		return client.Init(flag.Arg(1))
	case "test":
		if flag.NArg() != 3 {
			flag.Usage()
			os.Exit(2)
		}
		_, err = client.Test(flag.Arg(1), flag.Arg(2))
		return err
	case "submit":
		var contest, task string
		var testFlag bool
		switch flag.NArg() {
		case 3:
			contest, task = flag.Arg(1), flag.Arg(2)
		case 4:
			println(flag.Arg(1), flag.Arg(2), flag.Arg(3))
			switch {
			case flag.Arg(1) == "-t":
				contest, task = flag.Arg(2), flag.Arg(3)
				testFlag = true
			case flag.Arg(3) == "-t":
				contest, task = flag.Arg(1), flag.Arg(2)
				testFlag = true
			}
		}
		if contest == "" || task == "" {
			flag.Usage()
			os.Exit(2)
		}

		if testFlag {
			ok, err := client.Test(contest, task)
			if err != nil {
				return err
			}
			if !ok {
				os.Exit(1)
			}
		}

		if err = client.Login(); err != nil {
			return err
		}
		id, err := client.Submit(contest, task)
		if err != nil {
			return err
		}
		return client.WaitJudge(contest, id)
	default:
		flag.Usage()
		os.Exit(2)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
