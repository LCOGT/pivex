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
	delAuth := flag.Bool("d", false, "delete the authentication files being used for pivex")
	pivApiTok := flag.String(
		"p",
		"",
		fmt.Sprintf(
			"the Pivotal API token to be used, this token only needs to specified when the token is set for the first time and will be stored under %s after the first time it is set",
			credsPath))
	flag.Parse()

	piv := pivotal.New(*pivApiTok, logger)
	gs := export.New(credsPath, logger)

	piv.GetStories()

	if *delAuth {
		gs.DelAuth()
	}

	gs.Export(&piv.Intervals[0])
}
