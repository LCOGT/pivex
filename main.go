package main

import (
	"fmt"
	"github.com/gobuffalo/packr"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"os"
	"pivex/credentials"
	"pivex/export"
	"pivex/pivotal"
)

var (
	logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
)

func main() {
	box := packr.NewBox("./static")

	pivCreds := credentials.NewPivotal(logger)
	pivApiToken := flag.StringP(
		"pivotal-api-token",
		"p",
		"",
		fmt.Sprintf(
			"exports the given Pivotal API token, this token will be stored under %s and overwrite any existing token",
			pivCreds.Path))

	gSlideCreds := credentials.NewGoogleSlides(logger)
	oauthClientIdFile := flag.StringP(
		"google-token",
		"g",
		"exports the given Google API token JSON, this token will be stored under %s and overwrite any existing token",
		fmt.Sprintf("the file containing the Google API token JSON",
			gSlideCreds.Path))
	print(oauthClientIdFile)

	fCreate := flag.BoolP("force", "f", false, "overwrite an existing presentation")
	showVer := flag.BoolP("version", "v", false, "show the current version")

	flag.Parse()

	pivCreds.ApiToken = *pivApiToken

	pivCredsErr := pivCreds.Init()

	if pivCredsErr != nil {
		logger.Panicln(pivCredsErr)
	}

	piv := pivotal.New(pivCreds, logger)
	piv.GetIterations()

	gs := export.New(gSlideCreds, *fCreate, logger, piv.Iterations[0])

	if *showVer {
		ver := box.String("version")

		print("Version: " + ver)

		os.Exit(0)
	}

	gs.Export()
}

func copyFile(sourceFile string, destinationFile string) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}
}
