package main

import (
	"pivex/pivotal"
	"pivex/export"
	"flag"
	"fmt"
	"os"
	"log"
	"os/user"
)

var (
	logger    = log.New(os.Stdout, "logger: ", log.Lshortfile)
	credsPath = func() (credsPath string) {
		usr, err := user.Current()
		credsPath = fmt.Sprintf("%s/.pivex", usr.HomeDir)

		if _, err := os.Stat(credsPath); os.IsNotExist(err) {
			os.Mkdir(credsPath, 0600)
		}

		if err != nil {
			logger.Fatal(err)
		}

		return
	}()
)

func main() {
	delTok := flag.Bool("d", false, "delete the generated Google API token being used")
	pivApiTok := flag.String(
		"p",
		"",
		fmt.Sprintf(
			"the Pivotal API token to be used, this token only needs to specified when the token is set for the first time and will be stored under %s after the first time it is set",
			credsPath))
	fCreate := flag.Bool("f", false, "overwrite an existing presentation")
	flag.Parse()

	piv := pivotal.New(*pivApiTok, credsPath, logger)
	gs := export.New(credsPath, *fCreate, logger)

	if *delTok {
		gs.DelTok()

		os.Exit(0)
	}

	piv.GetStories()

	gs.Export(&piv.Intervals[0])
}
