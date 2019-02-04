package main

import (
	"fmt"
	"github.com/gobuffalo/packr"
	flag "github.com/spf13/pflag"
	"log"
	"os"
	"os/user"
	"pivex/export"
	"pivex/pivotal"
)

var (
	logger    = log.New(os.Stdout, "logger: ", log.Lshortfile)
	credsPath = func() (credsPath string) {
		usr, err := user.Current()
		credsPath = fmt.Sprintf("%s/.pivex", usr.HomeDir)

		if _, err := os.Stat(credsPath); os.IsNotExist(err) {
			os.Mkdir(credsPath, 0700)
		}

		if err != nil {
			logger.Fatal(err)
		}

		return
	}()
)

func main() {
	box := packr.NewBox("./static")

	delTok := flag.BoolP("delete-google-token", "d", false, "delete the generated Google API token being used")
	pivApiTok := flag.StringP(
		"pivotal-token",
		"p",
		"",
		fmt.Sprintf(
			"the Pivotal API token to be used, this token only needs to specified when the token is set for the first time and will be stored under %s after the first time it is set",
			credsPath))
	fCreate := flag.BoolP("force", "f", false, "overwrite an existing presentation")
	presName := flag.StringP("name", "n", "", "name of the presentation (default \"Sprint Demo [SPRINT NUMBER]\")")
	showVer := flag.BoolP("version", "v", false, "show the current version")
	groupEpic := flag.BoolP("group-epic", "e", false, "group sprint stories by epic on the same slide")

	flag.Parse()

	piv := pivotal.New(*pivApiTok, credsPath, logger)
	piv.GetIterations()
	//piv.GetEpics()

	gs := export.New(*presName, credsPath, *fCreate, *groupEpic, logger, piv.Iterations[0])

	if *delTok {
		gs.DelTok()

		os.Exit(0)
	}

	if *showVer {
		ver := box.String("version")

		print("Version: " + ver)

		os.Exit(0)
	}

	gs.Export()
}
